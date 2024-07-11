package main

import (
	"fmt"
	"net/http"
	"time"
)

func apiNextDate(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(r.URL.Query().Get("now"))
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
