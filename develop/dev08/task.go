package main

/*
=== Взаимодействие с ОС ===

Необходимо реализовать собственный шелл

встроенные команды: cd/pwd/echo/kill/ps
поддержать fork/exec команды
конвеер на пайпах

Реализовать утилиту netcat (nc) клиент
принимать данные из stdin и отправлять в соединение (tcp/udp)
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("⊳⊳⊳ ")
		input, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		Iterpreter(input)

	}
}

func Iterpreter(inp string) {
	cmds := strings.Split(inp, "|")
	first_cmd := strings.Fields(cmds[0])

	// Если не используется пайплайн, выполнить и вывести
	if len(cmds) <= 1 {
		out, _ := Execute(first_cmd)
		fmt.Println(out)
		return
	}

	// Выполнит первую команду
	out, ok := Execute(first_cmd)
	if !ok {
		return
	}

	for _, cmd := range cmds {
		parts := strings.Fields(cmd)
		if len(out) != 0 { // Если вывод первой команды содержит выходные данные
			parts = append(parts, out) // Добавить выходные данные как аргумент к след команде (хотя следовало бы направить в стандартный вход но я плохо задизайнил)
		}
		out, ok = Execute(parts)
		if !ok {
			return
		}
	}
}

func Execute(parts []string) (string, bool) {

	switch parts[0] {
	case "cd":
		if len(parts) < 2 {
			fmt.Fprintln(os.Stderr, "path required")
			return "", false
		}
		err := os.Chdir(parts[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		return "", true
	case "pwd":
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return "", false
		}
		return cwd, true
	case "echo":
		return strings.Join(parts[1:], " "), true
	case "kill":
		if len(parts) < 2 {
			fmt.Fprintln(os.Stderr, "pid required")
			return "", false
		}
		pid, _ := strconv.ParseUint(parts[1], 10, 32)

		process, err := os.FindProcess(int(pid))
		if err != nil {
			return "", false
		}

		err = process.Signal(syscall.SIGKILL)
		if err != nil {
			return "", false
		}
		fmt.Println("Процесс убит")
		return "", true
	case "ps":
		cmd := exec.Command("ps", "aux")
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return "", false
		}
		return string(out), true
	case "e":
		fmt.Println("Оболочка завершила работу")
		os.Exit(0)
		return "", false
	}
	return "wtf man", false
}
