package customcsv

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
)

type ParseError struct {
	Record int
	Column int
	Err    error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error on record %d, column %d: %v", e.Record, e.Column, e.Err)
}

type Reader struct {
	Delimiter              rune
	Quote                  rune
	SpecialRecordSeparator []rune
	r                      *bufio.Reader
	runeBuffer             []rune
	numRecoed              int
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		Delimiter:              ',',
		Quote:                  '"',
		SpecialRecordSeparator: nil, // 指定無しの場合は改行(\r,\n)が対象になる
		r:                      bufio.NewReader(r),
		numRecoed:              1,
	}
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
				return nil, &ParseError{Record: r.numRecoed, Column: len(field) + 1, Err: fmt.Errorf("quote is not closed")}
			}

			if len(record) == 0 && len(field) == 0 {
				return nil, err
			}

			record = append(record, string(field))
			r.numRecoed++
			return record, err
		}

		// レコードの終端は特殊なのでここで判定
		if !quoting {
			isRecordSeparator, err := r.judgeRecordSeparator(c)
			if err != nil {
				return nil, err
			}

			if isRecordSeparator {
				record = append(record, string(field))
				r.numRecoed++
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
			if len(field) == 0 {
				quotedField = true
				quoting = true
			} else {
				if !quotedField {
					return nil, &ParseError{Record: r.numRecoed, Column: len(field) + 1, Err: fmt.Errorf("bare quote in non quoted field")}
				}

				if !quoting {
					// クォートが2つ続いている
					field = append(field, c)
				}

				quoting = !quoting
			}
		default:
			if quotedField && !quoting {
				return nil, &ParseError{Record: r.numRecoed, Column: len(field) + 1, Err: fmt.Errorf("extraneous or missing quote in quoted field")}
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
			return r.runeBuffer, err
		}

		r.runeBuffer = append(r.runeBuffer, c)
	}

	return r.runeBuffer[0:n], nil
}

func (r *Reader) judgeRecordSeparator(c rune) (bool, error) {

	if r.SpecialRecordSeparator == nil {
		// 改行で判定
		if c == '\n' {
			return true, nil
		}

		if c == '\r' {
			next, err := r.peekRune(1)
			if err != nil {
				return false, err
			}

			if len(next) != 0 && next[0] == '\n' {
				// 次が\nならば、CRLFで終端として扱うために読み飛ばしておく
				r.readRune()
			}

			return true, nil
		}

	} else {
		// 指定した文字を行の区切りとして利用

		if c == r.SpecialRecordSeparator[0] {
			// 先頭の文字が同じ場合、2文字目以降も比較
			other, err := r.peekRune(len(r.SpecialRecordSeparator) - 1)
			if err != nil {
				return false, err
			}

			if reflect.DeepEqual(r.SpecialRecordSeparator[1:], other) {
				// 一致した場合には、読み飛ばしておく
				for i := 1; i < len(r.SpecialRecordSeparator); i++ {
					r.readRune()
				}

				return true, nil
			}
		}
	}

	return false, nil
}
