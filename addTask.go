package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/GalievRinat/go_final_project/model"
)

func apiAddTask(w http.ResponseWriter, r *http.Request) {
	var task model.Task
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
		jsonError(w, "Ошибка: пустой заголовок", err)
		return
	}

	if task.Date == "" {
		task.Date = Now.Format("20060102")
	}

	_, err = time.Parse("20060102", task.Date)

	if err != nil {
		jsonError(w, "Ошибка: неверный формат даты", err)
		return
	}

	if Now.Format("20060102") > task.Date {
		if task.Repeat == "" {
			task.Date = Now.Format("20060102")
		} else {
			task.Date, err = NextDate(Now, task.Date, task.Repeat)
			if err != nil {
				jsonError(w, "Ошибка даты/повторения", err)
				return
			}
		}
	}

	res, err := taskRepo.Add(task)
	if err != nil {
		fmt.Println(err)
		jsonError(w, "Ошибка добавления задачи в БД", err)
		return
	}
	id, err := res.LastInsertId()
	if err != nil {
		jsonError(w, "Ошибка получения ID добавленной задачи", err)
		return
	}

	answer, err := json.Marshal(map[string]int64{"id": id})
	if err != nil {
		fmt.Println("Ошибка генерации JSON для ID:", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err = w.Write(answer)
	if err != nil {
		fmt.Println("Ошибка записи данных в соединение:", err)
		return
	}
}
