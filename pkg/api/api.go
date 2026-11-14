package api

import (
	"net/http"
)

// Init реализует регистрацию обработчиков запросов
func Init() {

	http.HandleFunc("GET /api/nextdate", nextDayHandler)
	http.HandleFunc("POST /api/signin", authHandler)
	http.HandleFunc("GET /api/tasks", auth(tasksHandler))
	http.HandleFunc("POST /api/task/done", auth(DoneHandler))
	http.HandleFunc("POST /api/task", auth(AddTaskHandler))
	http.HandleFunc("GET /api/task", auth(getTaskHandler))
	http.HandleFunc("PUT /api/task", auth(UpdateTaskHandler))
	http.HandleFunc("DELETE /api/task", auth(DeleteHandler))
}
