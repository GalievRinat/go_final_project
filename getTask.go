package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

func apiGetTask(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	row := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id", sql.Named("id", id))

	var task Task
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		answer, _ := json.Marshal(map[string]string{"error": "Задача не найдена"})
		w.Write(answer)
		return
	}

	resp, err := json.Marshal(task)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		answer, _ := json.Marshal(map[string]string{"error": "Ошибка сериализации JSON"})
		w.Write(answer)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
