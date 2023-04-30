package main

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

type MockNtp struct {
	id int
}

// Типо непредвиденная ошибка
type TestErr struct{}

func (e TestErr) Error() string {
	return "Прострел бита лол"
}

// Подстановка заполненных структур ответов вместо веб запросов
func (m MockNtp) GetNtpTime() NtpResp {
	answers := []NtpResp{
		// Корректный ответ
		NtpResp{
			KissCode:    "7",
			Stratum:     7,
			ClockOffset: time.Duration(1 * time.Millisecond),
			err:         nil,
		},
		// Непредвиденная ошибка
		NtpResp{
			KissCode:    "7",
			Stratum:     7,
			ClockOffset: time.Duration(1 * time.Millisecond),
			err:         TestErr{},
		},
		// Ошибка перегруза запросами
		NtpResp{
			KissCode:    "7",
			Stratum:     0,
			ClockOffset: time.Duration(1 * time.Millisecond),
			err:         nil,
		},
	}

	return answers[m.id]
}

func TestCurrTime(t *testing.T) {
	var tests = []struct {
		name string
		want error
	}{
		{"Normal case", nil},
		{"Unknown error case", UnknownError{TestErr{}}},
		{"Overload error case", OverloadError{"7"}},
	}

	mock := MockNtp{-1}
	for _, test := range tests {
		mock.id++
		name := fmt.Sprintf("case: %s", test.name)
		t.Run(name, func(t *testing.T) {
			_, err := CurrTime(mock)
			if !errors.Is(err, test.want) {
				t.Errorf("want %T, got %T, .Error() -> %s", test.want, err, err.Error())
			}
		})
	}
}
