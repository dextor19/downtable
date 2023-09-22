package downtable

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	lazyQuotes        = true
	trimLeadingQuotes = true
)

func TestCSVFileInputAndGetMarkdownTableByte(t *testing.T) {
	start := time.Now()
	defer fmt.Println("took: ", time.Since(start))
	input := []byte("month, age, gender, name, blood type\nmay, 30, male, john, A\njune, 32, female, peter, O\njuly, 33, male, charles, B\nseptember, 34, alien, bob, AB\n")
	want := []byte("| month  | age  | gender  | name  | blood type    |\n|-------|-------|-------|-------|-------|\n| may    | 30    | male    | john    | A    |\n| june    | 32    | female    | peter    | O    |\n| july    | 33    | male    | charles    | B    |\n| september    | 34    | alien    | bob    | AB    |\n")
	mdt := NewMarkdownTable()
	tmpFile, err := ioutil.TempFile(t.TempDir(), t.Name())
	if err != nil {
		t.Errorf("got: %q", err)
	}
	if _, err := tmpFile.Write(input); err != nil {
		t.Errorf("got: %q", err)
	}
	if _, err := tmpFile.Seek(0, 0); err != nil {
		t.Errorf("got: %q", err)
	}
	os.Stdin = tmpFile
	mdt.AddTableFromCSVFile(os.Stdin, lazyQuotes, trimLeadingQuotes)
	got, err := mdt.GetMarkdownTable()
	if err != nil {
		t.Errorf("got: %q", err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("got: %q want %q", got, want)
	}
}

func TestCSVFileInputAndGetMarkdownTableString(t *testing.T) {
	start := time.Now()
	defer fmt.Println("took: ", time.Since(start))
	input := []byte("month, age, gender, name, blood type\nmay, 30, male, john, A\njune, 32, female, peter, O\njuly, 33, male, charles, B\nseptember, 34, alien, bob, AB\n")
	want := "| month  | age  | gender  | name  | blood type    |\n|-------|-------|-------|-------|-------|\n| may    | 30    | male    | john    | A    |\n| june    | 32    | female    | peter    | O    |\n| july    | 33    | male    | charles    | B    |\n| september    | 34    | alien    | bob    | AB    |\n"
	mdt := NewMarkdownTable()
	tmpFile, err := ioutil.TempFile(t.TempDir(), t.Name())
	if err != nil {
		t.Errorf("got: %q", err)
	}
	if _, err := tmpFile.Write(input); err != nil {
		t.Errorf("got: %q", err)
	}
	if _, err := tmpFile.Seek(0, 0); err != nil {
		t.Errorf("got: %q", err)
	}
	os.Stdin = tmpFile
	mdt.AddTableFromCSVFile(os.Stdin, lazyQuotes, trimLeadingQuotes)
	got, err := mdt.GetMarkdownTableString()
	if err != nil {
		t.Errorf("got: %q", err)
	}
	if strings.Compare(got, want) != 0 {
		t.Errorf("got: %q want %q", got, want)
	}
}

func TestJSONFileInputAndGetMarkdownTableByte(t *testing.T) {
	start := time.Now()
	defer fmt.Println("took: ", time.Since(start))
	input := jsonTable{
		Headers: []string{"month", "age", "gender", "name", "blood type"},
		Rows: [][]string{{"may", "30", "male", "john", "A"},
			{"june", "32", "female", "peter", "O"},
			{"july", "33", "male", "charles", "B"},
			{"september", "34", "alien", "bob", "AB"}},
	}
	jsonInput, err := json.Marshal(input)
	want := []byte("| month  | age  | gender  | name  | blood type    |\n|-------|-------|-------|-------|-------|\n| may    | 30    | male    | john    | A    |\n| june    | 32    | female    | peter    | O    |\n| july    | 33    | male    | charles    | B    |\n| september    | 34    | alien    | bob    | AB    |\n")
	mdt := NewMarkdownTable()
	tmpFile, err := ioutil.TempFile(t.TempDir(), t.Name())
	if err != nil {
		t.Errorf("got: %q", err)
	}
	if _, err := tmpFile.Write(jsonInput); err != nil {
		t.Errorf("got: %q", err)
	}
	if _, err := tmpFile.Seek(0, 0); err != nil {
		t.Errorf("got: %q", err)
	}
	os.Stdin = tmpFile
	mdt.AddTableFromJSONFile(os.Stdin)
	got, err := mdt.GetMarkdownTable()
	if err != nil {
		t.Errorf("got: %q", err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("got: %q want %q", got, want)
	}
}

func TestPrintMarkdownTable(t *testing.T) {
	start := time.Now()
	defer fmt.Println("took: ", time.Since(start))
	input := jsonTable{
		Headers: []string{"month", "age", "gender", "name", "blood type"},
		Rows: [][]string{{"may", "30", "male", "john", "A"},
			{"june", "32", "female", "peter", "O"},
			{"july", "33", "male", "charles", "B"},
			{"september", "34", "alien", "bob", "AB"}},
	}
	want := 293
	mdt := NewMarkdownTable()
	mdt.AddHeaders(input.Headers)
	mdt.AddRows(input.Rows)
	got, err := mdt.PrintMarkdownTable()
	if err != nil {
		t.Errorf("got: %q", err)
	}
	if got != want {
		t.Errorf("got: %d want %d", got, want)
	}
}

func TestPrintMarkdownTableNoRows(t *testing.T) {
	start := time.Now()
	defer fmt.Println("took: ", time.Since(start))
	input := jsonTable{
		Headers: []string{"month", "age", "gender", "name", "blood type"},
		Rows: [][]string{{"may", "30", "male", "john", "A"},
			{"june", "32", "female", "peter", "O"},
			{"july", "33", "male", "charles", "B"},
			{"september", "34", "alien", "bob", "AB"}},
	}
	want := 0
	mdt := NewMarkdownTable()
	mdt.AddHeaders(input.Headers)
	got, _ := mdt.PrintMarkdownTable()
	if got != want {
		t.Errorf("got: %d want %d", got, want)
	}
}

func TestPrintMarkdownTableNoHeaders(t *testing.T) {
	start := time.Now()
	defer fmt.Println("took: ", time.Since(start))
	input := jsonTable{
		Headers: []string{"month", "age", "gender", "name", "blood type"},
		Rows: [][]string{{"may", "30", "male", "john", "A"},
			{"june", "32", "female", "peter", "O"},
			{"july", "33", "male", "charles", "B"},
			{"september", "34", "alien", "bob", "AB"}},
	}
	want := 0
	mdt := NewMarkdownTable()
	mdt.AddRows(input.Rows)
	got, _ := mdt.PrintMarkdownTable()
	if got != want {
		t.Errorf("got: %d want %d", got, want)
	}
}

func TestPrintMarkdownTableEmptyRows(t *testing.T) {
	start := time.Now()
	defer fmt.Println("took: ", time.Since(start))
	input := jsonTable{
		Headers: []string{"month", "age", "gender", "name", "blood type"},
		Rows:    [][]string{{}, {}, {}, {}},
	}
	want := 0
	mdt := NewMarkdownTable()
	mdt.AddHeaders(input.Headers)
	mdt.AddRows(input.Rows)
	got, _ := mdt.PrintMarkdownTable()
	if got != want {
		t.Errorf("got: %d want %d", got, want)
	}
}

func TestAddTableEmptyStringMatrix(t *testing.T) {
	start := time.Now()
	defer fmt.Println("took: ", time.Since(start))
	input := [][]string{{}, {}, {}, {}}
	want := 0
	mdt := NewMarkdownTable()
	_ = mdt.AddTable(input)
	got, _ := mdt.PrintMarkdownTable()
	if got != want {
		t.Errorf("got: %d want %d", got, want)
	}
}

func TestAddTableOnlyHeaders(t *testing.T) {
	start := time.Now()
	defer fmt.Println("took: ", time.Since(start))
	input := [][]string{{"month", "age", "gender", "name", "blood type"}}
	want := 0
	mdt := NewMarkdownTable()
	_ = mdt.AddTable(input)
	got, _ := mdt.PrintMarkdownTable()
	if got != want {
		t.Errorf("got: %d want %d", got, want)
	}
}

func TestAddTableHeadersAndEmptyRows(t *testing.T) {
	start := time.Now()
	defer fmt.Println("took: ", time.Since(start))
	input := [][]string{{"month", "age", "gender", "name", "blood type"}, {}, {}, {}}
	want := 0
	mdt := NewMarkdownTable()
	_ = mdt.AddTable(input)
	got, _ := mdt.PrintMarkdownTable()
	if got != want {
		t.Errorf("got: %d want %d", got, want)
	}
}

func TestAddTableHeadersAndRows(t *testing.T) {
	start := time.Now()
	defer fmt.Println("took: ", time.Since(start))
	input := [][]string{{"month", "age", "gender", "name", "blood type"}, {"month", "age", "gender", "name", "blood type"}, {"month", "age", "gender", "name", "blood type"}}
	want := 214
	mdt := NewMarkdownTable()
	_ = mdt.AddTable(input)
	got, _ := mdt.PrintMarkdownTable()
	if got != want {
		t.Errorf("got: %d want %d", got, want)
	}
}

func TestAddTableFourHeadersAndRowsWithFiveItems(t *testing.T) {
	start := time.Now()
	defer fmt.Println("took: ", time.Since(start))
	input := [][]string{{"month", "age", "gender", "name"}, {"month", "age", "gender", "name", "blood type"}, {"month", "age", "gender", "name", "blood type"}}
	want := 0
	mdt := NewMarkdownTable()
	_ = mdt.AddTable(input)
	got, _ := mdt.PrintMarkdownTable()
	if got != want {
		t.Errorf("got: %d want %d", got, want)
	}
}

func TestAddTableFiveHeadersAndRowsWithFourItems(t *testing.T) {
	start := time.Now()
	defer fmt.Println("took: ", time.Since(start))
	input := [][]string{{"month", "age", "gender", "name", "blood type"}, {"month", "age", "gender", "name"}, {"month", "age", "gender", "name"}}
	want := 0
	mdt := NewMarkdownTable()
	_ = mdt.AddTable(input)
	got, _ := mdt.PrintMarkdownTable()
	if got != want {
		t.Errorf("got: %d want %d", got, want)
	}
}

func TestDeleteHeaders(t *testing.T) {
	start := time.Now()
	defer fmt.Println("took: ", time.Since(start))
	input := [][]string{{"month", "age", "gender", "name", "blood type"}, {"month", "age", "gender", "name", "blood type"}, {"month", "age", "gender", "name", "blood type"}}
	want := 0
	mdt := NewMarkdownTable()
	_ = mdt.AddTable(input)
	mdt.DeleteHeaders()
	got, _ := mdt.PrintMarkdownTable()
	if got != want {
		t.Errorf("got: %d want %d", got, want)
	}
}

func TestDeleteRows(t *testing.T) {
	start := time.Now()
	defer fmt.Println("took: ", time.Since(start))
	input := [][]string{{"month", "age", "gender", "name", "blood type"}, {"month", "age", "gender", "name", "blood type"}, {"month", "age", "gender", "name", "blood type"}}
	want := 0
	mdt := NewMarkdownTable()
	_ = mdt.AddTable(input)
	mdt.DeleteRows()
	got, _ := mdt.PrintMarkdownTable()
	if got != want {
		t.Errorf("got: %d want %d", got, want)
	}
}
