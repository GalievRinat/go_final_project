package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func apiGetTask(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	task, err := taskRepo.getbyID(id)
	if err != nil {
		fmt.Println(err)
		jsonError(w, "Задача не найдена", err)
		return
	}

	resp, err := json.Marshal(task)
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
