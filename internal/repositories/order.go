package repositories

import (
	"database/sql"
	"github.com/vllvll/diploma/internal/dictionaries"
	"github.com/vllvll/diploma/internal/types"
	"log"
	"time"
)

type Order struct {
	db *sql.DB
}

type OrderInterface interface {
	CreateOrder(number int64, userId int) error
	GetByNumber(number int64) (order types.Order, err error)
	GetOrdersByUser(userId int) (orders []types.Order, err error)
}

func NewOrderRepository(db *sql.DB) OrderInterface {
	return &Order{db: db}
}

func (o *Order) CreateOrder(number int64, userId int) error {
	_, err := o.db.Exec(
		"INSERT INTO orders (number, user_id, status, uploaded_at) VALUES ($1, $2, $3, $4) RETURNING id;",
		number,
		userId,
		dictionaries.OrderNew,
		time.Now(),
	)

	if err != nil {
		log.Printf("Error create order: %v", err)

		return err
	}

	return nil
}

func (o *Order) GetByNumber(number int64) (order types.Order, err error) {
	err = o.db.QueryRow("SELECT id, number, user_id, status, uploaded_at FROM orders WHERE number = $1 LIMIT 1", number).
		Scan(&order.Id, &order.Number, &order.UserId, &order.Status, &order.UploadedAt)
	if err != nil {
		return types.Order{}, err
	}

	return order, nil
}

func (o *Order) GetOrdersByUser(userId int) ([]types.Order, error) {
	var count int

	err := o.db.QueryRow("SELECT COUNT(*) as count FROM orders WHERE user_id = $1", userId).Scan(&count)
	if err != nil {
		log.Printf("Error get count orders by user: %v", err)

		return nil, err
	}

	rows, err := o.db.Query("SELECT id, number, user_id, status, uploaded_at FROM orders WHERE user_id = $1 ORDER BY uploaded_at", userId)
	if err != nil || rows.Err() != nil {
		log.Printf("Error get orders by user: %v", err)

		return nil, err
	}
	defer rows.Close()

	orders := make([]types.Order, 0, count)

	for rows.Next() {
		var order types.Order

		err = rows.Scan(&order.Id, &order.Number, &order.UserId, &order.Status, &order.UploadedAt)
		if err != nil {
			log.Printf("Error read order: %v", err)

			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}
