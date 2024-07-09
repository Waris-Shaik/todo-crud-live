package types

import "time"

type TodoStore interface {
	CreateTodo(Todo) error
	GetTodos(id int) ([]*Todo, error)
	GetTodoByID(id int) (*Todo, error)
	UpdateTodo(id int) error
	DeleteTodo(id int) error
}

type TodoPayload struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
}

type Todo struct {
	ID          int       `json:"_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	UserID      int       `json:"userID"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"-"`
}

type LoginUserPayload struct {
	Text     string `json:"text" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(User) (int, error)
}

type User struct {
	ID        int       `json:"_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	UserName  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
}

type RegisterUserPayload struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	UserName  string `json:"username" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
}
