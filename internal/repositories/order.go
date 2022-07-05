package repositories

import (
	"database/sql"
	"github.com/vllvll/diploma/internal/dictionaries"
	"github.com/vllvll/diploma/internal/types"
	"time"
)

type Order struct {
	db *sql.DB
}

type OrderInterface interface {
	CreateOrder(number string, userID int) error
	UpdateOrder(number string, status string, accrual float32) (userID int, err error)
	GetByNumber(number string) (order types.Order, err error)
	GetOrdersByUser(userID int) (orders []types.Order, err error)
}

func NewOrderRepository(db *sql.DB) OrderInterface {
	return &Order{db: db}
}

func (o *Order) CreateOrder(number string, userID int) error {
	_, err := o.db.Exec(
		"INSERT INTO orders (number, user_id, status, uploaded_at) VALUES ($1, $2, $3, $4) RETURNING id;",
		number,
		userID,
		dictionaries.OrderNew,
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (o *Order) UpdateOrder(number string, status string, accrual float32) (userID int, err error) {
	err = o.db.QueryRow("UPDATE orders SET status = $1, accrual = $2 WHERE number = $3 RETURNING user_id", status, accrual, number).
		Scan(&userID)
	if err != nil {
		return userID, err
	}

	return userID, err
}

func (o *Order) GetByNumber(number string) (order types.Order, err error) {
	err = o.db.QueryRow("SELECT id, number, user_id, status, uploaded_at FROM orders WHERE number = $1 LIMIT 1", number).
		Scan(&order.ID, &order.Number, &order.UserID, &order.Status, &order.UploadedAt)
	if err != nil {
		return types.Order{}, err
	}

	return order, nil
}

func (o *Order) GetOrdersByUser(userID int) ([]types.Order, error) {
	var count int

	err := o.db.QueryRow("SELECT COUNT(*) as count FROM orders WHERE user_id = $1", userID).Scan(&count)
	if err != nil {
		return nil, err
	}

	rows, err := o.db.Query("SELECT id, number, user_id, status, uploaded_at, accrual FROM orders WHERE user_id = $1 ORDER BY uploaded_at", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]types.Order, 0, count)

	for rows.Next() {
		var order types.Order

		err = rows.Scan(&order.ID, &order.Number, &order.UserID, &order.Status, &order.UploadedAt, &order.Accrual)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return orders, nil
}
