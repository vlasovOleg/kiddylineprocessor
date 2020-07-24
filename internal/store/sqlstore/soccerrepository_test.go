package sqlstore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vlasovoleg/kiddyLineProcessor/internal/store/sqlstore"
)

func TestSoccerRepository(t *testing.T) {
	db, _ := sqlstore.TestDB(t)
	store := sqlstore.New(db)

	err := store.SoccerRepository().UpdateCoefficient(1.75)
	assert.NoError(t, err)

	c, err := store.SoccerRepository().GetCoefficient()
	assert.NoError(t, err)
	assert.Equal(t, float32(1.75), c)
}
