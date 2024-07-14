package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/GalievRinat/go_final_project/task_repository"
	"github.com/go-chi/chi/v5"
	gotdotenv "github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

const dateFormat = "20060102"

func main() {
	err := gotdotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	dbFile := os.Getenv("TODO_DBFILE")
	taskRepo := task_repository.TaskRepository{}
	err = taskRepo.CreateRepo(dbFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer taskRepo.DB.Close()

	r := chi.NewRouter()

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "web"))
	FileServer(r, "/", filesDir)

	r.Get("/api/nextdate", apiNextDate(&taskRepo))
	r.Post("/api/task", apiAddTask(&taskRepo))
	r.Get("/api/tasks", apiGetTasks(&taskRepo))
	r.Get("/api/task", apiGetTask(&taskRepo))
	r.Put("/api/task", apiEditTask(&taskRepo))
	r.Post("/api/task/done", apiTaskDone(&taskRepo))
	r.Delete("/api/task", apiTaskDelete(&taskRepo))

	addr := fmt.Sprintf(":%s", os.Getenv("TODO_PORT"))
	fmt.Printf("Start web server on port [%s]\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		fmt.Printf("Start server error: %s", err.Error())
		return
	}
}
