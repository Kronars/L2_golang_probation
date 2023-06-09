Что выведет программа? Объяснить вывод программы. Объяснить как работают defer’ы и их порядок вызовов.

```go
package main

import (
	"fmt"
)


func test() (x int) {
	defer func() {
		x++
	}()
	x = 1
	return
}


func anotherTest() int {
	var x int
	defer func() {
		x++
	}()
	x = 1
	return x
}


func main() {
	fmt.Println(test())
	fmt.Println(anotherTest())
}
```

Ответ:
```go
Порядок выполнения:
1 Исполнение тела функции
2 Исполнение выражений в return
3 Исполнение defer
4 Возврат результата

Если в функции для возврата используются именованные аргументы, то в return вычисленная переменная 
будет передана будто по ссылке - изменения переменной в defer отразится на возвращаемом значении

Если в функции для возврата не используются именованные аргументы, то в return вычисленная переменная 
будет передана как значение - изменения в defer не отразятся на возвращаемом

Примеры:

func test() (y, x int) {
	defer func() {
		x++
		fmt.Println("defer - вычисляется вторым, после return, x == ", x)
		fmt.Println("происходит возврат значений из функции")
	}()
	x = 1
	return fmt.Println("return - вычисляется первым, перед defer, x == ", x), x
}


>>> test:
>>> return - вычисляется первым, перед defer, x ==  1
>>> defer - вычисляется вторым, после return, x ==  2
>>> происходит возврат значений из функции
>>> 2

func anotherTest() (int, int) {
	var x int
	defer func() { 
		x++
		fmt.Println("defer - вычисляется вторым, после return, x == ", x)
		fmt.Println("В return уже вычислено значение результата, возвращено будет оно - 1")
		}()
	x = 1
	return fmt.Println("return - вычисляется первым, перед defer, подготовленно значение для возврата - текущий х == ", x), x
}


>>> anotherTest:
>>> return - вычисляется первым, перед defer, подготовленно значение для возврата - текущий х ==  1
>>> defer - вычисляется вторым, после return, x ==  2
>>> В return уже вычислено значение результата, возвращено будет оно - 1
>>> 1
```
