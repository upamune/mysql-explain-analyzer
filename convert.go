package main

import (
	"fmt"
	"strconv"
	"strings"

	"mysql-explain-analyzer/model/input"
	"mysql-explain-analyzer/model/output"

	"github.com/hatena/godash"
)

func convert(input input.Explain) output.Result {
	var tables []output.Table
	if t := input.QueryBlock.Table; t.TableName != "" {
		tables = append(tables, convertTable(t))
	}
	if nestedLoop := input.QueryBlock.OrderingOperation.DuplicatesRemoval.NestedLoop; len(nestedLoop) > 0 {
		tables = append(tables, convertNestedLoopToTables(nestedLoop)...)
	}
	if nestedLoop := input.QueryBlock.OrderingOperation.NestedLoop; len(nestedLoop) > 0 {
		tables = append(tables, convertNestedLoopToTables(nestedLoop)...)
	}

	return output.Result{
		Tables: tables,
		FullTableScanTables: godash.Filter(tables, func(t output.Table, _ int) bool {
			return t.IsFullTableScans
		}),
		FullIndexScanTables: godash.Filter(tables, func(t output.Table, _ int) bool {
			return t.IsFullIndexScans
		}),
		HasAnyCommentsTables: godash.Filter(tables, func(t output.Table, _ int) bool {
			return t.Comment != ""
		}),
		Comment: analyzeOrderingOperationComment(input.QueryBlock.OrderingOperation),
	}
}

func analyzeOrderingOperationComment(orderingOperation input.OrderingOperation) string {
	var builder strings.Builder
	if orderingOperation.UsingTemporaryTable {
		_, _ = builder.WriteString("ソートのために一時テーブルを使用しています。データが設定したメモリのバッファーサイズに収まらないため非常に遅くなります。")
	}
	if orderingOperation.UsingFilesort {
		_, _ = builder.WriteString("ソートのためにファイルソートを使用しています。インデックスを用いてソートが行われていません。")
	}
	return builder.String()
}

func convertNestedLoopToTables(nestedLoops []input.NestedLoop) []output.Table {
	var tables []output.Table
	for _, nestedLoop := range nestedLoops {
		tables = append(tables, convertTable(nestedLoop.Table))
	}
	return tables
}

func convertTable(table input.Table) output.Table {
	ref := table.Ref
	if len(ref) == 0 {
		ref = nil
	}

	isFullTableScan, isFullIndexScan, comment := analyzeAccessType(table)
	return output.Table{
		Name:             table.TableName,
		AccessType:       table.AccessType,
		PossibleKeys:     table.PossibleKeys,
		Key:              stringPtr(table.Key),
		KeyLength:        fromStringToIntP(table.KeyLength),
		Ref:              ref,
		Rows:             table.RowsExaminedPerScan,
		Filtered:         fromStringToFloat64(table.Filtered),
		Scalability:      fmt.Sprintf("O(%s)", fromAccessTypeToScalability(table.AccessType)),
		IsFullTableScans: isFullTableScan,
		IsFullIndexScans: isFullIndexScan,
		Comment:          comment,
	}
}

const (
	AccessTypeAll        = "ALL"
	AccessTypeIndex      = "index"
	AccessTypeRef        = "ref"
	AccessTypeEqRef      = "eq_ref"
	AccessTypeConst      = "const"
	AccessTypeRange      = "range"
	AccessTypeIndexMerge = "index_merge"
)

func analyzeAccessType(t input.Table) (isFullTableScan bool, isFullIndexScan bool, comment string) {
	switch t.AccessType {
	case AccessTypeAll:
		isFullTableScan = true
	case AccessTypeIndex:
		if t.Key == "PRIMARY" {
			// InnoDB上のプライマリキーインデックススキャンは全テーブルスキャンと同等
			isFullTableScan = true
		} else {
			isFullIndexScan = true
		}
	case AccessTypeRange:
		comment = "インデックスを使って一定範囲の行にアクセスしています。"
	case AccessTypeRef:
		comment = "インデックスを使ってマッチする行にアクセスしています。"
	case AccessTypeEqRef:
		comment = "インデックスを使って最大1行にアクセスしています。"
	case AccessTypeConst:
		comment = "このテーブルは問い合わせの最初に1回読み込まれ、定数として扱われています。"
	case AccessTypeIndexMerge:
		comment = "MySQLが複数のインデックスを使用しています。"
	}

	if t.Key != "" {
		if t.Key == "PRIMARY" {
			comment += "MySQLがPRIMARY KEYを使用しています。"
		} else if t.AccessType == AccessTypeIndexMerge {
			comment += "使用されているインデックスとマージタイプは " + t.Key + " です。"
			if strings.HasPrefix(strings.ToLower(t.Key), "intersect") {
				comment += "マージタイプが'intersect'のようです。WHERE句が複雑な場合、遅くなる可能性があります!"
			}
		} else {
			comment += "MySQL が '" + t.Key + "' インデックスを使用しています。"
		}
	}

	return isFullTableScan, isFullIndexScan, comment
}

func fromAccessTypeToScalability(accessType string) string {
	switch accessType {
	case AccessTypeAll:
		return "N"
	case AccessTypeIndex:
		return "N"
	case AccessTypeRef:
		return "log N"
	case AccessTypeEqRef:
		return "log N"
	case AccessTypeConst:
		return "1"
	default:
		return "?"
	}
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func fromStringToFloat64(s string) float64 {
	return fromString(s, func(s string) (float64, error) {
		return strconv.ParseFloat(s, 64)
	})
}

func fromStringToIntP(s string) *int {
	return fromString(s, func(s string) (*int, error) {
		i, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		return &i, nil
	})
}

// fromString converts a string to a type T using the provided function.
// If the conversion fails, it returns the zero value of T.
func fromString[T any](s string, from func(s string) (T, error)) T {
	v, err := from(s)
	if err != nil {
		var v2 T
		return v2
	}
	return v
}
