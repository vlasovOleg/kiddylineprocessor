package sqlstore

type footballRepository struct {
	store *Store
}

func (br *footballRepository) GetCoefficient() (float32, error) {
	coefficient := float32(0)
	err := br.store.db.QueryRow("select * from football").Scan(&coefficient)
	return coefficient, err
}

func (br *footballRepository) UpdateCoefficient(coefficient float32) error {
	_, err := br.store.db.Exec(
		"UPDATE football SET coefficient = $1",
		coefficient,
	)
	return err
}
