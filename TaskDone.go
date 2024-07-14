package main

import (
	"fmt"
	"net/http"
	"time"
)

func apiTaskDone(w http.ResponseWriter, r *http.Request) {
	fmt.Println("TaskDone")
	id := r.URL.Query().Get("id")
	fmt.Println(id)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	task, err := taskRepo.getbyID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonError("Задача не найдена"))
		return
	}

	if task.Repeat == "" {
		err = taskRepo.Delete(task)
		if err != nil {
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
		_, err = taskRepo.Edit(task)
		if err != nil {
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
