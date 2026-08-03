package repository

import (
	"database/sql"
	"errors"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}
func (repo *OrderRepository) UpdateOrderStatus(orderID int, status string) error {
	query := `UPDATE orders SET status = $1 WHERE id = $2`

	res, err := repo.db.Exec(query, status, orderID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows updated")
	}

	return nil
}
