package main

import (
	"strconv"
	"testing"
)

func TestUnpackPositive(t *testing.T) {
	var tests = []struct {
		inp, out string
	}{
		{"b", "b"},
		{"b3", "bbb"},
		{`b\3\5`, "b35"},
		{`b\\`, `b\`},
		{`b\\2`, `b\\`},
	}

	for id, test := range tests {
		name_id := strconv.FormatInt(int64(id), 10)
		t.Run(name_id, func(t *testing.T) {
			res, _ := Unpack(test.inp)
			if res != test.out {
				t.Errorf("want: %s \t got: %s ", test.inp, res)
			}
		})
	}

}
