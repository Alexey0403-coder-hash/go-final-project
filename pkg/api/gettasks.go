package api

import (
	"encoding/json"
	"errors"
	"final_project/pkg/db"
	"log"
	"net/http"
)

const limit = 50

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// tasksHandler реализует вывод ближайших задач в JSON формате
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(limit)

	switch {
	case errors.Is(err, db.ErrInvalidLimit),
		errors.Is(err, db.ErrLimitExceeded):
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return

	case err != nil:
		log.Printf("ошибка получения задач: %v", err)
		writeJSONError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	tasksResp := TasksResp{
		Tasks: tasks,
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(tasksResp)
	if err != nil {
		log.Printf("ошибка отправки JSON: %v", err)
	}
}
