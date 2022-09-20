package main

import (
	"reflect"
	"testing"
)

func Test_parseFields(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		wantFs []Field
	}{
		{
			name:   "Newline in release notes",
			s:      "Release notes|line1\\nline2",
			wantFs: []Field{{Title: "Release notes", Value: "line1\nline2"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotFs := parseFields(tt.s); !reflect.DeepEqual(gotFs, tt.wantFs) {
				t.Errorf("parseFields() = %v, want %v", gotFs, tt.wantFs)
			}
		})
	}
}
