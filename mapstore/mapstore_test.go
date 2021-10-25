package mapstore

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapstoreCRUD(t *testing.T) {
	m := NewMapStore()

	key := "key"
	value := "value"
	newValue := "new value"

	v, err := m.get(key)
	require.Equal(t, notFoundErr, err)
	require.Zero(t, v)

	m.set(key, value)

	v, err = m.get(key)
	require.NoError(t, err)
	require.Equal(t, value, v)

	m.set(key, newValue)

	v, err = m.get(key)
	require.NoError(t, err)
	require.Equal(t, newValue, v)

	m.delete(key)

	v, err = m.get(key)
	require.Equal(t, notFoundErr, err)
	require.Zero(t, v)
}
