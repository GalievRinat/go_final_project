package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func apiEditTask(w http.ResponseWriter, r *http.Request) {
	fmt.Println("EditTask")
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("ID", task.ID)
	Now := time.Now()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if task.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonError("Ошибка: пустой заголовок"))
		return
	}

	if task.Date == "" {
		task.Date = Now.Format("20060102")
	}

	_, err = time.Parse("20060102", task.Date)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonError("Ошибка: неверный формат даты"))
		return
	}

	if Now.Format("20060102") > task.Date {
		if task.Repeat == "" {
			task.Date = Now.Format("20060102")
		} else {
			task.Date, err = NextDate(Now, task.Date, task.Repeat)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(jsonError("Ошибка даты/повторения"))
				return
			}
		}
	}

	res, err := db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))
	if err != nil {
		fmt.Println("Ошибка обновления задачи в БД:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonError("Ошибка обновления задачи в БД"))
		return
	}

	row_count, _ := res.RowsAffected()
	if row_count == 0 {
		fmt.Println("Задача не найдена:", task.ID)
		w.WriteHeader(http.StatusOK)
		w.Write(jsonError("Задача не найдена"))
		return
	}

	answer := []byte("{}")
	fmt.Println(string(answer))
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(answer)
}
