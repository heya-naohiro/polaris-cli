package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func captureStdout(f func()) string {
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	stdout := os.Stdout
	os.Stdout = w
	f()

	os.Stdout = stdout
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()

}

func TestSearch(t *testing.T) {
	d := t.TempDir()
	t.Logf("outdir: %s", d)

	s := NewSearch(d)

	s.AddDocument("testdata/samplefile1.txt")
	s.AddDocument("testdata/samplefile2.txt")
	s.AddDocument("testdata/samplefile3.txt")
	s.AddDocument("testdata/samplefile4.txt")

	out := captureStdout(func() {
		s.QueryPrint("猫")
	})

	// 吾輩は猫である
	if !strings.Contains(out, "samplefile2.txt") {
		t.Errorf("Query Print not contains 猫")
	}

}
