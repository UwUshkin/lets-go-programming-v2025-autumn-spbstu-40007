package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/UwUshkin/task-6/internal/db"
)

func TestNew(t *testing.T) {
	dbConn, _, _ := sqlmock.New()
	defer dbConn.Close()
	service := db.New(dbConn)
	require.NotNil(t, service)
	assert.Equal(t, dbConn, service.DB)
}

func TestDBService_GetNames(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Alice"))
		res, err := db.New(dbConn).GetNames()
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("empty", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}))
		res, err := db.New(dbConn).GetNames()
		assert.NoError(t, err)
		assert.Empty(t, res)
	})

	t.Run("query error", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()
		
		mock.ExpectQuery("SELECT name FROM users").WillReturnError(errors.New("q"))
		_, err := db.New(dbConn).GetNames()
		assert.Error(t, err)
	})

	t.Run("scan nil error", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))
		_, err := db.New(dbConn).GetNames()
		assert.Error(t, err)
	})

	t.Run("rows iteration error", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()
		
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Bob").
			RowError(1, errors.New("row error"))
		
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
		_, err := db.New(dbConn).GetNames()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rows error")
	})

	t.Run("rows close error", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()
		
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			CloseError(errors.New("close error"))
		
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
		
		_, err := db.New(dbConn).GetNames()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rows error")
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Alice"))
		res, err := db.New(dbConn).GetUniqueNames()
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("query error", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errors.New("q"))
		_, err := db.New(dbConn).GetUniqueNames()
		assert.Error(t, err)
	})

	t.Run("scan nil error", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))
		_, err := db.New(dbConn).GetUniqueNames()
		assert.Error(t, err)
	})

	t.Run("rows iteration error", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()
		
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Bob").
			RowError(1, errors.New("row error"))
		
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
		_, err := db.New(dbConn).GetUniqueNames()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rows error")
	})

	t.Run("rows close error", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()
		
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			CloseError(errors.New("close error"))
		
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
		
		_, err := db.New(dbConn).GetUniqueNames()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rows error")
	})
}
