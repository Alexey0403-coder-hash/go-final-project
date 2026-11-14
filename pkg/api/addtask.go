package api

import (
	"encoding/json"
	"final_project/pkg/db"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ErrResponse struct {
	Error string `json:"error"`
}

type LastID struct {
	Id string `json:"id"`
}

// AddTaskHandler обрабатывает запрос на добавление записи в БД и возвращает, в случае успеха, ее ID в формате JSON
func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var addTask db.Task

	err := json.NewDecoder(r.Body).Decode(&addTask)
	if err != nil {
		writeJSONError(w, "Ошибка декодирования запроса JSON")
		return
	}

	if strings.TrimSpace(addTask.Title) == "" {
		writeJSONError(w, "Ошибка запроса. Поле title не может быть пустым")
		return
	}

	err = checkDate(&addTask)
	if err != nil {
		writeJSONError(w, err.Error())
		return
	}

	lastIDint, err := db.AddTask(&addTask)
	if err != nil {
		writeJSONError(w, err.Error())
		return
	}

	lastIDstring := strconv.FormatInt(lastIDint, 10)

	last := LastID{
		Id: lastIDstring,
	}

	w.Header().Set("Content-type", "application/json")

	err = json.NewEncoder(w).Encode(last)
	if err != nil {
		log.Printf("Ошибка отправки JSON: %v", err)
		return
	}
}

// checkDate реализует проверку корректности полученной даты из запроса
func checkDate(addTask *db.Task) error {
	now := time.Now()

	var next string

	if addTask.Date == "" {
		addTask.Date = now.Format("20060102")
	}

	t, err := time.Parse("20060102", addTask.Date)
	if err != nil {
		return fmt.Errorf("Ошибка: некорректый формат даты: %w", err)
	}

	if addTask.Repeat != "" {
		next, err = NextDate(now, addTask.Date, addTask.Repeat)
		if err != nil {
			return err
		}
	}

	if afterNow(now, t) {
		if addTask.Repeat == "" {
			addTask.Date = now.Format("20060102")
		} else {
			addTask.Date = next
		}
	}
	return nil
}

// writeJSONError реализует сериализацию JSON с ошибкой
func writeJSONError(w http.ResponseWriter, textErr string) {
	w.Header().Set("Content-Type", "application/json")
	errResponse := ErrResponse{
		Error: textErr,
	}
	err := json.NewEncoder(w).Encode(errResponse)
	if err != nil {
		log.Printf("Ошибка отправки JSON: %v", err)
	}
}
