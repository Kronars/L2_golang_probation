package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port
go-telnet mysite.ru 8080
go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

func main() {
	// args := []string{"--pha=3s", "localhost", "8888"}
	args := os.Args[1:]
	host, port, timeout := Parse(args)

	conn, err := ConnectTCP(host, port, timeout)
	if err != nil {
		fmt.Println("Ошикба подключения: ", err)
	}
	defer conn.Close()
	fmt.Println("Успешное подключение, введите EXIT что бы закрыть соединение")

	// Запускаем горутину чтения данных из сокета
	go Listen(conn)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("⊳⊳⊳ ")
		input, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		if HasExit(input) {
			fmt.Println("\nСоединение закрыто")
			break
		}

		fmt.Fprintln(conn, input)
	}

}

func ConnectTCP(address string, port int, timeout time.Duration) (net.Conn, error) {
	addr := fmt.Sprintf("%s:%d", address, port)

	// Подключаемся к серверу с таймаутом
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return nil, err
	}

	// Возвращаем соединение
	return conn, nil
}

func HasExit(inp string) bool {
	if strings.Contains(inp, "EXIT") {
		return true
	}
	return false
}

func Listen(c net.Conn) {
	scanner := bufio.NewScanner(c)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if scanner.Err() != nil {
		fmt.Fprintf(os.Stderr, "Ошибка чтения из сокета: %v\n", scanner.Err())
		os.Exit(1)
	}

	// Если соединение закрыто, завершаем программу
	os.Exit(0)
}

func Parse(args []string) (string, int, time.Duration) {
	fmt.Println(args)
	if len(args) < 2 {
		fmt.Println("Недостаточно параметров")
		os.Exit(1)
	}
	host := args[len(args)-2]
	port, err_port := strconv.ParseUint(args[len(args)-1], 10, 32)
	if err_port != nil {
		fmt.Println("Неправильно передан порт", port)
		os.Exit(1)
	}

	timeout := time.Duration(1 * time.Second)
	if len(args) == 3 {
		str_timeout := strings.Split(args[len(args)-3], "=")[1]            // Вырезал значение таймаута. 3s
		rune_timeout := []rune(str_timeout)[:len(str_timeout)-1]           // Удалил последний символ. 3
		int_timeout, err := strconv.ParseInt(string(rune_timeout), 10, 32) // Распарсил в цифру
		if err != nil || int_timeout == 0 {
			fmt.Println("Неправильно передан параметр таймаута ", str_timeout)
			os.Exit(1)
		}
		timeout = time.Duration(int_timeout * int64(timeout.Seconds()))
	}

	return host, int(port), timeout
}
