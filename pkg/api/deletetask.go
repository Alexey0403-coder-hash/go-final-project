package api

import (
	"encoding/json"
	"errors"
	"final_project/pkg/db"
	"log"
	"net/http"
)

// DeleteHandler реализует удаление записи
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	err := db.DeleteTask(id)

	switch {
	case errors.Is(err, db.ErrTaskNotFound):
		writeJSONError(w, "задача не найдена", http.StatusNotFound)
		return

	case err != nil:
		log.Printf("ошибка удаления задачи: %v", err)
		writeJSONError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(struct{}{})
	if err != nil {
		log.Printf("ошибка отправки JSON: %v", err)
		return
	}
}
