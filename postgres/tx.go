package postgres

import (
	"github.com/hamster2020/gauth"
	"github.com/jmoiron/sqlx"
)

type Tx struct {
	*sqlx.Tx
}

func (tx Tx) InsideTx(fn func(gauth.Transaction) error) error {
	return fn(tx)
}

func (tx Tx) IsTx() {}
