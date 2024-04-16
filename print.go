package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"mysql-explain-analyzer/model/output"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

const (
	purple = lipgloss.Color("99")
	gray   = lipgloss.Color("245")
)

var descStyle = lipgloss.NewStyle()

func printTitle(w io.Writer, title string) {
	var style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		PaddingLeft(4).
		MarginTop(1).
		Width(len(title) + 8)
	fmt.Fprintln(w, style.Render(title))
}

func printInfo(w io.Writer, s string) {
	subtle := lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	infoStyle := lipgloss.NewStyle().
		Bold(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(subtle)
	fmt.Fprintln(w, infoStyle.Render(s))
}

func printDescription(w io.Writer, s string) {
	fmt.Fprintln(w, descStyle.Render(s))
}

func printExplainTable(w io.Writer, res output.Result) {
	printTitle(w, "Explain table")
	tableHeaders := []string{
		"Table",
		"Access type",
		"Possible indexes",
		"Index",
		"Index key length",
		"Ref",
		"Rows examined per scan",
		"Filtered",
		"Scalability",
	}

	var tableRows [][]string
	for _, t := range res.Tables {
		tableRows = append(tableRows, []string{
			t.Name,
			t.AccessType,
			strings.Join(t.PossibleKeys, ", "),
			strPtrToString(t.Key),
			intPtrToString(t.KeyLength),
			strings.Join(t.Ref, ", "),
			strconv.Itoa(t.Rows),
			fmt.Sprintf("%.2f", t.Filtered),
			t.Scalability,
		})
	}

	re := lipgloss.NewRenderer(w)
	var (
		HeaderStyle = re.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
		CellStyle   = re.NewStyle().Padding(0, 1).Width(14)
		RowStyle    = CellStyle.Copy().Foreground(gray)
		BorderStyle = lipgloss.NewStyle().Foreground(purple)
	)
	t := table.New().
		Width(180).
		Border(lipgloss.NormalBorder()).
		BorderRow(true).
		BorderStyle(BorderStyle).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == 0:
				return HeaderStyle
			default:
				return RowStyle
			}
		}).
		Headers(tableHeaders...).
		Rows(tableRows...)

	fmt.Fprintln(w, t)
}

func printFullTableScanComments(w io.Writer, output output.Result) {
	printTitle(w, "Are there any full table scans?")

	tables := output.FullTableScanTables
	if len(tables) == 0 {
		printInfo(w, "No")
		return
	}
	printInfo(w, "Yes")
	printDescription(w, "以下のテーブルに対して全テーブルスキャンが行われています。MySQLはこれらのテーブルの全ての行を読み込んでいます。")

	for _, t := range tables {
		var doc strings.Builder
		doc.WriteString(descStyle.Render("- Table `"))
		doc.WriteString(descStyle.Bold(true).Render(t.Name))
		doc.WriteString(descStyle.Render("` with "))
		doc.WriteString(descStyle.Bold(t.Rows > 10000).Render(fmt.Sprintf("%d", t.Rows)))
		doc.WriteString(descStyle.Render(" rows examined per scan."))
		fmt.Fprintln(w, doc.String())
	}
}

func printFullIndexScanComments(w io.Writer, output output.Result) {
	printTitle(w, "Are there any full index scans?")
	tables := output.FullIndexScanTables
	if len(tables) == 0 {
		printInfo(w, "No")
		return
	}
	printInfo(w, "Yes")
	printDescription(w, "次のテーブルにおいて、全インデックススキャンが発生しています。MySQLはこれらのテーブルの全インデックスを読み取っています。")
	for _, t := range tables {
		var doc strings.Builder
		doc.WriteString(descStyle.Render("- Table `"))
		doc.WriteString(descStyle.Bold(true).Render(t.Name))
		doc.WriteString(descStyle.Render("` with index "))
		doc.WriteString(descStyle.Bold(true).Render(strPtrToString(t.Key)))
		fmt.Fprintln(w, doc.String())
	}
}

func printAnythingElseComments(w io.Writer, output output.Result) {
	printTitle(w, "Is there anything else interesting?")
	tables := output.HasAnyCommentsTables
	if len(tables) == 0 {
		printInfo(w, "No")
		return
	}
	printInfo(w, "Yes")
	for _, t := range tables {
		fmt.Fprintln(w,
			descStyle.Render("- Table `")+
				descStyle.Bold(true).Render(t.Name)+
				descStyle.Render(fmt.Sprintf("`: %s", t.Comment)),
		)
	}
	for _, comment := range output.Comments {
		fmt.Fprintln(w, descStyle.Render(fmt.Sprintf("- %s", comment)))
	}
}

func strPtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func intPtrToString(i *int) string {
	if i == nil {
		return ""
	}
	return strconv.Itoa(*i)
}
