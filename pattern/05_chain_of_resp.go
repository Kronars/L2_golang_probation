package pattern

/*
	Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern
*/

import "fmt"

// Request - запрос
type Request struct {
	value int
}

// Handler - интерфейс обработчика
type Handler interface {
	SetNext(handler Handler) Handler
	Handle(request *Request)
}

// BaseHandler - базовый обработчик
type BaseHandler struct {
	next Handler
}

// SetNext - установка следующего обработчика
func (h *BaseHandler) SetNext(handler Handler) Handler {
	h.next = handler
	return handler
}

// Handle - обработка запроса
func (h *BaseHandler) Handle(request *Request) {
	if h.next != nil {
		h.next.Handle(request)
	}
}

// HandlerA - обработчик A
type HandlerA struct {
	BaseHandler
}

// Handle - обработка запроса
func (h *HandlerA) Handle(request *Request) {
	if request.value < 0 {
		fmt.Println("HandlerA: запрос обработан")
	} else {
		h.BaseHandler.Handle(request)
	}
}

// HandlerB - обработчик B
type HandlerB struct {
	BaseHandler
}

// Handle - обработка запроса
func (h *HandlerB) Handle(request *Request) {
	if request.value < 10 {
		fmt.Println("HandlerB: запрос обработан")
	} else {
		h.BaseHandler.Handle(request)
	}
}

// HandlerC - обработчик C
type HandlerC struct {
	BaseHandler
}

// Handle - обработка запроса
func (h *HandlerC) Handle(request *Request) {
	if request.value < 20 {
		fmt.Println("HandlerC: запрос обработан")
	} else {
		h.BaseHandler.Handle(request)
	}
}

func main() {
	// Создаем обработчики
	handlerA := &HandlerA{}
	handlerB := &HandlerB{}
	handlerC := &HandlerC{}

	// Связываем их в цепочку
	handlerA.SetNext(handlerB).SetNext(handlerC)

	// Отправляем запросы
	request1 := &Request{value: -5}
	request2 := &Request{value: 15}
	request3 := &Request{value: 25}

	handlerA.Handle(request1)
	handlerA.Handle(request2)
	handlerA.Handle(request3)
}

/*
Назначение:
	Позволяет избежать привязки отправителя запроса к его получателю,
	предоставляя возможность обработать запрос нескольким объектам
	Связывает объекты получателив цепочку, передаёт запрос по цепочке пока не будет обработан.

Достоинства:
	* Ослабление связности. Объект может не парится кто обработает его запрос
	упрощается взаимосвязь -> не требуется хранить ссылки на все объекты-получатели запроса

	* Гибкость распределения обязанностей между объектами
	Добавить или изменить обязанности по обработке можно включив в цепочку новых участников или изменив старых

Подводные:
	* Получение не гарантированно
	У запроса нет явного получателя -> нет гарантии что запрос будет обработан, достигнет конца цепочки и пропадёт
	Может быть необработан в случае неправильной конфигурации цепочки
*/
