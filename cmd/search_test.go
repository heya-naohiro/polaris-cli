package cmd

import "testing"

func TestSearch(t *testing.T) {
	d := t.TempDir()
	t.Logf("outdir: %s", d)

	s := NewSearch(d)

	s.AddDocument("testdata/samplefile1.txt")
	s.AddDocument("testdata/samplefile2.txt")
	s.AddDocument("testdata/samplefile3.txt")
	s.AddDocument("testdata/samplefile4.txt")

	s.QueryPrint("音楽")

}
