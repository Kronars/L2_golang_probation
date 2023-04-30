// package pattern
package main

/*
	Реализовать паттерн «строитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Builder_pattern
*/

import (
	"fmt"
	"strings"
)

// Интерфейс AbstractBuilder объявляет методы для поэтапного построения объекта.
type AbstractBuilder interface {
	SetWidth(width int)
	SetHeight(height int)
	SetRenderEngine(engine Render)
}

// Компонент строителя
type Render interface {
	Show()
}

// Default - представление, шаблон конструируемого объекта
// Конструирует из собранных компонентов конкретный объект
type Default struct {
	builder AbstractBuilder
}

func (d *Default) SetBuilder(builder AbstractBuilder) {
	d.builder = builder
}

func (d *Default) Construct() {
	d.builder.SetWidth(40)
	d.builder.SetHeight(10)
	d.builder.SetRenderEngine(&RenderDefault{20, 10, "◯"})
}

// Строитель, собирает из переданных объектов и значений объект
// Предоставляет метод для получения готового объекта
type ConcreteBuilder struct {
	product Product
}

func (b *ConcreteBuilder) SetWidth(width int) {
	b.product.Width = width
}

func (b *ConcreteBuilder) SetHeight(height int) {
	b.product.Height = height
}

func (b *ConcreteBuilder) SetRenderEngine(engine Render) {
	b.product.Render = engine
}

func (b *ConcreteBuilder) GetObj() Product {
	return b.product
}

// Движок рендера
type RenderDefault struct {
	w, h int
	char string
}

func (r *RenderDefault) Show() {
	var out string
	for i := 0; i < r.h; i++ {
		out += strings.Repeat(r.char, r.w)
		out += "\n"
	}
	fmt.Println(out)
}

// Product - собираемый объект.
type Product struct {
	Width  int
	Height int
	Render
}

func main() {
	builder := &ConcreteBuilder{} // Создание строителя
	artist := &Default{}          // Oбъект cобираемый по шаблону Default
	artist.SetBuilder(builder)    // Могут использоваться разные строители
	artist.Construct()
	creation := builder.GetObj()
	creation.Show()
}

/*
	Порождающий паттерн, отделяет создание сложного объекта от его представления
	Разделяет процесс создания сложных объектов на отдельные шаги, что позволяет контролировать процесс
	и создавать различные варианты объектов на основе одного и того же кода строительства

Применимость
 * процесс конструирования должен обеспечивать различные представления конструируемого объекта
 * алгоритм создания объекта не должен зависеть от компонентов из которых он собирается, стыкуются ли компоненты
*/
