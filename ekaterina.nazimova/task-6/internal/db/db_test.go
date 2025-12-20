package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/UwUshkin/task-6/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	errDBCustom  = errors.New("db query error")
	errRowCustom = errors.New("row error")
)

func TestNew(t *testing.T) {
	t.Parallel()
	mockDB, _, _ := sqlmock.New()
	defer mockDB.Close()

	service := db.New(mockDB)
	require.NotNil(t, service)
	assert.NotNil(t, service.DB)
}

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	setup := func() (*db.DBService, sqlmock.Sqlmock, func()) {
		mockDB, mock, _ := sqlmock.New()
		service := db.New(mockDB)
		return &service, mock, func() { mockDB.Close() }
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		service, mock, cleanup := setup()
		defer cleanup()

		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice")
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := service.GetNames()
		require.NoError(t, err)
		assert.Equal(t, []string{"Alice"}, names)
	})

	t.Run("empty_result", func(t *testing.T) {
		t.Parallel()
		service, mock, cleanup := setup()
		defer cleanup()

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}))

		names, err := service.GetNames()
		require.NoError(t, err)
		assert.Empty(t, names)
	})

	t.Run("query_error", func(t *testing.T) {
		t.Parallel()
		service, mock, cleanup := setup()
		defer cleanup()

		mock.ExpectQuery("SELECT name FROM users").WillReturnError(errDBCustom)

		_, err := service.GetNames()
		require.Error(t, err)
	})

	t.Run("scan_error", func(t *testing.T) {
		t.Parallel()
		service, mock, cleanup := setup()
		defer cleanup()

		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		_, err := service.GetNames()
		require.Error(t, err)
	})

	t.Run("rows_error", func(t *testing.T) {
		t.Parallel()
		service, mock, cleanup := setup()
		defer cleanup()

		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").RowError(0, errRowCustom)
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		_, err := service.GetNames()
		require.Error(t, err)
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		mockDB, mock, _ := sqlmock.New()
		defer mockDB.Close()
		service := db.New(mockDB)

		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice")
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := service.GetUniqueNames()
		require.NoError(t, err)
		assert.Equal(t, []string{"Alice"}, names)
	})

	t.Run("query_error", func(t *testing.T) {
		t.Parallel()
		mockDB, mock, _ := sqlmock.New()
		defer mockDB.Close()
		service := db.New(mockDB)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errDBCustom)
		_, err := service.GetUniqueNames()
		require.Error(t, err)
	})

	t.Run("scan_error", func(t *testing.T) {
		t.Parallel()
		mockDB, mock, _ := sqlmock.New()
		defer mockDB.Close()
		service := db.New(mockDB)

		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
		_, err := service.GetUniqueNames()
		require.Error(t, err)
	})

	t.Run("rows_error_final", func(t *testing.T) {
		t.Parallel()
		mockDB, mock, _ := sqlmock.New()
		defer mockDB.Close()
		service := db.New(mockDB)

		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").RowError(0, errRowCustom)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
		_, err := db.DBService(service).GetUniqueNames()
		require.Error(t, err)
	})
}
