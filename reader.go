package customcsv

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
)

type ParseError struct {
	Message string
	Record  int
	Column  int
}

func (e *ParseError) Error() string {
	if e.Column == 0 {
		return fmt.Sprintf("parse error on record %d: %v", e.Record, e.Message)
	} else {
		return fmt.Sprintf("parse error on record %d, column %d: %v", e.Record, e.Column, e.Message)
	}
}

type Reader struct {
	Delimiter              rune
	Quote                  rune
	SpecialRecordSeparator string
	VerifyFieldsPerRecord  bool
	FieldsPerRecord        int
	r                      *bufio.Reader
	runeBuffer             []rune
	numRecord              int
}

var utf8bom = []byte{0xEF, 0xBB, 0xBF}

func NewReader(r io.Reader) (*Reader, error) {

	br := bufio.NewReader(r)
	mark, err := br.Peek(len(utf8bom))

	if err != io.EOF && err != nil {
		return nil, err
	}

	if reflect.DeepEqual(mark, utf8bom) {
		// If there is a BOM, skip the BOM.
		br.Discard(len(utf8bom))
	}

	return &Reader{
		Delimiter:              ',',
		Quote:                  '"',
		SpecialRecordSeparator: "", // If not specified, a newline will be used as the record separator.
		VerifyFieldsPerRecord:  true,
		r:                      br,
		runeBuffer:             []rune{},
		numRecord:              1,
	}, nil
}

func (r *Reader) Read() ([]string, error) {

	quotedField := false
	quoting := false
	field := []rune{}
	record := []string{}

	for {

		c, err := r.readRune()
		if err != nil && err != io.EOF {
			return nil, err
		}

		if err == io.EOF {

			if quoting {
				return nil, &ParseError{Message: "quote is not closed", Record: r.numRecord, Column: len(record) + 1}
			}

			if len(record) == 0 && len(field) == 0 {
				return nil, err
			}

			record = append(record, string(field))
			if err := r.verifyRecord(record); err != nil {
				return nil, err
			}
			r.numRecord++
			return record, nil
		}

		// Judge the record separator first.
		if !quoting {
			isRecordSeparator, err := r.judgeRecordSeparator(c)
			if err != nil {
				return nil, err
			}

			if isRecordSeparator {
				record = append(record, string(field))
				if err := r.verifyRecord(record); err != nil {
					return nil, err
				}
				r.numRecord++
				return record, nil
			}
		}

		switch c {
		case r.Delimiter:
			if quoting {
				field = append(field, c)
			} else {
				record = append(record, string(field))
				field = []rune{}
				quotedField = false
				quoting = false
			}
		case r.Quote:
			if !quotedField && len(field) == 0 {
				quotedField = true
				quoting = true
			} else {
				if !quotedField {
					return nil, &ParseError{Message: "bare quote in non quoted field", Record: r.numRecord, Column: len(record) + 1}
				}

				if !quoting {
					// Escaped quote.
					field = append(field, c)
				}

				quoting = !quoting
			}
		default:
			if quotedField && !quoting {
				return nil, &ParseError{Message: "unescaped quote in quoted field", Record: r.numRecord, Column: len(record) + 1}
			}

			field = append(field, c)
		}
	}
}

func (r *Reader) ReadAll() ([][]string, error) {

	records := [][]string{}

	for {
		record, err := r.Read()
		if err == io.EOF {
			return records, nil
		}
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
}

func (r *Reader) readRune() (rune, error) {

	if len(r.runeBuffer) != 0 {
		c := r.runeBuffer[0]
		r.runeBuffer = r.runeBuffer[1:]
		return c, nil
	}

	c, _, err := r.r.ReadRune()
	return c, err
}

func (r *Reader) peekRune(n int) ([]rune, error) {

	for len(r.runeBuffer) < n {
		c, _, err := r.r.ReadRune()
		if err != nil && err != io.EOF {
			return nil, err
		}

		if err == io.EOF {
			// Peek does not error on EOF. It returns only what it can read.
			return r.runeBuffer, nil
		}

		r.runeBuffer = append(r.runeBuffer, c)
	}

	return r.runeBuffer[0:n], nil
}

func (r *Reader) judgeRecordSeparator(c rune) (bool, error) {

	if r.SpecialRecordSeparator == "" {
		// Newlines are record separators.
		if c == '\n' {
			return true, nil
		}

		if c == '\r' {
			next, err := r.peekRune(1)
			if err != nil {
				return false, err
			}

			if len(next) != 0 && next[0] == '\n' {
				// Use CR+LF as a single separator
				r.readRune()
			}

			return true, nil
		}

	} else {
		// The specified character is the record separator.
		if c == []rune(r.SpecialRecordSeparator)[0] {
			// If the first character is the same, the remaining characters are included in the comparison.
			remaining, err := r.peekRune(len(r.SpecialRecordSeparator) - 1)
			if err != nil {
				return false, err
			}

			if string(append([]rune{c}, remaining...)) == r.SpecialRecordSeparator {
				for i := 1; i < len(r.SpecialRecordSeparator); i++ {
					// Skip characters that are record separators.
					r.readRune()
				}

				return true, nil
			}
		}
	}

	return false, nil
}

func (r *Reader) verifyRecord(record []string) error {

	if !r.VerifyFieldsPerRecord {
		return nil
	}

	if r.FieldsPerRecord == 0 {
		// Keep the number of fields in the first record.
		r.FieldsPerRecord = len(record)
		return nil
	}

	if len(record) != r.FieldsPerRecord {
		return &ParseError{Message: "wrong number of fields", Record: r.numRecord}
	}

	return nil
}
