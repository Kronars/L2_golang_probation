package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	inp := `b\\2` // Ранил
	res, err := Unpack(inp)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

// Распаковка строки
func Unpack(str string) (string, error) {
	var res []rune // Заполняется в цикле

	if err := is_correct(str); err != nil {
		return "", err
	}

	for id := 0; id < len(str); id++ {
		char := rune(str[id])

		// Если текущий символ буква - добавить к результату
		if unicode.IsLetter(char) {
			res = append(res, char)
			continue
		}

		// Если пред символ слеш И текущий симв слеш ИЛИ цифра И последний симв в результате не слеш
		if str[id-1] == '\\' && (unicode.IsDigit(char) || char == '\\') && res[len(res)-1] != '\\' {
			res = append(res, char)
			continue
		}

		// Если цифра, распарсить, повторить n раз и добавить к резульатату
		if unicode.IsDigit(char) {
			amount, _ := strconv.ParseInt(string(char), 10, 32)
			sub_s := repeat(res[len(res)-1], amount-1)
			res = append(res, sub_s...)
		}
	}
	return string(res), nil
}

func repeat(c rune, amount int64) []rune {
	return []rune(strings.Repeat(string(c), int(amount)))
}

func is_correct(str string) error {
	if len(str) == 0 {
		return InccorectStringError{str}
	}

	if unicode.IsDigit(rune(str[0])) {
		return InccorectStringError{str}
	}

	if strings.Contains(str, "0") {
		return InccorectStringError{str}
	}

	return nil
}

type InccorectStringError struct {
	str string
}

func (e InccorectStringError) Error() string {
	return fmt.Sprintf(`The string has an incorrect format: %s
Restrictions: the first character is not a digit, no zeros, len > 0`, e.str)
}
