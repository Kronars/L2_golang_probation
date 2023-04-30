package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type Params struct {
	A, B, C, Count          int
	Ignore, InVert, F, Nums bool
}

type GREP struct {
	p        Params
	patterns []string
	rows     []string
}

func (g *GREP) Filter() string {
	rows := g.ignoreCase()                      // Обработка ignore-case
	core_idxs := g.findRows(rows)               // Обработка fixed, count
	shifted_idxs := g.shiftABC(core_idxs)       // Обработка A, B, C
	inverted_idxs := g.invertIdx(shifted_idxs)  // Обработка invert
	result_rows := g.collectRows(inverted_idxs) // Сбор строк по их индексам
	result_string := g.stickRows(result_rows)   // Склейка в единую строку, обработка line num
	return result_string
}

func NewGREP(p Params, pat []string, rows []string) *GREP {
	return &GREP{p, pat, rows}
}

func main() {
	params, patterns, path := Parse()

	if err := CheckPath(path); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	rows := ParseRows(path)

	grepa := NewGREP(params, patterns, rows)
	result := grepa.Filter()
	fmt.Println(result)
}

// Методы обработки

func (g *GREP) ignoreCase() []string {
	if !g.p.Ignore {
		return g.rows
	}
	copy_rows := make([]string, len(g.rows))
	copy(copy_rows, g.rows)

	for id, row := range copy_rows {
		copy_rows[id] = strings.ToLower(row)
	}
	for id, pat := range g.patterns {
		g.patterns[id] = strings.ToLower(pat)
	}
	return copy_rows
}

func (g *GREP) findRows(rows []string) []int {
	idxs := []int{}

	// Обработка ключа F fixed
	soft := func(a, b string) bool {
		return strings.Contains(a, b)
	}
	fixed := func(a, b string) bool {
		return a == b
	}

	comparator := soft
	if g.p.F {
		comparator = fixed
	}

	// Обработка ключа c count
	var found_count int
	enought := func() bool {
		if g.p.Count == 0 {
			return false
		}
		found_count++
		return found_count >= g.p.Count
	}

	// Сбор индексов строк содержащих паттерн
	for k, row := range rows {
		for _, pat := range g.patterns {
			if comparator(row, pat) {
				idxs = append(idxs, k)
				if enought() {
					return idxs
				}
			}
		}
	}
	return idxs
}

// Добавляет все индексы из сдвигов, сортирует, удаляет дубли
func (g *GREP) shiftABC(idxs []int) []int {
	if g.p.A == 0 && g.p.B == 0 && g.p.C == 0 {
		return idxs
	}

	copy_idxs := make([]int, len(idxs))
	copy(copy_idxs, idxs)

	for _, idx := range idxs {
		before_idx, after_idx := idx-g.p.A, idx+g.p.B
		if g.p.C > 0 {
			before_idx, after_idx = idx-g.p.C, idx+g.p.C
		}
		if before_idx <= 0 {
			before_idx = 0
		}
		if after_idx >= len(g.rows) {
			after_idx = len(g.rows) - 1
		}
		for i := before_idx; i <= after_idx; i++ {
			copy_idxs = append(copy_idxs, i) // Добавление индексов из диапазона
		}
	}

	// Избавление от дублей
	sort.Ints(copy_idxs)
	clean_idxs := []int{copy_idxs[0]}
	for prev, next := 0, 1; next < len(copy_idxs); prev, next = prev+1, next+1 {
		if copy_idxs[prev] == copy_idxs[next] {
			continue
		}
		clean_idxs = append(clean_idxs, copy_idxs[next])
	}
	return clean_idxs
}

func (g *GREP) invertIdx(idxs []int) []int {
	if !g.p.InVert {
		return idxs
	}
	out := make([]int, 0)
	contains := func(a int) bool {
		for _, i := range idxs {
			if i == a {
				return true
			}
		}
		return false
	}
	for i := 0; i < len(g.rows); i++ {
		if contains(i) {
			continue
		}
		out = append(out, i)
	}
	return out
}

// Собирает строки по индексам в единый срез
func (g *GREP) collectRows(idxs []int) []string {
	collection := []string{}
	for _, idx := range idxs {
		row := g.rows[idx]
		if g.p.Nums {
			row = fmt.Sprintf("%d: %s", idx, row)
		}
		collection = append(collection, row)
	}
	return collection
}

func (g *GREP) stickRows(rows []string) string {
	return strings.Join(rows, "\n")
}

// Всякая классическая фигня, парсинг командной строки, проверка пути, считывание строк файла

func Parse() (params Params, patterns []string, path string) {
	A_flag := flag.Int("A", 0, `"after" печатать +N строк после совпадения`)
	B_flag := flag.Int("B", 0, `"before" печатать +N строк до совпадения`)
	C_flag := flag.Int("C", 0, `"context" (A+B) печатать ±N строк вокруг совпадения`)
	c_flag := flag.Int("c", 0, `"count" (количество строк)`)
	i_flag := flag.Bool("i", false, `"ignore-case" (игнорировать регистр)`)
	v_flag := flag.Bool("v", false, `"invert" (вместо совпадения, исключать)`)
	F_flag := flag.Bool("F", false, `"fixed", точное совпадение со строкой, не паттерн`)
	n_flag := flag.Bool("n", false, `"line num", печатать номер строки`)

	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		fmt.Printf("Передано недостаточно параметров: %v\n", args)
		os.Exit(1)
	}
	patterns, path = args[:len(args)-1], args[len(args)-1]

	params = Params{*A_flag, *B_flag, *C_flag, *c_flag, *i_flag, *v_flag, *F_flag, *n_flag}
	return params, patterns, path
}

// Гарантирует что файл доступен и откроется
func CheckPath(path string) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("путь к файлу %s не сущесвует", err)
	}
	if info.IsDir() {
		return fmt.Errorf("путь является директорией - %s", path)
	}
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("файл недоступен, причина: %w", err)
	}
	defer file.Close()

	return nil
}

// Парсинг строк
func ParseRows(path string) []string {
	file, _ := os.Open(path)
	defer file.Close()

	var rows []string
	f_scan := bufio.NewScanner(file)
	for f_scan.Scan() {
		rows = append(rows, f_scan.Text())
	}

	return rows
}
