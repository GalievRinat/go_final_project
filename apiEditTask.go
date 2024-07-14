package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/GalievRinat/go_final_project/model"
	"github.com/GalievRinat/go_final_project/task_repository"
)

func apiEditTask(taskRepo *task_repository.TaskRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("EditTask")
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
		fmt.Println("ID", task.ID)
		Now := time.Now()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		if task.Title == "" {
			w.WriteHeader(http.StatusBadRequest)
			jsonError(w, "Ошибка: пустой заголовок", err)
			return
		}

		if task.Date == "" {
			task.Date = Now.Format(dateFormat)
		}

		_, err = time.Parse(dateFormat, task.Date)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonError(w, "Ошибка: неверный формат даты", err)
			return
		}

		if Now.Format(dateFormat) > task.Date {
			if task.Repeat == "" {
				task.Date = Now.Format(dateFormat)
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
}
