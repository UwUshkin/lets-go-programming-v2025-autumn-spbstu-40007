package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/UwUshkin/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

var errTest = errors.New("test error")

func TestDBService(t *testing.T) {
	t.Parallel()

	t.Run("GetNames", func(t *testing.T) {
		t.Parallel()

		t.Run("success", func(t *testing.T) {
			dbConn, mock, _ := sqlmock.New()
			defer dbConn.Close()
			mock.ExpectQuery("SELECT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Alice"))
			_, err := db.New(dbConn).GetNames()
			require.NoError(t, err)
		})

		t.Run("empty_result", func(t *testing.T) { // Добавляем пустой результат
			dbConn, mock, _ := sqlmock.New()
			defer dbConn.Close()
			mock.ExpectQuery("SELECT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}))
			res, err := db.New(dbConn).GetNames()
			require.NoError(t, err)
			require.Empty(t, res)
		})

		t.Run("query_error", func(t *testing.T) {
			dbConn, mock, _ := sqlmock.New()
			defer dbConn.Close()
			mock.ExpectQuery("SELECT name FROM users").WillReturnError(errTest)
			_, err := db.New(dbConn).GetNames()
			require.Error(t, err)
		})

		t.Run("scan_error", func(t *testing.T) {
			dbConn, mock, _ := sqlmock.New()
			defer dbConn.Close()
			// Несоответствие колонок вызывает ошибку Scan
			rows := sqlmock.NewRows([]string{"name", "id"}).AddRow("Alice", 1)
			mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			_, err := db.New(dbConn).GetNames()
			require.Error(t, err)
		})

		t.Run("rows_err", func(t *testing.T) {
			dbConn, mock, _ := sqlmock.New()
			defer dbConn.Close()
			rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").RowError(0, errTest)
			mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			_, err := db.New(dbConn).GetNames()
			require.Error(t, err)
		})
	})

	t.Run("GetUniqueNames", func(t *testing.T) {
		t.Parallel()

		t.Run("success", func(t *testing.T) {
			dbConn, mock, _ := sqlmock.New()
			defer dbConn.Close()
			mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Alice"))
			_, err := db.New(dbConn).GetUniqueNames()
			require.NoError(t, err)
		})

		t.Run("query_error", func(t *testing.T) {
			dbConn, mock, _ := sqlmock.New()
			defer dbConn.Close()
			mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errTest)
			_, err := db.New(dbConn).GetUniqueNames()
			require.Error(t, err)
		})

		t.Run("scan_error", func(t *testing.T) {
			dbConn, mock, _ := sqlmock.New()
			defer dbConn.Close()
			rows := sqlmock.NewRows([]string{"name", "id"}).AddRow("Alice", 1)
			mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			_, err := db.New(dbConn).GetUniqueNames()
			require.Error(t, err)
		})

		t.Run("rows_err", func(t *testing.T) {
			dbConn, mock, _ := sqlmock.New()
			defer dbConn.Close()
			rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").RowError(0, errTest)
			mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			_, err := db.New(dbConn).GetUniqueNames()
			require.Error(t, err)
		})
	})
}

func TestNew(t *testing.T) {
	t.Parallel()
	dbConn, _, _ := sqlmock.New()
	defer dbConn.Close()
	service := db.New(dbConn)
	require.NotNil(t, service)
}
