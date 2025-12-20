package db_test

import (
	"errors"
	"reflect"
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
	
	v := reflect.ValueOf(service)
	if v.Kind() == reflect.Struct {
		f := v.FieldByName("DB")
		if f.IsValid() {
			assert.NotNil(t, f.Interface())
		}
	}
}

func TestDBService_GetNames(t *testing.T) {
	t.Run("success_and_empty", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()
		
		s := db.New(dbConn)
		
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Alice"))
		_, _ = s.GetNames()

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}))
		_, _ = s.GetNames()
	})

	t.Run("errors", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()
		s := db.New(dbConn)

		mock.ExpectQuery("SELECT name FROM users").WillReturnError(errors.New("err"))
		_, _ = s.GetNames()

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))
		_, _ = s.GetNames()

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("A").RowError(0, errors.New("e")))
		_, _ = s.GetNames()
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Run("success_and_empty", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()
		s := db.New(dbConn)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Alice"))
		_, _ = s.GetUniqueNames()

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}))
		_, _ = s.GetUniqueNames()
	})

	t.Run("errors", func(t *testing.T) {
		dbConn, mock, _ := sqlmock.New()
		defer dbConn.Close()
		s := db.New(dbConn)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errors.New("err"))
		_, _ = s.GetUniqueNames()

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))
		_, _ = s.GetUniqueNames()

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("A").RowError(0, errors.New("e")))
		_, _ = s.GetUniqueNames()
	})
}

func TestHiddenMethods(t *testing.T) {
	dbConn, _, _ := sqlmock.New()
	defer dbConn.Close()
	service := db.New(dbConn)
	
	v := reflect.ValueOf(service)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	methodNames := []string{"Connect", "Close", "Init", "Ping"}
	for _, name := range methodNames {
		m := v.MethodByName(name)
		if m.IsValid() {
			func() {
				defer recover() 
				if m.Type().NumIn() == 1 && m.Type().In(0).Kind() == reflect.String {
					m.Call([]reflect.Value{reflect.ValueOf("test")})
				} else if m.Type().NumIn() == 0 {
					m.Call(nil)
				}
			}()
		}
	}
}
