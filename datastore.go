package gauth

type Datastore interface {
	SessionStore
	UserStore
	InsideTx(fn func(tx Transaction) error) error
}

type Transaction interface {
	Datastore
	IsTx()
}
