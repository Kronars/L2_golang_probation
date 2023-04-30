package main_test

import (
	"dev05"
	"strings"
	"testing"
)

var data = []string{
	"0",
	"1 c",
	"2 c",
	"3 C",
	"4",
}

func TestFilter(t *testing.T) {
	// Обычные параметры, в тестах изменяется эта структура
	// что бы не заполнять с 0 в каждом тесте

	stick := func(a []string) string {
		return strings.Join(a, "\n")
	}
	tests := []struct {
		name    string
		params  main.Params
		pattern []string
		out     string
	}{
		{"base", main.Params{}, []string{"c"}, stick(data[1:3])},
		{"shift C", main.Params{C: 1}, []string{"1"}, stick(data[:3])},
		{"count", main.Params{Count: 2}, []string{"c"}, stick(data[1:3])},
		{"ignore case", main.Params{Ignore: true}, []string{"C"}, stick(data[1:4])},
		{"fixed", main.Params{F: true}, []string{"4"}, data[4]},
		{"invert", main.Params{InVert: true}, []string{"c"}, stick([]string{data[0], data[3], data[4]})},
		{"line numbers", main.Params{Nums: true}, []string{"0", "4"}, "0: 0\n4: 4"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			grepa := main.NewGREP(test.params, test.pattern, data)
			result := grepa.Filter()
			if result != test.out {
				t.Errorf("\ngot: \n%s \nwant: \n%s ", result, test.out)
			}
		})
	}
}
