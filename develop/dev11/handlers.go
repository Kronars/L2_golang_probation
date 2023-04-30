package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// ------- Всякое веб серверное --------

const (
	HOST = "localhost"
	PORT = 8080
)

// ------- Схемы, сериализаторы, POST --------
// Ответ успеха с описанием (при успешных create upd del)
type ResultScheme struct {
	Result string `json:"result"`
}

func SerJsonResult(data string) []byte {
	rj, _ := json.Marshal(ResultScheme{data})
	return rj
}

// Ответ провала сервера - бизнес логики (удаление/обновление несуществующего)
type LogicErrScheme struct {
	Error string `json:"error"`
}

// {error: "not found event with date=%s id=%d"}
func SerLogicErr(user_id int, err_str string) string {
	err := fmt.Sprintf("Server error. user_id=%d Error: %s", user_id, err_str)
	rj, _ := json.Marshal(LogicErrScheme{err})
	return string(rj)
}

// Ответ провала клиента - валидации (отрицательный id, неформаттная дата)
type ValidErrScheme struct { // Дублируется с LogicErr, поздно заметил, рефакторить лень
	Error string `json:"error"`
}

// {error: "Validation error"}
func SerValidErr(field string) string {
	rj, _ := json.Marshal(ValidErrScheme{"Validation error: " + field})
	return string(rj)
}

// для десериализации
type PostBodyScheme struct {
	Data string `json:"data"`
}

// ------- Схемы, сериализаторы, GET --------

type EventScheme struct {
	Date  EventDate `json:"date"`
	Event EventDesc `json:"event"`
}

type EventListScheme struct {
	Result []EventScheme `json:"result"`
}

// Для сериализации событий одного дня
func SerOneDayEvents(events []EventDesc, date EventDate) []byte {
	e_schemes := make([]EventScheme, len(events))
	for id, event := range events {
		e_schemes[id] = EventScheme{date, event}
	}
	e_json, _ := json.Marshal(e_schemes)
	return e_json
}

// На вход приходят упорядоченные срезы событий и дат, позиции соответствуют
func SerManyDayEvents(events []EventDesc, date []EventDate) []byte {
	e_schemes := make([]EventScheme, len(events))
	for id, event := range events {
		e_schemes[id] = EventScheme{date[id], event}
	}
	e_json, _ := json.Marshal(e_schemes)
	return e_json
}

// ---------- Обработчики -----------

// Главная страница для красоты (но как надо работать она конечно же не будет)
func RootHandler(w http.ResponseWriter, r *http.Request) {
	desc := `
<!DOCTYPE html>
<html>
<head><title>dev11</title><style>body {background-color: rgb(18, 18, 18);color: white;text-align: center;font-family: monospace;}a {color: white;}</style></head>
<body>
	<h1>Доступные методы:</h1>
	<p><a href="%s:%d/create_event">/create_event</a> - POST</p>
	<p><a href="%s:%d/update_event">/update_event</a> - POST</p>
	<p><a href="%s:%d/delete_event">/delete_event</a> - POST</p>
	<p><a href="%s:%d/events_for_day">/events_for_day</a> - GET</p>
	<p><a href="%s:%d/events_for_week">/events_for_week</a> - GET</p>
	<p><a href="%s:%d/events_for_month">/events_for_month</a> - GET</p>
</body>
</html>
`
	desc = fmt.Sprintf(desc, HOST, PORT, HOST, PORT, HOST, PORT, HOST, PORT, HOST, PORT, HOST, PORT)
	fmt.Fprint(w, desc)
}

// ---------- POST Обработчики ----------

// * POST /create_event + user_id=3&date=2019-09-09 + body: { "data": "вставить текст" }
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка метода
	if r.Method != "POST" {
		fmt.Fprint(w, "Поддерживается только POST")
		return
	}
	defer r.Body.Close()

	// Преобразование и валидация параметров url
	date, user_id, valid := ParseValid_DateUser(r.URL.Query())
	if !valid {
		http.Error(w, SerValidErr("incorrect url params"), 400)
		return
	}

	// Чтение body
	body, err := ParseBody(r.Body)
	if err != nil {
		http.Error(w, SerValidErr(err.Error()), 400)
		return
	}

	// Сохранение данных
	user := data.GetOrCreateUser(user_id)
	user.EventCreate(date, string(body))

	// Ответ о успешной записи
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(SerJsonResult("Created"))
}

// * /update_event + ?user_id=3&date=2019-09-09&event_id=0 + body: { "data": "вставить текст" }
// * Ответ: { "result": "updated" }
func UpdHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка метода
	if r.Method != "POST" {
		fmt.Fprint(w, "Поддерживается только POST")
		return
	}
	defer r.Body.Close()

	// Преобразование и валидация параметров url
	user_id, date, event_id, valid := ParseValid_UserDateEvent_id(r.URL.Query())
	if !valid {
		http.Error(w, SerValidErr("incorrect url params"), 400)
		return
	}

	// Чтение body
	body, err := ParseBody(r.Body)
	if err != nil {
		http.Error(w, SerValidErr(err.Error()), 400)
		return
	}

	// Обновление данных
	user := data.GetOrCreateUser(user_id)
	err = user.EventUpd(date, body, event_id)

	if err != nil {
		http.Error(w, SerLogicErr(user_id, err.Error()), 503)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(SerJsonResult("Updated"))
}

// * /delete_event ?user_id=3&date=2019-09-09&event_id=0 + body: { "data": "вставить текст" }
// * Ответ: { "result": "deleted" }
func DelHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка метода
	if r.Method != "POST" {
		fmt.Fprint(w, "Поддерживается только POST")
		return
	}
	defer r.Body.Close()

	// Преобразование и валидация параметров url
	user_id, date, event_id, valid := ParseValid_UserDateEvent_id(r.URL.Query())
	if !valid {
		http.Error(w, SerValidErr("incorrect url params"), 400)
		return
	}

	// Обновление данных
	user := data.GetOrCreateUser(user_id)
	err := user.EventDel(date, event_id)

	if err != nil {
		http.Error(w, SerLogicErr(user_id, err.Error()), 503)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(SerJsonResult("Deleted"))
}

// -------- GET Обработчики --------

// * GET /events_for_day + ?user_id=3&date=2019-09-09
// * Ответ: { "result": [ {"date": 2007-09-09, "event": "a"}... ] }
func DayEventsHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка метода
	if r.Method != "GET" {
		fmt.Fprint(w, "Поддерживается только GET")
		return
	}
	defer r.Body.Close()

	// Преобразование и валидация параметров url
	date, user_id, valid := ParseValid_DateUser(r.URL.Query())
	if !valid {
		http.Error(w, SerValidErr("incorrect url params"), 400)
		return
	}

	// Получение пользователя поиск событий
	user := data.GetOrCreateUser(user_id)
	events := user.EventsDay(date)
	events_json := SerOneDayEvents(events, date)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(events_json)
}

// * GET /events_for_week + ?user_id=3&date=2019-09-09
// * Ответ: { "result": [ {"date": 2007-09-09, "event": "a"}... ] }
func WeekEventsHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка метода
	if r.Method != "GET" {
		fmt.Fprint(w, "Поддерживается только GET")
		return
	}
	defer r.Body.Close()

	// Преобразование и валидация параметров url
	date, user_id, valid := ParseValid_DateUser(r.URL.Query())
	if !valid {
		http.Error(w, SerValidErr("incorrect url params"), 400)
		return
	}

	// Получение пользователя поиск событий
	user := data.GetOrCreateUser(user_id)
	events, dates := user.EventsWeek(date)
	events_json := SerManyDayEvents(events, dates)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(events_json)
}

// * GET /events_for_month + ?user_id=3&date=2019-09-09
// * Ответ: { "result": [ {"date": 2007-09-09, "event": "a"}... ] }
func MonthEventsHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка метода
	if r.Method != "GET" {
		fmt.Fprint(w, "Поддерживается только GET")
		return
	}
	defer r.Body.Close()

	// Преобразование и валидация параметров url
	date, user_id, valid := ParseValid_DateUser(r.URL.Query())
	if !valid {
		http.Error(w, SerValidErr("incorrect url params"), 400)
		return
	}

	// Получение пользователя поиск событий
	user := data.GetOrCreateUser(user_id)
	events, dates := user.EventsMonth(date)
	events_json := SerManyDayEvents(events, dates)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(events_json)
}

// - Вспомогательное для распаковки -

// Чтение тела запроса
func ParseBody(r io.ReadCloser) (string, error) {
	raw_body, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	// Распаковка в строку
	var body_json PostBodyScheme
	if json.Unmarshal(raw_body, &body_json) != nil {
		return "", err
	}
	return body_json.Data, nil
}

// ---------- Валидация -------------

// Парсинг и валидация
func ParseValid_DateUser(u url.Values) (EventDate, int, bool) {
	usr := u.Get("user_id")
	dt := u.Get("date")
	if len(dt) == 0 || len(usr) == 0 {
		return EventDate(""), 0, false
	}

	usr_id, err := strconv.ParseUint(usr, 10, 32)
	if err != nil {
		return EventDate(""), 0, false
	}

	ev_tm, err := time.Parse("2006-01-02", dt)
	if err != nil {
		return EventDate(""), 0, false
	}

	return EventDate(FormatTime(ev_tm)), int(usr_id), true
}

// Для парсинга и валидации полей user_id и date используется уже существующий метод
func ParseValid_UserDateEvent_id(u url.Values) (int, EventDate, int, bool) {
	date, user_id, ok := ParseValid_DateUser(u)
	if !ok {
		return 0, EventDate(""), 0, false
	}

	ev_id := u.Get("event_id")
	if len(ev_id) == 0 {
		return 0, EventDate(""), 0, false
	}
	event_id, err := strconv.ParseUint(ev_id, 10, 32)
	if err != nil {
		return 0, EventDate(""), 0, false
	}

	return user_id, date, int(event_id), true
}
