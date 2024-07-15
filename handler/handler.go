package handler

import (
	"fmt"

	"github.com/GalievRinat/go_final_project/task_repository"
)

type Handler struct {
	taskRepo *task_repository.TaskRepository
}

func (handler *Handler) CreateHandler(dbFile string) error {
	handler.taskRepo = &task_repository.TaskRepository{}
	err := handler.taskRepo.CreateRepo(dbFile)
	return err
}

func (handler *Handler) CloseHandler() {
	err := handler.taskRepo.DB.Close()
	fmt.Println("Ошибка закрытия handler: ", err)
}
