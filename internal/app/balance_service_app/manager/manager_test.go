package manager

import (
	ac "Avito-Internship-Task/internal/app/balance_service_app/account/account_controller"
	"Avito-Internship-Task/internal/app/balance_service_app/account/account_repo"
	oc "Avito-Internship-Task/internal/app/balance_service_app/order/order_controller"
	"Avito-Internship-Task/internal/app/balance_service_app/order/order_repo"
	"Avito-Internship-Task/internal/app/balance_service_app/transaction"
	tc "Avito-Internship-Task/internal/app/balance_service_app/transaction/transaction_controller"
	"Avito-Internship-Task/internal/app/balance_service_app/transaction/transaction_repo"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"reflect"
	"testing"
	"time"
)

// TestRefillMoney проверяет, что сценарий пополнения денег отрабатывает без ошибок
func TestRefillMoneyWithExistsAccount(t *testing.T) {
	var (
		userID  int64   = 1
		sum     float64 = 200
		comment         = "Всё хорошо!"
	)

	// Подготовка БД для таблицы с аккаунтами
	accountDB, accountMock, createAccountDBErr := sqlmock.New()
	if createAccountDBErr != nil {
		t.Fatalf("cant create mock: %s", createAccountDBErr)
	}
	defer accountDB.Close()

	accountFirstRows := sqlmock.NewRows([]string{"amount"})
	var expectResult float64 = 209
	accountFirstRows.AddRow(expectResult)

	accountSecondRows := sqlmock.NewRows([]string{"amount"})
	accountSecondRows.AddRow(expectResult)

	accountMock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(userID).
		WillReturnRows(accountFirstRows)

	/*
		accountMock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
			WithArgs(userID).
			WillReturnRows(accountSecondRows)

	*/
	accountMock.ExpectExec("UPDATE balanceApp.accounts SET amount = amoumt +").
		WithArgs(sum, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Подготовка БД для таблицы с транзакциями
	transactionDB, transactionMock, createTransactDBErr := sqlmock.New()
	if createTransactDBErr != nil {
		t.Fatalf("cant create mock: %s", createTransactDBErr)
	}
	defer transactionDB.Close()

	newTransaction := transaction.Transaction{
		TransactionID:   0,
		UserID:          userID,
		TransactionType: transaction.Refill,
		Sum:             sum,
		Time:            time.Now(),
		ActionComments:  "зачислены средства на баланс",
		AddComments:     comment,
	}

	transactionMock.ExpectExec("INSERT INTO balanceApp.transactions").
		WithArgs(newTransaction.TransactionID, newTransaction.UserID, newTransaction.TransactionType,
			newTransaction.Sum, sqlmock.AnyArg(), newTransaction.ActionComments, newTransaction.AddComments).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Подготовка БД для таблицы с заказами
	orderDB, orderMock, createOrderErr := sqlmock.New()
	if createOrderErr != nil {
		t.Fatalf("cant create mock: %s", createOrderErr)
	}
	defer orderDB.Close()

	// Создание объектов
	accountRepo := account_repo.NewAccountRepo(accountDB)
	accountController := ac.CreateNewAccountController(accountRepo)

	orderRepo := order_repo.NewOrderRepo(orderDB)
	orderController := oc.CreateNewOrderController(orderRepo)

	transactionRepo := transaction_repo.NewTransactionRepo(transactionDB)
	transactionController := tc.CreateNewTransactionController(transactionRepo)

	testManager := CreateNewManager(accountController, orderController, transactionController)

	// Тест
	err := testManager.RefillBalance(userID, sum, comment)

	// Проверка
	if err != nil {
		t.Errorf("unexpected err: %v", err)
		return
	}

	accountSecRows := sqlmock.NewRows([]string{"amount"})
	lastExpect := expectResult + sum
	accountSecRows = accountSecRows.AddRow(lastExpect)

	accountMock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(userID).
		WillReturnRows(accountSecRows)

	balance, getBalanceErr := testManager.accountController.CheckBalance(userID)

	if getBalanceErr != nil {
		t.Errorf("unexpected err: %v", getBalanceErr)
		return
	}

	if expectationAccErr := accountMock.ExpectationsWereMet(); expectationAccErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationAccErr)
		return
	}

	if expectationOrderErr := orderMock.ExpectationsWereMet(); expectationOrderErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationOrderErr)
		return
	}

	if expectationTransactionsErr := transactionMock.ExpectationsWereMet(); expectationTransactionsErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationTransactionsErr)
		return
	}

	if !reflect.DeepEqual(balance, lastExpect) {
		t.Errorf("results not match, want %v, have %v", lastExpect, balance)
		return
	}
}

// TestRefillMoney проверяет, что сценарий пополнения денег отрабатывает без ошибок, если аккаунта не существует
func TestRefillMoneyWithNoExistsAccount(t *testing.T) {
	var (
		userID  int64   = 1
		sum     float64 = 200
		comment         = "Всё хорошо!"
	)

	// Подготовка БД для таблицы с аккаунтами
	accountDB, accountMock, createAccountDBErr := sqlmock.New()
	if createAccountDBErr != nil {
		t.Fatalf("cant create mock: %s", createAccountDBErr)
	}
	defer accountDB.Close()

	accountFirstRows := sqlmock.NewRows([]string{"amount"})
	accountSecondRows := sqlmock.NewRows([]string{"amount"})

	accountMock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(userID).
		WillReturnRows(accountFirstRows)

	accountMock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(userID).
		WillReturnRows(accountSecondRows)

	accountMock.ExpectExec("INSERT INTO balanceApp.accounts").
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	accountMock.ExpectExec("UPDATE balanceApp.accounts SET amount = amoumt +").
		WithArgs(sum, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Подготовка БД для таблицы с транзакциями
	transactionDB, transactionMock, createTransactDBErr := sqlmock.New()
	if createTransactDBErr != nil {
		t.Fatalf("cant create mock: %s", createTransactDBErr)
	}
	defer transactionDB.Close()

	newTransaction := transaction.Transaction{
		TransactionID:   0,
		UserID:          userID,
		TransactionType: transaction.Refill,
		Sum:             sum,
		Time:            time.Now(),
		ActionComments:  "зачислены средства на баланс",
		AddComments:     comment,
	}

	transactionMock.ExpectExec("INSERT INTO balanceApp.transactions").
		WithArgs(newTransaction.TransactionID, newTransaction.UserID, newTransaction.TransactionType,
			newTransaction.Sum, sqlmock.AnyArg(), newTransaction.ActionComments, newTransaction.AddComments).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Подготовка БД для таблицы с заказами
	orderDB, orderMock, createOrderErr := sqlmock.New()
	if createOrderErr != nil {
		t.Fatalf("cant create mock: %s", createOrderErr)
	}
	defer orderDB.Close()

	// Создание объектов
	accountRepo := account_repo.NewAccountRepo(accountDB)
	accountController := ac.CreateNewAccountController(accountRepo)

	orderRepo := order_repo.NewOrderRepo(orderDB)
	orderController := oc.CreateNewOrderController(orderRepo)

	transactionRepo := transaction_repo.NewTransactionRepo(transactionDB)
	transactionController := tc.CreateNewTransactionController(transactionRepo)

	testManager := CreateNewManager(accountController, orderController, transactionController)

	// Тест
	err := testManager.RefillBalance(userID, sum, comment)

	// Проверка
	if err != nil {
		t.Errorf("unexpected err: %v", err)
		return
	}

	accountSecRows := sqlmock.NewRows([]string{"amount"})
	accountSecRows = accountSecRows.AddRow(sum)

	accountMock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(userID).
		WillReturnRows(accountSecRows)

	balance, getBalanceErr := testManager.accountController.CheckBalance(userID)

	if getBalanceErr != nil {
		t.Errorf("unexpected err: %v", getBalanceErr)
		return
	}

	if expectationAccErr := accountMock.ExpectationsWereMet(); expectationAccErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationAccErr)
		return
	}

	if expectationOrderErr := orderMock.ExpectationsWereMet(); expectationOrderErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationOrderErr)
		return
	}

	if expectationTransactionsErr := transactionMock.ExpectationsWereMet(); expectationTransactionsErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationTransactionsErr)
		return
	}

	if !reflect.DeepEqual(balance, sum) {
		t.Errorf("results not match, want %v, have %v", sum, balance)
		return
	}
}
