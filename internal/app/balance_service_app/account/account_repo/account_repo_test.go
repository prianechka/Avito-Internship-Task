package account_repo

import (
	"reflect"
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

// go test -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html

func TestAddNewAccount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var newUserID int64 = 1

	mock.
		ExpectExec("INSERT INTO").
		WithArgs(newUserID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewAccountRepo(db)

	execErr := repo.AddNewAccount(newUserID)
	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}
}

func TestGetAccountAmount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var checkAccID int64 = 1

	rows := sqlmock.NewRows([]string{"amount"})
	expect := 293.46
	rows = rows.AddRow(expect)

	mock.
		ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE").
		WithArgs(checkAccID).
		WillReturnRows(rows).WillReturnError(nil)

	repo := NewAccountRepo(db)

	amount, execErr := repo.GetCurrentAmount(checkAccID)
	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}

	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	if !reflect.DeepEqual(amount, expect) {
		t.Errorf("results not match, want %v, have %v", expect, amount)
		return
	}
}

func TestChangeAmount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var changeAccID int64 = 1
	var delta float64 = 345

	mock.
		ExpectExec("UPDATE balanceApp.accounts SET amount = amoumt +").
		WithArgs(delta, changeAccID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rows := sqlmock.NewRows([]string{"amount"})
	expect := delta
	rows = rows.AddRow(expect)

	repo := NewAccountRepo(db)

	execErr := repo.ChangeAmount(changeAccID, delta)
	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}

	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	mock.
		ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE").
		WithArgs(changeAccID).
		WillReturnRows(rows).WillReturnError(nil)

	amount, getExecErr := repo.GetCurrentAmount(changeAccID)
	if getExecErr != nil {
		t.Errorf("unexpected err: %v", getExecErr)
		return
	}

	if !reflect.DeepEqual(amount, expect) {
		t.Errorf("results not match, want %v, have %v", expect, amount)
		return
	}
}
