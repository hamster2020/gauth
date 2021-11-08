package gauth

type Caller interface {
	Info() string
	HasRole(role Roles) bool
	HasAtLeastOneRole(roles Roles) bool
}

type EmptyCaller struct{}

func (e EmptyCaller) Info() string {
	return "empty caller"
}

func (e EmptyCaller) HasRole(role Roles) bool {
	return false
}

func (e EmptyCaller) HasAtLeastOneRole(roles Roles) bool {
	return false
}
