package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/UwUshkin/task-6/internal/db"
)

var (
	errDbQuery = errors.New("db error")
	errScan    = errors.New("scan error")
	errIter    = errors.New("post-iteration error")
)

func TestGetNames(t *testing.T) {
	t.Parallel()
	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	service := db.New(dbConn)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").AddRow("Bob")
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := service.GetNames()
		require.NoError(t, err)
		assert.Equal(t, []string{"Alice", "Bob"}, names)
	})

	t.Run("query_error", func(t *testing.T) {
		t.Parallel()
		mock.ExpectQuery("SELECT name FROM users").WillReturnError(errDbQuery)

		_, err := service.GetNames()
		require.Error(t, err)
	})

	t.Run("scan_error", func(t *testing.T) {
		t.Parallel()
		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		_, err := service.GetNames()
		require.Error(t, err)
	})

	t.Run("rows_err_after_loop", func(t *testing.T) {
		t.Parallel()
		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").RowError(0, errIter)
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		_, err := service.GetNames()
		require.Error(t, err)
	})
}

func TestGetUniqueNames(t *testing.T) {
	t.Parallel()
	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	service := db.New(dbConn)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		rows := sqlmock.NewRows([]string{"name"}).AddRow("UniqueName")
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := service.GetUniqueNames()
		require.NoError(t, err)
		assert.Equal(t, []string{"UniqueName"}, names)
	})

	t.Run("query_error", func(t *testing.T) {
		t.Parallel()
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errDbQuery)

		_, err := service.GetUniqueNames()
		require.Error(t, err)
	})

	t.Run("scan_error", func(t *testing.T) {
		t.Parallel()
		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		_, err := service.GetUniqueNames()
		require.Error(t, err)
	})

	t.Run("rows_err", func(t *testing.T) {
		t.Parallel()
		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").RowError(0, errIter)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		_, err := service.GetUniqueNames()
		require.Error(t, err)
	})
}
