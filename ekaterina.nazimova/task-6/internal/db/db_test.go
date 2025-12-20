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
	errTestDB    = errors.New("db error")
	errTestRow   = errors.New("row error")
	errTestClose = errors.New("close error")
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

	t.Run("success", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Alice"))
		
		res, err := db.New(dbConn).GetNames()
		require.NoError(t, err)
		assert.NotEmpty(t, res)
	})

	t.Run("all_errors", func(t *testing.T) {
		db1, m1, _ := sqlmock.New()
		defer db1.Close()
		m1.ExpectQuery("SELECT name FROM users").WillReturnError(errTestDB)
		_, err := db.New(db1).GetNames()
		require.Error(t, err)

		db2, m2, _ := sqlmock.New()
		defer db2.Close()
		m2.ExpectQuery("SELECT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))
		_, err = db.New(db2).GetNames()
		require.Error(t, err)

		db3, m3, _ := sqlmock.New()
		defer db3.Close()
		m3.ExpectQuery("SELECT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Alice").RowError(0, errTestRow))
		_, err = db.New(db3).GetNames()
		require.Error(t, err)

		db4, m4, _ := sqlmock.New()
		defer db4.Close()
		m4.ExpectQuery("SELECT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Alice").CloseError(errTestClose))
		_, _ = db.New(db4).GetNames()
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Alice"))
		
		res, err := db.New(dbConn).GetUniqueNames()
		require.NoError(t, err)
		assert.NotEmpty(t, res)
	})

	t.Run("all_errors", func(t *testing.T) {
		db1, m1, _ := sqlmock.New()
		defer db1.Close()
		m1.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errTestDB)
		_, err := db.New(db1).GetUniqueNames()
		require.Error(t, err)

		db2, m2, _ := sqlmock.New()
		defer db2.Close()
		m2.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))
		_, err = db.New(db2).GetUniqueNames()
		require.Error(t, err)

		db3, m3, _ := sqlmock.New()
		defer db3.Close()
		m3.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Alice").RowError(0, errTestRow).CloseError(errTestClose))
		_, err = db.New(db3).GetUniqueNames()
		require.Error(t, err)
	})
}
