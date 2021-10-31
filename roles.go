package gauth

import "fmt"

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

func RoleFromName(n string) (Roles, error) {
	roles, found := rolesAndNames[RoleName(n)]
	if !found {
		return Roles(0), fmt.Errorf("role with name %s not found", n)
	}

	return roles, nil
}

func RolesFromNames(ns []string) (Roles, error) {
	roles := Roles(0)
	for _, n := range ns {
		role, err := RoleFromName(n)
		if err != nil {
			return Roles(0), err
		}

		roles.AddRole(role)
	}

	return roles, nil
}
