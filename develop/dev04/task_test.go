package main_test

import (
	"dev04"
	"testing"
)

// Проверяет удаление дублей, сортировку, длинну
func TestCollectAnagrams(t *testing.T) {
	answer := map[string][]string{}
	answer["ab"] = []string{"ab", "ba"}

	t.Run("base", func(t *testing.T) {
		res := main.CollectAnagrams([]string{"ab", "ba", "ab", "a"})
		if len(res) == len(answer) {
			if val, ok := res["ab"]; ok {
				if len(res["ab"]) == len(answer["ab"]) {
					if val[0] == "ab" && val[1] == "ba" {
						return
					}
				}
			}
		}
		t.Errorf("\nwant: \n\t%#v \ngot: \n\t%#v", answer, res)
	})
}
