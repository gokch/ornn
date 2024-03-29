package db_mysql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConn(t *testing.T) {
	conn, err := New("localhost", "3306", "root", "1234", "test")
	require.NoError(t, err)

	err = conn.Raw().Ping()
	require.NoError(t, err)
}
