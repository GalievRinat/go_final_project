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
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonError("Задача не найдена"))
		return
	}

	resp, err := json.Marshal(task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonError("Ошибка сериализации JSON"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
