package main

// package pattern

/*
	Реализовать паттерн «состояние».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/State_pattern
*/

import "fmt"

// Контекст определяет интерфейс, представляющий объект,
// чье поведение изменяется в зависимости от состояния.
type Context struct {
	state State
}

// State определяет интерфейс состояний.
type State interface {
	Handle(ctx *Context)
	ExecState()
}

// ConcreteStateA и ConcreteStateB - это конкретные состояния,
// реализующие интерфейс State.
type ConcreteStateA struct{}

// Выполняет что то полезное и переводит в следующее состояние
func (s *ConcreteStateA) Handle(ctx *Context) {
	fmt.Println("Handle ConcreteStateA")
	s.ExecState()
	ctx.state = &ConcreteStateB{}
}

func (s *ConcreteStateA) ExecState() {
	fmt.Println("Executed state A, moving to state B")
}

type ConcreteStateB struct{}

func (s *ConcreteStateB) Handle(ctx *Context) {
	fmt.Println("Handle ConcreteStateB")
	s.ExecState()
	ctx.state = &ConcreteStateA{}
}

func (s *ConcreteStateB) ExecState() {
	fmt.Println("Executed state B, moving to state A")
}

func main() {
	// Создаем объект контекста с начальным состоянием ConcreteStateA.
	ctx := &Context{state: &ConcreteStateA{}}

	// Последовательно вызываем Handle на контексте, что приводит к переключению состояний.
	ctx.state.Handle(ctx)
	ctx.state.Handle(ctx)
}

/*
Суть:
	Изменение поведение объекта в зависимости от внутреннего состояния

Применимость:
	* Поведение объекта должно изменяться во время выполнения, поведение зависит от состояния

*/
