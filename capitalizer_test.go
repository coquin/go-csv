package gocsv

import (
	"bytes"
	"strings"
	"testing"
)

func TestCapitalizerRead(t *testing.T) {
	s := strings.NewReader("Foobar")
	buf := new(bytes.Buffer)
	capitalizerReader := NewCapitalizerReader(s)
	buf.ReadFrom(capitalizerReader)

	if buf.String() != "FOOBAR" {
		t.Log("CapitalizeReader should capitalize text string.")
		t.Error("expected", "FOOBAR", "got", buf.String())
	}
}
