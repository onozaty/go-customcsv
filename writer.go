package customcsv

import (
	"bufio"
	"io"
	"strings"
)

type Writer struct {
	Delimiter       rune
	Quote           rune
	AllQuotes       bool
	RecordSeparator string
	w               *bufio.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		Delimiter:       ',',
		Quote:           '"',
		AllQuotes:       false,
		RecordSeparator: "\r\n",
		w:               bufio.NewWriter(w),
	}
}

func (w *Writer) Write(record []string) error {

	for n, field := range record {
		if n > 0 {
			if _, err := w.w.WriteRune(w.Delimiter); err != nil {
				return err
			}
		}

		if !w.fieldNeedsQuotes(field) {
			if _, err := w.w.WriteString(field); err != nil {
				return err
			}
			continue
		}

		if _, err := w.w.WriteRune(w.Quote); err != nil {
			return err
		}

		if strings.ContainsRune(field, w.Quote) {
			escaped := strings.ReplaceAll(field, string(w.Quote), string([]rune{w.Quote, w.Quote}))
			if _, err := w.w.WriteString(escaped); err != nil {
				return err
			}
		} else {
			if _, err := w.w.WriteString(field); err != nil {
				return err
			}
		}

		if _, err := w.w.WriteRune(w.Quote); err != nil {
			return err
		}
	}

	_, err := w.w.WriteString(w.RecordSeparator)
	return err
}

func (w *Writer) Flush() error {
	return w.w.Flush()
}

func (w *Writer) WriteAll(records [][]string) error {
	for _, record := range records {
		err := w.Write(record)
		if err != nil {
			return err
		}
	}
	return w.w.Flush()
}

func (w *Writer) fieldNeedsQuotes(field string) bool {

	if w.AllQuotes {
		return true
	}

	return strings.ContainsRune(field, w.Delimiter) || strings.ContainsRune(field, w.Quote) ||
		strings.ContainsAny(field, w.RecordSeparator) || strings.ContainsAny(field, "\r\n")
}
