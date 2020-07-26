package sqlstore

type baseballRepository struct {
	store *Store
}

func (br *baseballRepository) GetCoefficient() (float32, error) {
	coefficient := float32(0)
	err := br.store.db.QueryRow("select * from baseball").Scan(&coefficient)
	return coefficient, err
}

func (br *baseballRepository) UpdateCoefficient(coefficient float32) error {
	_, err := br.store.db.Exec(
		"UPDATE baseball SET coefficient = $1",
		coefficient,
	)
	return err
}
