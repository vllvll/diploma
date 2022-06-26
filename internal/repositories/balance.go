package repositories

import (
	"database/sql"
	"github.com/vllvll/diploma/internal/types"
	"log"
	"time"
)

type Balance struct {
	db *sql.DB
}

type BalanceInterface interface {
	CreateBalance(userId int) error
	GetByUserId(userId int) (balance types.Balance, err error)
	UpdateBalance(userId int, number int64, sum int) (bool, error)
	GetWithdrawals(userId int) ([]types.Withdraw, error)
}

func NewBalanceRepository(db *sql.DB) BalanceInterface {
	return &Balance{db: db}
}

func (b *Balance) CreateBalance(userId int) error {
	_, err := b.db.Exec(
		"INSERT INTO balances (user_id, created_at, updated_at) VALUES ($1, $2, $3) RETURNING id;",
		userId,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		log.Printf("Error create balance: %v", err)

		return err
	}

	return nil
}

func (b *Balance) GetByUserId(userId int) (balance types.Balance, err error) {
	err = b.db.QueryRow("SELECT id, user_id, sum, created_at, updated_at FROM balances WHERE user_id = $1 LIMIT 1", userId).
		Scan(&balance.Id, &balance.UserId, &balance.Sum, &balance.CreatedAt, &balance.UpdatedAt)
	if err != nil {
		return types.Balance{}, err
	}

	return balance, nil
}

func (b *Balance) UpdateBalance(userId int, number int64, sum int) (bool, error) {
	tx, err := b.db.Begin()
	if err != nil {
		log.Printf("Error with open transaction: %v\n", err)

		return false, err
	}

	balance, err := b.GetByUserId(userId)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			log.Printf("Error with unable to rollback: %v", err)
		}

		return false, err
	}

	finalSum := balance.Sum - sum

	if finalSum < 0 {
		if err = tx.Rollback(); err != nil {
			log.Printf("Error with unable to rollback: %v", err)
		}

		return false, nil
	}

	_, err = tx.Exec("UPDATE balances SET sum = $1 WHERE user_id = $2", finalSum, userId)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			log.Printf("Error with unable to rollback: %v", err)
		}

		log.Printf("Error with update balance: %v", err)
		return false, err
	}

	_, err = tx.Exec(
		"INSERT INTO withdraw (user_id, number, sum, created_at) VALUES ($1, $2, $3, $4) RETURNING id;",
		userId,
		number,
		sum,
		time.Now(),
	)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			log.Printf("Error with unable to rollback: %v", err)
		}

		log.Printf("Error create withdraw: %v", err)
		return false, err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Erro with unable to commit: %v", err)

		return false, err
	}

	return true, nil
}

func (b *Balance) GetWithdrawals(userId int) ([]types.Withdraw, error) {
	var count int

	err := b.db.QueryRow("SELECT COUNT(*) as count FROM withdraw WHERE user_id = $1", userId).Scan(&count)
	if err != nil {
		log.Printf("Error get count withdrawals by user: %v", err)

		return nil, err
	}

	rows, err := b.db.Query("SELECT id, user_id, number, sum, created_at FROM withdraw WHERE user_id = $1 ORDER BY created_at", userId)
	if err != nil || rows.Err() != nil {
		log.Printf("Error get withdrawals by user: %v", err)

		return nil, err
	}
	defer rows.Close()

	withdrawals := make([]types.Withdraw, 0, count)

	for rows.Next() {
		var withdraw types.Withdraw

		err = rows.Scan(&withdraw.Id, &withdraw.UserId, &withdraw.Number, &withdraw.Sum, &withdraw.CreatedAt)
		if err != nil {
			log.Printf("Error read order: %v", err)

			return nil, err
		}

		withdrawals = append(withdrawals, withdraw)
	}

	return withdrawals, nil
}
