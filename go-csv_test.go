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
	var rec []string
	var err error

	csvReader := NewReader(strings.NewReader("foo,bar,baz\nnyan,cat,wat"))

	rec, err = csvReader.Read()
	expectedFirst := []string{"foo", "bar", "baz"}

	if !compareRecords(rec, expectedFirst) {
		t.Log("Read should read first string")
		t.Fatal("expected", expectedFirst, "got", rec)
	}
	if err != nil {
		t.Log("Read should not return error on success")
		t.Error("expected", nil, "got", err)
	}

	rec, err = csvReader.Read()
	expectedLast := []string{"nyan", "cat", "wat"}

	if !compareRecords(rec, expectedLast) {
		t.Log("Read should read first string")
		t.Fatal("expected", expectedLast, "got", rec)
	}
	if err != nil {
		t.Log("Read should not return error on success")
		t.Error("expected", nil, "got", err)
	}

	rec, err = csvReader.Read()
	if !compareRecords(rec, []string{}) {
		t.Log("Reading the end of the stream should return empty string")
		t.Fatal("expected", "(empty string)", "got", rec)
	}
	if err != io.EOF {
		t.Log("Reading the end of the stream should return EOF error")
		t.Error("expected", io.EOF, "got", err)
	}
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

	rec, err := csvReader.Read()
	expected := []string{"foo", "bar", "baz"}

	if !compareRecords(rec, expected) {
		t.Log("Read should read string with custom fields delimeter")
		t.Fatal("expected", expected, "got", rec)
	}
	if err != nil {
		t.Log("Read should not return error on success")
		t.Error("expected", nil, "got", err)
	}
}

func TestReadQuotes(t *testing.T) {
	var rec []string

	csvReader := NewReader(strings.NewReader("foo,\"bar\",baz\nnyan,\"Tyrion \"\"Imp\"\" Lannister\",wat"))

	rec, _ = csvReader.Read()
	expectedFirst := []string{"foo", "bar", "baz"}

	if !compareRecords(rec, expectedFirst) {
		t.Log("Read should unquote qouted fields")
		t.Error("expected", expectedFirst, "got", rec)
	}

	rec, _ = csvReader.Read()
	expectedLast := []string{"nyan", "Tyrion \"Imp\" Lannister", "wat"}

	if !compareRecords(rec, expectedLast) {
		t.Log("Read should unescape doubled quotes")
		t.Error("expected", expectedLast, "got", rec)
	}
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
