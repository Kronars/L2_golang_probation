package pattern

import "fmt"

/*
	Реализовать паттерн «фабричный метод».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Factory_method_pattern
*/

type Shape interface {
	draw() string
}

// интерфейс создания объектов
type ShapeFactory interface {
	create() Shape
}

type Circle struct {
	radius float64
}

func (c *Circle) draw() string {
	return fmt.Sprintf("Drawing a circle with radius %f", c.radius)
}

type CircleFactory struct{}

func (cf *CircleFactory) create() Shape {
	return &Circle{radius: 1.0}
}

type Rectangle struct {
	width, height float64
}

func (r *Rectangle) draw() string {
	return fmt.Sprintf("Drawing a rectangle with width %f and height %f", r.width, r.height)
}

type RectangleFactory struct{}

func (rf *RectangleFactory) create() Shape {
	return &Rectangle{width: 1.0, height: 2.0}
}

func main() {
	circleFactory := &CircleFactory{}
	circle := circleFactory.create()
	fmt.Println(circle.draw())

	rectangleFactory := &RectangleFactory{}
	rectangle := rectangleFactory.create()
	fmt.Println(rectangle.draw())
}

/*
Суть:
	Порождающий паттерн проектирования, предоставляет интерфейс создания объектов,
	позволяет подклассам выбрать класс для создания.
	Предоставляет инструменты для создания объектов, не указывая конкретного класса объекта.

Применимость:
	* Когда есть общий интерфейс создания объектов, но требуемые для создания классы определяются в рантайме
	Иными словами
	* Когда нужно создавать объекты в зависимости от других объектов или условий

Достоинства:
	hooks - появляется точка входа для контроля/переопределения/дополнения создания всех объектов определённого класса, +гибкость
	соединение параллельных иерархий - когда класс делегирует часть своих обязанностей другому классу, который не является его производным == паралелен
	в создателя можно добавить операцию-создание паралельного класса
*/
