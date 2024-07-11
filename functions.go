package main

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("Wrong format REPEAT (empty)")
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
			return "", errors.New("Wrong format REPEAT (no days count)")
		}
		days, err := strconv.Atoi(words[1])

		if err != nil {
			return "", err
		}

		if days > 400 {
			return "", errors.New("Wrong format REPEAT (days > 400)")
		}
		start_time = start_time.AddDate(0, 0, days)
		for start_time.Before(now) {
			start_time = start_time.AddDate(0, 0, days)
		}
	} else {
		return "", errors.New("Wrong format REPEAT")
	}

	str_time := start_time.Format("20060102")

	return str_time, err
}
