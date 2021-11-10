package gauth

type Datastore interface {
	SessionStore
	UserStore
}
