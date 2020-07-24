package sqlstore

// FootballRepository ...
type FootballRepository struct {
	store *Store
}

// GetCoefficient ...
func (br *FootballRepository) GetCoefficient() (float32, error) {
	coefficient := float32(0)
	err := br.store.db.QueryRow("select * from football").Scan(&coefficient)
	return coefficient, err
}

// UpdateCoefficient ...
func (br *FootballRepository) UpdateCoefficient(coefficient float32) error {
	_, err := br.store.db.Exec(
		"UPDATE football SET coefficient = $1",
		coefficient,
	)
	return err
}
