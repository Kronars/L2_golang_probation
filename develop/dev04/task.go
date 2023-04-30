package main

import (
	"fmt"
	"sort"
	"strings"
)

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func NewAnamSet() AnamSet {
	return AnamSet{map[string][]string{}}
}

type AnamSet struct {
	set map[string][]string
}

func (a *AnamSet) AddFirst(word string) {
	a.set[word] = []string{word}
}

func (a *AnamSet) Append(word, key string) {
	a.set[key] = append(a.set[key], word)
}

// Удаляет дупликаты, сортирует, удаляет короткие
func (a *AnamSet) GetSet() map[string][]string {
	for key, words := range a.set {
		if len(words) == 1 { // Множества из одного элемента не должны попасть в результат
			delete(a.set, key)
			continue
		}
		sort.Strings(words)

		to_delete := []int{}
		for next := 1; next < len(words); next++ {
			prev := next - 1
			if words[prev] == words[next] {
				to_delete = append(to_delete, prev) // Что бы не изменять список по которомy итерирование
			}
		}

		for _, id := range to_delete {
			words = append(words[:id], words[id+1:]...) // Создаётся новый срез
		}
		a.set[key] = words // Перезапись
	}
	return a.set
}

// Проверяет есть ли в множестве анаграммы переданного слова
// Проверяет функцией IsAnagram переданное слово с каждым словoм в множестве
// true: в множестве есть анаграмма, возвращает её ключ false: анаграмм нет
func (a *AnamSet) Check(word_a string) (bool, string) {
	for key, words := range a.set {
		for _, word_b := range words {
			if IsAnagram(word_a, word_b) {
				return true, key // В множестве найдена анаграмма, вернуть ключ по которому записанна
			}
		}
	}
	return false, word_a // В множестве нет анаграммы переданного слова, вернуть его же
}

// false: разные или одинаковые слова true: являются анаграммами
func IsAnagram(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	b_copy := []rune(b)

	for _, letter_a := range a {
		// Если в слове b нет буквы из слова a
		if idx := strings.Index(b, string(letter_a)); idx == -1 {
			return false
		} else { // Удалить найденную букву из слова b, что бы не найти её снова
			b_copy = []rune(strings.Replace(string(b_copy), string(a[idx]), "", 1))
		}
	}
	return true
}

func CollectAnagrams(arr []string) map[string][]string {
	set := NewAnamSet()
	for _, word := range arr {
		word = strings.ToLower(word)
		if contains, key := set.Check(word); contains {
			set.Append(word, key)
		} else {
			set.AddFirst(word)
		}
	}
	return set.GetSet()
}

func main() {
	arr := []string{"пятак", "пятка", "тяпка", "типы", "типа", "пати", "типа"}
	res := CollectAnagrams(arr)
	fmt.Println(res)
}
