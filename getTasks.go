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
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonError("Ошибка получения списка задач"))
		return
	}

	fmt.Println(tasks)
	if tasks == nil {
		w.WriteHeader(http.StatusBadRequest)
		tasks = make([]Task, 0)
	}
	resp, err := json.Marshal(map[string][]Task{"tasks": tasks})

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonError("Ошибка сериализации JSON"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
