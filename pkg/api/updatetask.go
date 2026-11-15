package api

import (
	"encoding/json"
	"errors"
	"final_project/pkg/db"
	"log"
	"net/http"
	"strings"
)

// UpdateTaskHandler реализует обновление задачи
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var newTask db.Task

	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		writeJSONError(w, "невозможно прочитать JSON", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(newTask.ID) == "" {
		writeJSONError(w, "передан пустой ID", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(newTask.Title) == "" {
		writeJSONError(w, "передана пустая задача", http.StatusBadRequest)
		return
	}

	err = checkDate(&newTask)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.UpdateTask(&newTask)
	if err != nil {
		switch {
		case errors.Is(err, db.ErrTaskNotFound):
			writeJSONError(w, "задача не найдена", http.StatusNotFound)
			return

		case err != nil:
			log.Printf("ошибка обновления задачи: %v", err)
			writeJSONError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(struct{}{})
	if err != nil {
		log.Printf("ошибка отправки JSON: %v", err)
		return
	}
}
