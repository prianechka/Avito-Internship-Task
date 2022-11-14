package manager

import (
	ac "Avito-Internship-Task/internal/app/balance_service_app/account/account_controller"
	"Avito-Internship-Task/internal/app/balance_service_app/account/account_repo"
	"Avito-Internship-Task/internal/app/balance_service_app/order"
	oc "Avito-Internship-Task/internal/app/balance_service_app/order/order_controller"
	"Avito-Internship-Task/internal/app/balance_service_app/order/order_repo"
	"Avito-Internship-Task/internal/app/balance_service_app/report"
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

// TestBuyServiceSuccess проверяет, что сценарий покупки сервиса проходит успешно.
func TestBuyServiceSuccess(t *testing.T) {
	var (
		userID    int64   = 1
		orderID   int64   = 1
		serviceID int64   = 1
		sum       float64 = 200
		balance   float64 = 400
		comment           = "Всё хорошо!"
	)

	curTime := time.Now()

	myOrder := order.Order{
		OrderID:      orderID,
		UserID:       userID,
		ServiceID:    serviceID,
		OrderCost:    sum,
		CreatingTime: curTime,
		Comment:      comment,
		OrderState:   order.REGISTRATED,
	}

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
	accountThirdRows := sqlmock.NewRows([]string{"amount"})
	accountThirdRows.AddRow(balance)

	accountMock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(userID).
		WillReturnRows(accountFirstRows)

	accountMock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(userID).
		WillReturnRows(accountSecondRows)

	accountMock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(userID).
		WillReturnRows(accountThirdRows)

	accountMock.ExpectExec("UPDATE balanceApp.accounts SET amount = amount +").
		WithArgs(-sum, userID).
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
		TransactionType: transaction.Buy,
		Sum:             sum,
		Time:            time.Now(),
		ActionComments:  "куплена услуга: " + order.Types[serviceID],
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

	firstOrderRow := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})
	secondOrderRow := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})
	thirdOrderRows := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})
	fourthOrderRows := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})

	expect := myOrder
	thirdOrderRows.AddRow(expect.OrderID, expect.UserID, expect.ServiceID, expect.OrderCost,
		expect.CreatingTime, expect.Comment, expect.OrderState)
	fourthOrderRows.AddRow(expect.OrderID, expect.UserID, expect.ServiceID, expect.OrderCost,
		expect.CreatingTime, expect.Comment, expect.OrderState)

	orderMock.ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(firstOrderRow).WillReturnError(nil)

	orderMock.ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(secondOrderRow).WillReturnError(nil)

	orderMock.
		ExpectExec("INSERT INTO balanceApp.orders").
		WithArgs(orderID, userID, serviceID, sum, sqlmock.AnyArg(), comment, order.REGISTRATED).
		WillReturnResult(sqlmock.NewResult(1, 1))

	orderMock.ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(thirdOrderRows).WillReturnError(nil)

	orderMock.ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(fourthOrderRows).WillReturnError(nil)

	orderMock.ExpectExec("UPDATE balanceApp.orders SET orderState = ").
		WithArgs(order.RESERVED, orderID, userID, serviceID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Создание объектов
	accountRepo := account_repo.NewAccountRepo(accountDB)
	accountController := ac.CreateNewAccountController(accountRepo)

	orderRepo := order_repo.NewOrderRepo(orderDB)
	orderController := oc.CreateNewOrderController(orderRepo)

	transactionRepo := transaction_repo.NewTransactionRepo(transactionDB)
	transactionController := tc.CreateNewTransactionController(transactionRepo)

	testManager := CreateNewManager(accountController, orderController, transactionController)

	// Тест
	err := testManager.BuyService(userID, orderID, serviceID, sum, comment)

	// Проверка
	if err != nil {
		t.Errorf("unexpected err: %v", err)
		return
	}

	accountLastRows := sqlmock.NewRows([]string{"amount"})
	accountLastRows = accountLastRows.AddRow(balance - sum)

	accountMock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(userID).
		WillReturnRows(accountLastRows)

	curBalance, getBalanceErr := testManager.accountController.CheckBalance(userID)

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

	if !reflect.DeepEqual(balance-sum, curBalance) {
		t.Errorf("results not match, want %v, have %v", balance-sum, curBalance)
		return
	}
}

// TestBuyServiceNotAccExist проверяет, что покупки не произойдёт, если аккаунта не существует.
func TestBuyServiceNotAccExist(t *testing.T) {
	var (
		userID    int64   = 1
		orderID   int64   = 1
		serviceID int64   = 1
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
		WithArgs(userID).
		WillReturnRows(accountFirstRows)

	// Подготовка БД для таблицы с транзакциями
	transactionDB, transactionMock, createTransactDBErr := sqlmock.New()
	if createTransactDBErr != nil {
		t.Fatalf("cant create mock: %s", createTransactDBErr)
	}
	defer transactionDB.Close()

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
	err := testManager.BuyService(userID, orderID, serviceID, sum, comment)

	// Проверка
	if err != ac.AccountNotExistErr {
		t.Errorf("unexpected err: %v", err)
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
}

// TestBuyServiceNotEnoughMoneyErr проверяет, что покупка не произойдёт, если баланс пользователя меньше нужного.
func TestBuyServiceNotEnoughMoneyErr(t *testing.T) {
	var (
		userID    int64   = 1
		orderID   int64   = 1
		serviceID int64   = 1
		sum       float64 = 400
		balance   float64 = 200
		comment           = "Всё хорошо!"
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
	err := testManager.BuyService(userID, orderID, serviceID, sum, comment)

	// Проверка
	if err != ac.NotEnoughMoneyErr {
		t.Errorf("unexpected err: %v", err)
		return
	}

	accountLastRows := sqlmock.NewRows([]string{"amount"})
	accountLastRows = accountLastRows.AddRow(balance)

	accountMock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(userID).
		WillReturnRows(accountLastRows)

	curBalance, getBalanceErr := testManager.accountController.CheckBalance(userID)

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

	if !reflect.DeepEqual(balance, curBalance) {
		t.Errorf("results not match, want %v, have %v", balance, curBalance)
		return
	}
}

// TestAcceptBuySuccess проверяет, что подтверждение покупки услуги проходит успешно.
func TestAcceptBuySuccess(t *testing.T) {
	var (
		userID    int64   = 1
		orderID   int64   = 1
		serviceID int64   = 1
		sum       float64 = 200
		comment           = "Всё хорошо!"
	)

	curTime := time.Now()

	myOrder := order.Order{
		OrderID:      orderID,
		UserID:       userID,
		ServiceID:    serviceID,
		OrderCost:    sum,
		CreatingTime: curTime,
		Comment:      comment,
		OrderState:   order.RESERVED,
	}

	// Подготовка БД для таблицы с аккаунтами
	accountDB, accountMock, createAccountDBErr := sqlmock.New()
	if createAccountDBErr != nil {
		t.Fatalf("cant create mock: %s", createAccountDBErr)
	}
	defer accountDB.Close()

	// Подготовка БД для таблицы с транзакциями
	transactionDB, transactionMock, createTransactDBErr := sqlmock.New()
	if createTransactDBErr != nil {
		t.Fatalf("cant create mock: %s", createTransactDBErr)
	}
	defer transactionDB.Close()

	// Подготовка БД для таблицы с заказами
	orderDB, orderMock, createOrderErr := sqlmock.New()
	if createOrderErr != nil {
		t.Fatalf("cant create mock: %s", createOrderErr)
	}
	defer orderDB.Close()

	firstOrderRow := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})
	secondOrderRow := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})

	thirdOrderRow := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})

	expect := myOrder
	firstOrderRow.AddRow(expect.OrderID, expect.UserID, expect.ServiceID, expect.OrderCost,
		expect.CreatingTime, expect.Comment, expect.OrderState)
	secondOrderRow.AddRow(expect.OrderID, expect.UserID, expect.ServiceID, expect.OrderCost,
		expect.CreatingTime, expect.Comment, expect.OrderState)
	thirdOrderRow.AddRow(expect.OrderID, expect.UserID, expect.ServiceID, expect.OrderCost,
		expect.CreatingTime, expect.Comment, expect.OrderState)

	orderMock.ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(firstOrderRow).WillReturnError(nil)

	orderMock.ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(secondOrderRow).WillReturnError(nil)

	orderMock.ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(thirdOrderRow).WillReturnError(nil)

	orderMock.ExpectExec("UPDATE balanceApp.orders SET orderState = ").
		WithArgs(order.FINISHED, orderID, userID, serviceID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Создание объектов
	accountRepo := account_repo.NewAccountRepo(accountDB)
	accountController := ac.CreateNewAccountController(accountRepo)

	orderRepo := order_repo.NewOrderRepo(orderDB)
	orderController := oc.CreateNewOrderController(orderRepo)

	transactionRepo := transaction_repo.NewTransactionRepo(transactionDB)
	transactionController := tc.CreateNewTransactionController(transactionRepo)

	testManager := CreateNewManager(accountController, orderController, transactionController)

	// Тест
	err := testManager.AcceptBuy(userID, orderID, serviceID)

	// Проверка
	if err != nil {
		t.Errorf("unexpected err: %v", err)
		return
	}

	lastOrderRow := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})
	lastOrderRow.AddRow(expect.OrderID, expect.UserID, expect.ServiceID, expect.OrderCost,
		expect.CreatingTime, expect.Comment, order.FINISHED)

	orderMock.ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(lastOrderRow).WillReturnError(nil)

	curOrder, getOrderErr := testManager.orderController.GetOrder(orderID, userID, serviceID)

	if getOrderErr != nil {
		t.Errorf("unexpected err: %v", getOrderErr)
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

	if !reflect.DeepEqual(curOrder.OrderState, int64(order.FINISHED)) {
		t.Errorf("results not match, want %v, have %v", order.FINISHED, curOrder.OrderState)
		return
	}
}

// TestAcceptBuyError проверяет, что если заказа не существует, то не получится подтвердить покупку.
func TestAcceptBuyError(t *testing.T) {
	var (
		userID    int64 = 1
		orderID   int64 = 1
		serviceID int64 = 1
	)

	// Подготовка БД для таблицы с аккаунтами
	accountDB, accountMock, createAccountDBErr := sqlmock.New()
	if createAccountDBErr != nil {
		t.Fatalf("cant create mock: %s", createAccountDBErr)
	}
	defer accountDB.Close()

	// Подготовка БД для таблицы с транзакциями
	transactionDB, transactionMock, createTransactDBErr := sqlmock.New()
	if createTransactDBErr != nil {
		t.Fatalf("cant create mock: %s", createTransactDBErr)
	}
	defer transactionDB.Close()

	// Подготовка БД для таблицы с заказами
	orderDB, orderMock, createOrderErr := sqlmock.New()
	if createOrderErr != nil {
		t.Fatalf("cant create mock: %s", createOrderErr)
	}
	defer orderDB.Close()

	firstOrderRow := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})

	orderMock.ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(firstOrderRow).WillReturnError(nil)

	// Создание объектов
	accountRepo := account_repo.NewAccountRepo(accountDB)
	accountController := ac.CreateNewAccountController(accountRepo)

	orderRepo := order_repo.NewOrderRepo(orderDB)
	orderController := oc.CreateNewOrderController(orderRepo)

	transactionRepo := transaction_repo.NewTransactionRepo(transactionDB)
	transactionController := tc.CreateNewTransactionController(transactionRepo)

	testManager := CreateNewManager(accountController, orderController, transactionController)

	// Тест
	err := testManager.AcceptBuy(userID, orderID, serviceID)

	// Проверка
	if err != oc.OrderNotFound {
		t.Errorf("unexpected err: %v", err)
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
}

// TestAcceptBuyWrongStatusError проверяет, что если статус неверный, подтвердить покупку не получится.
func TestAcceptBuyWrongStatusError(t *testing.T) {
	var (
		userID    int64   = 1
		orderID   int64   = 1
		serviceID int64   = 1
		sum       float64 = 200
		comment           = "Всё хорошо!"
	)

	curTime := time.Now()

	myOrder := order.Order{
		OrderID:      orderID,
		UserID:       userID,
		ServiceID:    serviceID,
		OrderCost:    sum,
		CreatingTime: curTime,
		Comment:      comment,
		OrderState:   order.REGISTRATED,
	}

	// Подготовка БД для таблицы с аккаунтами
	accountDB, accountMock, createAccountDBErr := sqlmock.New()
	if createAccountDBErr != nil {
		t.Fatalf("cant create mock: %s", createAccountDBErr)
	}
	defer accountDB.Close()

	// Подготовка БД для таблицы с транзакциями
	transactionDB, transactionMock, createTransactDBErr := sqlmock.New()
	if createTransactDBErr != nil {
		t.Fatalf("cant create mock: %s", createTransactDBErr)
	}
	defer transactionDB.Close()

	// Подготовка БД для таблицы с заказами
	orderDB, orderMock, createOrderErr := sqlmock.New()
	if createOrderErr != nil {
		t.Fatalf("cant create mock: %s", createOrderErr)
	}
	defer orderDB.Close()

	firstOrderRow := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})
	secondOrderRow := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})

	thirdOrderRow := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})

	expect := myOrder
	firstOrderRow.AddRow(expect.OrderID, expect.UserID, expect.ServiceID, expect.OrderCost,
		expect.CreatingTime, expect.Comment, expect.OrderState)
	secondOrderRow.AddRow(expect.OrderID, expect.UserID, expect.ServiceID, expect.OrderCost,
		expect.CreatingTime, expect.Comment, expect.OrderState)
	thirdOrderRow.AddRow(expect.OrderID, expect.UserID, expect.ServiceID, expect.OrderCost,
		expect.CreatingTime, expect.Comment, expect.OrderState)

	orderMock.ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(firstOrderRow).WillReturnError(nil)

	orderMock.ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(secondOrderRow).WillReturnError(nil)

	orderMock.ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(thirdOrderRow).WillReturnError(nil)

	// Создание объектов
	accountRepo := account_repo.NewAccountRepo(accountDB)
	accountController := ac.CreateNewAccountController(accountRepo)

	orderRepo := order_repo.NewOrderRepo(orderDB)
	orderController := oc.CreateNewOrderController(orderRepo)

	transactionRepo := transaction_repo.NewTransactionRepo(transactionDB)
	transactionController := tc.CreateNewTransactionController(transactionRepo)

	testManager := CreateNewManager(accountController, orderController, transactionController)

	// Тест
	err := testManager.AcceptBuy(userID, orderID, serviceID)

	// Проверка
	if err != oc.WrongStateError {
		t.Errorf("unexpected err: %v", err)
		return
	}

	lastOrderRow := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})
	lastOrderRow.AddRow(expect.OrderID, expect.UserID, expect.ServiceID, expect.OrderCost,
		expect.CreatingTime, expect.Comment, order.REGISTRATED)

	orderMock.ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(lastOrderRow).WillReturnError(nil)

	curOrder, getOrderErr := testManager.orderController.GetOrder(orderID, userID, serviceID)

	if getOrderErr != nil {
		t.Errorf("unexpected err: %v", getOrderErr)
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

	if !reflect.DeepEqual(curOrder.OrderState, myOrder.OrderState) {
		t.Errorf("results not match, want %v, have %v", myOrder.OrderState, curOrder.OrderState)
		return
	}
}

// TestRefuseServiceSuccess проверяет, что возврат произойдёт успешно.
func TestRefuseServiceSuccess(t *testing.T) {
	var (
		userID    int64   = 1
		orderID   int64   = 1
		serviceID int64   = 1
		sum       float64 = 200
		balance   float64 = 400
		comment           = "Всё хорошо!"
	)

	curTime := time.Now()

	myOrder := order.Order{
		OrderID:      orderID,
		UserID:       userID,
		ServiceID:    serviceID,
		OrderCost:    sum,
		CreatingTime: curTime,
		Comment:      comment,
		OrderState:   order.RESERVED,
	}

	// Подготовка БД для таблицы с аккаунтами
	accountDB, accountMock, createAccountDBErr := sqlmock.New()
	if createAccountDBErr != nil {
		t.Fatalf("cant create mock: %s", createAccountDBErr)
	}
	defer accountDB.Close()

	accountFirstRows := sqlmock.NewRows([]string{"amount"})
	accountFirstRows.AddRow(balance)

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
		TransactionType: transaction.Return,
		Sum:             sum,
		Time:            time.Now(),
		ActionComments:  "возврат за услугу: " + order.Types[serviceID],
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

	firstOrderRow := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})
	secondOrderRow := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})

	thirdOrderRow := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})

	expect := myOrder
	firstOrderRow.AddRow(expect.OrderID, expect.UserID, expect.ServiceID, expect.OrderCost,
		expect.CreatingTime, expect.Comment, expect.OrderState)
	secondOrderRow.AddRow(expect.OrderID, expect.UserID, expect.ServiceID, expect.OrderCost,
		expect.CreatingTime, expect.Comment, expect.OrderState)
	thirdOrderRow.AddRow(expect.OrderID, expect.UserID, expect.ServiceID, expect.OrderCost,
		expect.CreatingTime, expect.Comment, expect.OrderState)

	orderMock.ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(firstOrderRow).WillReturnError(nil)

	orderMock.ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(secondOrderRow).WillReturnError(nil)

	orderMock.ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(thirdOrderRow).WillReturnError(nil)

	orderMock.ExpectExec("UPDATE balanceApp.orders SET orderState = ").
		WithArgs(order.RETURNED, orderID, userID, serviceID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Создание объектов
	accountRepo := account_repo.NewAccountRepo(accountDB)
	accountController := ac.CreateNewAccountController(accountRepo)

	orderRepo := order_repo.NewOrderRepo(orderDB)
	orderController := oc.CreateNewOrderController(orderRepo)

	transactionRepo := transaction_repo.NewTransactionRepo(transactionDB)
	transactionController := tc.CreateNewTransactionController(transactionRepo)

	testManager := CreateNewManager(accountController, orderController, transactionController)

	// Тест
	err := testManager.RefuseBuy(userID, orderID, serviceID, comment)

	// Проверка
	if err != nil {
		t.Errorf("unexpected err: %v", err)
		return
	}

	accountLastRows := sqlmock.NewRows([]string{"amount"})
	accountLastRows = accountLastRows.AddRow(balance + sum)

	accountMock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(userID).
		WillReturnRows(accountLastRows)

	curBalance, getBalanceErr := testManager.accountController.CheckBalance(userID)

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

	if !reflect.DeepEqual(balance+sum, curBalance) {
		t.Errorf("results not match, want %v, have %v", balance+sum, curBalance)
		return
	}
}

// TestRefuseServiceWrongStatusError проверяет, что возврат не произойдёт, если статус неправильный.
func TestRefuseServiceWrongStatusError(t *testing.T) {
	var (
		userID    int64   = 1
		orderID   int64   = 1
		serviceID int64   = 1
		sum       float64 = 200
		balance   float64 = 400
		comment           = "Всё хорошо!"
	)

	curTime := time.Now()

	myOrder := order.Order{
		OrderID:      orderID,
		UserID:       userID,
		ServiceID:    serviceID,
		OrderCost:    sum,
		CreatingTime: curTime,
		Comment:      comment,
		OrderState:   order.REGISTRATED,
	}

	// Подготовка БД для таблицы с аккаунтами
	accountDB, accountMock, createAccountDBErr := sqlmock.New()
	if createAccountDBErr != nil {
		t.Fatalf("cant create mock: %s", createAccountDBErr)
	}
	defer accountDB.Close()

	// Подготовка БД для таблицы с транзакциями
	transactionDB, transactionMock, createTransactDBErr := sqlmock.New()
	if createTransactDBErr != nil {
		t.Fatalf("cant create mock: %s", createTransactDBErr)
	}
	defer transactionDB.Close()

	// Подготовка БД для таблицы с заказами
	orderDB, orderMock, createOrderErr := sqlmock.New()
	if createOrderErr != nil {
		t.Fatalf("cant create mock: %s", createOrderErr)
	}
	defer orderDB.Close()

	firstOrderRow := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})
	secondOrderRow := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})

	thirdOrderRow := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})

	expect := myOrder
	firstOrderRow.AddRow(expect.OrderID, expect.UserID, expect.ServiceID, expect.OrderCost,
		expect.CreatingTime, expect.Comment, expect.OrderState)
	secondOrderRow.AddRow(expect.OrderID, expect.UserID, expect.ServiceID, expect.OrderCost,
		expect.CreatingTime, expect.Comment, expect.OrderState)
	thirdOrderRow.AddRow(expect.OrderID, expect.UserID, expect.ServiceID, expect.OrderCost,
		expect.CreatingTime, expect.Comment, expect.OrderState)

	orderMock.ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(firstOrderRow).WillReturnError(nil)

	orderMock.ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(secondOrderRow).WillReturnError(nil)

	orderMock.ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(thirdOrderRow).WillReturnError(nil)

	// Создание объектов
	accountRepo := account_repo.NewAccountRepo(accountDB)
	accountController := ac.CreateNewAccountController(accountRepo)

	orderRepo := order_repo.NewOrderRepo(orderDB)
	orderController := oc.CreateNewOrderController(orderRepo)

	transactionRepo := transaction_repo.NewTransactionRepo(transactionDB)
	transactionController := tc.CreateNewTransactionController(transactionRepo)

	testManager := CreateNewManager(accountController, orderController, transactionController)

	// Тест
	err := testManager.RefuseBuy(userID, orderID, serviceID, comment)

	// Проверка
	if err != oc.WrongStateError {
		t.Errorf("unexpected err: %v", err)
		return
	}

	accountLastRows := sqlmock.NewRows([]string{"amount"})
	accountLastRows = accountLastRows.AddRow(balance)

	accountMock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(userID).
		WillReturnRows(accountLastRows)

	curBalance, getBalanceErr := testManager.accountController.CheckBalance(userID)

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

	if !reflect.DeepEqual(balance, curBalance) {
		t.Errorf("results not match, want %v, have %v", balance+sum, curBalance)
		return
	}
}

// TestRefuseServiceOrderNoExistError проверяет, что возврат не произойдёт, если заказа не существует.
func TestRefuseServiceOrderNoExistError(t *testing.T) {
	var (
		userID    int64 = 1
		orderID   int64 = 1
		serviceID int64 = 1
		comment         = "Всё хорошо!"
	)

	// Подготовка БД для таблицы с аккаунтами
	accountDB, accountMock, createAccountDBErr := sqlmock.New()
	if createAccountDBErr != nil {
		t.Fatalf("cant create mock: %s", createAccountDBErr)
	}
	defer accountDB.Close()

	// Подготовка БД для таблицы с транзакциями
	transactionDB, transactionMock, createTransactDBErr := sqlmock.New()
	if createTransactDBErr != nil {
		t.Fatalf("cant create mock: %s", createTransactDBErr)
	}
	defer transactionDB.Close()

	// Подготовка БД для таблицы с заказами
	orderDB, orderMock, createOrderErr := sqlmock.New()
	if createOrderErr != nil {
		t.Fatalf("cant create mock: %s", createOrderErr)
	}
	defer orderDB.Close()

	firstOrderRow := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})

	orderMock.ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(firstOrderRow).WillReturnError(nil)

	// Создание объектов
	accountRepo := account_repo.NewAccountRepo(accountDB)
	accountController := ac.CreateNewAccountController(accountRepo)

	orderRepo := order_repo.NewOrderRepo(orderDB)
	orderController := oc.CreateNewOrderController(orderRepo)

	transactionRepo := transaction_repo.NewTransactionRepo(transactionDB)
	transactionController := tc.CreateNewTransactionController(transactionRepo)

	testManager := CreateNewManager(accountController, orderController, transactionController)

	// Тест
	err := testManager.RefuseBuy(userID, orderID, serviceID, comment)

	// Проверка
	if err != oc.OrderNotFound {
		t.Errorf("unexpected err: %v", err)
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
}

// TestGetBalance проверяет, что проверка баланса работает корректно.
func TestGetBalance(t *testing.T) {
	var (
		userID  int64   = 1
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

	if expectationOrderErr := orderMock.ExpectationsWereMet(); expectationOrderErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationOrderErr)
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
		srcUserID     int64   = 1
		dstUserID     int64   = 2
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

	if expectationOrderErr := orderMock.ExpectationsWereMet(); expectationOrderErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationOrderErr)
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
		srcUserID int64   = 1
		dstUserID int64   = 2
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

	if expectationOrderErr := orderMock.ExpectationsWereMet(); expectationOrderErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationOrderErr)
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
		srcUserID     int64   = 1
		dstUserID     int64   = 2
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

	if expectationOrderErr := orderMock.ExpectationsWereMet(); expectationOrderErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationOrderErr)
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

// TestGetFinanceReportSuccess показывает, что программа способна сделать финансовый отчёт по заданию.
func TestGetFinanceReportSuccess(t *testing.T) {
	var (
		year  int64 = 2022
		month int64 = 2
		url         = "report.csv"
	)

	// Подготовка БД для таблицы с аккаунтами
	accountDB, accountMock, createAccountDBErr := sqlmock.New()
	if createAccountDBErr != nil {
		t.Fatalf("cant create mock: %s", createAccountDBErr)
	}
	defer accountDB.Close()

	// Подготовка БД для таблицы с транзакциями
	transactionDB, transactionMock, createTransactDBErr := sqlmock.New()
	if createTransactDBErr != nil {
		t.Fatalf("cant create mock: %s", createTransactDBErr)
	}
	defer transactionDB.Close()

	// Подготовка БД для таблицы с заказами
	orderDB, orderMock, createOrderErr := sqlmock.New()
	if createOrderErr != nil {
		t.Fatalf("cant create mock: %s", createOrderErr)
	}
	defer orderDB.Close()

	rows := sqlmock.NewRows([]string{"serviceType", "sum"})
	allServicesReport := []report.FinanceReport{{1, 100}, {2, 150}}
	for _, service := range allServicesReport {
		rows.AddRow(service.ServiceType, service.Sum)
	}

	orderMock.ExpectQuery("SELECT serviceType ").
		WithArgs(month, year).WillReturnRows(rows)

	// Создание объектов
	accountRepo := account_repo.NewAccountRepo(accountDB)
	accountController := ac.CreateNewAccountController(accountRepo)

	orderRepo := order_repo.NewOrderRepo(orderDB)
	orderController := oc.CreateNewOrderController(orderRepo)

	transactionRepo := transaction_repo.NewTransactionRepo(transactionDB)
	transactionController := tc.CreateNewTransactionController(transactionRepo)

	testManager := CreateNewManager(accountController, orderController, transactionController)

	// Тест
	err := testManager.GetFinanceReport(month, year, url)

	// Проверка
	if err != nil {
		t.Errorf("unexpected err: %v", err)
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
}
