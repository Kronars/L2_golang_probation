package main

import (
	"reflect"
	"testing"
)

func Test_getLinks(t *testing.T) {
	tests := []struct {
		name     string
		htmlFile string
		want     []string
	}{
		{"base", "test_links", []string{"http 1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getLinks(tt.htmlFile); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getLinks() = %v, want %v", got, tt.want)
			}
		})
	}
}
