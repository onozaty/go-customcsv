# go-customcsv

This is a library for customizing the CSV format.

You can customize the following.

* (Reader/Writer) Format characters
    * Delimiter (default: `,`)
    * Quote (default: `"`)
    * Record separator (default: `\r\n`)
* (Writer) Always quote (default: `false`)
* (Reader) Verify the number of fields per record (default: `true`)

In Reader, the head BOM will be automatically skipped.

## Usage

### Reader

```go
f, err := os.Open("input.csv")
if err != nil {
	return err
}

r, err := customcsv.NewReader(f)
if err != nil {
	return err
}

for {
	record, err := r.Read()
	if err == io.EOF {
		break
	}
	if err != nil {
    	return err
	}

    fmt.Println(record)
}
```

In `Reader`, the following items can be customized.

```go
type Reader struct {
	// Delimiter is the field delimiter.
	// It is set to default comma (',') by NewReader.
	Delimiter rune

	// Quote is the field quote character.
	// It is set to default double quote ('"') by NewReader.
	Quote rune

	// SpecialRecordSeparator is the special record separator.
	// If not specified, a newline ('\n' '\r' '\r\n') will be used as the record separator.
	SpecialRecordSeparator string

	// FieldsPerRecord is the number of expected fields per record.
	// FieldsPerRecord > 0 : Checks for the specified value.
	// FieldsPerRecord = 0 : Check by the number of fields in the first record.
	// FieldsPerRecord < 0 : No check.
	FieldsPerRecord int
}
```

Set this after `NewReader()`.

```go
r, err := customcsv.NewReader(f)
if err != nil {
	return err
}

// Customize format
r.Delimiter = ';'
r.SpecialRecordSeparator = "|"
```
