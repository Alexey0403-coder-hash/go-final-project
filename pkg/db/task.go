package db

import (
	"errors"
	"fmt"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

var (
	ErrInvalidLimit  = errors.New("некорректное значение limit")
	ErrLimitExceeded = errors.New("превышен максимальный limit")
	ErrTaskNotFound  = errors.New("задача не найдена")
)

// AddTask реализует добавление задачи в БД и, при успешном добавлении, возвращает ее ID
func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	result, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("ошибка добавления записи в БД: %w", err)
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("ошибка получения ID добавленной записи: %w", err)
	}
	return lastID, nil
}

// Tasks реализует вывод отсортированных задач в виде слайса
func Tasks(limit int) ([]*Task, error) {

	tasks := make([]*Task, 0)

	if limit <= 0 {
		return nil, ErrInvalidLimit
	}

	if limit > 50 {
		return nil, ErrLimitExceeded
	}

	rowsQuery := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?`

	rows, err := db.Query(rowsQuery, limit)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения записей: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		task := new(Task)

		err := rows.Scan(
			&task.ID,
			&task.Date,
			&task.Title,
			&task.Comment,
			&task.Repeat,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка прочтения данных: %w", err)
		}

		tasks = append(tasks, task)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("ошибка при переборе строк: %w", err)
	}

	return tasks, nil
}

// GetTask реализует получение записи из БД по ее ID
func GetTask(id string) (*Task, error) {
	var task Task

	row := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id=?`

	err := db.QueryRow(row, id).Scan(
		&task.ID,
		&task.Date,
		&task.Title,
		&task.Comment,
		&task.Repeat,
	)

	if err != nil {
		return nil, fmt.Errorf("ошибка чтения из БД: %w", err)
	}

	return &task, nil
}

// UpdateTask обновляет задачу
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	result, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("ошибка изменения записи в БД: %w", err)
	}
	rowsAff, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества измененных строк: %w", err)
	}

	if rowsAff == 0 {
		return ErrTaskNotFound
	}
	return nil
}

// DeleteTask реализует удаление записи из БД
func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления записи из БД: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества измененных строк: %w", err)
	}

	if rows == 0 {
		return ErrTaskNotFound
	}
	return nil
}

// UpdateDate реализует обновление даты выполнения задачи в БД
func UpdateDate(next string, id string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	result, err := db.Exec(query, next, id)
	if err != nil {
		return fmt.Errorf("ошибка изменения даты в БД: %w", err)
	}
	rowsAff, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества измененных строк: %w", err)
	}

	if rowsAff == 0 {
		return ErrTaskNotFound
	}
	return nil
}
