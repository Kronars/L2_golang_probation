package main_test

import (
	"dev06"
	"reflect"
	"testing"
)

var t_rows = []string{
	"a	1 sep 2",
	"b	1 sep 2",
	"c	1 sep 2",
	"d	1 se p 2",
}

var t_cols = [][]string{
	{"a", "1 sep 2"},
	{"b", "1 sep 2"},
	{"c", "1 sep 2"},
	{"d", "1 se p 2"},
}

var t_col = `1 sep 2
1 sep 2
1 sep 2
1 se p 2
`

var t_out_of_bounds = `a	1 sep 2
b	1 sep 2
c	1 sep 2
d	1 se p 2
`

func TestCUT_SplitColumns(t *testing.T) {
	tests := []struct {
		name string
		c    main.CUT
		want [][]string
	}{
		{"base", main.CUT{main.Params{Delimeter: "	"}, t_rows}, t_cols},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.SplitColumns(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CUT.SplitColumns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCUT_Cut(t *testing.T) {
	tests := []struct {
		name string
		c    main.CUT
		want string
	}{
		{"base", main.CUT{main.Params{Delimeter: "	", Fields: []int{1}}, t_rows}, t_col},
		{"out of bounds", main.CUT{main.Params{Delimeter: "	", Fields: []int{3}}, t_rows}, t_out_of_bounds},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Cut(); got != tt.want {
				t.Errorf("CUT.Cut() = %v, want %v", got, tt.want)
			}
		})
	}
}
