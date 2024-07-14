package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	gotdotenv "github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

var taskRepo TaskRepository

func main() {
	err := gotdotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	dbFile := os.Getenv("TODO_DBFILE")
	err = taskRepo.createRepo(dbFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer taskRepo.db.Close()

	r := chi.NewRouter()

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "web"))
	FileServer(r, "/", filesDir)

	r.Get("/api/nextdate", apiNextDate)
	r.Post("/api/task", apiAddTask)
	r.Get("/api/tasks", apiGetTasks)
	r.Get("/api/task", apiGetTask)
	r.Put("/api/task", apiEditTask)
	r.Post("/api/task/done", apiTaskDone)
	r.Delete("/api/task", apiTaskDelete)

	addr := fmt.Sprintf(":%s", os.Getenv("TODO_PORT"))
	fmt.Printf("Start web server on port [%s]\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		fmt.Printf("Start server error: %s", err.Error())
		return
	}
}
