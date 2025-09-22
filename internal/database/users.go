package database

import (
	"context"
	"database/sql"
	"time"
)

type UserModel struct {
	DB *sql.DB
}
type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

func (m *UserModel) Insert(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `INSERT INTO users (name, email, password) VALUES (?, ?, ?) RETURNING id`
	return m.DB.QueryRowContext(ctx, query, user.Name, user.Email, user.Password).Scan(&user.Id)
}

func (m *UserModel) getUser(query string, args ...interface{}) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User
	err:= m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (m *UserModel) GetById(id int) (*User, error) {
	query := `SELECT * FROM users WHERE id = ?`
	return m.getUser(query, id)
}

func (m *UserModel) GetByEmail(email string) (*User, error) {
	query := `SELECT * FROM users WHERE email = ?`
	return m.getUser(query, email)
}
