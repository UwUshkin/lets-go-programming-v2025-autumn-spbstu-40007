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
	errDBTest  = errors.New("database query error")
	errRowTest = errors.New("row iteration error")
)

func TestNew(t *testing.T) {
	t.Parallel()
	mockDB, _, _ := sqlmock.New()
	defer mockDB.Close()

	service := db.New(mockDB)
	require.NotNil(t, service)
	assert.Equal(t, mockDB, service.DB)
}

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").AddRow("Bob")
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := db.New(dbConn).GetNames()
		require.NoError(t, err)
		assert.Equal(t, []string{"Alice", "Bob"}, names)
	})

	t.Run("query error", func(t *testing.T) {
		t.Parallel()
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()
		mock.ExpectQuery("SELECT name FROM users").WillReturnError(errDBTest)

		res, err := db.New(dbConn).GetNames()
		require.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("scan error", func(t *testing.T) {
		t.Parallel()
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))

		res, err := db.New(dbConn).GetNames()
		require.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("rows error", func(t *testing.T) {
		t.Parallel()
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Alice").RowError(0, errRowTest))

		res, err := db.New(dbConn).GetNames()
		require.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Alice"))

		names, err := db.New(dbConn).GetUniqueNames()
		require.NoError(t, err)
		assert.Equal(t, []string{"Alice"}, names)
	})

	t.Run("query error", func(t *testing.T) {
		t.Parallel()
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errDBTest)

		res, err := db.New(dbConn).GetUniqueNames()
		require.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("scan error", func(t *testing.T) {
		t.Parallel()
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))

		res, err := db.New(dbConn).GetUniqueNames()
		require.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("rows error", func(t *testing.T) {
		t.Parallel()
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Alice").RowError(0, errRowTest))

		res, err := db.New(dbConn).GetUniqueNames()
		require.Error(t, err)
		assert.Nil(t, res)
	})
}
