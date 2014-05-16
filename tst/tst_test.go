package tst

import (
	"strings"
	"testing"
)

const (
	lorem = "Lorem ipsum dolor sit amet consectetur adipiscing elit " +
		"Donec tempus placerat quam nec vehicula " +
		"Etiam eleifend justo vitae lorem tincidunt feugiat sed ut augue " +
		"Nullam hendrerit ultrices risus " +
		"Aliquam elementum nunc ac arcu gravida a dictum mauris semper " +
		"Nulla sit amet lectus sit amet enim pellentesque dictum at vitae metus " +
		"Praesent tincidunt lacus non cursus sodales ligula nisl porttitor leo vitae mollis metus orci a massa " +
		"Nulla facilisi " +
		"Nullam rutrum pharetra sagittis " +
		"Phasellus rutrum nisi odio eu venenatis odio dapibus et " +
		"Praesent ultricies erat vitae euismod pellentesque " +
		"Phasellus tempus elit et felis tempus non consectetur sapien accumsan " +
		"Pellentesque elementum metus ultrices neque dictum tempor " +
		"Etiam rhoncus lacinia luctus " +
		"Vestibulum tempus metus vel justo feugiat sagittis " +
		"Fusce tempus odio eu leo tristique pretium " +
		"Fusce non dolor lectus"
)

var (
	words = strings.Fields(lorem)
	wlen  = len(words)
)

func TestInsertGet(t *testing.T) {
	tst := New()
	if tst == nil {
		t.Fatal("New() is nil but it should be *TST")
	}

	for i := 0; i < wlen; i++ {
		if !tst.Insert(words[i], words[i]) {
			t.Fatalf("failed to insert %q (%d of %d)", words[i], i+1, wlen)
			break
		}
	}

	if tst.Size == 0 {
		t.Fatalf("TST size is 0 but expected <= %d", wlen)
	}

	// test values at start, middle and end of list
	for _, key := range []string{"Lorem", "ipsum", "dolor", "pretium", "vel", "tristique", "cursus", "enim", "rutrum"} {
		ret := tst.Get(key)
		if ret != key {
			t.Errorf("expected %q but got %q", key, ret)
		}
	}

	// test unexistent values
	for _, key := range []string{"horse", "tennis", "automobile"} {
		ret := tst.Get(key)
		if ret != nil {
			t.Errorf("expected nil but got %q", ret)
		}
	}

	// test partial matches (should all return nil)
	for _, key := range []string{"Lore", "curs", "tristi"} {
		ret := tst.Get(key)
		if ret != nil {
			t.Errorf("expected nil but got %q", ret)
		}
	}

	// test wildcard matches
	if ret := tst.Match("temp__"); strings.Join(ret, " ") != "tempor tempus" {
		t.Errorf("wildcard match failed for %q got %q", "temp__", ret)
	}

	// test prefix matching
	if ret := tst.Prefix("temp"); strings.Join(ret, " ") != "tempor tempus" {
		t.Errorf("prefix match failed for %q got %q", "temp", ret)
	}
}

// TODO: benchmarks
