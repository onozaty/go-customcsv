package customcsv

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestNewReader(t *testing.T) {

	s := `a,"b","c,d"
,"","g
"
`

	r, err := NewReader(strings.NewReader(s))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	// record:1
	{
		record, err := r.Read()
		if err != nil {
			t.Fatal("failed test\n", err)
		}

		if !reflect.DeepEqual(record, []string{"a", "b", "c,d"}) {
			t.Fatal("failed test\n", record)
		}
	}

	// record:2
	{
		record, err := r.Read()
		if err != nil {
			t.Fatal("failed test\n", err)
		}

		if !reflect.DeepEqual(record, []string{"", "", "g\n"}) {
			t.Fatal("failed test\n", record)
		}
	}

	_, err = r.Read()
	if err != io.EOF {
		t.Fatal("failed test\n", err)
	}
}

func TestNewReader_NonAscii(t *testing.T) {

	s := `あ,"日本語,한글"
`

	r, err := NewReader(strings.NewReader(s))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	// record:1
	{
		record, err := r.Read()
		if err != nil {
			t.Fatal("failed test\n", err)
		}

		if !reflect.DeepEqual(record, []string{"あ", "日本語,한글"}) {
			t.Fatal("failed test\n", record)
		}
	}

	_, err = r.Read()
	if err != io.EOF {
		t.Fatal("failed test\n", err)
	}
}

func TestNewReader_RecordEndEOF(t *testing.T) {

	s := "a,b"

	r, err := NewReader(strings.NewReader(s))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	// record:1
	{
		record, err := r.Read()
		if err != nil {
			t.Fatal("failed test\n", err)
		}

		if !reflect.DeepEqual(record, []string{"a", "b"}) {
			t.Fatal("failed test\n", record)
		}
	}

	_, err = r.Read()
	if err != io.EOF {
		t.Fatal("failed test\n", err)
	}
}

func TestNewReader_QuotedField(t *testing.T) {

	s := `"a","","a,b","""","""""","a""b","1
2"`

	r, err := NewReader(strings.NewReader(s))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	// record:1
	{
		record, err := r.Read()
		if err != nil {
			t.Fatal("failed test\n", err)
		}

		if !reflect.DeepEqual(record, []string{"a", "", "a,b", `"`, `""`, `a"b`, "1\n2"}) {
			t.Fatal("failed test\n", record)
		}
	}

	_, err = r.Read()
	if err != io.EOF {
		t.Fatal("failed test\n", err)
	}
}

func TestNewReader_QuoteNotClose(t *testing.T) {

	s := `a,b,c
d,e,"f,
`

	r, err := NewReader(strings.NewReader(s))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	r.Read()
	_, err = r.Read()
	if err == nil || err.Error() != "parse error on record 2, column 3: quote is not closed" {
		t.Fatal("failed test\n", err)
	}
}

func TestNewReader_BareQuote(t *testing.T) {

	s := `a,b",c`

	r, err := NewReader(strings.NewReader(s))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	_, err = r.Read()
	if err == nil || err.Error() != "parse error on record 1, column 2: bare quote in non quoted field" {
		t.Fatal("failed test\n", err)
	}
}

func TestNewReader_UnescapedQuote(t *testing.T) {

	s := `a,b,"c"d"`

	r, err := NewReader(strings.NewReader(s))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	_, err = r.Read()
	if err == nil || err.Error() != "parse error on record 1, column 3: unescaped quote in quoted field" {
		t.Fatal("failed test\n", err)
	}
}

func TestNewReader_RecordSeparator_NewLine(t *testing.T) {

	s := "a,b\n" +
		"c,d\r\n" +
		"e,f\r" +
		"g,h"

	r, err := NewReader(strings.NewReader(s))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	// record:1 LF
	{
		record, err := r.Read()
		if err != nil {
			t.Fatal("failed test\n", err)
		}

		if !reflect.DeepEqual(record, []string{"a", "b"}) {
			t.Fatal("failed test\n", record)
		}
	}

	// record:2 CRLF
	{
		record, err := r.Read()
		if err != nil {
			t.Fatal("failed test\n", err)
		}

		if !reflect.DeepEqual(record, []string{"c", "d"}) {
			t.Fatal("failed test\n", record)
		}
	}

	// record:3 CR
	{
		record, err := r.Read()
		if err != nil {
			t.Fatal("failed test\n", err)
		}

		if !reflect.DeepEqual(record, []string{"e", "f"}) {
			t.Fatal("failed test\n", record)
		}
	}

	// record:4 EOF
	{
		record, err := r.Read()
		if err != nil {
			t.Fatal("failed test\n", err)
		}

		if !reflect.DeepEqual(record, []string{"g", "h"}) {
			t.Fatal("failed test\n", record)
		}
	}

	_, err = r.Read()
	if err != io.EOF {
		t.Fatal("failed test\n", err)
	}
}

func TestNewReader_RecordSeparator_NewLine_LastLF(t *testing.T) {

	s := "a,b\r"

	r, err := NewReader(strings.NewReader(s))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	// record:1 LF
	{
		record, err := r.Read()
		if err != nil {
			t.Fatal("failed test\n", err)
		}

		if !reflect.DeepEqual(record, []string{"a", "b"}) {
			t.Fatal("failed test\n", record)
		}
	}

	_, err = r.Read()
	if err != io.EOF {
		t.Fatal("failed test\n", err)
	}
}

func TestNewReader_SpecialRecordSeparator_OneChar(t *testing.T) {

	s := `a,b|c,d|e
,"f|"`

	r, err := NewReader(strings.NewReader(s))
	r.SpecialRecordSeparator = "|"
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	// record:1
	{
		record, err := r.Read()
		if err != nil {
			t.Fatal("failed test\n", err)
		}

		if !reflect.DeepEqual(record, []string{"a", "b"}) {
			t.Fatal("failed test\n", record)
		}
	}

	// record:2
	{
		record, err := r.Read()
		if err != nil {
			t.Fatal("failed test\n", err)
		}

		if !reflect.DeepEqual(record, []string{"c", "d"}) {
			t.Fatal("failed test\n", record)
		}
	}

	// record:3
	{
		record, err := r.Read()
		if err != nil {
			t.Fatal("failed test\n", err)
		}

		if !reflect.DeepEqual(record, []string{"e\n", "f|"}) {
			t.Fatal("failed test\n", record)
		}
	}

	_, err = r.Read()
	if err != io.EOF {
		t.Fatal("failed test\n", err)
	}
}

func TestNewReader_SpecialRecordSeparator_MultiChar(t *testing.T) {

	s := `a,b[RS]c
,[RS[RS]d,"[RS]"""[RS]`

	r, err := NewReader(strings.NewReader(s))
	r.SpecialRecordSeparator = "[RS]"
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	// record:1
	{
		record, err := r.Read()
		if err != nil {
			t.Fatal("failed test\n", err)
		}

		if !reflect.DeepEqual(record, []string{"a", "b"}) {
			t.Fatal("failed test\n", record)
		}
	}

	// record:2
	{
		record, err := r.Read()
		if err != nil {
			t.Fatal("failed test\n", err)
		}

		if !reflect.DeepEqual(record, []string{"c\n", "[RS"}) {
			t.Fatal("failed test\n", record)
		}
	}

	// record:3
	{
		record, err := r.Read()
		if err != nil {
			t.Fatal("failed test\n", err)
		}

		if !reflect.DeepEqual(record, []string{"d", "[RS]\""}) {
			t.Fatal("failed test\n", record)
		}
	}

	_, err = r.Read()
	if err != io.EOF {
		t.Fatal("failed test\n", err)
	}
}

func TestNewReader_Delimiter(t *testing.T) {

	s := "a\tb\n" +
		"\"\t\"\t"

	r, err := NewReader(strings.NewReader(s))
	r.Delimiter = '\t'
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	// record:1
	{
		record, err := r.Read()
		if err != nil {
			t.Fatal("failed test\n", err)
		}

		if !reflect.DeepEqual(record, []string{"a", "b"}) {
			t.Fatal("failed test\n", record)
		}
	}

	// record:2
	{
		record, err := r.Read()
		if err != nil {
			t.Fatal("failed test\n", err)
		}

		if !reflect.DeepEqual(record, []string{"\t", ""}) {
			t.Fatal("failed test\n", record)
		}
	}

	_, err = r.Read()
	if err != io.EOF {
		t.Fatal("failed test\n", err)
	}
}

func TestNewReader_Quote(t *testing.T) {

	s := `'a','b'
'''''',''
`

	r, err := NewReader(strings.NewReader(s))
	r.Quote = '\''
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	// record:1
	{
		record, err := r.Read()
		if err != nil {
			t.Fatal("failed test\n", err)
		}

		if !reflect.DeepEqual(record, []string{"a", "b"}) {
			t.Fatal("failed test\n", record)
		}
	}

	// record:2
	{
		record, err := r.Read()
		if err != nil {
			t.Fatal("failed test\n", err)
		}

		if !reflect.DeepEqual(record, []string{"''", ""}) {
			t.Fatal("failed test\n", record)
		}
	}

	_, err = r.Read()
	if err != io.EOF {
		t.Fatal("failed test\n", err)
	}
}

func TestNewReader_VerifyFieldsPerRecord_True(t *testing.T) {

	s := `a,b
c,d,e
`

	r, err := NewReader(strings.NewReader(s))
	r.VerifyFieldsPerRecord = true
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	// record:1
	{
		record, err := r.Read()
		if err != nil {
			t.Fatal("failed test\n", err)
		}

		if !reflect.DeepEqual(record, []string{"a", "b"}) {
			t.Fatal("failed test\n", record)
		}
	}

	// record:2
	{
		// VerifyFieldsPerRecord = true
		// If the number of fields is different from the first record, an error will occur.
		_, err := r.Read()
		if err == nil || err.Error() != "parse error on record 2: wrong number of fields" {
			t.Fatal("failed test\n", err)
		}
	}
}

func TestNewReader_VerifyFieldsPerRecord_True_EOF(t *testing.T) {

	s := `a,b
c,d,e`

	r, err := NewReader(strings.NewReader(s))
	r.VerifyFieldsPerRecord = true
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	// record:1
	{
		record, err := r.Read()
		if err != nil {
			t.Fatal("failed test\n", err)
		}

		if !reflect.DeepEqual(record, []string{"a", "b"}) {
			t.Fatal("failed test\n", record)
		}
	}

	// record:2
	{
		// VerifyFieldsPerRecord = true
		// If the number of fields is different from the first record, an error will occur.
		_, err := r.Read()
		if err == nil || err.Error() != "parse error on record 2: wrong number of fields" {
			t.Fatal("failed test\n", err)
		}
	}
}

func TestNewReader_VerifyFieldsPerRecord_True_FixNum(t *testing.T) {

	s := `a,b
c,d
`

	r, err := NewReader(strings.NewReader(s))
	r.VerifyFieldsPerRecord = true
	r.FieldsPerRecord = 1
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	// record:1
	{
		// VerifyFieldsPerRecord = true and FieldsPerRecord = 1
		// The number of fields is different from the number specified in "FieldsPerRecord", so an error occurs.
		_, err := r.Read()
		if err == nil || err.Error() != "parse error on record 1: wrong number of fields" {
			t.Fatal("failed test\n", err)
		}
	}
}

func TestNewReader_VerifyFieldsPerRecord_False(t *testing.T) {

	s := `a,b
c,d,e
`

	r, err := NewReader(strings.NewReader(s))
	r.VerifyFieldsPerRecord = false
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	// record:1
	{
		record, err := r.Read()
		if err != nil {
			t.Fatal("failed test\n", err)
		}

		if !reflect.DeepEqual(record, []string{"a", "b"}) {
			t.Fatal("failed test\n", record)
		}
	}

	// record:2
	{
		// VerifyFieldsPerRecord = false
		// If the number of fields is different from the first record, no error will occur.
		record, err := r.Read()
		if err != nil {
			t.Fatal("failed test\n", err)
		}

		if !reflect.DeepEqual(record, []string{"c", "d", "e"}) {
			t.Fatal("failed test\n", record)
		}
	}

	_, err = r.Read()
	if err != io.EOF {
		t.Fatal("failed test\n", err)
	}
}

func TestNewReader_ReadAll(t *testing.T) {

	s := `a,"b","c,d"
,"","g
"
`

	r, err := NewReader(strings.NewReader(s))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	records, err := r.ReadAll()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	expected := [][]string{
		{"a", "b", "c,d"},
		{"", "", "g\n"},
	}

	if !reflect.DeepEqual(records, expected) {
		t.Fatal("failed test\n", records)
	}
}

func TestNewReader_WithBOM(t *testing.T) {

	s := "\uFEFFa,b"

	r, err := NewReader(strings.NewReader(s))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	// record:1
	{
		record, err := r.Read()
		if err != nil {
			t.Fatal("failed test\n", err)
		}

		// Verify that the BOM has been removed.
		if !reflect.DeepEqual(record, []string{"a", "b"}) {
			t.Fatal("failed test\n", record)
		}
	}

	_, err = r.Read()
	if err != io.EOF {
		t.Fatal("failed test\n", err)
	}
}

func TestNewReader_Empty(t *testing.T) {

	s := ""

	r, err := NewReader(strings.NewReader(s))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	_, err = r.Read()
	if err != io.EOF {
		t.Fatal("failed test\n", err)
	}
}

func TestNewReader_OnlyLF(t *testing.T) {

	s := "\n"

	r, err := NewReader(strings.NewReader(s))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	// record:1
	{
		record, err := r.Read()
		if err != nil {
			t.Fatal("failed test\n", err)
		}

		if !reflect.DeepEqual(record, []string{""}) {
			t.Fatal("failed test\n", record)
		}
	}

	_, err = r.Read()
	if err != io.EOF {
		t.Fatal("failed test\n", err)
	}
}
