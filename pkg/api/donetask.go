package api

import (
	"encoding/json"
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

	if err != nil {
		writeJSONError(w, err.Error())
		return
	}

	if task.Repeat == "" {
		err := db.DeleteTask(id)
		if err != nil {
			writeJSONError(w, err.Error())
			return
		}
	} else {
		now := time.Now()
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSONError(w, err.Error())
			return
		}

		err = db.UpdateDate(next, id)
		if err != nil {
			writeJSONError(w, err.Error())
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(struct{}{})
	if err != nil {
		log.Printf("Ошибка отправки JSON: %v", err)
	}
}
