package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/GalievRinat/go_final_project/handler"
	"github.com/go-chi/chi/v5"
	gotdotenv "github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

func main() {
	err := gotdotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	dbFile := os.Getenv("TODO_DBFILE")
	handler, err := handler.NewHandler(dbFile)
	if err != nil {
		fmt.Println("Ошибка создания handler: ", err)
		return
	}
	defer handler.CloseHandler()

	r := chi.NewRouter()

	workDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Ошибка обнаружения рабочей директории: ", err)
		return
	}

	filesDir := http.Dir(filepath.Join(workDir, "web"))
	handler.FileServer(r, "/", filesDir)

	r.Get("/api/nextdate", handler.ApiNextDate)
	r.Post("/api/task", handler.ApiAddTask)
	r.Get("/api/tasks", handler.ApiGetTasks)
	r.Get("/api/task", handler.ApiGetTask)
	r.Put("/api/task", handler.ApiEditTask)
	r.Post("/api/task/done", handler.ApiTaskDone)
	r.Delete("/api/task", handler.ApiTaskDelete)

	addr := fmt.Sprintf(":%s", os.Getenv("TODO_PORT"))
	fmt.Printf("Start web server on port [%s]\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		fmt.Printf("Start server error: %s", err.Error())
		return
	}
}
