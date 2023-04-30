package pattern

/*
	Реализовать паттерн «комманда».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Command_pattern
*/

import "fmt"

// Command interface
type Command interface {
	Execute()
}

// Receiver struct
type Receiver struct {
	message string
}

// SetMessage method of Receiver
func (r *Receiver) SetMessage(msg string) {
	r.message = msg
}

// PrintMessage method of Receiver
func (r *Receiver) PrintMessage() {
	fmt.Println(r.message)
}

// ConcreteCommand struct
type ConcreteCommand struct {
	receiver *Receiver
}

// Execute method of ConcreteCommand
func (cc *ConcreteCommand) Execute() {
	cc.receiver.PrintMessage()
}

// Invoker struct
type Invoker struct {
	commands []Command
}

// AddCommand method of Invoker
func (i *Invoker) AddCommand(c Command) {
	i.commands = append(i.commands, c)
}

// ExecuteCommands method of Invoker
func (i *Invoker) ExecuteCommands() {
	for _, c := range i.commands {
		c.Execute()
	}
}

func main() {
	// Create a new Receiver
	r := &Receiver{}

	// Create a new ConcreteCommand and set its Receiver
	cc1 := &ConcreteCommand{r}

	// Set the message of the Receiver
	r.SetMessage("Hello world!")

	// Create a new Invoker and add the ConcreteCommand to its commands list
	i := &Invoker{}
	i.AddCommand(cc1)

	// Execute the ConcreteCommand using the Invoker
	i.ExecuteCommands()
}

/*
Суть и применимость:
Eдиный интерфейс для описания всех типов операций над системой
Упрощение добавления в систему поддержки новой операции -> достаточно реализовать предлагаемый интерфейс.
Kаждая операция представляется самостоятельным объектом инкапсулирующим некоторый набор дополнительных свойств
Система приобретает возможно выполнять дополнительный набор действий над запросами (объектами)
Протоколирование, отмена предыдущего действия повторение последующего и т.д.

Примеры использования
 1. desktop-приложение c возможностями отмены и повторения действий пользователя (undo/redo)
 2. Сетевые распределенные системы использующие запросы в виде объектов в качестве основного примитива инициализации каких-либо операций
 3. Системы с поддержкой асинхронных вызовов, инкапсулирующие обратный вызов в виде callback-объекта


*/
