package input

type Explain struct {
	QueryBlock QueryBlock `json:"query_block"`
}

type QueryBlock struct {
	SelectID          int               `json:"select_id"`
	CostInfo          TotalCostInfo     `json:"cost_info"`
	OrderingOperation OrderingOperation `json:"ordering_operation"`
	Table             Table             `json:"table"`
}

type TotalCostInfo struct {
	QueryCost string `json:"query_cost"`
}

type OrderingOperation struct {
	UsingTemporaryTable bool              `json:"using_temporary_table"`
	UsingFilesort       bool              `json:"using_filesort"`
	DuplicatesRemoval   DuplicatesRemoval `json:"duplicates_removal"`
	SortCostInfo        SortCostInfo      `json:"cost_info"`
	NestedLoop          []NestedLoop      `json:"nested_loop"`
}

type DuplicatesRemoval struct {
	UsingTemporaryTable bool         `json:"using_temporary_table"`
	UsingFilesort       bool         `json:"using_filesort"`
	NestedLoop          []NestedLoop `json:"nested_loop"`
}

type NestedLoop struct {
	Table Table `json:"table"`
}

type SortCostInfo struct {
	SortCost string `json:"sort_cost"`
}

type CostInfo struct {
	ReadCost        string `json:"read_cost"`
	EvalCost        string `json:"eval_cost"`
	PrefixCost      string `json:"prefix_cost"`
	DataReadPerJoin string `json:"data_read_per_join"`
}

type Table struct {
	TableName           string   `json:"table_name"`
	AccessType          string   `json:"access_type"`
	PossibleKeys        []string `json:"possible_keys"`
	Key                 string   `json:"key"`
	UsedKeyParts        []string `json:"used_key_parts"`
	KeyLength           string   `json:"key_length"`
	RowsExaminedPerScan int      `json:"rows_examined_per_scan"`
	RowsProducedPerJoin int      `json:"rows_produced_per_join"`
	Filtered            string   `json:"filtered"`
	CostInfo            CostInfo `json:"cost_info"`
	UsedColumns         []string `json:"used_columns"`
	AttachedCondition   string   `json:"attached_condition"`
	Ref                 []string `json:"ref"`
}
