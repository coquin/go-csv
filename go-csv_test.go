package gocsv

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestNewReader(t *testing.T) {
	// TODO: write it
}

func TestRead(t *testing.T) {
	var rec []string
	var err error

	csvReader := NewReader(strings.NewReader("foo,bar,baz\nnyan,cat,wat\npeow,,"))

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
	expectedSecond := []string{"nyan", "cat", "wat"}

	if !compareRecords(rec, expectedSecond) {
		t.Log("Read should read first string")
		t.Fatal("expected", expectedSecond, "got", rec)
	}

	rec, err = csvReader.Read()
	expectedLast := []string{"peow", "", ""}
	if !compareRecords(rec, expectedLast) {
		t.Log("Read should read empty fields as empty strings")
		t.Fatal("expected", expectedLast, "got", rec)
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

	csvReader := NewReader(strings.NewReader("\"foo\",\"\"bar\"\",\"\"\"baz\"\"\"\nnyan,Tyrion \"\"Imp\"\" Lannister,wat"))

	rec, _ = csvReader.Read()
	expectedFirst := []string{"foo", "\"bar\"", "\"baz\""}

	if !compareRecords(rec, expectedFirst) {
		t.Log("Read should unquote qouted fields and escape double quotes")
		t.Error("expected", expectedFirst, "got", rec)
	}

	rec, _ = csvReader.Read()
	expectedLast := []string{"nyan", "Tyrion \"Imp\" Lannister", "wat"}

	if !compareRecords(rec, expectedLast) {
		t.Log("Read should unescape doubled quotes")
		t.Error("expected", expectedLast, "got", rec)
	}
}

func TestReadSpecialCharsInQuotedFields(t *testing.T) {
	testStr := []string{
		"Track,\"Daddy, Brother, Lover, Little Boy\",\"Lean Into It\",1991,Mr. Big",
		"\"\"\"Daddy, Brother, Lover, Little Boy\"\"\"",
	}
	csvReader := NewReader(strings.NewReader(strings.Join(testStr, "\n")))

	rec, _ := csvReader.Read()
	expectedFirst := []string{"Track", "Daddy, Brother, Lover, Little Boy", "Lean Into It", "1991", "Mr. Big"}

	if !compareRecords(rec, expectedFirst) {
		t.Log("Read should keep special characters (commas) in quoted fields")
		t.Error("expectedFirst", expectedFirst, "got", rec)
	}

	rec, _ = csvReader.Read()
	expectedLast := []string{"\"Daddy, Brother, Lover, Little Boy\""}

	if !compareRecords(rec, expectedLast) {
		t.Log("Read should keep special characters (commas) in quoted fields")
		t.Error("expectedLast", expectedLast, "got", rec)
	}
}

func TestReadFromFile(t *testing.T) {
	sample, err := os.Open("test-csv/1.csv")

	if err != nil {
		t.Fatal("Error during reading sample CSV file:", err)
	}

	defer sample.Close()

	csvReader := NewReader(sample)

	expected := [][]string{
		{"Jaime", "Lannister", "Kingslayer", ""},
		{"Robert", "Baratheon", "", "dead"},
		{"Eddard", "Stark", "", "dead"},
		{"Tyrion", "Lannister", "Imp", ""},
		{"John", "Snow", "You know nothing", "dead (GRM you're bastard!)"},
		{"Daenerys", "Targarien", "\"Stormborn\", \"Unburnt\"", ""},
	}

	for idx := 0; ; idx++ {
		rec, err := csvReader.Read()

		if err == io.EOF {
			break
		}
		if !compareRecords(rec, expected[idx]) {
			t.Log("Should read records from file correctly")
			t.Fatal("expected", expected[idx], "got", rec)
		}
		if err != nil {
			t.Log("Read should not return error on success")
			t.Error("expected", nil, "got", err)
		}
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
		t.Log("Write should append new string")
		t.Fatal("expected", "foo,bar,baz\\nnyan,cat,wat", "got", buf.String())
	}
}

func TestWriteCustomComma(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	csvWriter := NewWriter(buf)
	csvWriter.Comma = '|'

	csvWriter.Write([]string{"foo", "bar", "baz"})
	if buf.String() != "foo|bar|baz" {
		t.Log("Write should write array of strings using custom delimiter")
		t.Fatal("expected", "foo|bar|baz", "got", buf.String())
	}
}

func TestWriteQuotesEscaped(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	csvWriter := NewWriter(buf)

	csvWriter.Write([]string{"nyan", "Tyrion \"Imp\" Lannister", "wat"})
	expected := "nyan,Tyrion \"\"Imp\"\" Lannister,wat"

	if buf.String() != expected {
		t.Log("Write should escape quotes inside field and wrap field in quotes")
		t.Fatal("expected", expected, "got", buf.String())
	}
}

func TestWriteSpecialCharacters(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	csvWriter := NewWriter(buf)

	csvWriter.Write([]string{"foo, bar", "baz"})
	expected := "\"foo, bar\",baz"

	if buf.String() != expected {
		t.Log("Write should put field with special characters in quotes")
		t.Fatal("expected", expected, "got", buf.String())
	}
}

func TestWriteSpecialCharactersQuoted(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	csvWriter := NewWriter(buf)

	csvWriter.Write([]string{"foo \"bar, baz\"", "quix"})
	expected := "foo \"\"bar, baz\"\",quix"

	if buf.String() != expected {
		t.Log("Write should not put field with special characters in quotes special chars already in quotes (should escape with pairs of quotes instead)")
		t.Fatal("expected", expected, "got", buf.String())
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
