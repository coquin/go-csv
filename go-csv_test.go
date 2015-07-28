package gocsv

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestNewReader(t *testing.T) {
	// TODO: write it
}

func TestRead(t *testing.T) {
	csvReader := NewReader(strings.NewReader("foo,bar,baz\nnyan,cat,wat"))

	testReadStr(t, csvReader, []string{"foo", "bar", "baz"})
	testReadStr(t, csvReader, []string{"nyan", "cat", "wat"})
	testReadEnd(t, csvReader)
}

func TestReadInLoop(t *testing.T) {
	var str string

	csvReader := NewReader(strings.NewReader("foo,bar,baz\nnyan,cat,wat"))

	for {
		rec, err := csvReader.Read()
		str += strings.Join(rec, "")

		if err == io.EOF {
			break
		}
	}

	if str != "foobarbaznyancatwat" {
		t.Log("Should read till the end of reader in loop")
		t.Fatal("expected", "foobarbaznyancatwat", "got", str)
	}
}

func TestReadCustomComma(t *testing.T) {
	csvReader := NewReader(strings.NewReader("foo|bar|baz"))
	csvReader.Comma = '|'

	testReadStr(t, csvReader, []string{"foo", "bar", "baz"})
}

func TestWrite(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	csvWriter := NewWriter(buf)

	csvWriter.Write([]string{"foo", "bar", "baz"})
	if buf.String() != "foo,bar,baz" {
		t.Log("Write should write array of strings")
		t.Fatal("expected", "foo,bar,baz", "got", buf.String())
	}

	csvWriter.Write([]string{"nyan", "cat", "wat"})
	if buf.String() != "foo,bar,baz\nnyan,cat,wat" {
		t.Log("Write should append array of strings")
		t.Fatal("expected", "foo,bar,baz\\nnyan,cat,wat", "got", buf.String())
	}
}

func testReadStr(t *testing.T, r *Reader, expected []string) {
	rec, err := r.Read()

	if !compareRecords(rec, expected) {
		t.Log("Read should read a string")
		t.Fatal("expected", expected, "got", rec)
	} else if err != nil {
		t.Log("Read should not return error on success")
		t.Error("expected", nil, "got", err)
	}
}

func testReadEnd(t *testing.T, r *Reader) {
	rec, err := r.Read()

	if !compareRecords(rec, []string{}) {
		t.Log("Reading the end of the stream should return empty string")
		t.Fatal("expected", "(empty string)", "got", rec)
	}
	if err != io.EOF {
		t.Log("Reading the end of the stream should return EOF error")
		t.Error("expected", io.EOF, "got", err)
	}
}

func compareRecords(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
