package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/GalievRinat/go_final_project/model"
)

const dateFormat = "20060102"

func (handler *Handler) ApiAddTask(w http.ResponseWriter, r *http.Request) {
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
		task.Date = Now.Format(dateFormat)
	}

	_, err = time.Parse(dateFormat, task.Date)

	if err != nil {
		jsonError(w, "Ошибка: неверный формат даты", err)
		return
	}

	if Now.Format(dateFormat) > task.Date {
		if task.Repeat == "" {
			task.Date = Now.Format(dateFormat)
		} else {
			task.Date, err = NextDate(Now, task.Date, task.Repeat)
			if err != nil {
				jsonError(w, "Ошибка даты/повторения", err)
				return
			}
		}
	}

	id, err := handler.taskRepo.Add(task)
	if err != nil {
		fmt.Println(err)
		jsonError(w, "Ошибка добавления задачи в БД", err)
		return
	}

	answer, err := json.Marshal(map[string]int64{"id": id})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
