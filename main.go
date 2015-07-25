package gocsv

import (
	"io"
	"io/ioutil"
)

type Reader struct {
	b     []byte
	pos   int
	Comma rune
}

func (r *Reader) Read() (record string, err error) {
	if r.pos >= len(r.b) {
		return "", io.EOF
	}

	delim := byte('\n')
	i := r.pos

	for ; i < len(r.b); i++ {
		if r.b[i] == delim {
			break
		}
	}

	record = string(r.b[r.pos:i])
	err = nil

	r.pos = i + 1
	return
}

func NewReader(r io.Reader) *Reader {
	bytes, _ := ioutil.ReadAll(r)
	return &Reader{bytes, 0, ','}
}
