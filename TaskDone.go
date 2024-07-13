package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

func apiTaskDone(w http.ResponseWriter, r *http.Request) {
	fmt.Println("TaskDone")
	id := r.URL.Query().Get("id")
	fmt.Println(id)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	row := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id", sql.Named("id", id))

	var task Task
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		fmt.Println("Задача не найдена:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonError("Задача не найдена"))
		return
	}

	if task.Repeat == "" {
		_, err = db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
		if err != nil {
			fmt.Println("Ошибка удаления:", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(jsonError("Ошибка удаления"))
			return
		}
	} else {
		task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(jsonError("Ошибка даты/повторения"))
			return
		}
		_, err = db.Exec("UPDATE scheduler SET date = :date WHERE id = :id",
			sql.Named("date", task.Date),
			sql.Named("id", task.ID))
		if err != nil {
			fmt.Println("Ошибка обновления задачи в БД:", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(jsonError("Ошибка обновления задачи в БД"))
			return
		}
	}

	answer := []byte("{}")
	fmt.Println(string(answer))
	w.WriteHeader(http.StatusOK)
	w.Write(answer)
}
