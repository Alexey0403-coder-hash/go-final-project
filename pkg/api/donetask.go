package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"final_project/pkg/db"
	"log"
	"net/http"
	"time"
)

// DoneHandler реализует отметку о выполнении задачи. Если повторяется - передвигает дату выполнения, если одноразовая - удаляется
func DoneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	var task *db.Task

	task, err := db.GetTask(id)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		writeJSONError(w, "задача не найдена", http.StatusNotFound)
		return

	case err != nil:
		log.Printf("ошибка получения задачи: %v", err)
		writeJSONError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	if task.Repeat == "" {
		err := db.DeleteTask(id)
		if err != nil {
			switch {
			case errors.Is(err, db.ErrTaskNotFound):
				writeJSONError(w, "задача не найдена", http.StatusNotFound)
				return

			default:
				log.Printf("ошибка удаления задачи: %v", err)
				writeJSONError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
				return
			}
		}
	} else {
		now := time.Now()
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = db.UpdateDate(next, id)
		switch {
		case errors.Is(err, db.ErrTaskNotFound):
			writeJSONError(w, "задача не найдена", http.StatusNotFound)
			return

		case err != nil:
			log.Printf("ошибка удаления задачи: %v", err)
			writeJSONError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(struct{}{})
	if err != nil {
		log.Printf("ошибка отправки JSON: %v", err)
	}
}
