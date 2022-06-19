package repositories

import (
	"database/sql"
	"log"
	"time"
)

type User struct {
	db *sql.DB
}

type UserInterface interface {
	IsExists(login string) (isExist bool, err error)
	CreateUser(login string, password string) (id int, err error)
	GetUserByLoginAndHash()
}

func NewUserRepository(db *sql.DB) UserInterface {
	return &User{db: db}
}

func (u *User) IsExists(login string) (isExist bool, err error) {
	var count int

	err = u.db.QueryRow("SELECT COUNT(*) FROM users WHERE login = $1 LIMIT 1", login).Scan(&count)
	if err != nil {
		log.Printf("Error with get counter count: %v", err)

		return true, err
	}

	return count > 0, nil
}

func (u *User) CreateUser(login string, password string) (id int, err error) {
	err = u.db.QueryRow(
		"INSERT INTO users (login, password_hash, created_at) VALUES ($1, $2, $3) RETURNING id;",
		login,
		password,
		time.Now(),
	).Scan(&id)

	if err != nil {
		log.Printf("Error create user: %v", err)

		return 0, err
	}

	return id, nil
}

func (u *User) GetUserByLoginAndHash() {

}
