package main

import (
	"encoding/json"
	"fmt"
)

func jsonError(text string) []byte {
	answer, err := json.Marshal(map[string]string{"error": text})
	if err != nil {
		fmt.Println("Ошибка генерации JSON для ошибки:", err)
		return []byte("")
	}
	fmt.Println("Ошибка:", text)
	return answer
}