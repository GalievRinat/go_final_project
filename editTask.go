package main

import (
	"bytes"
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
		jsonError(w, "Ошибка: пустой заголовок", err)
		return
	}

	if task.Date == "" {
		task.Date = Now.Format("20060102")
	}

	_, err = time.Parse("20060102", task.Date)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonError(w, "Ошибка: неверный формат даты", err)
		return
	}

	if Now.Format("20060102") > task.Date {
		if task.Repeat == "" {
			task.Date = Now.Format("20060102")
		} else {
			task.Date, err = NextDate(Now, task.Date, task.Repeat)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				jsonError(w, "Ошибка даты/повторения", err)
				return
			}
		}
	}

	res, err := taskRepo.Edit(task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonError(w, "Ошибка обновления задачи в БД", err)
		return
	}

	row_count, _ := res.RowsAffected()
	if row_count == 0 {
		jsonError(w, "Задача не найдена", err)
		//w.WriteHeader(http.StatusOK)
		return
	}

	resp := []byte("{}")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		fmt.Println("Ошибка записи данных в соединение:", err)
		return
	}
}
