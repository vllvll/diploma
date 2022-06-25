package repositories

import (
	"database/sql"
	"log"
	"time"
)

type Token struct {
	db *sql.DB
}

type TokenInterface interface {
	IsExists(token string) bool
	CreateToken(token string, userId int) error
	GetUserIdByToken(token string) (userId int, err error)
}

func NewTokenRepository(db *sql.DB) TokenInterface {
	return &Token{db: db}
}

func (t *Token) IsExists(token string) bool {
	var count int

	err := t.db.QueryRow("SELECT 1 FROM tokens WHERE token = $1 ORDER BY last_login LIMIT 1", token).Scan(&count)
	if err != nil {
		return false
	}

	return true
}

func (t *Token) CreateToken(token string, userId int) error {
	_, err := t.db.Exec(
		"INSERT INTO tokens (token, user_id, last_login) VALUES ($1, $2, $3)",
		token,
		userId,
		time.Now(),
	)

	if err != nil {
		log.Printf("Error create token: %v", err)

		return err
	}

	return nil
}

func (t *Token) GetUserIdByToken(token string) (userId int, err error) {
	err = t.db.QueryRow("SELECT user_id FROM tokens WHERE token = $1 ORDER BY last_login LIMIT 1", token).Scan(&userId)
	if err != nil {
		return 0, err
	}

	return userId, nil
}
