package summarizer

import (
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/JesusIslam/tldr"
	"github.com/dgraph-io/badger"
)

// PersistentSummarizer wraps a K/V db and stores/retrieves summarized text snippets
type PersistentSummarizer struct {
	dbpath        string
	sentenceCount int
	db            *badger.DB
}

// NewPersistentSummarizer creates a persistent summarizer
func NewPersistentSummarizer(dbpath string, sentenceCount int) (*PersistentSummarizer, error) {
	var err error
	psz := PersistentSummarizer{dbpath: dbpath, sentenceCount: sentenceCount}
	opts := badger.DefaultOptions(dbpath)
	opts.SyncWrites = true
	opts.Dir = dbpath
	psz.db, err = badger.Open(opts)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			log.Println("Summarizer Database housekeeping")
			var ll sync.Mutex
			ll.Lock()
			psz.db.RunValueLogGC(1.0)
			ll.Unlock()
			time.Sleep(10 * time.Minute)
		}
	}()

	return &psz, nil
}

// SummarizeAndStore uses a key to point to a summary on disk
func (psz *PersistentSummarizer) SummarizeAndStore(key string, value string) (string, error) {
	var summary string
	var err error
	if summary, err = Summarize(value, psz.sentenceCount); err != nil {
		return "", errors.New("Error summarizing text: " + err.Error())
	}

	err = psz.db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), []byte(summary))
		if err != nil {
			return err
		}
		return nil
	})
	return summary, err
}

// Fetch a summary based on key, returns error if the key doesn't exist
func (psz *PersistentSummarizer) Fetch(key string) (string, error) {
	var payload []byte
	err := psz.db.View(func(txn *badger.Txn) error {

		item, err := txn.Get([]byte(key))
		if err == badger.ErrKeyNotFound {
			return nil
		}
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}
		err = item.Value(func(val []byte) error {
			payload = val
			return nil

		})
		return err
	})
	return string(payload), err
}

// Summarize - standalone lexrank summarization wrapper
func Summarize(body string, sentenceCount int) (string, error) {
	bag := tldr.New()
	summary, err := bag.Summarize(body, sentenceCount)
	if err != nil {
		return "", err
	}
	summ := strings.Join(summary, " ")
	return summ, nil
}
