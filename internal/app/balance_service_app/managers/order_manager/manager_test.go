package order_manager

import (
	ac "Avito-Internship-Task/internal/app/balance_service_app/account/account_controller"
	"Avito-Internship-Task/internal/app/balance_service_app/account/account_repo"
	"Avito-Internship-Task/internal/app/balance_service_app/order"
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

// TestBuyServiceSuccess проверяет, что сценарий покупки сервиса проходит успешно.
func TestBuyServiceSuccess(t *testing.T) {
	var (
		userID    int     = 1
		orderID   int     = 1
		serviceID int     = 1
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

	testManager := CreateNewOrderManager(accountController, orderController, transactionController)

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
		userID    int     = 1
		orderID   int     = 1
		serviceID int     = 1
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

	testManager := CreateNewOrderManager(accountController, orderController, transactionController)

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
		userID    int     = 1
		orderID   int     = 1
		serviceID int     = 1
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

	testManager := CreateNewOrderManager(accountController, orderController, transactionController)

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
		userID    int     = 1
		orderID   int     = 1
		serviceID int     = 1
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

	testManager := CreateNewOrderManager(accountController, orderController, transactionController)

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

	if !reflect.DeepEqual(curOrder.OrderState, int(order.FINISHED)) {
		t.Errorf("results not match, want %v, have %v", order.FINISHED, curOrder.OrderState)
		return
	}
}

// TestAcceptBuyError проверяет, что если заказа не существует, то не получится подтвердить покупку.
func TestAcceptBuyError(t *testing.T) {
	var (
		userID    int = 1
		orderID   int = 1
		serviceID int = 1
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

	testManager := CreateNewOrderManager(accountController, orderController, transactionController)

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
		userID    int     = 1
		orderID   int     = 1
		serviceID int     = 1
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

	testManager := CreateNewOrderManager(accountController, orderController, transactionController)

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
		userID    int     = 1
		orderID   int     = 1
		serviceID int     = 1
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

	testManager := CreateNewOrderManager(accountController, orderController, transactionController)

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
		userID    int     = 1
		orderID   int     = 1
		serviceID int     = 1
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

	testManager := CreateNewOrderManager(accountController, orderController, transactionController)

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
		userID    int = 1
		orderID   int = 1
		serviceID int = 1
		comment       = "Всё хорошо!"
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

	testManager := CreateNewOrderManager(accountController, orderController, transactionController)

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
