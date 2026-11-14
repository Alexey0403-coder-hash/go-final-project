package api

import (
	"encoding/json"
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
		writeJSONError(w, "Невозможно прочитать JSON")
		return
	}

	if strings.TrimSpace(newTask.ID) == "" {
		writeJSONError(w, "Передан пустой ID")
		return
	}

	if strings.TrimSpace(newTask.Title) == "" {
		writeJSONError(w, "Передана пустая задача")
		return
	}

	err = checkDate(&newTask)
	if err != nil {
		writeJSONError(w, err.Error())
		return
	}

	err = db.UpdateTask(&newTask)
	if err != nil {
		writeJSONError(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(struct{}{})
	if err != nil {
		log.Printf("Ошибка отправки JSON: %v", err)
		return
	}
}
