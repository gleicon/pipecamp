package searchengine

type SearchEngine interface {
	CreateOrOpenIndex(string) error
	Close()
	docSearch()
	Query()
	skipName()
	AddDocuments()
	indexDocument()
}

/*
MetaDocument holds the meta document info and summary to be indexed
*/
type MetaDocument struct {
	Id      string `json:"id"`
	Path    string `json:"path"`
	Body    string `json:"body"`
	Summary string `json:"summary"`
}

// Result is a query result holder, using id as reference and optionally the summary
type Result struct {
	Id      string  `json:"id"`
	Score   float64 `json:"score"`
	Summary string  `json:"summary"`
}

// QueryResults holds a slice of returns
type QueryResults struct {
	Query    string   `json:"query"`
	CutScore float64  `json:"cutscore"`
	Results  []Result `json:"results"`
}
