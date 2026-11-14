package api

import (
	"encoding/json"
	"final_project/pkg/db"
	"log"
	"net/http"
)

// DeleteHandler реализует удаление записи
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	err := db.DeleteTask(id)
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
