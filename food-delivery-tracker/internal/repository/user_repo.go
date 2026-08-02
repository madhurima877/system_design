package repository

import (
	"database/sql"
	"system_design/food-delivery-tracker/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) CreateUser(user model.User) (int, error) {
	query := `
		INSERT INTO users (name, email)
		VALUES ($1, $2)
		RETURNING id
	`

	var id int

	err := repo.db.QueryRow(query, user.Name, user.Email).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
