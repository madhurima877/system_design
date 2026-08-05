package repository

import "database/sql"

type URLRepository struct {
	db *sql.DB
}

func NewURLRepository(db *sql.DB) *URLRepository {
	return &URLRepository{db: db}
}
func (r *URLRepository) CreateShortURL(originalURL, shortCode string) error {
	_, err := r.db.Exec(
		`INSERT INTO urls (original_url, short_code) VALUES ($1, $2)`,
		originalURL,
		shortCode,
	)
	if err != nil {
		return err
	}
	return nil
}
func (r *URLRepository) GetOriginalURL(shortCode string) (string, error) {
	var url string
	query := `SELECT original_url FROM urls WHERE short_code=$1`

	err := r.db.QueryRow(query, shortCode).Scan(&url)
	if err != nil {
		return "", err
	}
	return url, nil
}
