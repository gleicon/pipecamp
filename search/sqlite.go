package searchengine

import (
	"database/sql"
	"log"
	"os"

	"github.com/gleicon/pipecamp/summarizer"
)

type SQLiteSearchEngine struct {
	dbHandler  *sql.DB
	dbPath     string
	summarizer *summarizer.PersistentSummarizer
}

func NewSQLiteSearchEngine(dbPath string, summarizer *summarizer.PersistentSummarizer) *SQLiteSearchEngine {
	se := SQLiteSearchEngine{}
	se.dbPath = dbPath
	se.CreateOrOpenIndex()
	se.summarizer = summarizer
	return &se
}

func (se *SQLiteSearchEngine) CreateOrOpenIndex() error {
	var err error
	// if the db doesn't exists create it and setup schema
	// if the db exists, open and test if the fts5 vtable exists
	if _, err := os.Stat(se.dbPath); err == nil {
		se.dbHandler, err = sql.Open("sqlite3", se.dbPath)
		if err != nil {
			return err
		}
		return se.testFTS5Schema()
	}
	// create and apply schema
	se.dbHandler, err = sql.Open("sqlite3", se.dbPath)
	if err == nil {
		return err
	}
	return se.createFTS5Schema()
}

func (se *SQLiteSearchEngine) testFTS5Schema() error {
	sqlStmt := "SELECT * FROM documents LIMIT 1"
	_, err := se.dbHandler.Query(sqlStmt)

	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return err
	}
	return nil
}

func (se *SQLiteSearchEngine) createFTS5Schema() error {

	sqlStmt := "CREATE VIRTUAL TABLE documents USING fts5(body, summary)"
	_, err := se.dbHandler.Query(sqlStmt)

	//rows, err := se.dbHandler.Query("SELECT * FROM user WHERE id = ?", id)
	//_, err = db.Exec(sqlStmt)

	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return err
	}
	return nil
}

// CloseIndex safely closes an open index
func (se *SQLiteSearchEngine) Close() {
	se.dbHandler.Close()
}

func (se *SQLiteSearchEngine) docSearch(input string) (string, error) {
	return "", nil
}

func (se *SQLiteSearchEngine) Query(query string) (*QueryResults, error) {
	return nil, nil
}
