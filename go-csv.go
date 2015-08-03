package gocsv

import (
	"io"
	"io/ioutil"
	"regexp"
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

	// "        Match any double quote
	// (?:      Containing noncapturing group
	//    [^\"] of non-double quotes
	//    |     OR
	//    \"\"  pair of consecutive double quotes
	// )+       at least one match
	// \"       ending with double quote (i.e. any letters between two quotes, plus pairs of consecutive double quotes)
	// |        OR
	// [^       NOT
	//    ,\"   comma (r.Comma) followed by double quote
	// ]*       zero or more matches
	splitRgx := regexp.MustCompile("\"(?:[^\"]|\"\")+\"|[^" + string(r.Comma) + "\"]*")
	trimQuotesRgx := regexp.MustCompile("^\"|\"$")

	str := string(r.b[r.pos:i])
	record = splitRgx.FindAllString(str, -1)
	err = nil

	for idx, recStr := range record {
		recStr = trimQuotesRgx.ReplaceAllLiteralString(recStr, "")
		recStr = strings.Replace(recStr, "\"\"", "\"", -1)
		record[idx] = recStr
	}

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
	var str string

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
