package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()
	dbConn, _, _ := sqlmock.New()
	defer dbConn.Close()
	service := New(dbConn)
	assert.NotNil(t, service)
	assert.Equal(t, dbConn, service.DB)
}

func TestDBService_GetNames(t *testing.T) {
	dbConn, mock, _ := sqlmock.New()
	defer dbConn.Close()
	s := New(dbConn)

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Alice"))
	_, _ = s.GetNames()

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}))
	_, _ = s.GetNames()

	mock.ExpectQuery("SELECT name FROM users").WillReturnError(errors.New("e"))
	_, _ = s.GetNames()

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))
	_, _ = s.GetNames()

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("A").RowError(0, errors.New("e")))
	_, _ = s.GetNames()
}

func TestDBService_GetUniqueNames(t *testing.T) {
	dbConn, mock, _ := sqlmock.New()
	defer dbConn.Close()
	s := New(dbConn)

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Bob"))
	_, _ = s.GetUniqueNames()

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}))
	_, _ = s.GetUniqueNames()

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errors.New("e"))
	_, _ = s.GetUniqueNames()

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))
	_, _ = s.GetUniqueNames()

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("B").RowError(0, errors.New("e")))
	_, _ = s.GetUniqueNames()
}
