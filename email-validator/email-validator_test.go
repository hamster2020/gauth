package emailvalidator

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	token := "2b1e810090b21cab8a8753ec6bd1f091b4b01255302e5c5314263898f3fd3d04d49f54b43333a30279114b37a53b6f40"
	eValidator := NewEmailValidator(token)

	cases := []struct {
		name  string
		email string

		expErr error
	}{
		{
			name:   "invalid email",
			email:  "bogus@email.com",
			expErr: errors.New("The provided email 'bogus@email.com' is invalid"),
		},
		{
			name:  "valid email",
			email: "test@gmail.com",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := eValidator.Validate(tc.email)
			require.Equal(t, tc.expErr, err)
		})
	}
}
