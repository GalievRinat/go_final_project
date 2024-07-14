package main

import (
	"fmt"
	"net/http"
)

func apiTaskDelete(w http.ResponseWriter, r *http.Request) {
	fmt.Println("TaskDone")
	id := r.URL.Query().Get("id")
	fmt.Println(id)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	task, err := taskRepo.getbyID(id)
	if err != nil {
		jsonError(w, "Задача не найдена", err)
		return
	}

	err = taskRepo.Delete(task)
	if err != nil {
		jsonError(w, "Ошибка удаления", err)
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
