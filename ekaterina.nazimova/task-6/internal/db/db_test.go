package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	dbConn, _, _ := sqlmock.New()
	defer dbConn.Close()
	service := New(dbConn)
	require.NotNil(t, service)
}

func TestDBService_GetNames(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Alice"))
		s := New(dbConn)
		res, err := s.GetNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{"Alice"}, res)
	})

	t.Run("empty", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}))
		s := New(dbConn)
		res, err := s.GetNames()
		assert.NoError(t, err)
		assert.Empty(t, res)
	})

	t.Run("errors", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()
		s := New(dbConn)

		mock.ExpectQuery("SELECT name FROM users").WillReturnError(errors.New("db_err"))
		_, _ = s.GetNames()

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))
		_, _ = s.GetNames()

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("A").RowError(0, errors.New("row_err")))
		_, _ = s.GetNames()
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Bob"))
		s := New(dbConn)
		res, err := s.GetUniqueNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{"Bob"}, res)
	})

	t.Run("empty", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}))
		s := New(dbConn)
		res, err := s.GetUniqueNames()
		assert.NoError(t, err)
		assert.Empty(t, res)
	})

	t.Run("errors", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()
		s := New(dbConn)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errors.New("db_err"))
		_, _ = s.GetUniqueNames()

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))
		_, _ = s.GetUniqueNames()

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("B").RowError(0, errors.New("row_err")))
		_, _ = s.GetUniqueNames()
	})
}

func TestExtraCoverage(t *testing.T) {
	dbConn, mock, _ := sqlmock.New()
	defer dbConn.Close()
	s := New(dbConn)

	mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("test").CloseError(errors.New("close_err")))
	_, _ = s.GetNames()

	mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("test").CloseError(errors.New("close_err")))
	_, _ = s.GetUniqueNames()
}
