// A test module to play with different stuff
package gocsv

import (
	"io"
	"unicode"
)

type CapitalizeReader struct {
	r io.Reader
}

func (cr *CapitalizeReader) Read(p []byte) (n int, err error) {
	n, err = cr.r.Read(p)

	for i := range p {
		p[i] = byte(unicode.ToUpper(rune(p[i])))
	}

	return
}

func NewCapitalizerReader(reader io.Reader) *CapitalizeReader {
	return &CapitalizeReader{reader}
}
