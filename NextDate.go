package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func apiNextDate(w http.ResponseWriter, r *http.Request) {
	now, err := time.Parse("20060102", r.URL.Query().Get("now"))
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
	w.Write([]byte(result))
}

func jsonError(text string) []byte {
	answer, err := json.Marshal(map[string]string{"error": text})
	if err != nil {
		fmt.Println("Ошибка генерации JSON для ошибки:", err)
		return []byte("")
	}
	return answer
}
