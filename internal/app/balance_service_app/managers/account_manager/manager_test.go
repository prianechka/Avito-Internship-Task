package account_manager

import (
	ac "Avito-Internship-Task/internal/app/balance_service_app/account/account_controller"
	"Avito-Internship-Task/internal/app/balance_service_app/account/account_repo"
	"Avito-Internship-Task/internal/app/balance_service_app/transaction"
	tc "Avito-Internship-Task/internal/app/balance_service_app/transaction/transaction_controller"
	"Avito-Internship-Task/internal/app/balance_service_app/transaction/transaction_repo"
	"fmt"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"reflect"
	"testing"
	"time"
)

// TestRefillMoney проверяет, что сценарий пополнения денег отрабатывает без ошибок
func TestRefillMoneyWithExistsAccount(t *testing.T) {
	var (
		userID          = 1
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

	accountMock.ExpectExec("UPDATE balanceApp.accounts SET amount = amount +").
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

	// Создание объектов
	accountRepo := account_repo.NewAccountRepo(accountDB)
	accountController := ac.CreateNewAccountController(accountRepo)

	transactionRepo := transaction_repo.NewTransactionRepo(transactionDB)
	transactionController := tc.CreateNewTransactionController(transactionRepo)

	testManager := CreateNewAccountManager(accountController, transactionController)

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
		userID          = 1
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

	accountMock.ExpectExec("UPDATE balanceApp.accounts SET amount = amount +").
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

	// Создание объектов
	accountRepo := account_repo.NewAccountRepo(accountDB)
	accountController := ac.CreateNewAccountController(accountRepo)

	transactionRepo := transaction_repo.NewTransactionRepo(transactionDB)
	transactionController := tc.CreateNewTransactionController(transactionRepo)

	testManager := CreateNewAccountManager(accountController, transactionController)

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

	if expectationTransactionsErr := transactionMock.ExpectationsWereMet(); expectationTransactionsErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationTransactionsErr)
		return
	}

	if !reflect.DeepEqual(balance, sum) {
		t.Errorf("results not match, want %v, have %v", sum, balance)
		return
	}
}

// TestGetBalance проверяет, что проверка баланса работает корректно.
func TestGetBalance(t *testing.T) {
	var (
		userID          = 1
		balance float64 = 200
	)

	// Подготовка БД для таблицы с аккаунтами
	accountDB, accountMock, createAccountDBErr := sqlmock.New()
	if createAccountDBErr != nil {
		t.Fatalf("cant create mock: %s", createAccountDBErr)
	}
	defer accountDB.Close()

	accountFirstRows := sqlmock.NewRows([]string{"amount"})
	accountFirstRows.AddRow(balance)
	accountSecondRows := sqlmock.NewRows([]string{"amount"})
	accountSecondRows.AddRow(balance)

	accountMock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(userID).
		WillReturnRows(accountFirstRows)

	accountMock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(userID).
		WillReturnRows(accountSecondRows)

	// Подготовка БД для таблицы с транзакциями
	transactionDB, transactionMock, createTransactDBErr := sqlmock.New()
	if createTransactDBErr != nil {
		t.Fatalf("cant create mock: %s", createTransactDBErr)
	}
	defer transactionDB.Close()

	// Создание объектов
	accountRepo := account_repo.NewAccountRepo(accountDB)
	accountController := ac.CreateNewAccountController(accountRepo)

	transactionRepo := transaction_repo.NewTransactionRepo(transactionDB)
	transactionController := tc.CreateNewTransactionController(transactionRepo)

	testManager := CreateNewAccountManager(accountController, transactionController)

	// Тест
	curBalance, err := testManager.GetUserBalance(userID)

	// Проверка
	if err != nil {
		t.Errorf("unexpected err: %v", err)
		return
	}

	if expectationAccErr := accountMock.ExpectationsWereMet(); expectationAccErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationAccErr)
		return
	}

	if expectationTransactionsErr := transactionMock.ExpectationsWereMet(); expectationTransactionsErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationTransactionsErr)
		return
	}

	if !reflect.DeepEqual(balance, curBalance) {
		t.Errorf("results not match, want %v, have %v", balance, curBalance)
		return
	}
}

// TestTransferSuccess проверяет, что перевод средств между двумя аккаунтами прошёл успешно
func TestTransferSuccess(t *testing.T) {
	var (
		srcUserID             = 1
		dstUserID             = 2
		sum           float64 = 200
		balanceFirst  float64 = 400
		balanceSecond float64 = 200
		comment               = "Всё хорошо!"
	)

	// Подготовка БД для таблицы с аккаунтами
	accountDB, accountMock, createAccountDBErr := sqlmock.New()
	if createAccountDBErr != nil {
		t.Fatalf("cant create mock: %s", createAccountDBErr)
	}
	defer accountDB.Close()

	accountFirstRows := sqlmock.NewRows([]string{"amount"})
	accountFirstRows.AddRow(balanceFirst)
	accountSecondRows := sqlmock.NewRows([]string{"amount"})
	accountSecondRows.AddRow(balanceSecond)
	accountThirdRows := sqlmock.NewRows([]string{"amount"})
	accountThirdRows.AddRow(balanceFirst)
	accountFourthRows := sqlmock.NewRows([]string{"amount"})
	accountFourthRows.AddRow(balanceFirst)
	accountFifthRows := sqlmock.NewRows([]string{"amount"})
	accountFifthRows.AddRow(balanceFirst)

	accountMock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(srcUserID).
		WillReturnRows(accountFirstRows)

	accountMock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(dstUserID).
		WillReturnRows(accountSecondRows)

	accountMock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(srcUserID).
		WillReturnRows(accountThirdRows)

	accountMock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(srcUserID).
		WillReturnRows(accountFourthRows)

	accountMock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(srcUserID).
		WillReturnRows(accountFifthRows)

	accountMock.ExpectExec("UPDATE balanceApp.accounts SET amount = amount +").
		WithArgs(-sum, srcUserID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	accountMock.ExpectExec("UPDATE balanceApp.accounts SET amount = amount +").
		WithArgs(sum, dstUserID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Подготовка БД для таблицы с транзакциями
	transactionDB, transactionMock, createTransactDBErr := sqlmock.New()
	if createTransactDBErr != nil {
		t.Fatalf("cant create mock: %s", createTransactDBErr)
	}
	defer transactionDB.Close()

	newTransaction := transaction.Transaction{
		TransactionID:   0,
		UserID:          srcUserID,
		TransactionType: transaction.Transfer,
		Sum:             sum,
		Time:            time.Now(),
		ActionComments:  "перевод пользователю: " + fmt.Sprintf("%d", dstUserID),
		AddComments:     comment,
	}

	newTransaction1 := transaction.Transaction{
		TransactionID:   1,
		UserID:          dstUserID,
		TransactionType: transaction.Transfer,
		Sum:             sum,
		Time:            time.Now(),
		ActionComments:  "перевод от пользователя: " + fmt.Sprintf("%d", srcUserID),
		AddComments:     comment,
	}

	transactionMock.ExpectExec("INSERT INTO balanceApp.transactions").
		WithArgs(newTransaction.TransactionID, newTransaction.UserID, newTransaction.TransactionType,
			newTransaction.Sum, sqlmock.AnyArg(), newTransaction.ActionComments, newTransaction.AddComments).
		WillReturnResult(sqlmock.NewResult(1, 1))

	transactionMock.ExpectExec("INSERT INTO balanceApp.transactions").
		WithArgs(newTransaction1.TransactionID, newTransaction1.UserID, newTransaction1.TransactionType,
			newTransaction1.Sum, sqlmock.AnyArg(), newTransaction1.ActionComments, newTransaction1.AddComments).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Создание объектов
	accountRepo := account_repo.NewAccountRepo(accountDB)
	accountController := ac.CreateNewAccountController(accountRepo)

	transactionRepo := transaction_repo.NewTransactionRepo(transactionDB)
	transactionController := tc.CreateNewTransactionController(transactionRepo)

	testManager := CreateNewAccountManager(accountController, transactionController)

	// Тест
	err := testManager.Transfer(srcUserID, dstUserID, sum, comment)

	// Проверка
	if err != nil {
		t.Errorf("unexpected err: %v", err)
		return
	}

	accountLastRowsSrc := sqlmock.NewRows([]string{"amount"})
	accountLastRowsSrc = accountLastRowsSrc.AddRow(balanceFirst - sum)

	accountLastRowsDst := sqlmock.NewRows([]string{"amount"})
	accountLastRowsDst = accountLastRowsSrc.AddRow(balanceSecond + sum)

	accountMock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(srcUserID).
		WillReturnRows(accountLastRowsSrc)

	accountMock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(dstUserID).
		WillReturnRows(accountLastRowsDst)

	curBalanceSrc, getBalanceSrcErr := testManager.accountController.CheckBalance(srcUserID)

	if getBalanceSrcErr != nil {
		t.Errorf("unexpected err: %v", getBalanceSrcErr)
		return
	}

	curBalanceDst, getBalanceDstErr := testManager.accountController.CheckBalance(dstUserID)

	if getBalanceDstErr != nil {
		t.Errorf("unexpected err: %v", getBalanceDstErr)
		return
	}

	if expectationAccErr := accountMock.ExpectationsWereMet(); expectationAccErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationAccErr)
		return
	}

	if expectationTransactionsErr := transactionMock.ExpectationsWereMet(); expectationTransactionsErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationTransactionsErr)
		return
	}

	if !reflect.DeepEqual(balanceFirst-sum, curBalanceSrc) {
		t.Errorf("results not match, want %v, have %v", balanceFirst-sum, curBalanceSrc)
		return
	}

	if !reflect.DeepEqual(balanceSecond+sum, curBalanceDst) {
		t.Errorf("results not match, want %v, have %v", balanceSecond+sum, curBalanceDst)
		return
	}
}

// TestTransferAccNotExistError проверяет, что если передан несуществующий аккаунт, то вернётся ошибка
func TestTransferAccNotExistError(t *testing.T) {
	var (
		srcUserID         = 1
		dstUserID         = 2
		sum       float64 = 200
		comment           = "Всё хорошо!"
	)

	// Подготовка БД для таблицы с аккаунтами
	accountDB, accountMock, createAccountDBErr := sqlmock.New()
	if createAccountDBErr != nil {
		t.Fatalf("cant create mock: %s", createAccountDBErr)
	}
	defer accountDB.Close()

	accountFirstRows := sqlmock.NewRows([]string{"amount"})

	accountMock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(srcUserID).
		WillReturnRows(accountFirstRows)

	// Подготовка БД для таблицы с транзакциями
	transactionDB, transactionMock, createTransactDBErr := sqlmock.New()
	if createTransactDBErr != nil {
		t.Fatalf("cant create mock: %s", createTransactDBErr)
	}
	defer transactionDB.Close()

	// Создание объектов
	accountRepo := account_repo.NewAccountRepo(accountDB)
	accountController := ac.CreateNewAccountController(accountRepo)

	transactionRepo := transaction_repo.NewTransactionRepo(transactionDB)
	transactionController := tc.CreateNewTransactionController(transactionRepo)

	testManager := CreateNewAccountManager(accountController, transactionController)

	// Тест
	err := testManager.Transfer(srcUserID, dstUserID, sum, comment)

	// Проверка
	if err != ac.AccountNotExistErr {
		t.Errorf("unexpected err: %v", err)
		return
	}

	if expectationAccErr := accountMock.ExpectationsWereMet(); expectationAccErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationAccErr)
		return
	}

	if expectationTransactionsErr := transactionMock.ExpectationsWereMet(); expectationTransactionsErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationTransactionsErr)
		return
	}
}

// TestTransferNotEnoughMoneyError проверяет, что если недостаточно денег, то перевод не произойдёт
func TestTransferNotEnoughMoneyError(t *testing.T) {
	var (
		srcUserID             = 1
		dstUserID             = 2
		sum           float64 = 500
		balanceFirst  float64 = 400
		balanceSecond float64 = 200
		comment               = "Всё хорошо!"
	)

	// Подготовка БД для таблицы с аккаунтами
	accountDB, accountMock, createAccountDBErr := sqlmock.New()
	if createAccountDBErr != nil {
		t.Fatalf("cant create mock: %s", createAccountDBErr)
	}
	defer accountDB.Close()

	accountFirstRows := sqlmock.NewRows([]string{"amount"})
	accountFirstRows.AddRow(balanceFirst)
	accountSecondRows := sqlmock.NewRows([]string{"amount"})
	accountSecondRows.AddRow(balanceSecond)
	accountThirdRows := sqlmock.NewRows([]string{"amount"})
	accountThirdRows.AddRow(balanceFirst)
	accountFourthRows := sqlmock.NewRows([]string{"amount"})
	accountFourthRows.AddRow(balanceFirst)

	accountMock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(srcUserID).
		WillReturnRows(accountFirstRows)

	accountMock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(dstUserID).
		WillReturnRows(accountSecondRows)

	accountMock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(srcUserID).
		WillReturnRows(accountThirdRows)

	accountMock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(srcUserID).
		WillReturnRows(accountFourthRows)

	// Подготовка БД для таблицы с транзакциями
	transactionDB, transactionMock, createTransactDBErr := sqlmock.New()
	if createTransactDBErr != nil {
		t.Fatalf("cant create mock: %s", createTransactDBErr)
	}
	defer transactionDB.Close()

	// Создание объектов
	accountRepo := account_repo.NewAccountRepo(accountDB)
	accountController := ac.CreateNewAccountController(accountRepo)

	transactionRepo := transaction_repo.NewTransactionRepo(transactionDB)
	transactionController := tc.CreateNewTransactionController(transactionRepo)

	testManager := CreateNewAccountManager(accountController, transactionController)

	// Тест
	err := testManager.Transfer(srcUserID, dstUserID, sum, comment)

	// Проверка
	if err != ac.NotEnoughMoneyErr {
		t.Errorf("unexpected err: %v", err)
		return
	}

	accountLastRowsSrc := sqlmock.NewRows([]string{"amount"})
	accountLastRowsSrc = accountLastRowsSrc.AddRow(balanceFirst)

	accountLastRowsDst := sqlmock.NewRows([]string{"amount"})
	accountLastRowsDst = accountLastRowsSrc.AddRow(balanceSecond)

	accountMock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(srcUserID).
		WillReturnRows(accountLastRowsSrc)

	accountMock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(dstUserID).
		WillReturnRows(accountLastRowsDst)

	curBalanceSrc, getBalanceSrcErr := testManager.accountController.CheckBalance(srcUserID)

	if getBalanceSrcErr != nil {
		t.Errorf("unexpected err: %v", getBalanceSrcErr)
		return
	}

	curBalanceDst, getBalanceDstErr := testManager.accountController.CheckBalance(dstUserID)

	if getBalanceDstErr != nil {
		t.Errorf("unexpected err: %v", getBalanceDstErr)
		return
	}

	if expectationAccErr := accountMock.ExpectationsWereMet(); expectationAccErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationAccErr)
		return
	}

	if expectationTransactionsErr := transactionMock.ExpectationsWereMet(); expectationTransactionsErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationTransactionsErr)
		return
	}

	if !reflect.DeepEqual(balanceFirst, curBalanceSrc) {
		t.Errorf("results not match, want %v, have %v", balanceFirst-sum, curBalanceSrc)
		return
	}

	if !reflect.DeepEqual(balanceSecond, curBalanceDst) {
		t.Errorf("results not match, want %v, have %v", balanceSecond+sum, curBalanceDst)
		return
	}
}
