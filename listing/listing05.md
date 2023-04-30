Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

Ответ:
```
Происходит ситуация схожая с третьим листингом
Функция возвращает ссылку на структуру реализующую интерфейс error
При возврате из функции переменная является структурой
Но в main, результат test преобразуется в переменную err интерфейсного типа error
Интерфейсная переменная получает тип customError и не получает никаких данных описывающих струтуру
Операция сравнения интерфейсной переменной err с nil, пытается сравнить данные переменной которых нет
поэтому сравнение возвращает true 

```
