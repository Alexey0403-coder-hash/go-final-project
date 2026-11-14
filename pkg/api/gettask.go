package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"final_project/pkg/db"
	"log"
	"net/http"
	"strings"
)

// getTaskHandler реализует поиск записи по ID и отправляет ее клиенту
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		writeJSONError(w, "Ошибка: передан пустой ID")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			writeJSONError(w, "Задачи с указанным ID не существует")
			return
		default:
			writeJSONError(w, "Внутренняя ошибка сервера")
			log.Print(err)
			return
		}
	}

	w.Header().Set("Content-Type", "applications/json")

	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		log.Printf("Ошибка отправки JSON: %v", err)
	}
}
