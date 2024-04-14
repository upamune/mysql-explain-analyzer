package output

type Result struct {
	Tables []Table

	FullTableScanTables  []Table
	FullIndexScanTables  []Table
	HasAnyCommentsTables []Table
}

type Table struct {
	Name         string
	AccessType   string
	PossibleKeys []string
	Key          *string
	KeyLength    *int
	Ref          []string
	Rows         int
	Filtered     float64
	Scalability  string

	IsFullTableScans bool
	IsFullIndexScans bool
	Comment          string
}
