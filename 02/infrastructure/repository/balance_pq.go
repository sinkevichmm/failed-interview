package repository

import (
	"database/sql"
	"log"

	"failed-interview/02/entity"
)

type BalancePG struct {
	db *sql.DB
}

// NewbalancePG create new repository
func NewbalancePG(db *sql.DB) *BalancePG {
	return &BalancePG{
		db: db,
	}
}

// Get a balance
func (r *BalancePG) Get(idFrom int, idTo int) (bb []*entity.Balance, err error) {
	stmt, err := r.db.Prepare(`select account, balance from balance where account in ($1,$2)`)
	if err != nil {
		return nil, err
	}

	var b entity.Balance

	rows, err := stmt.Query(idFrom, idTo)

	if err != nil {
		return nil, err
	}

	bb = make([]*entity.Balance, 0, 2)

	for rows.Next() {
		err = rows.Scan(&b.ID, &b.Value)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		bb = append(bb, &b)
	}

	err = rows.Close()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = rows.Err()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return bb, nil
}

// List balances
func (r *BalancePG) List() ([]*entity.Balance, error) {
	stmt, err := r.db.Prepare(`select account, balance from balance`)
	if err != nil {
		return nil, err
	}

	var balances []*entity.Balance

	rows, err := stmt.Query()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var b entity.Balance
		err = rows.Scan(&b.ID, &b.Value)

		if err != nil {
			log.Println(err)
			return nil, err
		}

		balances = append(balances, &b)
	}

	err = rows.Close()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = rows.Err()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return balances, nil
}

// Update a balance
func (r *BalancePG) Update(idFrom int, idTo int, value int) error {
	_, err := r.db.Exec(
		"UPDATE balance SET balance = (case when account = $1 then balance - $2 when account = $3 then balance + $4 end)", idFrom, value, idTo, value)
	if err != nil {
		return err
	}

	return nil
}
