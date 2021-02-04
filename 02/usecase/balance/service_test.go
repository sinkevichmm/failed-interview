package balance

import (
	"failed-interview/02/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_List(t *testing.T) {
	repo := newInmem()
	repo.create(&entity.Balance{ID: 1, Value: 100})
	repo.create(&entity.Balance{ID: 2, Value: 100})
	m := NewService(repo)

	l, err := m.ListBalances()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(l))
}

func Test_Get(t *testing.T) {
	repo := newInmem()
	repo.create(&entity.Balance{ID: 1, Value: 100})
	repo.create(&entity.Balance{ID: 2, Value: 100})
	m := NewService(repo)

	t.Run("Update not found", func(t *testing.T) {
		err := m.UpdateBalance(0, 0, 50)
		assert.ErrorIs(t, entity.ErrNotFound, err)
	})
	t.Run("Update not found from", func(t *testing.T) {
		err := m.UpdateBalance(0, 2, 50)
		assert.ErrorIs(t, entity.ErrIDFromNotFound, err)
	})

	t.Run("Update not found to", func(t *testing.T) {
		err := m.UpdateBalance(1, 0, 50)
		assert.ErrorIs(t, entity.ErrIDToNotFound, err)
	})

	t.Run("Update ok", func(t *testing.T) {
		err := m.UpdateBalance(1, 2, 50)
		assert.Nil(t, err)

		l, err := m.ListBalances()
		assert.Nil(t, err)
		assert.Equal(t, 2, len(l))

		for i := 0; i < len(l); i++ {
			if l[i].ID == 1 {
				assert.Equal(t, 50, l[i].Value)
			}
			if l[i].ID == 2 {
				assert.Equal(t, 150, l[i].Value)
			}
		}
	})
}
