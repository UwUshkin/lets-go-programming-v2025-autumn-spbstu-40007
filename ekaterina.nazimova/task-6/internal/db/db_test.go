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
	errDatabase = errors.New("database failure")
	errRowScan  = errors.New("scan failure")
)

func TestNew(t *testing.T) {
	t.Parallel()
	mockDB, _, _ := sqlmock.New()
	defer mockDB.Close()

	service := db.New(mockDB)
	require.NotNil(t, service)
	// Проверяем, что поле DB инициализировано (это закроет ветку конструктора)
	assert.NotNil(t, service.DB)
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	t.Run("success_with_data", func(t *testing.T) {
		t.Parallel()
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		service := db.New(dbConn)
		rows := sqlmock.NewRows([]string{"name"}).AddRow("User1").AddRow("User2")
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		res, err := service.GetNames()
		require.NoError(t, err)
		assert.Equal(t, []string{"User1", "User2"}, res)
	})

	t.Run("success_empty", func(t *testing.T) {
		t.Parallel()
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		service := db.New(dbConn)
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}))

		res, err := service.GetNames()
		require.NoError(t, err)
		assert.Empty(t, res)
	})

	t.Run("query_error", func(t *testing.T) {
		t.Parallel()
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		service := db.New(dbConn)
		mock.ExpectQuery("SELECT name FROM users").WillReturnError(errDatabase)

		_, err := service.GetNames()
		require.Error(t, err)
	})

	t.Run("scan_error", func(t *testing.T) {
		t.Parallel()
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		service := db.New(dbConn)
		// Передаем nil там, где ожидается строка — это гарантирует вызов блока rows.Scan error
		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		_, err := service.GetNames()
		require.Error(t, err)
	})

	t.Run("rows_iteration_error", func(t *testing.T) {
		t.Parallel()
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		service := db.New(dbConn)
		rows := sqlmock.NewRows([]string{"name"}).AddRow("User1").RowError(0, errRowScan)
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		_, err := service.GetNames()
		require.Error(t, err)
	})
}

func TestGetUniqueNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		service := db.New(dbConn)
		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice")
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		res, err := service.GetUniqueNames()
		require.NoError(t, err)
		assert.Contains(t, res, "Alice")
	})

	t.Run("query_error", func(t *testing.T) {
		t.Parallel()
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		service := db.New(dbConn)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errDatabase)

		_, err := service.GetUniqueNames()
		require.Error(t, err)
	})

	t.Run("scan_error", func(t *testing.T) {
		t.Parallel()
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		service := db.New(dbConn)
		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		_, err := service.GetUniqueNames()
		require.Error(t, err)
	})

	t.Run("rows_err_final", func(t *testing.T) {
		t.Parallel()
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		service := db.New(dbConn)
		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").RowError(0, errRowScan)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		_, err := service.GetUniqueNames()
		require.Error(t, err)
	})
}
