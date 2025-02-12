package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateClient(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Inicia transação
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO clients").
		WithArgs("Test Client", "12345678901", "PERSON", false).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
}
