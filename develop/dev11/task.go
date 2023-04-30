package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

/*
=== HTTP server ===

Реализовать HTTP сервер для работы с календарем. В рамках задания необходимо работать строго со стандартной HTTP библиотекой.
В рамках задания необходимо:
	1. ✅ Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.
	2. ✅ Реализовать вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event.
	3. ✅ Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные функции и объекты доменной области.
	4. ❌ Реализовать middleware для логирования запросов
Методы API: POST /create_event POST /update_event POST /delete_event GET /events_for_day GET /events_for_week GET /events_for_month
Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09).
В GET методах параметры передаются через queryString, в POST через тело запроса.
В результате каждого запроса должен возвращаться JSON документ содержащий либо {"result": "..."} в случае успешного выполнения метода,
либо {"error": "..."} в случае ошибки бизнес-логики.

В рамках задачи необходимо:
	1. ✅ Реализовать все методы.
	2. ✅ Бизнес логика НЕ должна зависеть от кода HTTP сервера.
	3. ✅ В случае ошибки бизнес-логики сервер должен возвращать HTTP 503. В случае ошибки входных данных (невалидный int например) сервер должен возвращать HTTP 400. В случае остальных ошибок сервер должен возвращать HTTP 500. Web-сервер должен запускаться на порту указанном в конфиге и выводить в лог каждый обработанный запрос.
	4. ✅ Код должен проходить проверки go vet и golint.
*/

// ---------- Структуры бизнес логики ----------

// Штамп даты для удобства
type Stamp struct {
	day, month, year int
}

// Описание события представляет собой одну строку, в теле POST передаётся по ключу event
type EventDesc string

// Индекс события, строка в формате год-месяц-дата
type EventDate string

// Конструктор EventDate, форматирует строку к единому формату
func NewEventDate(d Stamp) EventDate {
	return EventDate(FormatStamp(d)) // time.DateOnly
}

// Форматирование EventDate из штампа
func FormatStamp(d Stamp) string {
	return time.Date(d.year, time.Month(d.month), d.day, 1, 1, 1, 1, time.Local).Format("2006-01-02")
}

// Форматирование EventDate из объекта time.Time
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02")
}

type Calendar struct {
	events map[EventDate][]EventDesc // Ключ в формате EventDate.Format
}

func NewCalendar() Calendar {
	e := make(map[EventDate][]EventDesc, 10)
	return Calendar{e}
}

type User struct {
	id int
	Calendar
}

func NewUser(id int) User {
	return User{id, NewCalendar()}
}

// Структура для хранения данных
type Data struct {
	users map[int]User
}

// Получение/создание пользователя
func (d *Data) GetOrCreateUser(user_id int) *User {
	user, ok := data.users[user_id]
	if ok {
		return &user
	}
	user = NewUser(user_id)
	data.users[user_id] = user
	return &user
}

// Для удобства доступа данные хранятся в глобальном неймспейсе
var data Data = Data{make(map[int]User, 16)} // TODO хорошо бы блочить запись мутексами

func main() {
	http.HandleFunc("/", RootHandler)
	http.HandleFunc("/create_event", CreateHandler)
	http.HandleFunc("/update_event", UpdHandler)
	http.HandleFunc("/delete_event", DelHandler)
	http.HandleFunc("/events_for_day", DayEventsHandler)
	http.HandleFunc("/events_for_week", WeekEventsHandler)
	http.HandleFunc("/events_for_month", MonthEventsHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// ---------- Бизнес логика ----------

// Ошибка обращения к несуществующим событиям в Update и Delete
type ErrNotFound struct {
	date string
	id   int
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("not found event with date=%s event_id=%d", e.date, e.id)
}

// Добавление события
func (c *Calendar) EventCreate(dt EventDate, eb string) {
	event_body := EventDesc(eb)
	if _, exists := c.events[dt]; exists { // Если в этот день уже есть события - добавить новое
		c.events[dt] = append(c.events[dt], event_body)
		return
	}
	c.events[dt] = []EventDesc{event_body}
}

// Обновление события
func (c *Calendar) EventUpd(ev_date EventDate, desc string, id int) error {
	ev_desc := EventDesc(desc)

	// Если сущестуют события на переданную дату
	if evs, ok := c.events[ev_date]; ok {
		// Если есть событие по указанному id
		if id >= len(evs) || id < 0 {
			return ErrNotFound{string(ev_date), id}
		}
		c.events[ev_date][id] = ev_desc
		return nil
	}
	return ErrNotFound{string(ev_date), id}
}

// Удаление события
func (c *Calendar) EventDel(ev_date EventDate, id int) error {
	// Проверка сущестуют ли события на переданную дату
	if evs, ok := c.events[ev_date]; ok {
		// Если есть событие по указанному id
		if id >= len(evs) || id < 0 {
			return ErrNotFound{string(ev_date), id}
		}
		c.events[ev_date] = append(evs[:id], evs[id+1:]...) // Вырезаение события по айди
		return nil
	}
	return ErrNotFound{string(ev_date), id}
}

// Поиск и получение событий за день
func (c *Calendar) EventsDay(date EventDate) []EventDesc {
	if calendar_events, ok := c.events[date]; ok {
		return calendar_events
	}
	return []EventDesc{}
}

// Bсе события за рабочую неделю ВС-СБ (привторится что я американец)
func (c *Calendar) EventsWeek(dw EventDate) ([]EventDesc, []EventDate) {
	date_day, _ := time.Parse("2006-01-02", string(dw)) // Перевод в time.Time, ! dw уже отвалидирован

	// Рассчёт даты начала недели
	day_dur := 24 * time.Hour
	week_start := date_day.Add(-time.Duration(day_dur * time.Duration(date_day.Weekday()))) //  дата начала недели == день недели - (номер дня недели * сутки)

	// Создание среза из ключей-дат недели
	week_dates := make([]EventDate, 7)
	for i := 0; i < 7; i++ {
		day := week_start.Add(time.Duration(time.Duration(i) * day_dur))
		week_dates[i] = EventDate(FormatTime(day))
	}

	// Поиск событий по ключам из среза дат
	var events []EventDesc
	var events_dates []EventDate
	for _, day := range week_dates {
		if ev, ok := c.events[day]; ok {
			events = append(events, ev...)
			events_dates = append(events_dates, day)
		}
		continue
	}

	return events, events_dates
}

// Поиск и получение событий за месяц
func (c *Calendar) EventsMonth(dm EventDate) ([]EventDesc, []EventDate) {
	d, _ := time.Parse("2006-01-02", string(dm)) // Перевод в time.Time, ! dm уже отвалидирован
	day_dur := 24 * time.Hour
	// Получение первого дня месяца
	first_day := time.Date(d.Year(), d.Month(), 1, 1, 0, 0, 0, time.Local)

	// Прибавление 1 месяца и вычитание 1 дня
	last_day := first_day.AddDate(0, 1, -1)

	// Вычисление количества дней в месяце
	daysInMonth := last_day.Day()

	// Создание среза из ключей-дат недели
	month_days := make([]EventDate, daysInMonth)
	for i := 0; i < daysInMonth; i++ {
		day := first_day.Add(time.Duration(time.Duration(i) * day_dur))
		month_days[i] = EventDate(FormatTime(day))
	}

	// Поиск событий по ключам из среза дат
	var events []EventDesc
	var events_dates []EventDate
	for _, day := range month_days {
		if ev, ok := c.events[day]; ok {
			events = append(events, ev...)
			events_dates = append(events_dates, day)
		}
		continue
	}

	return events, events_dates
}
