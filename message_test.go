package main

import (
	"encoding/json"
	"reflect"
	"strings"
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

// Test_Message_unfurlFieldsAreNeverOmitted guards against reintroducing `omitempty` on
// UnfurlLinks/UnfurlMedia. With omitempty, a false value (the exact value a caller sets to turn
// unfurling off) would be dropped from the JSON instead of sent, since false is also the zero
// value omitempty checks for.
func Test_Message_unfurlFieldsAreNeverOmitted(t *testing.T) {
	msg := Message{Channel: "#general", Text: "hello", UnfurlLinks: false, UnfurlMedia: false}

	b, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	got := string(b)
	if !strings.Contains(got, `"unfurl_links":false`) {
		t.Errorf("expected marshaled JSON to explicitly contain \"unfurl_links\":false, got: %s", got)
	}
	if !strings.Contains(got, `"unfurl_media":false`) {
		t.Errorf("expected marshaled JSON to explicitly contain \"unfurl_media\":false, got: %s", got)
	}
}
