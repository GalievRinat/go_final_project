package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

func apiGetTasks(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 50")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		answer, _ := json.Marshal(map[string]string{"error": "Ошибка чтения из БД"})
		w.Write(answer)
		return
	}
	defer rows.Close()
	var tasks []Task
	for rows.Next() {
		var task Task

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			answer, _ := json.Marshal(map[string]string{"error": "Ошибка парсинга строк БД"})
			w.Write(answer)
			return
		}

		tasks = append(tasks, task)
	}

	fmt.Println(tasks)
	if tasks == nil {
		w.WriteHeader(http.StatusBadRequest)
		tasks = make([]Task, 0)
	}
	resp, err := json.Marshal(map[string][]Task{"tasks": tasks})

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		answer, _ := json.Marshal(map[string]string{"error": "Ошибка сериализации JSON"})
		w.Write(answer)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
