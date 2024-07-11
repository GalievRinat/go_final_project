package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

func apiTaskDelete(w http.ResponseWriter, r *http.Request) {
	fmt.Println("TaskDone")
	//id := chi.URLParam(r, "id")
	id := r.URL.Query().Get("id")
	fmt.Println(id)

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		fmt.Println("DB error:", err)
		return
	}
	defer db.Close()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	row := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id", sql.Named("id", id))

	var task Task
	err = row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		fmt.Println("Задача не найдена:", err)
		w.WriteHeader(http.StatusBadRequest)
		answer, _ := json.Marshal(map[string]string{"error": "Задача не найдена"})
		w.Write(answer)
		return
	}

	_, err = db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
	if err != nil {
		fmt.Println("Ошибка удаления:", err)
		w.WriteHeader(http.StatusBadRequest)
		answer, _ := json.Marshal(map[string]string{"error": "Ошибка удаления"})
		w.Write(answer)
		return
	}

	//answer, _ := json.Marshal()
	answer := []byte("{}")
	fmt.Println(string(answer))
	w.WriteHeader(http.StatusOK)
	w.Write(answer)
}