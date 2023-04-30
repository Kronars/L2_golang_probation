package main

import (
	"reflect"
	"testing"
)

// Календарь с одним событием и пустой
var date1 = NewEventDate(Stamp{3, 9, 2007})
var event1 = []EventDesc{"я", "календарь"}
var c1 = Calendar{map[EventDate][]EventDesc{date1: event1}}
var c2 = Calendar{make(map[EventDate][]EventDesc, 2)}

// События на разных неделях и месяцах
var event2 = []EventDesc{"Второе событие недели"}
var event3 = []EventDesc{"Событие другой недели"}
var event4 = []EventDesc{"Событие другого месяца"}
var date2 = NewEventDate(Stamp{4, 9, 2007})
var date3 = NewEventDate(Stamp{29, 9, 2007})
var date4 = NewEventDate(Stamp{29, 10, 2007})
var c3 = Calendar{map[EventDate][]EventDesc{date1: event1, date2: event2, date3: event3, date4: event4}}

func TestCalendar_EventsDay(t *testing.T) {
	tests := []struct {
		name        string
		c           *Calendar
		date_search EventDate
		want        []EventDesc
	}{
		{"base", &c1, date1, event1},
		{"unknown", &c1, NewEventDate(Stamp{8, 3, 2007}), []EventDesc{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.EventsDay(tt.date_search); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Calendar.EventsDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalendar_EventCreate(t *testing.T) {
	tests := []struct {
		name      string
		c         *Calendar
		dt_create EventDate
		eb        string
		want      []EventDesc
	}{
		{"base", &c2, date1, "Событие 1", []EventDesc{"Событие 1"}},
		{"append", &c2, date1, "Событие 2", []EventDesc{"Событие 1", "Событие 2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.EventCreate(tt.dt_create, tt.eb)

			if got := tt.c.events[tt.dt_create]; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Calendar.EventCreate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalendar_EventsWeek(t *testing.T) {
	tests := []struct {
		name string
		c    *Calendar
		date EventDate
		want []EventDesc
	}{
		{"base", &c3, NewEventDate(Stamp{3, 9, 2007}), []EventDesc{"я", "календарь", "Второе событие недели"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := tt.c.EventsWeek(tt.date); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Calendar.EventsWeek() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalendar_EventsMonth(t *testing.T) {
	tests := []struct {
		name string
		c    *Calendar
		dm   EventDate
		want []EventDesc
	}{
		{"base", &c3, NewEventDate(Stamp{3, 9, 2007}), []EventDesc{"я", "календарь", "Второе событие недели", "Событие другой недели"}}, // не содержит события другого месяца!
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := tt.c.EventsMonth(tt.dm); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Calendar.EventsMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}
