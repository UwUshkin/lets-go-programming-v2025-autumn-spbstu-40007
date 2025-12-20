package db

import (
	"errors"
	"reflect"
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

	val := reflect.ValueOf(service)
	if val.Kind() == reflect.Struct {
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			if field.Type().String() == "*sql.DB" {
				assert.Equal(t, dbConn, field.Interface())
			}
		}
	}
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
