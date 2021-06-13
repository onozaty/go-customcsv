package customcsv

import (
	"bufio"
	"bytes"
	"testing"
)

func TestNewWriter(t *testing.T) {

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	cw := NewWriter(w)

	if err := cw.Write([]string{"1", "2", "3", "4"}); err != nil {
		t.Fatal("failed test\n", err)
	}
	if err := cw.Write([]string{"", ",", "\"", "\"x\""}); err != nil {
		t.Fatal("failed test\n", err)
	}
	if err := cw.Write([]string{"\r", "\n", "\r\n", "a\r\nb\nc"}); err != nil {
		t.Fatal("failed test\n", err)
	}

	if err := cw.Flush(); err != nil {
		t.Fatal("failed test\n", err)
	}

	result := b.String()

	expect := "1,2,3,4\r\n" +
		",\",\",\"\"\"\",\"\"\"x\"\"\"\r\n" +
		"\"\r\",\"\n\",\"\r\n\",\"a\r\nb\nc\"\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestNewWriter_Delimiter(t *testing.T) {

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	cw := NewWriter(w)
	cw.Delimiter = ';'

	if err := cw.Write([]string{"1", "2", "3"}); err != nil {
		t.Fatal("failed test\n", err)
	}
	if err := cw.Write([]string{"", "a;b;", ";"}); err != nil {
		t.Fatal("failed test\n", err)
	}
	if err := cw.Write([]string{"\r", "\n", "\r\n"}); err != nil {
		t.Fatal("failed test\n", err)
	}

	if err := cw.Flush(); err != nil {
		t.Fatal("failed test\n", err)
	}

	result := b.String()

	expect := "1;2;3\r\n" +
		";\"a;b;\";\";\"\r\n" +
		"\"\r\";\"\n\";\"\r\n\"\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestNewWriter_Quote(t *testing.T) {

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	cw := NewWriter(w)
	cw.Quote = '\''

	if err := cw.Write([]string{"1", "2"}); err != nil {
		t.Fatal("failed test\n", err)
	}
	if err := cw.Write([]string{"", "'a'"}); err != nil {
		t.Fatal("failed test\n", err)
	}
	if err := cw.Write([]string{"\r", "\r\n"}); err != nil {
		t.Fatal("failed test\n", err)
	}

	if err := cw.Flush(); err != nil {
		t.Fatal("failed test\n", err)
	}

	result := b.String()

	expect := "1,2\r\n" +
		",'''a'''\r\n" +
		"'\r','\r\n'\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestNewWriter_AllQuotes(t *testing.T) {

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	cw := NewWriter(w)
	cw.AllQuotes = true

	if err := cw.Write([]string{"1", "2", "3"}); err != nil {
		t.Fatal("failed test\n", err)
	}
	if err := cw.Write([]string{"", ",", "\""}); err != nil {
		t.Fatal("failed test\n", err)
	}
	if err := cw.Write([]string{"\r", "\n", "a\rb\n"}); err != nil {
		t.Fatal("failed test\n", err)
	}

	if err := cw.Flush(); err != nil {
		t.Fatal("failed test\n", err)
	}

	result := b.String()

	expect := "\"1\",\"2\",\"3\"\r\n" +
		"\"\",\",\",\"\"\"\"\r\n" +
		"\"\r\",\"\n\",\"a\rb\n\"\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestNewWriter_RecordSeparator(t *testing.T) {

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	cw := NewWriter(w)
	cw.RecordSeparator = "|"

	if err := cw.Write([]string{"1", "2"}); err != nil {
		t.Fatal("failed test\n", err)
	}
	if err := cw.Write([]string{"|", ","}); err != nil {
		t.Fatal("failed test\n", err)
	}
	if err := cw.Write([]string{"\r", "\n"}); err != nil {
		t.Fatal("failed test\n", err)
	}

	if err := cw.Flush(); err != nil {
		t.Fatal("failed test\n", err)
	}

	result := b.String()

	expect := "1,2|" +
		"\"|\",\",\"|" +
		"\"\r\",\"\n\"|"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestNewWriter_WriteAll(t *testing.T) {

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	cw := NewWriter(w)

	err := cw.WriteAll(
		[][]string{
			{"1", "2", "3", "4"},
			{"", ",", "\"", "\"x\""},
			{"\r", "\n", "\r\n", "a\r\nb\nc"},
		})

	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := b.String()

	expect := "1,2,3,4\r\n" +
		",\",\",\"\"\"\",\"\"\"x\"\"\"\r\n" +
		"\"\r\",\"\n\",\"\r\n\",\"a\r\nb\nc\"\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}
