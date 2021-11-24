package postgres

import (
	"testing"

	"github.com/hamster2020/gauth"
	"github.com/stretchr/testify/require"
)

func newTestDB(t *testing.T) (DB, func()) {
	cfg, err := gauth.NewConfig()
	require.NoError(t, err)
	require.NotZero(t, cfg)
	require.NotZero(t, cfg.DBURL)

	testSchema, err := gauth.RandomHex(8)
	require.NoError(t, err)

	db, err := NewDB(testSchema, cfg.DBURL)
	require.NoError(t, err)
	require.NotZero(t, db)

	return db, func() {
		require.NoError(t, db.DropSchema())
		require.NoError(t, db.Close())
	}
}

func TestNewDB(t *testing.T) {
	db, cleanup := newTestDB(t)
	defer cleanup()

	require.NotZero(t, db)
	require.NoError(t, db.Ping())
}
