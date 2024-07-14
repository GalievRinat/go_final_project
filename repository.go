package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

type TaskRepository struct {
	db *sql.DB
}

func (taskRepo *TaskRepository) createRepo(dbFile string) error {
	fmt.Printf("DB on file [%s]\n", dbFile)
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
		fmt.Println("DB file not exist")
	}

	taskRepo.db, err = sql.Open("sqlite", dbFile)
	fmt.Println(taskRepo.db)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if install {
		fmt.Println("Create DB table")
		_, err := taskRepo.db.Exec("CREATE TABLE scheduler ( id INTEGER PRIMARY KEY AUTOINCREMENT, date VARCHAR (8), title VARCHAR(255), comment TEXT, repeat VARCHAR(128))")

		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

func (taskRepo *TaskRepository) getbyID(taskID string) (Task, error) {
	var task Task
	row := taskRepo.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id", sql.Named("id", taskID))
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	return task, err
}

func (taskRepo *TaskRepository) getAll() ([]Task, error) {
	var tasks []Task
	rows, err := taskRepo.db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 50")
	if err != nil {
		fmt.Println(err)
		return []Task{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			fmt.Println(err)
			return []Task{}, err
		}
		tasks = append(tasks, task)
	}
	return tasks, err
}

func (taskRepo *TaskRepository) Delete(task Task) error {
	fmt.Println("Удаление задачи", task.ID)
	_, err := taskRepo.db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", task.ID))
	return err
}

func (taskRepo *TaskRepository) Edit(task Task) (sql.Result, error) {
	fmt.Println("Изменение задачи", task.ID)
	res, err := taskRepo.db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))
	return res, err
}

func (taskRepo *TaskRepository) Add(task Task) (sql.Result, error) {
	fmt.Println("Добавление задачи", task.Title)
	res, err := taskRepo.db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	return res, err
}
