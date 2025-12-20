package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/UwUshkin/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

var (
	errDb            = errors.New("db error")
	errPostIteration = errors.New("post-iteration error")
)

func TestGetNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		dbConn, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbConn.Close()

		service := db.New(dbConn)
		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").AddRow("Bob")
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := service.GetNames()
		require.NoError(t, err)
		require.Equal(t, []string{"Alice", "Bob"}, names)
	})

	t.Run("query_error", func(t *testing.T) {
		t.Parallel()
		dbConn, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbConn.Close()

		service := db.New(dbConn)
		mock.ExpectQuery("SELECT name FROM users").WillReturnError(errDb)

		_, err = service.GetNames()
		require.Error(t, err)
	})
}

func TestGetUniqueNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		dbConn, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbConn.Close()

		service := db.New(dbConn)
		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice")
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := service.GetUniqueNames()
		require.NoError(t, err)
		require.Contains(t, names, "Alice")
	})

	t.Run("scan_error", func(t *testing.T) {
		t.Parallel()
		dbConn, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbConn.Close()

		service := db.New(dbConn)
		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").RowError(0, errPostIteration)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		_, err = service.GetUniqueNames()
		require.Error(t, err)
	})
}
