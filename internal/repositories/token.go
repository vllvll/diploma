package repositories

import (
	"database/sql"
	"time"
)

type Token struct {
	db *sql.DB
}

type TokenInterface interface {
	IsExists(token string) bool
	CreateToken(token string, userID int) error
	GetUserIDByToken(token string) (userID int, err error)
}

func NewTokenRepository(db *sql.DB) TokenInterface {
	return &Token{db: db}
}

func (t *Token) IsExists(token string) bool {
	var count int
	err := t.db.QueryRow("SELECT 1 FROM tokens WHERE token = $1 ORDER BY last_login LIMIT 1", token).Scan(&count)

	return err == nil
}

func (t *Token) CreateToken(token string, userID int) error {
	_, err := t.db.Exec(
		"INSERT INTO tokens (token, user_id, last_login) VALUES ($1, $2, $3)",
		token,
		userID,
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (t *Token) GetUserIDByToken(token string) (userID int, err error) {
	err = t.db.QueryRow("SELECT user_id FROM tokens WHERE token = $1 ORDER BY last_login LIMIT 1", token).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
