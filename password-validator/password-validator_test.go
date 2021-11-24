package passwordvalidator

import (
	"fmt"
	"testing"

	"github.com/hamster2020/gauth"
	"github.com/stretchr/testify/require"
)

func TestIsPasswordBreached(t *testing.T) {
	cfg, err := gauth.NewConfig()
	require.NoError(t, err)
	pwdValidator := NewPasswordValidator(cfg.PwnedPasswordsURL)
	fmt.Println(cfg.PwnedPasswordsURL)

	cases := []struct {
		name     string
		password string

		expErr error
	}{
		{
			name:     "invalid password",
			password: "P@ssword",
			expErr:   breachedPasswordErr,
		},
		{
			name:     "valid password",
			password: "Af^hl45G6&",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := pwdValidator.isPasswordBreached(tc.password)
			require.Equal(t, tc.expErr, err)
		})
	}
}
