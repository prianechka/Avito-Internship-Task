package transaction_repo

import (
	"Avito-Internship-Task/internal/app/balance_service_app/transaction"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"reflect"
	"testing"
	"time"
)

func TestAddNewTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var newTransaction transaction.Transaction

	mock.
		ExpectExec("INSERT INTO balanceApp.transactions").
		WithArgs(newTransaction.TransactionID, newTransaction.UserID, newTransaction.TransactionType,
			newTransaction.Sum, newTransaction.Time, newTransaction.ActionComments, newTransaction.AddComments).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewTransactionRepo(db)

	execErr := repo.AddNewTransaction(newTransaction)
	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}
}

func TestGetAllTransactions(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	curTime := time.Now()

	rows := sqlmock.NewRows([]string{"transactionID", "userID", "transactionType", "sum",
		"time", "actionComment", "addComment"})
	expect := []transaction.Transaction{{1, 1, 1, 100, curTime, "Good", "Good"},
		{2, 2, 2, 200, curTime, "Bad", "Bad"}}
	for _, transact := range expect {
		rows = rows.AddRow(transact.TransactionID, transact.UserID, transact.TransactionType,
			transact.Sum, transact.Time, transact.ActionComments, transact.AddComments)
	}

	mock.
		ExpectQuery(MySQLGetAllTransactions{}.GetString()).
		WillReturnRows(rows).WillReturnError(nil)

	repo := NewTransactionRepo(db)

	transact, execErr := repo.GetAllTransactions()
	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}

	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	if !reflect.DeepEqual(transact, expect) {
		t.Errorf("results not match, want %v, have %v", expect, transact)
		return
	}
}

func TestGetUserTransactions(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	curTime := time.Now()
	var userID int64 = 1

	rows := sqlmock.NewRows([]string{"transactionID", "userID", "transactionType", "sum",
		"time", "actionComment", "addComment"})
	expect := []transaction.Transaction{{1, userID, 1, 100, curTime, "Good", "Good"},
		{2, userID, 2, 200, curTime, "Bad", "Bad"}}
	for _, transact := range expect {
		rows = rows.AddRow(transact.TransactionID, transact.UserID, transact.TransactionType,
			transact.Sum, transact.Time, transact.ActionComments, transact.AddComments)
	}

	mock.
		ExpectQuery("SELECT transactionID, userID, transactionType, sum, time," +
			" actionComments, addComments FROM balanceApp.transaction WHERE userID").
		WillReturnRows(rows).WillReturnError(nil)

	repo := NewTransactionRepo(db)

	transact, execErr := repo.GetUserTransactions(userID)
	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}

	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	if !reflect.DeepEqual(transact, expect) {
		t.Errorf("results not match, want %v, have %v", expect, transact)
		return
	}
}

func TestGetTransactionByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	curTime := time.Now()
	var orderID int64 = 1

	rows := sqlmock.NewRows([]string{"transactionID", "userID", "transactionType", "sum",
		"time", "actionComment", "addComment"})
	expect := transaction.Transaction{orderID, 1, 1, 200, curTime, "", ""}

	rows.AddRow(expect.TransactionID, expect.UserID, expect.TransactionType, expect.Sum, expect.Time,
		expect.ActionComments, expect.AddComments)

	mock.
		ExpectQuery("SELECT transactionID, userID, transactionType, sum, time," +
			" actionComments, addComments FROM balanceApp.transaction WHERE transactionID").
		WillReturnRows(rows).WillReturnError(nil)

	repo := NewTransactionRepo(db)

	transact, execErr := repo.GetTransactionByID(orderID)
	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}

	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	if !reflect.DeepEqual(transact, expect) {
		t.Errorf("results not match, want %v, have %v", expect, transact)
		return
	}
}
