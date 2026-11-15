package main

import (
	"bytes"
	"testing"
)

func TestConvert(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  string
	}{
		{"1.5 GiB", "1610612736 bytes\n"},
		{"1 GB", "1000000000 bytes\n"},
		{"1 KiB", "1024 bytes\n"},
		{"0 B", "0 bytes\n"},
	} {
		t.Run(tc.input, func(t *testing.T) {
			var out bytes.Buffer
			if err := run([]string{tc.input}, &out); err != nil {
				t.Fatal(err)
			}
			if out.String() != tc.want {
				t.Fatalf("got %q, want %q", out.String(), tc.want)
			}
		})
	}
}

func TestInvalidInput(t *testing.T) {
	for _, args := range [][]string{nil, {"1 GB", "2 GB"}, {"oops"}, {"-1 GiB"}, {"1 potato"}} {
		var out bytes.Buffer
		if err := run(args, &out); err == nil {
			t.Fatalf("expected error for %q", args)
		}
		if out.Len() != 0 {
			t.Fatalf("unexpected output for %q: %q", args, out.String())
		}
	}
}
