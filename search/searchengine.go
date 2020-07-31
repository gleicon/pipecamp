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

/*
MetaDocument holds the meta document info and summary
*/
type MetaDocument struct {
	ID      string `json:"id"` // ID is the document path
	Body    string `json:"body"`
	Summary string `json:"summary"`
}

// Result is a query result holder
type Result struct {
	ID    string  `json:"id"`
	Score float64 `json:"score"`
}

// QueryResults holds a slice of returns
type QueryResults struct {
	Query    string   `json:"query"`
	CutScore float64  `json:"cutscore"`
	Results  []Result `json:"results"`
}

/*
SearchEngine represents a generic search engine wrapping bleve
*/
type SearchEngine struct {
	index      bleve.Index
	indexPath  string
	summarizer *summarizer.PersistentSummarizer
}

// NewSearchEngine creates a new search engine
func NewSearchEngine(indexPath string, summarizer *summarizer.PersistentSummarizer) *SearchEngine {
	se := SearchEngine{}
	se.CreateOrOpenIndex(indexPath)
	se.summarizer = summarizer
	return &se
}

// CreateOrOpenIndex sets up the search index structure
func (se *SearchEngine) CreateOrOpenIndex(indexPath string) error {
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
func (se *SearchEngine) CloseIndex() {
	se.index.Close()
}

func (se *SearchEngine) docSearch(input string) (*bleve.SearchResult, error) {

	/*query := bleve.NewBooleanQuery(nil, nil, nil)
	for _, term := range strings.Split(input, " ") {
		query.AddShould(bleve.NewFuzzyQuery(term))
		query.AddShould(bleve.NewTermQuery(term))
	}*/

	query := bleve.NewMatchQuery(input)
	searchRequest := bleve.NewSearchRequest(query)
	fmt.Println(searchRequest)
	return se.index.Search(searchRequest)

}

// Query , fetch results, rank and aggregate
func (se *SearchEngine) Query(query string) (*QueryResults, error) {
	searchResult, err := se.docSearch(query)
	if err != nil {
		return nil, err
	}
	results := []Result{}
	for _, doc := range searchResult.Hits {
		if doc.Score > 0.1 {
			results = append(results, Result{ID: doc.ID, Score: doc.Score})
		}
	}
	res := QueryResults{Query: query, Results: results, CutScore: 0.1}
	return &res, nil
}

func (se *SearchEngine) skipName(name string) bool {
	return strings.HasSuffix(name, ".jpg") ||
		strings.HasSuffix(name, ".jpeg") ||
		strings.HasSuffix(name, ".png") ||
		strings.HasSuffix(name, ".gif") ||
		strings.HasSuffix(name, ".gifv")
}

/*
AddDocuments creates the needed meta info around a file, reads it and holds its summary
*/
func (se *SearchEngine) AddDocuments(filePath string) error {
	err := filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if se.skipName(path) {
			return nil
		}
		mf := MetaDocument{}
		mf.ID = path

		fmt.Printf("\radding '\x1b[34;1m%v\x1b[0m' ... ", mf.ID)

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

func (se *SearchEngine) indexDocument(mf *MetaDocument) error {

	body, err := ioutil.ReadFile(mf.ID)

	if err != nil {
		return err
	}
	// read and fill body
	mf.Body = string(body)
	// create a summary
	if mf.Summary, err = se.summarizer.SummarizeAndStore(mf.ID, mf.Body); err != nil {
		return err
	}
	// index
	if err := se.index.Index(mf.ID, mf); err != nil {
		return err
	}

	fmt.Println("OK")
	return nil
}
