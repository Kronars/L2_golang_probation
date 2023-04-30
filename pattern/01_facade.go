package main

/*
	Реализовать паттерн «фасад».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Facade_pattern
*/

import (
	"fmt"
	"time"
)

// Подсистема #1
type Database struct {
	IsConnected bool
}

func (db *Database) Connect() {
	db.IsConnected = true
}

func (db *Database) Disconnect() {
	db.IsConnected = false
}

// Подсистема #2
type Validator struct{}

func (v *Validator) Validate(input string) bool {
	return input != ""
}

// Подсистема #3
type Processor struct{}

func (p *Processor) AddMeta(input string) string {
	prefix := time.Now().Format("Mon 15:03") + "\t"
	return prefix + input
}

// Фасад
type API struct {
	data      []string
	db        *Database
	validator *Validator
	proc      *Processor
}

func NewAPI() *API {
	return &API{
		db:        &Database{},
		validator: &Validator{},
		proc:      &Processor{},
	}
}

func (api *API) Connect() {
	api.db.Connect()
}

func (api *API) Disconnect() {
	api.db.Disconnect()
}

func (api *API) SaveData(data string) error {
	if !api.validator.Validate(data) {
		return fmt.Errorf("invalid data")
	}

	// Обработка
	data = api.proc.AddMeta(data)
	// Типо сохранение данных в базу данных
	api.data = append(api.data, data)
	return nil
}

func (api *API) Get(id byte) string {
	return api.data[id]
}

func main() {
	api := NewAPI()

	api.Connect()
	defer api.Disconnect()

	err := api.SaveData("some data")
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println(api.Get(0))
	// >>> Fri 23:11       some data
}

/*
	Фасад скрывает сложность взаимодействия между компонентами
	Предоставляет более удобный и понятный интерфейс для использования в клиентском коде

Применимость:
 * Предоставление простого интерфейса к сложной подсистеме
 * Наличие множества зависимостей между клиентами и структурами-классами реализации
 * Требуется разложить подсистему на отдельные уровни
*/
