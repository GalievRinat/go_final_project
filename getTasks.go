package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func apiGetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	tasks, err := taskRepo.getAll()
	if err != nil {
		jsonError(w, "Ошибка получения списка задач", err)
		return
	}

	fmt.Println(tasks)
	if tasks == nil {
		tasks = make([]Task, 0)
	}

	resp, err := json.Marshal(map[string][]Task{"tasks": tasks})
	if err != nil {
		jsonError(w, "Ошибка сериализации JSON", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		fmt.Println("Ошибка записи данных в соединение:", err)
		return
	}
}
