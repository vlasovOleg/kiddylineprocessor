package sqlstore

type soccerRepository struct {
	store *Store
}

func (br *soccerRepository) GetCoefficient() (float32, error) {
	coefficient := float32(0)
	err := br.store.db.QueryRow("select * from soccer").Scan(&coefficient)
	return coefficient, err
}

func (br *soccerRepository) UpdateCoefficient(coefficient float32) error {
	_, err := br.store.db.Exec(
		"UPDATE soccer SET coefficient = $1",
		coefficient,
	)
	return err
}
