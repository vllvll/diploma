package repositories

import (
	"database/sql"
	"github.com/vllvll/diploma/internal/types"
	"log"
	"time"
)

type User struct {
	db *sql.DB
}

type UserInterface interface {
	IsExists(login string) bool
	CreateUser(login string, password string) (id int, err error)
	GetUserHashByLogin(login string) (user types.User, err error)
	GetUserById(userId int) (user types.User, err error)
}

func NewUserRepository(db *sql.DB) UserInterface {
	return &User{db: db}
}

func (u *User) IsExists(login string) (isExist bool) {
	var count int

	err := u.db.QueryRow("SELECT 1 FROM users WHERE login = $1 LIMIT 1", login).Scan(&count)
	if err != nil {
		return false
	}

	return true
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

func (u *User) GetUserHashByLogin(login string) (user types.User, err error) {
	err = u.db.QueryRow("SELECT id, login, password_hash FROM users WHERE login = $1 LIMIT 1", login).Scan(&user.Id, &user.Login, &user.Hash)
	if err != nil {
		return types.User{}, err
	}

	return user, nil
}

func (u *User) GetUserById(userId int) (user types.User, err error) {
	err = u.db.QueryRow("SELECT id, login, password_hash FROM users WHERE id = $1 LIMIT 1", userId).Scan(&user.Id, &user.Login, &user.Hash)
	if err != nil {
		return types.User{}, err
	}

	return user, nil
}
