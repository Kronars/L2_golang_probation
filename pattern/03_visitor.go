package pattern

/*
	Реализовать паттерн «посетитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Visitor_pattern
*/
import "fmt"

type Element interface {
	Accept(visitor Visitor)
}

type ConcreteElementA struct {
	Name string
}

func (e *ConcreteElementA) Accept(visitor Visitor) {
	visitor.VisitElementA(e)
}

type ConcreteElementB struct {
	Id int
}

func (e *ConcreteElementB) Accept(visitor Visitor) {
	visitor.VisitElementB(e)
}

type Visitor interface {
	VisitElementA(element *ConcreteElementA)
	VisitElementB(element *ConcreteElementB)
}

type ConcreteVisitor struct{}

func (v *ConcreteVisitor) VisitElementA(element *ConcreteElementA) {
	fmt.Printf("Visited ConcreteElementA with name %s\n", element.Name)
}

func (v *ConcreteVisitor) VisitElementB(element *ConcreteElementB) {
	fmt.Printf("Visited ConcreteElementB with id %d\n", element.Id)
}

func main() {
	elements := []Element{&ConcreteElementA{Name: "foo"}, &ConcreteElementB{Id: 42}}

	visitor := &ConcreteVisitor{}
	for _, e := range elements {
		e.Accept(visitor)
	}
}

/*
Паттерн проектирования "посетитель" позволяет добавлять новые операции для объектов, не добавляя новых методов структурам.
Используется в случаях, когда необходимо обойти множество объектов разных типов и выполнить для каждого некоторую операцию,
при этом изменять эти объекты нежелательно или невозможно.

Применимость:
 * Множество структур
 	* Необходимы операции зависящие от типа структуры
 * Нежелательно изменять исходные структуры

Используется когда необходимо обойти множество объектов разных структур и выполнить для каждого из них некоторую операцию, зависящую от структуры
Основная идея паттерна заключается в разделении поведения-операций над объектами, и структуры данных на две отдельные структуры.
Объекты структуры данных предоставляют интерфейс "принять посетителя", он позволяет посетителю выполнить свою операцию для данного объекта.
*/
