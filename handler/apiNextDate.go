package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("wrong format repeat (empty)")
	}
	start_time, err := time.Parse("20060102", date)
	words := strings.Fields(repeat)

	if words[0] == "y" {
		start_time = start_time.AddDate(1, 0, 0)
		for start_time.Before(now) {
			start_time = start_time.AddDate(1, 0, 0)
		}

	} else if words[0] == "d" {
		if len(words) == 1 {
			return "", errors.New("wrong format repeat (no days count)")
		}
		days, err := strconv.Atoi(words[1])

		if err != nil {
			return "", err
		}

		if days > 400 {
			return "", errors.New("wrong format repeat (days > 400)")
		}
		start_time = start_time.AddDate(0, 0, days)
		for start_time.Before(now) {
			start_time = start_time.AddDate(0, 0, days)
		}
	} else {
		return "", errors.New("wrong format repeat")
	}

	str_time := start_time.Format("20060102")

	return str_time, err
}

func (handler *Handler) ApiNextDate(w http.ResponseWriter, r *http.Request) {
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
