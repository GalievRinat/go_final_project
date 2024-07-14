package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/GalievRinat/go_final_project/task_repository"
)

func apiNextDate(taskRepo *task_repository.TaskRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now, err := time.Parse(dateFormat, r.URL.Query().Get("now"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		date := r.URL.Query().Get("date")
		repeat := r.URL.Query().Get("repeat")
		s, err := NextDate(now, date, repeat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		result := fmt.Sprintf(s)
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(result))
		if err != nil {
			fmt.Println("Ошибка записи данных в соединение:", err)
			return
		}
	}
}
