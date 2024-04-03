package common

import (
	"testing"
)

func TestIsPrefixMatched(t *testing.T) {
	tests := []struct {
		name     string
		want     bool
		path     string
		prefixes []string
	}{
		{
			name:     "Matched",
			want:     true,
			path:     "/api/checkout/?a=1",
			prefixes: []string{"/api/checkout", "/api/2/match"},
		},
		{
			name:     "Not matched",
			want:     false,
			path:     "/api/checkout/?a=1",
			prefixes: []string{"/api/2/match"},
		},
	}
	for _, tt := range tests {
		isMatched := IsPrefixMatched(tt.path, tt.prefixes)
		if tt.want != isMatched {
			t.Errorf("isMatch want %v but got %v", tt.want, isMatched)
		}
	}
}
