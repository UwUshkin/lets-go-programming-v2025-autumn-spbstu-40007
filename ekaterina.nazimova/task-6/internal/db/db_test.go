package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetNames(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	service := New(db)

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").AddRow("Bob")
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := service.GetNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{"Alice", "Bob"}, names)
	})

	t.Run("query_error", func(t *testing.T) {
		mock.ExpectQuery("SELECT name FROM users").WillReturnError(errors.New("db error"))
		
		_, err := service.GetNames()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db query")
	})

	t.Run("scan_error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		_, err := service.GetNames()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rows scanning")
	})

	t.Run("rows_err_after_loop", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").RowError(0, errors.New("post-iteration error"))
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		_, err := service.GetNames()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rows error")
	})
}

func TestGetUniqueNames(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	service := New(db)

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow("UniqueName")
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := service.GetUniqueNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{"UniqueName"}, names)
	})

	t.Run("query_error", func(t *testing.T) {
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errors.New("db error"))
		
		_, err := service.GetUniqueNames()
		assert.Error(t, err)
	})

	t.Run("scan_error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		_, err := service.GetUniqueNames()
		assert.Error(t, err)
	})
    
	t.Run("rows_err", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").RowError(0, errors.New("err"))
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		_, err := service.GetUniqueNames()
		assert.Error(t, err)
	})
}