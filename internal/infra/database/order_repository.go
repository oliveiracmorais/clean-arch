package database

import (
	"database/sql"

	"github.com/oliveiracmorais/clean-arch/internal/entity"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{Db: db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	stmt, err := r.Db.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) GetTotal() (int, error) {
	var total int
	err := r.Db.QueryRow("Select count(*) from orders").Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *OrderRepository) List() ([]*entity.Order, error) {
	rows, err := r.Db.Query("SELECT id, price, tax, final_price FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	orders := []*entity.Order{}
	for rows.Next() {
		var id string
		var price, tax, finalPrice float64
		if err := rows.Scan(&id, &price, &tax, &finalPrice); err != nil {
			return nil, err
		}
		order, err := entity.NewOrder(id, price, tax)
		if err != nil {
			return nil, err
		}
		order.CalculateFinalPrice()
		orders = append(orders, order)
	}
	return orders, nil
}
