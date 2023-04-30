// package main

package pattern

/*
	Реализовать паттерн «стратегия».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Strategy_pattern
*/

import (
	"fmt"
	"time"
)

// Интерфейс стратегии определяет метод, который должны реализовать все конкретные стратегии.
type SortStrategy interface {
	Sort(arr []int) []int
}

// Конкретная стратегия для сортировки пузырьком.
type BubbleSortStrategy struct{}

func (s *BubbleSortStrategy) Sort(arr []int) []int {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
	return arr
}

// Конкретная стратегия для сортировки выбором.
type SelectionSortStrategy struct{}

func (s *SelectionSortStrategy) Sort(arr []int) []int {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		minIndex := i
		for j := i + 1; j < n; j++ {
			if arr[j] < arr[minIndex] {
				minIndex = j
			}
		}
		arr[i], arr[minIndex] = arr[minIndex], arr[i]
	}
	return arr
}

// Контекст, который использует стратегию для сортировки.
type SortContext struct {
	strategy SortStrategy
}

func (c *SortContext) SetStrategy(strategy SortStrategy) {
	c.strategy = strategy
}

func (c *SortContext) Sort(arr []int) []int {
	return c.strategy.Sort(arr)
}

func main() {
	// Создаем контекст с дефолтной стратегией - сортировкой пузырьком.
	context := &SortContext{&BubbleSortStrategy{}}

	// Сортируем массив с помощью текущей стратегии.
	start := time.Now()
	arr := []int{5, 2, 1, 8, 4}
	fmt.Printf("took: %dms	bubble sort: %v\n", time.Since(start).Milliseconds(), context.Sort(arr)) // [1 2 4 5 8]

	// Меняем стратегию на сортировку выбором.
	context.SetStrategy(&SelectionSortStrategy{})
	start = time.Now()
	fmt.Printf("took: %dms	selection sort: %v\n", time.Since(start).Milliseconds(), context.Sort(arr)) // [1 2 4 5 8]
}

/*
Суть:
	Инкапсуляция семейства алгоритмов, обеспечение их взаимозаменяемости
	Изменение алгоритмов независимо от клиентов пользователей

Применимость:
	* наличие множества родственных классов отличающихся только поведением
	* наличие множества алгоритмов для решения одной задачи
	* в классе определено много вариантов поведения представленных развлетвлёнными условными опеаторами
														-> проще перенести в отдельные классы стратегии

*/
