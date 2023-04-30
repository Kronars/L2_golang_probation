package main

import (
	"fmt"
	"os"
	"time"

	ntp "github.com/beevik/ntp"
)

/*
=== Базовая задача ===

Создать программу печатающую точное время с использованием NTP библиотеки.Инициализировать как go module.
Использовать библиотеку https://github.com/beevik/ntp.
Написать программу печатающую текущее время / точное время с использованием этой библиотеки.

Программа должна быть оформлена с использованием как go module.
Программа должна корректно обрабатывать ошибки библиотеки: распечатывать их в STDERR и возвращать ненулевой код выхода в OS.
Программа должна проходить проверки go vet и golint.
*/

// Всякие возможные ошибки
type UnknownError struct {
	err error
}

func (e UnknownError) Error() string {
	return fmt.Sprintf("Unpredicted error: %s", e.err.Error())
}

type OverloadError struct {
	code string
}

func (e OverloadError) Error() string {
	return fmt.Sprintf("kiss of death: too many requests to the server in a short period of time. KissCode: %s", e.code)
}

func main() {
	// Создание пустой структуры ответа
	var reply NtpResp
	// Выполнение запроса и наполнение структуры
	time, err := CurrTime(reply)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}
	fmt.Println(time)
}

// Испльзуемые параметры ответа
type NtpResp struct {
	KissCode    string
	Stratum     uint8
	ClockOffset time.Duration
	err         error
}

// Выполнение запроса
func (r NtpResp) GetNtpTime() NtpResp {
	resp, err := ntp.Query("ntp7.ntp-servers.net")

	r.KissCode = resp.KissCode
	r.Stratum = resp.Stratum
	r.ClockOffset = resp.ClockOffset
	r.err = err

	return r
}

// Интерфейс для возможности мокать ответ
type NtpResponser interface {
	GetNtpTime() NtpResp
}

// Обработка запроса
func CurrTime(r NtpResponser) (time.Time, error) {
	response := r.GetNtpTime()
	if response.err != nil {
		return time.Now(), UnknownError{response.err}
	}

	if response.Stratum == 0 {
		return time.Now(), OverloadError{response.KissCode}
	}

	time := time.Now().Add(response.ClockOffset)
	return time, nil
}
