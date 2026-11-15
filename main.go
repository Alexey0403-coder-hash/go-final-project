package main

import (
	"final_project/pkg/api"
	"final_project/pkg/db"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Println("файл .env не найден, используются системные переменные окружения")
	}
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}
	_, err := db.Init(dbFile)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := db.Close()
		if err != nil {
			log.Printf("ошибка закрытия БД: %v", err)
		}
	}()

	api.Init()

	dir := http.Dir("./web")
	handler := http.FileServer(dir)
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	http.Handle("/", handler)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Printf("ошибка HTTP-сервера: %v", err)
	}
}
