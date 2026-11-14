package api

import (
	"encoding/json"
	"final_project/pkg/db"
	"log"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// tasksHandler реализует вывод ближайших задач в JSON формате
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50)
	if err != nil {
		writeJSONError(w, err.Error())
		return
	}

	tasksResp := TasksResp{
		Tasks: tasks,
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(tasksResp)
	if err != nil {
		log.Printf("Ошибка отправки JSON: %v", err)
	}
}
