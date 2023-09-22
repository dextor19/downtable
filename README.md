# Downtable

package downtable generates markdown tables from csv and json files. there is methods for modifying the table data before generating the markdown table.

## Installation

`go get -u github.com/dextor19/downtable`

## Usage

### Parsing CSV File

Providing an CSV file as a input for the markdown table you need to use the `AddTableFromCSVFile(*os.File, lazyQuotes: bool, trimLeadingQuotes: bool)` method inside of the `MarkdownTable` interface. When using providing a CSV file you need to enable formatting options based on the type of CSV file, options `lazyQuotes` " double quotes are allowed in csv fields and `trimLeadingQuotes` leading white spaces in the csv field is ignored.

### Parsing JSON Files

JSON files are also able to provided as input for markdown tables, using the `MarkdownTable` method `AddTableFromJSONFile(*os.File)`.

Parse a JSON file it requires a specific format:

```json
{
    "Headers": [
        "header1",
        "header2",
        "header3",
    ],
    "Rows": [
        [
            "row1item1",
            "row1item2",
            "row1item3",
        ],
        [
            "row2item1",
            "row2item2",
            "row3item3",
        ]
    ]
}
```

### Output Markdown Table String

there are three way of exporting a Markdown table, with the `MarkdownTable` interface you are able to output the markdown table in `os.Stdout`, `string` or `[]byte` type.

#### Stdout

`PrintMarkdownTable()` prints a single string of the markdown table via the standard output.

#### String

`GetMarkdownTableString()` outputs a single string of the markdown table.

#### Byte

`GetMarkdownTable()` outputs a byte array containing a single string of the markdown table.

## Examples

### CSV File

```golang
package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"

    "github.com/dextor19/downtable"
)

func main() {
    input := []byte("header1, header2, header3, header4, header5\nrow1item1, row1item2, row1item3, row1item4, row1item5\nrow2item1, row2item2, row2item3, row2item4, row2item5\nrow3item1, row3item2, row3item3, row3item4, row3item5\n")
    mdt := downtable.NewMarkdownTable()
    tmpFile, err := ioutil.TempFile(os.TempDir(), "temp_example.csv")
    if err != nil {
        log.Fatal(err)
    }
    if _, err := tmpFile.Write(input); err != nil {
        log.Fatal(err)
    }
    if _, err := tmpFile.Seek(0, 0); err != nil {
        log.Fatal(err)
    }
    mdt.AddTableFromCSVFile(tmpFile, true, true)
    _, err = mdt.PrintMarkdownTable()
    if err != nil {
        log.Fatal(err)
    }
}
```

### JSON File

```golang
package main

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "os"

    "github.com/dextor19/downtable"
)

type Table struct {
    Headers []string
    Rows    [][]string
}

func main() {
    input := Table{
        Headers: []string{"header1", "header2", "header3", "header4", "header5"},
        Rows: [][]string{{"row1item1", "row1item2", "row1item3", "row1item4", "row1item5"},
            {"row2item1", "row2item2", "row2item3", "row2item4", "row2item5"},
            {"row3item1", "row3item2", "row3item3", "row3item4", "row3item5"}},
    }
    jsonInput, err := json.Marshal(input)
    if err != nil {
        log.Fatal(err)
    }
    mdt := downtable.NewMarkdownTable()
    tmpFile, err := ioutil.TempFile(os.TempDir(), "temp_example.json")
    if err != nil {
        log.Fatal(err)
    }
    if _, err := tmpFile.Write(jsonInput); err != nil {
        log.Fatal(err)
    }
    if _, err := tmpFile.Seek(0, 0); err != nil {
        log.Fatal(err)
    }
    mdt.AddTableFromJSONFile(tmpFile)
    _, err = mdt.PrintMarkdownTable()
    if err != nil {
        log.Fatal(err)
    }
}
```

### Output

```bash
| header1  | header2  | header3  | header4  | header5    |
|-------|-------|-------|-------|-------|
| row1item1    | row1item2    | row1item3    | row1item4    | row1item5    |
| row2item1    | row2item2    | row2item3    | row2item4    | row2item5    |
| row3item1    | row3item2    | row3item3    | row3item4    | row3item5    |
```

| header1  | header2  | header3  | header4  | header5    |
|-------|-------|-------|-------|-------|
| row1item1    | row1item2    | row1item3    | row1item4    | row1item5    |
| row2item1    | row2item2    | row2item3    | row2item4    | row2item5    |
| row3item1    | row3item2    | row3item3    | row3item4    | row3item5    |
