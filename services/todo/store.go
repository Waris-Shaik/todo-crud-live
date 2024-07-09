package todo

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Waris-Shaik/todo/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateTodo(todo types.Todo) error {
	_, err := s.db.Exec("INSERT INTO todo (title, description, userID) VALUES (?,?, ?)", todo.Title, todo.Description, todo.UserID)
	if err != nil {
		log.Println("Error in QUERY:", err)
		return fmt.Errorf("something went wrong")
	}

	return nil
}

func (s *Store) GetTodos(id int) ([]*types.Todo, error) {
	rows, err := s.db.Query("SELECT * FROM todo WHERE userID IN(?)", id)
	if err != nil {
		log.Println("Error in QUERY:", err)
		return nil, err
	}

	todos := make([]*types.Todo, 0)
	for rows.Next() {
		todo, err := scanRowsIntoTodo(rows)
		if err != nil {
			log.Println("Error in rows.Next():", err)
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func scanRowsIntoTodo(rows *sql.Rows) (*types.Todo, error) {
	todo := new(types.Todo)

	err := rows.Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Status,
		&todo.UserID,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		log.Println("Error in scanRowsIntoTodo:", err)
		return nil, err
	}
	return todo, nil
}

func (s *Store) GetTodoByID(id int) (*types.Todo, error) {
	rows, err := s.db.Query("SELECT * FROM todo WHERE id = ?", id)
	if err != nil {
		log.Println("Error in QUERY:", err)
		return nil, err
	}

	todo := new(types.Todo)
	for rows.Next() {
		todo, err = scanRowsIntoTodo(rows)
		if err != nil {
			log.Println("Error in rows.Next()", err)
			return nil, err
		}
	}

	if todo.ID == 0 {
		log.Println("todo not found:", todo)
		return nil, fmt.Errorf("todo not found")
	}

	return todo, nil
}

func (s *Store) UpdateTodo(id int) error {

	currentStatus, err := getCurrentTodoStatus(id, s.db)
	if err != nil {
		log.Println("Error getting current task status", err)
		return fmt.Errorf("something went wrong")
	}
	newStatus := "pending"
	if currentStatus == "pending" {
		newStatus = "completed"
	}

	_, err = s.db.Exec("UPDATE todo SET status = ? WHERE id = ?", newStatus, id)

	if err != nil {
		log.Println("Error updating task status:", err)
		return fmt.Errorf("something went wrong")
	}
	return nil
}

func getCurrentTodoStatus(id int, db *sql.DB) (string, error) {
	var status string
	err := db.QueryRow("SELECT status FROM todo WHERE id = ?", id).Scan(&status)
	if err != nil {
		log.Println("Error getting current task status from database:", err)
		return "", err
	}
	return status, nil
}
func (s *Store) DeleteTodo(id int) error {
	_, err := s.db.Exec("DELETE FROM todo WHERE id = ?", id)
	if err != nil {
		log.Println("Error in EXEC:", err)
		return fmt.Errorf("something went wrong")
	}
	return nil
}
