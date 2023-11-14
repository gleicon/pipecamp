package searchengine

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gleicon/pipecamp/summarizer"

	"github.com/blevesearch/bleve"
)

type BleveSearchEngine struct {
	index      bleve.Index
	indexPath  string
	summarizer *summarizer.PersistentSummarizer
}

// NewSearchEngine creates a new search engine
func NewBleveSearchEngine(indexPath string, summarizer *summarizer.PersistentSummarizer) *BleveSearchEngine {
	se := BleveSearchEngine{}
	se.CreateOrOpenIndex(indexPath)
	se.summarizer = summarizer
	return &se
}

// CreateOrOpenIndex sets up the search index structure
func (se *BleveSearchEngine) CreateOrOpenIndex(indexPath string) error {
	mapping := bleve.NewIndexMapping()
	var err error
	se.indexPath = indexPath
	se.index, err = bleve.New(indexPath, mapping)

	if err == bleve.ErrorIndexPathExists {
		se.index, err = bleve.Open(indexPath)
	}
	return err
}

// CloseIndex safely closes an open index
func (se *BleveSearchEngine) Close() {
	se.index.Close()
}

func (se *BleveSearchEngine) docSearch(input string) (*bleve.SearchResult, error) {

	query := bleve.NewMatchQuery(input)
	searchRequest := bleve.NewSearchRequest(query)
	fmt.Println(searchRequest)
	return se.index.Search(searchRequest)

}

// Query , fetch results, rank and aggregate
func (se *BleveSearchEngine) Query(query string) (*QueryResults, error) {
	searchResult, err := se.docSearch(query)
	if err != nil {
		return nil, err
	}
	results := []Result{}
	for _, doc := range searchResult.Hits {
		if doc.Score > 0.1 {
			results = append(results, Result{Id: doc.ID, Score: doc.Score})
		}
	}
	res := QueryResults{Query: query, Results: results, CutScore: 0.1}
	return &res, nil
}

func (se *BleveSearchEngine) skipName(name string) bool {
	return strings.HasSuffix(name, ".jpg") ||
		strings.HasSuffix(name, ".jpeg") ||
		strings.HasSuffix(name, ".png") ||
		strings.HasSuffix(name, ".gif") ||
		strings.HasSuffix(name, ".gifv")
}

/*
AddDocuments creates the needed meta info around a file, reads it and holds its summary
*/
func (se *BleveSearchEngine) AddDocuments(filePath string) error {
	err := filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if se.skipName(path) {
			return nil
		}
		mf := MetaDocument{}
		mf.Id = path

		fmt.Printf("\radding '\x1b[34;1m%v\x1b[0m' ... ", mf.Id)

		if err := se.indexDocument(&mf); err != nil {
			return err
		}
		return nil
	})
	if err != nil {

		fmt.Printf("Error indexing: %#v\n", err)
	}
	return nil
}

func (se *BleveSearchEngine) indexDocument(mf *MetaDocument) error {

	body, err := ioutil.ReadFile(mf.Id)

	if err != nil {
		return err
	}
	// read and fill body
	mf.Body = string(body)
	// create a summary
	if mf.Summary, err = se.summarizer.SummarizeAndStore(mf.Id, mf.Body); err != nil {
		return err
	}
	// index
	if err := se.index.Index(mf.Id, mf); err != nil {
		return err
	}

	fmt.Println("OK")
	return nil
}
