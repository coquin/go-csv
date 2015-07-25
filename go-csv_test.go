package gocsv

import (
	"io"
	"strings"
	"testing"
)

func TestNewReader(t *testing.T) {
	// TODO: write it
}

func TestRead(t *testing.T) {
	csvReader := NewReader(strings.NewReader("foo\nbar\nbaz"))

	testReadStr(t, csvReader, "foo")
	testReadStr(t, csvReader, "bar")
	testReadStr(t, csvReader, "baz")
	testReadEnd(t, csvReader)
}

func testReadStr(t *testing.T, r *Reader, expected string) {
	rec, err := r.Read()

	if rec != expected {
		t.Log("Read should read a string")
		t.Fatal("expected", expected, "got", rec)
	} else if err != nil {
		t.Log("Read should not return error on success")
		t.Error("expected", nil, "got", err)
	}
}

func testReadEnd(t *testing.T, r *Reader) {
	rec, err := r.Read()

	if rec != "" {
		t.Log("Reading the end of the stream should return empty string")
		t.Fatal("expected", "(empty string)", "got", rec)
	} else if err != io.EOF {
		t.Log("Reading the end of the stream should EOF error")
		t.Error("expected", io.EOF, "got", err)
	}
}
