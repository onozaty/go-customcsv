package customcsv

import (
	"bufio"
	"bytes"
	"testing"
)

func TestNewCsvWriter(t *testing.T) {

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	cw := NewWriter(w)

	cw.Write([]string{"1", "2", "3", "4"})
	cw.Write([]string{"", ",", "\"", "\"x\""})
	cw.Write([]string{"\r", "\n", "\r\n", "a\r\nb\nc"})

	cw.Flush()
	result := b.String()

	expect := "1,2,3,4\r\n" +
		",\",\",\"\"\"\",\"\"\"x\"\"\"\r\n" +
		"\"\r\",\"\n\",\"\r\n\",\"a\r\nb\nc\"\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestNewCsvWriter_Delimiter(t *testing.T) {

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	cw := NewWriter(w)
	cw.Delimiter = ';'

	cw.Write([]string{"1", "2", "3"})
	cw.Write([]string{"", "a;b;", ";"})
	cw.Write([]string{"\r", "\n", "\r\n"})

	cw.Flush()
	result := b.String()

	expect := "1;2;3\r\n" +
		";\"a;b;\";\";\"\r\n" +
		"\"\r\";\"\n\";\"\r\n\"\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestNewCsvWriter_Quote(t *testing.T) {

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	cw := NewWriter(w)
	cw.Quote = '\''

	cw.Write([]string{"1", "2"})
	cw.Write([]string{"", "'a'"})
	cw.Write([]string{"\r", "\r\n"})

	cw.Flush()
	result := b.String()

	expect := "1,2\r\n" +
		",'''a'''\r\n" +
		"'\r','\r\n'\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestNewCsvWriter_AllQuotes(t *testing.T) {

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	cw := NewWriter(w)
	cw.AllQuotes = true

	cw.Write([]string{"1", "2", "3"})
	cw.Write([]string{"", ",", "\""})
	cw.Write([]string{"\r", "\n", "a\rb\n"})

	cw.Flush()
	result := b.String()

	expect := "\"1\",\"2\",\"3\"\r\n" +
		"\"\",\",\",\"\"\"\"\r\n" +
		"\"\r\",\"\n\",\"a\rb\n\"\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestNewCsvWriter_RecordSeparator(t *testing.T) {

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	cw := NewWriter(w)
	cw.RecordSeparator = "|"

	cw.Write([]string{"1", "2"})
	cw.Write([]string{"|", ","})
	cw.Write([]string{"\r", "\n"})

	cw.Flush()
	result := b.String()

	expect := "1,2|" +
		"\"|\",\",\"|" +
		"\"\r\",\"\n\"|"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}
