package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type Params struct {
	k       string
	sep     string
	n, r, u bool
	path    string
}

type Sorter struct {
	p Params
}

func PrepareParams(k_key, sep_key *string, n_key, r_key, u_key *bool, path []string) Params {
	if len(path) == 0 {
		fmt.Println("Не передан путь к файлу для сортировки")
		os.Exit(1)
	}
	return Params{*k_key, *sep_key, *n_key, *r_key, *u_key, path[0]}
}

func MakeSorter(p Params) Sorter {
	return Sorter{p}
}

func (s *Sorter) Sort() (string, error) {
	if err := s.checkPath(); err != nil {
		return "", fmt.Errorf("ошибка доступа: %w", err)
	}

	rows := s.parseRows()
	splited := s.splitRows(rows)
	sorted := s.sortSplitedRows(splited)    // ключи k n r
	concated := s.ConcatSplitedRows(sorted) // ключ d
	return concated, nil
}

func main() {
	sep_key := flag.String("sep", " ", "разделительный символ для сортировки")
	k_key := flag.String("k", "0.0", `указание колонки для сортировки, формат: Н.С 
	Н-номер колонки, С-номер символа с котрого начать учитывать символы для сортировки`)
	n_key := flag.Bool("n", false, "сортировать по числовому значению")
	r_key := flag.Bool("r", false, "сортировать в обратном порядке")
	u_key := flag.Bool("u", false, "не выводить повторяющиеся строки")

	flag.Parse()

	path := flag.Args()

	params := PrepareParams(k_key, sep_key, n_key, r_key, u_key, path)
	sorter := MakeSorter(params)

	out, err := sorter.Sort()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(out)
}

// Гарантирует что файл доступен и откроется
func (s *Sorter) checkPath() error {
	info, err := os.Stat(s.p.path)
	if os.IsNotExist(err) {
		return fmt.Errorf("путь к файлу %s не сущесвует", err)
	}
	if info.IsDir() {
		return fmt.Errorf("путь является директорией - %s", s.p.path)
	}
	file, err := os.Open(s.p.path)
	if err != nil {
		return fmt.Errorf("файл недоступен, причина: %w", err)
	}
	defer file.Close()

	return nil
}

// Парсинг строк
func (s *Sorter) parseRows() []string {
	var rows []string
	file, _ := os.Open(s.p.path)
	defer file.Close()

	f_scan := bufio.NewScanner(file)
	for f_scan.Scan() {
		rows = append(rows, f_scan.Text())
	}

	return rows
}

// Разбиение строк по символу разделителю
func (s *Sorter) splitRows(rows []string) [][]string {
	var res [][]string
	for _, row := range rows {
		splited := strings.Split(row, s.p.sep)
		res = append(res, splited)
	}
	return res
}

// Выполнение сортировки
func (s *Sorter) sortSplitedRows(rows [][]string) [][]string {
	col_numb, char_offset := s.getK()
	data := &Rows{rows, col_numb, char_offset, s.p}

	if s.p.r {
		sort.Sort(sort.Reverse(data))
	} else {
		sort.Sort(data)
	}

	return data.rows
}

// Получение параметра k
func (s *Sorter) getK() (int, int) {
	res := strings.Split(s.p.k, ".")
	if len(res) != 2 {
		fmt.Printf("Неправильно передан параметр k - %s\nСмотри --help\n", s.p.k)
		os.Exit(1)
	}
	column, err_1 := strconv.ParseUint(res[0], 10, 32)
	offset, err_2 := strconv.ParseUint(res[1], 10, 32)
	if err_1 != nil || err_2 != nil {
		fmt.Printf(`Переданы неккоректные значения в параметр k - %s
		Ошибка: %s или %s\n`, s.p.k, err_1, err_2)
		os.Exit(1)
	}
	return int(column), int(offset)
}

// Структура - строка, для стандартной либы sort
type Rows struct {
	rows                [][]string
	column, char_offset int
	p                   Params
}

func (r *Rows) Len() int {
	return len(r.rows)
}

// Сравнение при сортировке, меньше ли iый элемент чем jый
func (r *Rows) Less(i, j int) bool {
	// Проверки на захождение за длинну столбцов и символов, такие элементы считаются наименьшими

	if r.column >= len(r.rows[i]) || r.column >= len(r.rows[j]) {
		return false
	}

	val_l, val_r := r.rows[i][r.column], r.rows[j][r.column]
	if r.char_offset >= len(val_l) || r.char_offset >= len(val_r) {
		return false
	}

	// Если сортировка по числовому значению, сравнить как числа
	if r.p.n {
		val_l, _ := strconv.ParseInt(string([]rune(val_l)[r.char_offset:]), 10, 32) // строка -> массив рун -> срез -> строка -> парсинг в число
		val_r, _ := strconv.ParseInt(string([]rune(val_r)[r.char_offset:]), 10, 32)
		return int(val_l) < int(val_r)
	} // Если сортировка по строкам, сравнить как строки

	if strings.Compare(val_l[r.char_offset:], val_r[r.char_offset:]) <= 0 {
		return true
	} else {
		return false
	}
}

func (r *Rows) Swap(i, j int) {
	r.rows[i], r.rows[j] = r.rows[j], r.rows[i]
}

// Склейка результаотв сортировки в единую строку
func (s *Sorter) ConcatSplitedRows(rows [][]string) string {
	var out string

	for _, columns := range rows {
		row := strings.Join(columns, s.p.sep)
		if s.p.u {
			if strings.Contains(out, row) {
				continue
			}
		}
		out += row
		out += "\n"
	}
	return out
}
