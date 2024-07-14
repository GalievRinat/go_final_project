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
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonError("Задача не найдена"))
		return
	}

	err = taskRepo.Delete(task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonError("Ошибка удаления"))
		return
	}

	answer := []byte("{}")
	fmt.Println(string(answer))
	w.WriteHeader(http.StatusOK)
	w.Write(answer)
}
