package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GalievRinat/go_final_project/model"
)

func (handler *Handler) ApiGetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	tasks, err := handler.taskRepo.GetAll()
	if err != nil {
		jsonError(w, "Ошибка получения списка задач", err)
		return
	}

	fmt.Println(tasks)
	if tasks == nil {
		tasks = make([]model.Task, 0)
	}

	resp, err := json.Marshal(map[string][]model.Task{"tasks": tasks})
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
