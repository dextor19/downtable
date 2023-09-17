package mdtable

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type MarkdownTable interface {
	AddHeader(string)
	AddHeaders([]string)
	DeleteHeaders()
	AddRowItem(string, int)
	AddRow([]string) error
	AddRows([][]string)
	DeleteRows()
	AddTable([][]string, error) error
	AddJSONFileTable(*os.File) error
	GetMarkdownTableString() (string, error)
	GetMarkdownTable() ([]byte, error)
	PrintMarkdownTable() (int, error)
}

type table struct {
	headers []string
	rows    [][]string
}

func (t *table) AddHeader(header string) {
	t.headers = append(t.headers, header)
}

func (t *table) AddHeaders(newHeaders []string) {
	t.headers = newHeaders
}

func (t *table) DeleteHeaders() {
	t.headers = []string{}
}

func (t *table) AddRowItem(rowItem string, rowIndex int) {
	t.rows[rowIndex] = append(t.rows[rowIndex], rowItem)
}

func (t *table) AddRow(newRow []string) error {
	if len(newRow) != len(t.headers) {
		return fmt.Errorf("AddRow new row length %d is greater than current headers length: %d", len(newRow), len(t.headers))
	}
	t.rows = append(t.rows, newRow)
	return nil
}

func (t *table) AddRows(newRows [][]string) {
	t.rows = newRows
}

func (t *table) DeleteRows() {
	t.rows = [][]string{}
}

func (t *table) AddTable(rows [][]string, err error) error {
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		return fmt.Errorf("AddTable provided rows are empty")
	}
	if len(rows) < 2 {
		return fmt.Errorf("AddTable provided rows are less than 2, AddTable requires 1 row for headers and another row for items")
	}
	t.AddHeaders(rows[0])

	emptyRowCount := 0

	for _, row := range rows[1:] {
		if len(row) == 0 {
			emptyRowCount++
			continue
		}
		t.AddRow(row)
	}
	return nil

}

type jsonTable struct {
	Headers []string
	Rows    [][]string
}

func (t *table) AddJSONFileTable(file *os.File) error {
	jsonData, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	var jt jsonTable
	err = json.Unmarshal(jsonData, &jt)
	if err != nil {
		return err
	}
	t.AddHeaders(jt.Headers)
	t.AddRows(jt.Rows)
	return nil
}

func (t *table) GetMarkdownTableString() (string, error) {
	if len(t.headers) == 0 {
		return "", fmt.Errorf("markdown table struct headers are empty")
	}
	if len(t.rows) == 0 {
		return "", fmt.Errorf("markdown table struct rows has no rows")
	}
	var markdown string
	markdown = markdown + "| " + strings.Join(t.headers, "  | ") + "    |\n"

	var headerSeparator string
	for i := 0; i < len(t.headers); i++ {
		headerSeparator = headerSeparator + "|-------"
	}
	headerSeparator = headerSeparator + "|\n"
	markdown = markdown + headerSeparator

	emptyRowCount := 0

	for _, row := range t.rows {
		if len(row) == 0 {
			emptyRowCount++
			continue
		}
		markdown = markdown + "| " + strings.Join(row, "    | ") + "    |\n"
	}
	if emptyRowCount == len(t.rows) {
		return "", fmt.Errorf("markdown table struct rows have no items, empty row count: %d", emptyRowCount)
	}
	return markdown, nil
}

func (t *table) PrintMarkdownTable() (int, error) {
	mdtString, err := t.GetMarkdownTableString()
	if err != nil {
		return 0, err
	}
	numberOfBytes, err := fmt.Printf("%s", mdtString)
	if err != nil {
		return numberOfBytes, err
	}
	return numberOfBytes, nil
}

func (t *table) GetMarkdownTable() ([]byte, error) {
	mdtString, err := t.GetMarkdownTableString()
	if err != nil {
		return nil, err
	}
	return []byte(mdtString), nil
}

func NewMarkdownTable() MarkdownTable {
	var mdt MarkdownTable
	mdt = &table{}
	return mdt
}

func WithCSVFile(file *os.File, lazyQuotes bool, trimLeadingSpace bool) ([][]string, error) {
	csvReader := csv.NewReader(file)
	csvReader.LazyQuotes = lazyQuotes
	csvReader.TrimLeadingSpace = trimLeadingSpace
	csvRows, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	return csvRows, nil
}
