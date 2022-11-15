package transaction_controller

import (
	"Avito-Internship-Task/internal/app/balance_service_app/order"
	"Avito-Internship-Task/internal/app/balance_service_app/transaction"
	"Avito-Internship-Task/internal/app/balance_service_app/transaction/transaction_repo"
	"fmt"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"reflect"
	"testing"
	"time"
)

// TestAddNewRecordBuyService проверяет, что запись была добавлена корректно
func TestAddNewRecordBuyService(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var (
		userID    int     = 1
		sum       float64 = 100
		serviceID int     = 1
		comment           = "Все успешно!"
	)

	newTransaction := transaction.Transaction{
		TransactionID:   0,
		UserID:          userID,
		TransactionType: transaction.Buy,
		Sum:             sum,
		Time:            time.Now(),
		ActionComments:  "куплена услуга: " + order.Types[serviceID],
		AddComments:     comment,
	}

	mock.
		ExpectExec("INSERT INTO balanceApp.transactions").
		WithArgs(newTransaction.TransactionID, newTransaction.UserID, newTransaction.TransactionType,
			newTransaction.Sum, sqlmock.AnyArg(), newTransaction.ActionComments, newTransaction.AddComments).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := transaction_repo.NewTransactionRepo(db)
	controller := CreateNewTransactionController(repo)

	execErr := controller.AddNewRecordBuyService(userID, sum, serviceID, comment)
	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	rows := sqlmock.NewRows([]string{"transactionID", "userID", "transactionType", "sum",
		"time", "actionComment", "addComment"})

	rows.AddRow(newTransaction.TransactionID, newTransaction.UserID, newTransaction.TransactionType,
		newTransaction.Sum, newTransaction.Time, newTransaction.ActionComments, newTransaction.AddComments)

	mock.
		ExpectQuery("SELECT transactionID, userID, transactionType, sum, time," +
			" actionComments, addComments FROM balanceApp.transactions WHERE transactionID").
		WillReturnRows(rows).WillReturnError(nil)

	getTransact, getError := controller.GetTransactionByID(newTransaction.TransactionID)

	if getError != nil {
		t.Errorf("unexpected err: %v", getError)
		return
	}
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	getTransact.Time = newTransaction.Time

	if !reflect.DeepEqual(getTransact, newTransaction) {
		t.Errorf("results not match, want %v, have %v", getTransact, newTransaction)
		return
	}
}

// TestAddNewRecordReturnService проверяет, что запись была добавлена корректно
func TestAddNewRecordReturnService(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var (
		userID    int     = 1
		sum       float64 = 100
		serviceID int     = 1
		comment           = "Все успешно!"
	)

	newTransaction := transaction.Transaction{
		TransactionID:   0,
		UserID:          userID,
		TransactionType: transaction.Return,
		Sum:             sum,
		Time:            time.Now(),
		ActionComments:  "возврат за услугу: " + order.Types[serviceID],
		AddComments:     comment,
	}

	mock.
		ExpectExec("INSERT INTO balanceApp.transactions").
		WithArgs(newTransaction.TransactionID, newTransaction.UserID, newTransaction.TransactionType,
			newTransaction.Sum, sqlmock.AnyArg(), newTransaction.ActionComments, newTransaction.AddComments).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := transaction_repo.NewTransactionRepo(db)
	controller := CreateNewTransactionController(repo)

	execErr := controller.AddNewRecordReturnService(userID, sum, serviceID, comment)
	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	rows := sqlmock.NewRows([]string{"transactionID", "userID", "transactionType", "sum",
		"time", "actionComment", "addComment"})

	rows.AddRow(newTransaction.TransactionID, newTransaction.UserID, newTransaction.TransactionType,
		newTransaction.Sum, newTransaction.Time, newTransaction.ActionComments, newTransaction.AddComments)

	mock.
		ExpectQuery("SELECT transactionID, userID, transactionType, sum, time," +
			" actionComments, addComments FROM balanceApp.transactions WHERE transactionID").
		WillReturnRows(rows).WillReturnError(nil)

	getTransact, getError := controller.GetTransactionByID(newTransaction.TransactionID)

	if getError != nil {
		t.Errorf("unexpected err: %v", getError)
		return
	}
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	newTransaction.Time = getTransact.Time

	if !reflect.DeepEqual(getTransact, newTransaction) {
		t.Errorf("results not match, want %v, have %v", getTransact, newTransaction)
		return
	}
}

// TestAddNewRecordRefillBalance проверяет, что запись была добавлена корректно
func TestAddNewRecordRefillBalance(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var (
		userID  int     = 1
		sum     float64 = 100
		comment         = "Все успешно!"
	)

	newTransaction := transaction.Transaction{
		TransactionID:   0,
		UserID:          userID,
		TransactionType: transaction.Refill,
		Sum:             sum,
		Time:            time.Now(),
		ActionComments:  "зачислены средства на баланс",
		AddComments:     comment,
	}

	mock.
		ExpectExec("INSERT INTO balanceApp.transactions").
		WithArgs(newTransaction.TransactionID, newTransaction.UserID, newTransaction.TransactionType,
			newTransaction.Sum, sqlmock.AnyArg(), newTransaction.ActionComments, newTransaction.AddComments).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := transaction_repo.NewTransactionRepo(db)
	controller := CreateNewTransactionController(repo)

	execErr := controller.AddNewRecordRefillBalance(userID, sum, comment)
	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	rows := sqlmock.NewRows([]string{"transactionID", "userID", "transactionType", "sum",
		"time", "actionComment", "addComment"})

	rows.AddRow(newTransaction.TransactionID, newTransaction.UserID, newTransaction.TransactionType,
		newTransaction.Sum, newTransaction.Time, newTransaction.ActionComments, newTransaction.AddComments)

	mock.
		ExpectQuery("SELECT transactionID, userID, transactionType, sum, time," +
			" actionComments, addComments FROM balanceApp.transactions WHERE transactionID").
		WillReturnRows(rows).WillReturnError(nil)

	getTransact, getError := controller.GetTransactionByID(newTransaction.TransactionID)

	if getError != nil {
		t.Errorf("unexpected err: %v", getError)
		return
	}
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	newTransaction.Time = getTransact.Time

	if !reflect.DeepEqual(getTransact, newTransaction) {
		t.Errorf("results not match, want %v, have %v", getTransact, newTransaction)
		return
	}
}

// TestAddNewRecordTransferTo проверяет, что запись была добавлена корректно
func TestAddNewRecordTransferTo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var (
		srcUserID int     = 1
		dstUserID int     = 2
		sum       float64 = 100
		comment           = "Все успешно!"
	)

	newTransaction := transaction.Transaction{
		TransactionID:   0,
		UserID:          srcUserID,
		TransactionType: transaction.Transfer,
		Sum:             sum,
		Time:            time.Now(),
		ActionComments:  "перевод пользователю: " + fmt.Sprintf("%d", dstUserID),
		AddComments:     comment,
	}

	mock.
		ExpectExec("INSERT INTO balanceApp.transactions").
		WithArgs(newTransaction.TransactionID, newTransaction.UserID, newTransaction.TransactionType,
			newTransaction.Sum, sqlmock.AnyArg(), newTransaction.ActionComments, newTransaction.AddComments).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := transaction_repo.NewTransactionRepo(db)
	controller := CreateNewTransactionController(repo)

	execErr := controller.AddNewRecordTransferTo(srcUserID, dstUserID, sum, comment)
	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	rows := sqlmock.NewRows([]string{"transactionID", "userID", "transactionType", "sum",
		"time", "actionComment", "addComment"})

	rows.AddRow(newTransaction.TransactionID, newTransaction.UserID, newTransaction.TransactionType,
		newTransaction.Sum, newTransaction.Time, newTransaction.ActionComments, newTransaction.AddComments)

	mock.
		ExpectQuery("SELECT transactionID, userID, transactionType, sum, time," +
			" actionComments, addComments FROM balanceApp.transactions WHERE transactionID").
		WillReturnRows(rows).WillReturnError(nil)

	getTransact, getError := controller.GetTransactionByID(newTransaction.TransactionID)

	if getError != nil {
		t.Errorf("unexpected err: %v", getError)
		return
	}
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	newTransaction.Time = getTransact.Time

	if !reflect.DeepEqual(getTransact, newTransaction) {
		t.Errorf("results not match, want %v, have %v", getTransact, newTransaction)
		return
	}
}

// TestGetUserTransactions проверяет, что запись была добавлена корректно
func TestGetUserTransactions(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var (
		userID    int     = 1
		serviceID int     = 1
		sum       float64 = 100
		comment           = "Хорошо"
		orderBy           = "id"
		limit             = 2
		offset            = 0
	)

	newTransactions := []transaction.Transaction{{
		TransactionID:   0,
		UserID:          userID,
		TransactionType: transaction.Refill,
		Sum:             sum,
		Time:            time.Now(),
		ActionComments:  "зачислены средства на баланс",
		AddComments:     comment,
	}, {TransactionID: 1,
		UserID:          userID,
		TransactionType: transaction.Buy,
		Sum:             sum,
		Time:            time.Now(),
		ActionComments:  "куплена услуга: " + order.Types[serviceID],
		AddComments:     comment},
	}

	rows := sqlmock.NewRows([]string{"transactionID", "userID", "transactionType", "sum",
		"time", "actionComment", "addComment"})

	for _, newTransaction := range newTransactions {
		rows.AddRow(newTransaction.TransactionID, newTransaction.UserID, newTransaction.TransactionType,
			newTransaction.Sum, newTransaction.Time, newTransaction.ActionComments, newTransaction.AddComments)
	}

	mock.ExpectQuery("SELECT transactionID, userID, transactionType, sum, time," +
		" actionComments, addComments FROM balanceApp.transactions WHERE userID").
		WillReturnRows(rows).WillReturnError(nil)

	repo := transaction_repo.NewTransactionRepo(db)
	controller := CreateNewTransactionController(repo)

	userTransacts, execErr := controller.GetUserTransactions(userID, orderBy, limit, offset)
	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	for i := range userTransacts {
		userTransacts[i].Time = newTransactions[i].Time
	}

	if !reflect.DeepEqual(userTransacts, newTransactions) {
		t.Errorf("results not match, want %v, have %v", newTransactions, userTransacts)
		return
	}
}
