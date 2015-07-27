package gocsv

import (
	"io"
	"io/ioutil"
	"strings"
)

type Reader struct {
	b     []byte
	pos   int
	Comma rune
}

func (r *Reader) Read() (record []string, err error) {
	if r.pos >= len(r.b) {
		return []string{}, io.EOF
	}

	delim := byte('\n')
	i := r.pos

	for ; i < len(r.b); i++ {
		if r.b[i] == delim {
			break
		}
	}

	str := string(r.b[r.pos:i])
	record = strings.Split(str, string(r.Comma))
	err = nil

	r.pos = i + 1
	return
}

func NewReader(r io.Reader) *Reader {
	bytes, _ := ioutil.ReadAll(r)
	return &Reader{bytes, 0, ','}
}

type Writer struct {
	w     io.Writer
	pos   int
	Comma rune
}

func (w *Writer) Write(record []string) error {
	str := ""

	if w.pos > 0 {
		str = "\n"
	}

	str += strings.Join(record, string(w.Comma))
	_, err := w.w.Write([]byte(str))

	w.pos++
	return err
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w, 0, ','}
}
