package task_repository

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/GalievRinat/go_final_project/model"
	_ "modernc.org/sqlite"
)

type TaskRepository struct {
	DB *sql.DB
}

func (taskRepo *TaskRepository) CreateRepo(dbFile string) error {
	fmt.Printf("DB on file [%s]\n", dbFile)
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
		fmt.Println("DB file not exist")
	}

	taskRepo.DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if install {
		fmt.Println("Create DB table")
		_, err := taskRepo.DB.Exec("CREATE TABLE scheduler ( id INTEGER PRIMARY KEY AUTOINCREMENT, date VARCHAR (8), title VARCHAR(255), comment TEXT, repeat VARCHAR(128))")

		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func (taskRepo *TaskRepository) GetbyID(taskID string) (model.Task, error) {
	var task model.Task
	row := taskRepo.DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id", sql.Named("id", taskID))
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		fmt.Println(err)
		return model.Task{}, err
	}
	return task, nil
}

func (taskRepo *TaskRepository) GetAll() ([]model.Task, error) {
	var tasks []model.Task
	rows, err := taskRepo.DB.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 50")
	if err != nil {
		fmt.Println(err)
		return []model.Task{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var task model.Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			fmt.Println(err)
			return []model.Task{}, err
		}
		tasks = append(tasks, task)
	}
	if rows.Err() != nil {
		return []model.Task{}, rows.Err()
	}
	return tasks, nil
}

func (taskRepo *TaskRepository) Delete(task model.Task) error {
	fmt.Println("Удаление задачи", task.ID)
	_, err := taskRepo.DB.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", task.ID))
	return err
}

func (taskRepo *TaskRepository) Edit(task model.Task) (int64, error) {
	fmt.Println("Изменение задачи", task.ID)
	res, err := taskRepo.DB.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))
	if err != nil {
		return 0, err
	}
	row_count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return row_count, err
}

func (taskRepo *TaskRepository) Add(task model.Task) (int64, error) {
	fmt.Println("Добавление задачи", task.Title)
	res, err := taskRepo.DB.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, err
}
