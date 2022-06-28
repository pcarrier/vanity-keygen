package main

import (
	"bytes"
	"testing"
)

func Test_incr(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		out  []byte
	}{
		{
			name: "zero",
			in:   []byte{0, 0},
			out:  []byte{0, 1},
		},
		{
			name: "second digit",
			in:   []byte{1, 255},
			out:  []byte{2, 0},
		},
		{
			name: "overflow",
			in:   []byte{255, 255},
			out:  []byte{0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			incr(&tt.in)
			if bytes.Compare(tt.in, tt.out) != 0 {
				t.Fail()
			}
		})
	}
}
