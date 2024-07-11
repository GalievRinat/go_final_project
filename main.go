package main

import (
	"database/sql"
	"fmt"
	"log"
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

var dbFile string

func main() {

	//fmt.Println(NextDate(time.Now(), "20250701", "y"))

	err := gotdotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	dbFile = os.Getenv("TODO_DBFILE")
	fmt.Printf("DB on file [%s]\n", dbFile)
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
		fmt.Println("DB file not exist")
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	if install {
		fmt.Println("Create DB table")
		_, err := db.Exec("CREATE TABLE scheduler ( id INTEGER PRIMARY KEY AUTOINCREMENT, date VARCHAR (8), title VARCHAR(255), comment TEXT, repeat VARCHAR(128))")

		if err != nil {
			log.Fatal(err)
		}
	}

	r := chi.NewRouter()
	//r := http.NewServeMux()
	//webDir := http.Dir("./web")
	//webFs := http.FileServer(webDir)
	//r.Handle("/", webFs)

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "web"))
	FileServer(r, "/", filesDir)

	r.Get("/api/nextdate", apiNextDate)
	r.Post("/api/task", apiAddTask)
	r.Get("/api/tasks", apiGetTasks)
	r.Get("/api/task", apiGetTask)
	r.Put("/api/task", apiEditTask)
	r.Post("/api/task/done", apiTaskDone)

	//r.Get("/tasks", handlers.GetTasks)
	//r.Post("/tasks", handlers.PostTask)
	//r.Get("/tasks/{id}", handlers.GetTaskByID)
	//r.Delete("/tasks/{id}", handlers.DeleteTask)
	port := fmt.Sprintf(":%s", os.Getenv("TODO_PORT"))
	fmt.Printf("Start web server on port [%s]\n", port)
	if err := http.ListenAndServe(port, r); err != nil {
		fmt.Printf("Start server error: %s", err.Error())
		return
	}
}
