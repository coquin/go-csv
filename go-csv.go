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

	strComma := string(r.Comma)

	// Regular expression to split records by commas or quotes using following rules:
	// 1) Match quoted fields with quoted (with pairs of double quotes) words:
	//    "         match single starting double quote
	//      [^,"]*  followed with any non-comma (r.Comma) and double quote characters
	//      ""      followed with pair of opening double quotes
	//        .+    followed with at least one character
	//      ""      followed with pair of closing double quotes
	//      [^,"]*  followed with any non-comma (r.Comma) and double quote characters
	//    "         finished with single ending double quote
	//    |         OR
	// 2) Match any string between commas, bypassing commas between pairs of double quotes:
	//    [^,\"]*   match any non-comma (r.Comma) and double quote characters
	//    ""        followed with pair of opening double quotes
	//      [^\"]+  followed with at least one non double quote character
	//    ""        followed with pair of closing double quotes
	//    [^,\"]*   finished with any non-comma (r.Comma) and double quote characters
	//    |         OR
	// 3) Match any string between commas
	//    [^,]*
	splitRgx := regexp.MustCompile("\"[^" + strComma + "\"]*\"\".+\"\"[^" + strComma + "\"]*\"|[^" + strComma + "\"]*\"\"[^\"]+\"\"[^" + strComma + "\"]*|\"[^\"]+\"|[^" + strComma + "]*")
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
	specialCharsGroup := string(w.Comma)
	specialCharsRgx := regexp.MustCompile("^[^\"]*[" + specialCharsGroup + "][^\"]*$")

	for idx, r := range record {
		if specialCharsRgx.MatchString(r) {
			r = "\"" + r + "\""
		} else {
			r = strings.Replace(r, "\"", "\"\"", -1)
		}

		record[idx] = r
	}

	str := strings.Join(record, string(w.Comma))

	if w.pos > 0 {
		str = "\n" + str
		w.pos++
	}

	w.pos += len(str)

	_, err := w.w.Write([]byte(str))

	return err
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w, 0, ','}
}
