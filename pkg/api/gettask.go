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
		writeJSONError(w, "ошибка: передан пустой ID", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			writeJSONError(w, "задачи с указанным ID не существует", http.StatusNotFound)
			return
		default:
			writeJSONError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
			log.Print(err)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		log.Printf("ошибка отправки JSON: %v", err)
	}
}
