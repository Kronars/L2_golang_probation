package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type Params struct {
	Fields    []int
	Delimeter string
	Separated bool
}

type CUT struct {
	Conf Params
	Rows []string
}

func (c *CUT) Cut() string {
	rows_columns := c.SplitColumns()            // Разбивает на колонки, отрабатывает флаг -s
	rows_columns = c.ChooseFields(rows_columns) // Выбирает столбец, если превышает границы - добавляет строку целиком и идёт дальше
	cuted_rows := c.StickTogether(rows_columns) // Объединяет результат в единую строку
	return cuted_rows
}

func main() {
	params, rows := Parse()

	cutter := CUT{params, rows}
	fmt.Println(cutter.Cut())
}

// Разбивает на колонки, если поднят флаг s скипает строки без разделителя
func (c *CUT) SplitColumns() [][]string {
	var r_c [][]string                     // сюда складывается результат
	splitter := func(id int, row string) { // Дефолтное поведение
		cols := strings.Split(row, c.Conf.Delimeter)
		r_c = append(r_c, cols)
	}
	// Поведение при поднятом флаге Separated
	if c.Conf.Separated { // Если результат Split - срез с одной строкой - строка не содержит разделителя надо скипать
		splitter = func(id int, row string) {
			cols := strings.Split(row, c.Conf.Delimeter)
			if len(cols) == 1 {
				return
			}
			r_c[id] = cols
		}
	}

	for id, row := range c.Rows {
		splitter(id, row)
	}
	return r_c
}

// Выбирает переданные столбцы
func (c *CUT) ChooseFields(rows [][]string) [][]string {
	res := make([][]string, len(rows))
	for id, row := range rows {
		for _, column := range c.Conf.Fields {
			if column >= len(row) {
				res[id] = append(res[id], strings.Join(row, c.Conf.Delimeter))
				break
			}
			res[id] = append(res[id], row[column])
		}
	}
	return res
}

// Склеивание результата в одну строку
func (c *CUT) StickTogether(rows_n_cols [][]string) string {
	var res string
	for _, row := range rows_n_cols {
		res += strings.Join(row, c.Conf.Delimeter)
		res += "\n"
	}
	return res
}

// Всякое io

// Парсинг параметров вызова, поддерживаются короткие и длинные верии параметров
func Parse() (Params, []string) {
	var fields string
	fields_usage := "обязательный аргумент, выбрать поля (колонки)"
	flag.StringVar(&fields, "f", "-1", fields_usage+" (короткая версия)")
	flag.StringVar(&fields, "fields", "-1", fields_usage)

	var delimiter string
	del_usage := "использовать другой разделитель"
	flag.StringVar(&delimiter, "d", "	", del_usage+" (короткая версия)")
	flag.StringVar(&delimiter, "delimeter", "	", del_usage)

	var is_sep bool
	sep_usage := "только строки с разделителем"
	flag.BoolVar(&is_sep, "s", false, sep_usage+" (короткая версия)")
	flag.BoolVar(&is_sep, "separated", false, sep_usage)

	flag.Parse()

	if fields == "-1" {
		fmt.Println("Необходимо указать параметр -f")
		os.Exit(1)
	}

	input := ScanStdIn()

	if len(input) == 0 {
		fmt.Println("Программа не получила текст на вход")
		os.Exit(1)
	}

	params := PrepareParams(fields, delimiter, is_sep)

	return params, input
}

// Чтение со стандартнго входа
func ScanStdIn() []string {
	var out []string
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		out = append(out, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "ошибка считывания стандартного входа:", err)
	}
	return out
}

// Парсинг параметра fields и создание параметров
func PrepareParams(fields, delimiter string, is_sep bool) Params {
	f := strings.Split(fields, ",")
	var f_list []int
	for _, v := range f {
		column, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			fmt.Println("Некорректный формат аргумента field. Передайте цифры столбцов через запятую без пробелов. f - ", column)
		}
		f_list = append(f_list, int(column))
	}

	return Params{f_list, delimiter, is_sep}
}
