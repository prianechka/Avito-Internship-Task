package service_handler

import (
	ac "Avito-Internship-Task/internal/app/balance_service_app/account/account_controller"
	"Avito-Internship-Task/internal/app/balance_service_app/account/account_repo"
	"Avito-Internship-Task/internal/app/balance_service_app/handlers/models"
	"Avito-Internship-Task/internal/app/balance_service_app/handlers/service_handler/request_models"
	"Avito-Internship-Task/internal/app/balance_service_app/manager"
	"Avito-Internship-Task/internal/app/balance_service_app/order"
	oc "Avito-Internship-Task/internal/app/balance_service_app/order/order_controller"
	"Avito-Internship-Task/internal/app/balance_service_app/order/order_repo"
	"Avito-Internship-Task/internal/app/balance_service_app/transaction"
	tc "Avito-Internship-Task/internal/app/balance_service_app/transaction/transaction_controller"
	"Avito-Internship-Task/internal/app/balance_service_app/transaction/transaction_repo"
	"bytes"
	"encoding/json"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestHandlerGetAccountBalance проверяет, что сервер правильно отвечает на запрос по количеству денег на балансе
func TestHandlerBuyServiceSuccess(t *testing.T) {
	var (
		userID             int     = 1
		orderID            int     = 1
		serviceID          int     = 1
		sum                float64 = 200
		balance            float64 = 400
		comment                    = "Всё хорошо!"
		expectedStatusCode         = http.StatusOK
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

	accountRepo := account_repo.NewAccountRepo(accountDB)
	accountController := ac.CreateNewAccountController(accountRepo)

	orderRepo := order_repo.NewOrderRepo(orderDB)
	orderController := oc.CreateNewOrderController(orderRepo)

	transactionRepo := transaction_repo.NewTransactionRepo(transactionDB)
	transactionController := tc.CreateNewTransactionController(transactionRepo)

	serverManager := manager.CreateNewManager(accountController, orderController, transactionController)

	serviceHandler := CreateServiceHandler(serverManager)
	ts := httptest.NewServer(http.HandlerFunc(serviceHandler.BuyService))
	defer ts.Close()

	bodyParams := request_models.BuyServiceMessage{
		UserID:    userID,
		OrderID:   orderID,
		ServiceID: serviceID,
		Sum:       sum,
		Comment:   comment,
	}
	reqBody, _ := json.Marshal(bodyParams)

	searcherReq, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(reqBody))

	r, _ := ts.Client().Do(searcherReq)

	// Проверка
	if err != nil {
		t.Errorf("unexpected err: %v", err)
		return
	}

	msg := models.ShortResponseMessage{}
	body, _ := ioutil.ReadAll(r.Body)

	unmarshalError := json.Unmarshal(body, &msg)
	if unmarshalError != nil {
		t.Errorf("unexpected error: %v", unmarshalError)
		return
	}

	if r.StatusCode != expectedStatusCode {
		t.Errorf("unexpected status code: %d %v", r.StatusCode, msg.Comment)
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

// TestHandlerBuyServiceNotAccExist проверяет, что покупки не произойдёт, если аккаунта не существует, и сервер вернет 400
func TestHandlerBuyServiceNotAccExist(t *testing.T) {
	var (
		userID             int     = 1
		orderID            int     = 1
		serviceID          int     = 1
		sum                float64 = 200
		comment                    = "Всё хорошо!"
		expectedStatusCode         = http.StatusUnauthorized
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

	accountRepo := account_repo.NewAccountRepo(accountDB)
	accountController := ac.CreateNewAccountController(accountRepo)

	orderRepo := order_repo.NewOrderRepo(orderDB)
	orderController := oc.CreateNewOrderController(orderRepo)

	transactionRepo := transaction_repo.NewTransactionRepo(transactionDB)
	transactionController := tc.CreateNewTransactionController(transactionRepo)

	serverManager := manager.CreateNewManager(accountController, orderController, transactionController)

	serviceHandler := CreateServiceHandler(serverManager)
	ts := httptest.NewServer(http.HandlerFunc(serviceHandler.BuyService))
	defer ts.Close()

	bodyParams := request_models.BuyServiceMessage{
		UserID:    userID,
		OrderID:   orderID,
		ServiceID: serviceID,
		Sum:       sum,
		Comment:   comment,
	}
	reqBody, _ := json.Marshal(bodyParams)

	searcherReq, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(reqBody))

	r, _ := ts.Client().Do(searcherReq)

	// Проверка
	if err != nil {
		t.Errorf("unexpected err: %v", err)
		return
	}

	msg := models.ShortResponseMessage{}
	body, _ := ioutil.ReadAll(r.Body)

	unmarshalError := json.Unmarshal(body, &msg)
	if unmarshalError != nil {
		t.Errorf("unexpected error: %v", unmarshalError)
		return
	}

	if r.StatusCode != expectedStatusCode {
		t.Errorf("unexpected status code: %d %v", r.StatusCode, msg.Comment)
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

// TestHandlerBuyServiceNotEnoughMoneyErr проверяет, что покупки не произойдёт, если аккаунта не существует, и сервер вернет 400
func TestHandlerBuyServiceNotEnoughMoneyErr(t *testing.T) {
	var (
		userID             int     = 1
		orderID            int     = 1
		serviceID          int     = 1
		sum                float64 = 400
		balance            float64 = 200
		comment                    = "Всё хорошо!"
		expectedStatusCode         = http.StatusUnprocessableEntity
	)

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

	accountRepo := account_repo.NewAccountRepo(accountDB)
	accountController := ac.CreateNewAccountController(accountRepo)

	orderRepo := order_repo.NewOrderRepo(orderDB)
	orderController := oc.CreateNewOrderController(orderRepo)

	transactionRepo := transaction_repo.NewTransactionRepo(transactionDB)
	transactionController := tc.CreateNewTransactionController(transactionRepo)

	serverManager := manager.CreateNewManager(accountController, orderController, transactionController)

	serviceHandler := CreateServiceHandler(serverManager)
	ts := httptest.NewServer(http.HandlerFunc(serviceHandler.BuyService))
	defer ts.Close()

	bodyParams := request_models.BuyServiceMessage{
		UserID:    userID,
		OrderID:   orderID,
		ServiceID: serviceID,
		Sum:       sum,
		Comment:   comment,
	}
	reqBody, _ := json.Marshal(bodyParams)

	searcherReq, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(reqBody))

	r, _ := ts.Client().Do(searcherReq)

	// Проверка
	if err != nil {
		t.Errorf("unexpected err: %v", err)
		return
	}

	msg := models.ShortResponseMessage{}
	body, _ := ioutil.ReadAll(r.Body)

	unmarshalError := json.Unmarshal(body, &msg)
	if unmarshalError != nil {
		t.Errorf("unexpected error: %v", unmarshalError)
		return
	}

	if r.StatusCode != expectedStatusCode {
		t.Errorf("unexpected status code: %d %v", r.StatusCode, msg.Comment)
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

// TestHandlerAcceptBuySuccess проверяет, что подтверждение покупки осуществлено корректно
func TestHandlerAcceptBuySuccess(t *testing.T) {
	var (
		userID             int     = 1
		orderID            int     = 1
		serviceID          int     = 1
		sum                float64 = 200
		comment                    = "Всё хорошо!"
		expectedStatusCode         = http.StatusOK
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

	accountRepo := account_repo.NewAccountRepo(accountDB)
	accountController := ac.CreateNewAccountController(accountRepo)

	orderRepo := order_repo.NewOrderRepo(orderDB)
	orderController := oc.CreateNewOrderController(orderRepo)

	transactionRepo := transaction_repo.NewTransactionRepo(transactionDB)
	transactionController := tc.CreateNewTransactionController(transactionRepo)

	serverManager := manager.CreateNewManager(accountController, orderController, transactionController)

	serviceHandler := CreateServiceHandler(serverManager)
	ts := httptest.NewServer(http.HandlerFunc(serviceHandler.AcceptService))
	defer ts.Close()

	bodyParams := request_models.AcceptServiceMessage{
		UserID:    userID,
		OrderID:   orderID,
		ServiceID: serviceID,
	}
	reqBody, _ := json.Marshal(bodyParams)

	searcherReq, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(reqBody))

	r, _ := ts.Client().Do(searcherReq)

	// Проверка
	if err != nil {
		t.Errorf("unexpected err: %v", err)
		return
	}

	msg := models.ShortResponseMessage{}
	body, _ := ioutil.ReadAll(r.Body)

	unmarshalError := json.Unmarshal(body, &msg)
	if unmarshalError != nil {
		t.Errorf("unexpected error: %v", unmarshalError)
		return
	}

	if r.StatusCode != expectedStatusCode {
		t.Errorf("unexpected status code: %d %v", r.StatusCode, msg.Comment)
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

// TestHandlerAcceptBuyError проверяет, что если заказа не существует, то вернется 404.
func TestHandlerAcceptBuyError(t *testing.T) {
	var (
		userID             int = 1
		orderID            int = 1
		serviceID          int = 1
		expectedStatusCode     = http.StatusNotFound
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

	accountRepo := account_repo.NewAccountRepo(accountDB)
	accountController := ac.CreateNewAccountController(accountRepo)

	orderRepo := order_repo.NewOrderRepo(orderDB)
	orderController := oc.CreateNewOrderController(orderRepo)

	transactionRepo := transaction_repo.NewTransactionRepo(transactionDB)
	transactionController := tc.CreateNewTransactionController(transactionRepo)

	serverManager := manager.CreateNewManager(accountController, orderController, transactionController)

	serviceHandler := CreateServiceHandler(serverManager)
	ts := httptest.NewServer(http.HandlerFunc(serviceHandler.AcceptService))
	defer ts.Close()

	bodyParams := request_models.AcceptServiceMessage{
		UserID:    userID,
		OrderID:   orderID,
		ServiceID: serviceID,
	}
	reqBody, _ := json.Marshal(bodyParams)

	searcherReq, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(reqBody))

	r, _ := ts.Client().Do(searcherReq)

	// Проверка
	if err != nil {
		t.Errorf("unexpected err: %v", err)
		return
	}

	msg := models.ShortResponseMessage{}
	body, _ := ioutil.ReadAll(r.Body)

	unmarshalError := json.Unmarshal(body, &msg)
	if unmarshalError != nil {
		t.Errorf("unexpected error: %v", unmarshalError)
		return
	}

	if r.StatusCode != expectedStatusCode {
		t.Errorf("unexpected status code: %d %v", r.StatusCode, msg.Comment)
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

// TestHandlerAcceptBuyWrongStatusError проверяет, что если статус неверный, то вернется 403.
func TestHandlerAcceptBuyWrongStatusError(t *testing.T) {
	var (
		userID             int     = 1
		orderID            int     = 1
		serviceID          int     = 1
		sum                float64 = 200
		comment                    = "Всё хорошо!"
		expectedStatusCode         = http.StatusForbidden
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

	accountRepo := account_repo.NewAccountRepo(accountDB)
	accountController := ac.CreateNewAccountController(accountRepo)

	orderRepo := order_repo.NewOrderRepo(orderDB)
	orderController := oc.CreateNewOrderController(orderRepo)

	transactionRepo := transaction_repo.NewTransactionRepo(transactionDB)
	transactionController := tc.CreateNewTransactionController(transactionRepo)

	serverManager := manager.CreateNewManager(accountController, orderController, transactionController)

	serviceHandler := CreateServiceHandler(serverManager)
	ts := httptest.NewServer(http.HandlerFunc(serviceHandler.AcceptService))
	defer ts.Close()

	bodyParams := request_models.AcceptServiceMessage{
		UserID:    userID,
		OrderID:   orderID,
		ServiceID: serviceID,
	}
	reqBody, _ := json.Marshal(bodyParams)

	searcherReq, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(reqBody))

	r, _ := ts.Client().Do(searcherReq)

	// Проверка
	if err != nil {
		t.Errorf("unexpected err: %v", err)
		return
	}

	msg := models.ShortResponseMessage{}
	body, _ := ioutil.ReadAll(r.Body)

	unmarshalError := json.Unmarshal(body, &msg)
	if unmarshalError != nil {
		t.Errorf("unexpected error: %v", unmarshalError)
		return
	}

	if r.StatusCode != expectedStatusCode {
		t.Errorf("unexpected status code: %d %v", r.StatusCode, msg.Comment)
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

// TestHandlerRefuseServiceSuccess проверяет, что возврат произойдёт успешно и вернется 200.
func TestHandlerRefuseServiceSuccess(t *testing.T) {
	var (
		userID             int     = 1
		orderID            int     = 1
		serviceID          int     = 1
		sum                float64 = 200
		balance            float64 = 400
		comment                    = "Всё хорошо!"
		expectedStatusCode         = http.StatusOK
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
			newTransaction.Sum, sqlmock.AnyArg(), newTransaction.ActionComments, sqlmock.AnyArg()).
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

	accountRepo := account_repo.NewAccountRepo(accountDB)
	accountController := ac.CreateNewAccountController(accountRepo)

	orderRepo := order_repo.NewOrderRepo(orderDB)
	orderController := oc.CreateNewOrderController(orderRepo)

	transactionRepo := transaction_repo.NewTransactionRepo(transactionDB)
	transactionController := tc.CreateNewTransactionController(transactionRepo)

	serverManager := manager.CreateNewManager(accountController, orderController, transactionController)

	serviceHandler := CreateServiceHandler(serverManager)
	ts := httptest.NewServer(http.HandlerFunc(serviceHandler.RefuseService))
	defer ts.Close()

	bodyParams := request_models.AcceptServiceMessage{
		UserID:    userID,
		OrderID:   orderID,
		ServiceID: serviceID,
	}
	reqBody, _ := json.Marshal(bodyParams)

	searcherReq, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(reqBody))

	r, _ := ts.Client().Do(searcherReq)

	// Проверка
	if err != nil {
		t.Errorf("unexpected err: %v", err)
		return
	}

	msg := models.ShortResponseMessage{}
	body, _ := ioutil.ReadAll(r.Body)

	unmarshalError := json.Unmarshal(body, &msg)
	if unmarshalError != nil {
		t.Errorf("unexpected error: %v", unmarshalError)
		return
	}

	if r.StatusCode != expectedStatusCode {
		t.Errorf("unexpected status code: %d %v", r.StatusCode, msg.Comment)
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

// TestHandlerRefuseServiceWrongStatusError проверяет, что если у заказа неправильный статус, то вернётся 403.
func TestHandlerRefuseServiceWrongStatusError(t *testing.T) {
	var (
		userID             int     = 1
		orderID            int     = 1
		serviceID          int     = 1
		sum                float64 = 200
		comment                    = "Всё хорошо!"
		expectedStatusCode         = http.StatusForbidden
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

	accountRepo := account_repo.NewAccountRepo(accountDB)
	accountController := ac.CreateNewAccountController(accountRepo)

	orderRepo := order_repo.NewOrderRepo(orderDB)
	orderController := oc.CreateNewOrderController(orderRepo)

	transactionRepo := transaction_repo.NewTransactionRepo(transactionDB)
	transactionController := tc.CreateNewTransactionController(transactionRepo)

	serverManager := manager.CreateNewManager(accountController, orderController, transactionController)

	serviceHandler := CreateServiceHandler(serverManager)
	ts := httptest.NewServer(http.HandlerFunc(serviceHandler.RefuseService))
	defer ts.Close()

	bodyParams := request_models.AcceptServiceMessage{
		UserID:    userID,
		OrderID:   orderID,
		ServiceID: serviceID,
	}
	reqBody, _ := json.Marshal(bodyParams)

	searcherReq, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(reqBody))

	r, _ := ts.Client().Do(searcherReq)

	// Проверка
	if err != nil {
		t.Errorf("unexpected err: %v", err)
		return
	}

	msg := models.ShortResponseMessage{}
	body, _ := ioutil.ReadAll(r.Body)

	unmarshalError := json.Unmarshal(body, &msg)
	if unmarshalError != nil {
		t.Errorf("unexpected error: %v", unmarshalError)
		return
	}

	if r.StatusCode != expectedStatusCode {
		t.Errorf("unexpected status code: %d %v", r.StatusCode, msg.Comment)
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

// TestHandlerRefuseServiceOrderNoExistError проверяет, что если аккаунта не существует, то вернется 404.
func TestHandlerRefuseServiceOrderNoExistError(t *testing.T) {
	var (
		userID             int = 1
		orderID            int = 1
		serviceID          int = 1
		expectedStatusCode     = http.StatusNotFound
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

	accountRepo := account_repo.NewAccountRepo(accountDB)
	accountController := ac.CreateNewAccountController(accountRepo)

	orderRepo := order_repo.NewOrderRepo(orderDB)
	orderController := oc.CreateNewOrderController(orderRepo)

	transactionRepo := transaction_repo.NewTransactionRepo(transactionDB)
	transactionController := tc.CreateNewTransactionController(transactionRepo)

	serverManager := manager.CreateNewManager(accountController, orderController, transactionController)

	serviceHandler := CreateServiceHandler(serverManager)
	ts := httptest.NewServer(http.HandlerFunc(serviceHandler.RefuseService))
	defer ts.Close()

	bodyParams := request_models.AcceptServiceMessage{
		UserID:    userID,
		OrderID:   orderID,
		ServiceID: serviceID,
	}
	reqBody, _ := json.Marshal(bodyParams)

	searcherReq, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(reqBody))

	r, _ := ts.Client().Do(searcherReq)

	// Проверка
	if err != nil {
		t.Errorf("unexpected err: %v", err)
		return
	}

	msg := models.ShortResponseMessage{}
	body, _ := ioutil.ReadAll(r.Body)

	unmarshalError := json.Unmarshal(body, &msg)
	if unmarshalError != nil {
		t.Errorf("unexpected error: %v", unmarshalError)
		return
	}

	if r.StatusCode != expectedStatusCode {
		t.Errorf("unexpected status code: %d %v", r.StatusCode, msg.Comment)
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
