package gauth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRoles(t *testing.T) {
	emptyRoles := Roles(0)
	r := &emptyRoles

	for _, role := range rolesAndNames {
		require.False(t, r.HasRole(role))
	}

	for _, role := range rolesAndNames {
		r.AddRole(role)
		require.True(t, r.HasRole(role))
		r.DropRole(role)
		require.False(t, r.HasRole(role))
	}
}
