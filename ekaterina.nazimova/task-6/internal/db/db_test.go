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

	t.Run("errors", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()
		
		mock.ExpectQuery("SELECT name FROM users").WillReturnError(errors.New("q"))
		_, err := db.New(dbConn).GetNames()
		assert.Error(t, err)

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))
		_, err = db.New(dbConn).GetNames()
		assert.Error(t, err)

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("A").RowError(0, errors.New("e")))
		_, err = db.New(dbConn).GetNames()
		assert.Error(t, err)
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

	t.Run("errors", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errors.New("q"))
		_, err := db.New(dbConn).GetUniqueNames()
		assert.Error(t, err)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))
		_, err = db.New(dbConn).GetUniqueNames()
		assert.Error(t, err)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("A").RowError(0, errors.New("e")))
		_, err = db.New(dbConn).GetUniqueNames()
		assert.Error(t, err)
	})
}
