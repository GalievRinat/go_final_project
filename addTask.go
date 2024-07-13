package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func apiAddTask(w http.ResponseWriter, r *http.Request) {
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

	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonError("Ошибка добавления задачи в БД"))
		return
	}

	id, _ := res.LastInsertId()

	answer, err := json.Marshal(map[string]int64{"id": id})
	if err != nil {
		fmt.Println("Ошибка генерации JSON для ошибки:", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(answer)
}
