package main

import (
	"strings"
	"testing"
)

func TestSort(t *testing.T) {
	sorted_cases := map[string]string{
		"base": `1 АА пре1
2 АБ пре2
3 АВ пре3
4 АГ пре4
5 Б пре5
5 А пре5
`,
		"lexem": `5 А пре5
1 АА пре1
2 АБ пре2
3 АВ пре3
4 АГ пре4
5 Б пре5
`,
		"number": `5 А пре5
5 Б пре5
4 АГ пре4
3 АВ пре3
2 АБ пре2
1 АА пре1
`,
	}

	var tests = []struct {
		name, out string
		params    Params
	}{
		{"base case", sorted_cases["base"],
			Params{
				k:    "0.0",
				sep:  " ",
				n:    false,
				r:    false,
				u:    false,
				path: "sample_for_tests",
			},
		},
		{"lexem case", sorted_cases["lexem"],
			Params{
				k:    "1.0",
				sep:  " ",
				n:    false,
				r:    false,
				u:    false,
				path: "sample_for_tests",
			},
		},
		{"number case", sorted_cases["number"],
			Params{
				k:    "2.3",
				sep:  " ",
				n:    true,
				r:    true,
				u:    false,
				path: "sample_for_tests",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sorter := MakeSorter(test.params)
			out, _ := sorter.Sort()
			if strings.Compare(test.out, out) != 0 {
				t.Errorf("want: \n%s\n got: \n%s", test.out, out)
			}
		})
	}

}
