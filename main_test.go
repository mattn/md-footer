package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestSimple(t *testing.T) {
	files, err := filepath.Glob("testdir/*.in")
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(files)

	for _, file := range files {
		t.Log(file)
		f, err := os.Open(file)
		if err != nil {
			t.Fatal(err)
		}
		var out bytes.Buffer
		err = convert(f, &out)
		f.Close()
		if err != nil {
			t.Fatal(err)
		}

		b, err := ioutil.ReadFile(file[:len(file)-3] + ".out")
		if err != nil {
			t.Fatal(err)
		}
		want := string(b)
		got := out.String()

		if want != got {
			t.Fatalf("\nwant: %v\nbut: %v\n", want, got)
		}
	}
}
