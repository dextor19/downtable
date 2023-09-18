/*
package downtable generates markdown tables from csv and json files.
there is methods for modifying the table data before generating the markdown table.

to parse an csv file as a input for the markdown table you need to use the [WithCSVFile] function inside of the [AddTable] method on the [MarkdownTable] interface.
csv data needs to have the first row be the headers and the subsequent rows will be general rows within the table.

  "header1, header2,header3,\nrow1item1, row1item2, row1item3\n"

when using providing a csv file you need to enable formatting options based on the type of csv file,
options `lazyQuotes` " double quotes are allowed in csv fields and `trimLeadingQuotes` leading white spaces in the csv field is ignored.

json files are also able to provided as input for markdown tables, using the [MarkdownTable] method [AddJSONFileTable].

to parse json files the package requires the following format:

  {
    "Headers": ["header1"],
    "Rows": [["row1item1"], ["row2item1"]]
  }

main idea is to use the [MarkdownTable] interface to parse array strings and produce a string of markdown that represents a table.


*/
package downtable

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// MarkdownTable interface defines all the methods that consumers of this package
// can use to generate a markdown table string.
type MarkdownTable interface {
	// AddHeader adds a single string to the end of the headers string array
	AddHeader(string)

	// AddHeaders replaces the array of strings in table instance with new array of strings
	// that will be used as headers for the markdown table.
	AddHeaders([]string)

	// DeleteHeaders removes all items inside the array of strings for table.headers
	DeleteHeaders()

	// AddRowItem adds a single items to one row in the matrix of arrays for table.rows.
	// requires an input specifying which row in the matrix this string will be appended too.
	AddRowItem(string, int)

	// AddRow will append an addition row to the matrix of strings for table.rows.
	AddRow([]string) error

	// AddRows will replace all matrix of strings with a new matrix for table.rows.
	AddRows([][]string)

	// DeleteRows will delete the matrix of strings in table.rows.
	DeleteRows()

	// AddTable takes a matrix of strings as input and will use the first row in the matrix
	// as the headers list of strings and replace existing headers. all other rows will be
	// appended to the matrix of strings for table.rows.
	//
	// AddTable requires that the headers and rows have the same number of items in the string array
	// otherwise error will occur saying that header and row arrays need to have the same length.
	AddTable([][]string, error) error

	// AddJSONFileTable takes a file pointer as input requires the file associated with the
	// pointer contains a json formatted object and has a specific Headers and Rows structure.
	// this method will replace existing headers in the table instance.
	AddJSONFileTable(*os.File) error

	// GetMarkdownTableString a single string of the markdown table via the standard output.
	GetMarkdownTableString() (string, error)

	// GetMarkdownTableString outputs a single string of the markdown table.
	GetMarkdownTable() ([]byte, error)

	// GetMarkdownTable outputs a byte array containing a single string of the markdown table.
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

// NewMarkdownTable initiates a empty instance of a [table] struct with a [MarkdownTable] interface
// that allows you to modify the data within the [table] instance using getter and setter methods or
// formatted data files.
func NewMarkdownTable() MarkdownTable {
	var mdt MarkdownTable
	mdt = &table{}
	return mdt
}

// WithCSVFile takes a file pointer with csv reader options as parameters to the function and
// outputs a matrix of strings which the first row in the matrix is the headers and the other
// rows are row items for the markdown table.
//
// the matrix of strings output is meant for the MarkdownTable.AddTable() method which populates [table]
// struct instances using matrix of strings.
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
