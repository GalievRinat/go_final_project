package handler

import (
	"fmt"
	"net/http"
	"time"
)

func (handler *Handler) ApiTaskDone(w http.ResponseWriter, r *http.Request) {
	fmt.Println("TaskDone")
	id := r.URL.Query().Get("id")
	fmt.Println(id)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	task, err := handler.taskRepo.GetbyID(id)
	if err != nil {
		jsonError(w, "Задача не найдена", err)
		return
	}

	if task.Repeat == "" {
		err = handler.taskRepo.Delete(task)
		if err != nil {
			jsonError(w, "Ошибка удаления", err)
			return
		}
	} else {
		task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			jsonError(w, "Ошибка даты/повторения", err)
			return
		}
		_, err = handler.taskRepo.Edit(task)
		if err != nil {
			jsonError(w, "Ошибка обновления задачи в БД", err)
			return
		}
	}

	resp := []byte("{}")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		fmt.Println("Ошибка записи данных в соединение:", err)
		return
	}

}
