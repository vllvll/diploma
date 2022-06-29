package repositories

import (
	"database/sql"
	"github.com/vllvll/diploma/internal/types"
	"time"
)

type Balance struct {
	db *sql.DB
}

type BalanceInterface interface {
	CreateBalance(userID int) error
	UpdateBalance(userID int, sum float32) error
	GetSumAndWithdrawals(userID int) (balance types.ResponseBalance, err error)
	GetByUserID(userID int) (balance types.Balance, err error)
	AddWithdraw(userID int, number string, sum float32) (bool, error)
	GetWithdrawals(userID int) ([]types.Withdraw, error)
}

func NewBalanceRepository(db *sql.DB) BalanceInterface {
	return &Balance{db: db}
}

func (b *Balance) CreateBalance(userID int) error {
	_, err := b.db.Exec(
		"INSERT INTO balances (user_id, created_at, updated_at) VALUES ($1, $2, $3) RETURNING id;",
		userID,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (b *Balance) UpdateBalance(userID int, sum float32) error {
	_, err := b.db.Exec("UPDATE balances SET sum = sum + $1 WHERE user_id = $2", sum, userID)
	if err != nil {
		return err
	}

	return nil
}

func (b *Balance) GetSumAndWithdrawals(userID int) (balance types.ResponseBalance, err error) {
	err = b.db.QueryRow("SELECT b.sum, COALESCE((SELECT SUM(sum) FROM withdraw WHERE user_id = $1 GROUP BY user_id), 0) AS withdrawn  FROM balances b WHERE b.user_id = $1", userID).
		Scan(&balance.Current, &balance.Withdrawn)

	if err != nil {
		return types.ResponseBalance{}, err
	}

	return balance, nil
}

func (b *Balance) GetByUserID(userID int) (balance types.Balance, err error) {
	err = b.db.QueryRow("SELECT id, user_id, sum, created_at, updated_at FROM balances WHERE user_id = $1 LIMIT 1", userID).
		Scan(&balance.ID, &balance.UserID, &balance.Sum, &balance.CreatedAt, &balance.UpdatedAt)
	if err != nil {
		return types.Balance{}, err
	}

	return balance, nil
}

func (b *Balance) AddWithdraw(userID int, number string, sum float32) (bool, error) {
	tx, err := b.db.Begin()
	if err != nil {
		return false, err
	}

	balance, err := b.GetByUserID(userID)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return false, err
		}

		return false, err
	}

	finalSum := balance.Sum - sum

	if finalSum < 0 {
		if err = tx.Rollback(); err != nil {
			return false, err
		}

		return false, nil
	}

	_, err = tx.Exec("UPDATE balances SET sum = $1 WHERE user_id = $2", finalSum, userID)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return false, err
		}

		return false, err
	}

	_, err = tx.Exec(
		"INSERT INTO withdraw (user_id, number, sum, created_at) VALUES ($1, $2, $3, $4) RETURNING id;",
		userID,
		number,
		sum,
		time.Now(),
	)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return false, err
		}

		return false, err
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return true, nil
}

func (b *Balance) GetWithdrawals(userID int) ([]types.Withdraw, error) {
	var count int

	err := b.db.QueryRow("SELECT COUNT(*) as count FROM withdraw WHERE user_id = $1", userID).Scan(&count)
	if err != nil {
		return nil, err
	}

	rows, err := b.db.Query("SELECT id, user_id, number, sum, created_at FROM withdraw WHERE user_id = $1 ORDER BY created_at", userID)
	if err != nil || rows.Err() != nil {
		return nil, err
	}
	defer rows.Close()

	withdrawals := make([]types.Withdraw, 0, count)

	for rows.Next() {
		var withdraw types.Withdraw

		err = rows.Scan(&withdraw.ID, &withdraw.UserID, &withdraw.Number, &withdraw.Sum, &withdraw.CreatedAt)
		if err != nil {
			return nil, err
		}

		withdrawals = append(withdrawals, withdraw)
	}

	return withdrawals, nil
}
