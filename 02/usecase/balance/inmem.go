package balance

import (
	"failed-interview/02/entity"
)

// inmem in memory repo
type inmem struct {
	m map[int]*entity.Balance
}

// newInmem create new repository
func newInmem() *inmem {
	var m = map[int]*entity.Balance{}

	return &inmem{
		m: m,
	}
}

func (r *inmem) create(b *entity.Balance) {
	r.m[b.ID] = b
}

// Get a balances
func (r *inmem) Get(idFrom int, idTo int) (bb []*entity.Balance, err error) {
	bb = make([]*entity.Balance, 0, 2)

	if from, ok := r.m[idFrom]; ok {
		bb = append(bb, from)
	}

	if to, ok := r.m[idTo]; ok {
		bb = append(bb, to)
	}

	return bb, err
}

// List balances
func (r *inmem) List() (bb []*entity.Balance, err error) {
	for _, v := range r.m {
		bb = append(bb, v)
	}

	return bb, err
}

// Update a balances
func (r *inmem) Update(idFrom int, idTo int, value int) error {
	r.m[idFrom].Value -= value
	r.m[idTo].Value += value

	return nil
}
