package user

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

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM user WHERE email = ? OR username = ?", email, email)
	if err != nil {
		log.Println("Error in QUERY:", err)
		return nil, fmt.Errorf("something went wrong")
	}

	user := new(types.User)

	for rows.Next() {
		user, err = scanRowIntoUser(rows)
		if err != nil {
			log.Println("Error in rows.Next()", err)
			return nil, err
		}
	}

	if user.ID == 0 {
		log.Println("User not found:", user)
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func scanRowIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.UserName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		log.Println("Error in scanningRowIntoUser:", err)
		return nil, err
	}
	return user, nil
}

func (s *Store) GetUserByID(id int) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM user WHERE id = ?", id)
	if err != nil {
		log.Println("Error in QUERY:", err)
		return nil, err
	}

	user := new(types.User)

	for rows.Next() {
		user, err = scanRowIntoUser(rows)
		if err != nil {
			log.Println("Error in rows.Next()", err)
			return nil, err
		}
	}

	if user.ID == 0 {
		log.Println("User not found:", user)
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *Store) CreateUser(user types.User) (int, error) {
	result, err := s.db.Exec("INSERT INTO user (first_name, last_name, username, email, password) VALUES (?,?,?,?,?)", user.FirstName, user.LastName, user.UserName, user.Email, user.Password)
	if err != nil {
		log.Println("Error in QUERY:", err)
		return 0, fmt.Errorf("something went wrong")
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v\n", err)
		return 0, err
	}
	log.Printf("User created with ID: %d", id)
	return int(id), nil
}
