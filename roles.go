package gauth

type Roles uint8

const (
	RolesAdmin Roles = 1 << iota
	RolesBase
)

func (r *Roles) HasRole(role Roles) bool {
	return *r&role != 0
}

func (r *Roles) AddRole(role Roles) {
	*r |= role
}

func (r *Roles) DropRole(role Roles) {
	*r &^= role
}

type RoleName string

const (
	RoleNameAdmin = "admin"
	RoleNameBase  = "base"
)

var rolesAndNames = map[RoleName]Roles{
	RoleNameAdmin: RolesAdmin,
	RoleNameBase:  RolesBase,
}
