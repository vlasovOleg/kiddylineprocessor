package sqlstore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vlasovoleg/kiddyLineProcessor/internal/store/sqlstore"
)

func TestFootballRepository(t *testing.T) {
	db, _ := sqlstore.TestDB(t)
	store := sqlstore.New(db)

	err := store.FootballRepository().UpdateCoefficient(1.75)
	assert.NoError(t, err)

	c, err := store.FootballRepository().GetCoefficient()
	assert.NoError(t, err)
	assert.Equal(t, float32(1.75), c)
}
