package sqlstore

// BaseballRepository ...
type BaseballRepository struct {
	store *Store
}

// GetCoefficient ...
func (br *BaseballRepository) GetCoefficient() (float32, error) {
	coefficient := float32(0)
	err := br.store.db.QueryRow("select * from baseball").Scan(&coefficient)
	return coefficient, err
}

// UpdateCoefficient ...
func (br *BaseballRepository) UpdateCoefficient(coefficient float32) error {
	_, err := br.store.db.Exec(
		"UPDATE baseball SET coefficient = $1",
		coefficient,
	)
	return err
}
