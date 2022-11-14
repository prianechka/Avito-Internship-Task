package account_handler

import (
	ac "Avito-Internship-Task/internal/app/balance_service_app/account/account_controller"
	"Avito-Internship-Task/internal/app/balance_service_app/account/account_repo"
	"Avito-Internship-Task/internal/app/balance_service_app/handlers/account_handler/messages"
	"Avito-Internship-Task/internal/app/balance_service_app/handlers/response"
	"Avito-Internship-Task/internal/app/balance_service_app/manager"
	oc "Avito-Internship-Task/internal/app/balance_service_app/order/order_controller"
	"Avito-Internship-Task/internal/app/balance_service_app/order/order_repo"
	"Avito-Internship-Task/internal/app/balance_service_app/transaction"
	tc "Avito-Internship-Task/internal/app/balance_service_app/transaction/transaction_controller"
	"Avito-Internship-Task/internal/app/balance_service_app/transaction/transaction_repo"
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

// TestGetAccountBalance проверяет, что сервер правильно отвечает на запрос по количеству денег на балансе
func TestGetAccountBalance(t *testing.T) {

	// Подготовка БД к тестам
	var (
		userID             int64   = 1
		balance            float64 = 200
		expectedStatusCode         = http.StatusOK
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

	accountRepo := account_repo.NewAccountRepo(accountDB)
	accountController := ac.CreateNewAccountController(accountRepo)

	orderRepo := order_repo.NewOrderRepo(orderDB)
	orderController := oc.CreateNewOrderController(orderRepo)

	transactionRepo := transaction_repo.NewTransactionRepo(transactionDB)
	transactionController := tc.CreateNewTransactionController(transactionRepo)

	serverManager := manager.CreateNewManager(accountController, orderController, transactionController)

	accountHandler := CreateAccountHandler(serverManager)
	ts := httptest.NewServer(http.HandlerFunc(accountHandler.GetBalance))
	defer ts.Close()

	searcherReq, err := http.NewRequest("GET", ts.URL+fmt.Sprintf("?id=%d", userID), nil)
	//searcherReq := httptest.NewRequest("GET", "/balance/{id}", nil)
	r, err := ts.Client().Do(searcherReq)

	// Проверка
	if err != nil {
		t.Errorf("unexpected err: %v", err)
		return
	}

	msg := response.BalanceResponseMessage{}
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

	if !reflect.DeepEqual(balance, msg.Balance) {
		t.Errorf("results not match, want %v, have %v", balance, msg.Balance)
		return
	}
}

// TestGetAccountBalance проверяет, что аккаунт неправильный, то выдаст 400
func TestGetAccountBalanceBadAccount(t *testing.T) {

	// Подготовка БД к тестам
	var (
		userID             int64 = 1
		expectedStatusCode       = http.StatusBadRequest
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

	accountHandler := CreateAccountHandler(serverManager)
	ts := httptest.NewServer(http.HandlerFunc(accountHandler.GetBalance))
	defer ts.Close()

	searcherReq, err := http.NewRequest("GET", ts.URL+fmt.Sprintf("?id=%d", userID), nil)
	//searcherReq := httptest.NewRequest("GET", "/balance/{id}", nil)
	r, err := ts.Client().Do(searcherReq)

	// Проверка
	if err != nil {
		t.Errorf("unexpected err: %v", err)
		return
	}

	msg := response.BalanceResponseMessage{}
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

// TestRefillAccountSuccess проверяет, что сервер правильно отвечает на запрос по пополнению счёта
func TestRefillAccountSuccess(t *testing.T) {
	var (
		userID             int64   = 1
		sum                float64 = 200
		comment                    = "Всё хорошо!"
		expectedStatusCode         = http.StatusOK
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

	serverManager := manager.CreateNewManager(accountController, orderController, transactionController)

	accountHandler := CreateAccountHandler(serverManager)
	ts := httptest.NewServer(http.HandlerFunc(accountHandler.RefillBalance))
	defer ts.Close()

	bodyParams := messages.RefillParams{
		UserID:  userID,
		Sum:     sum,
		Comment: comment,
	}

	reqBody, _ := json.Marshal(bodyParams)

	searcherReq, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(reqBody))

	r, _ := ts.Client().Do(searcherReq)

	// Проверка
	if err != nil {
		t.Errorf("unexpected err: %v", err)
		return
	}

	msg := response.ShortResponseMessage{}
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

	if !reflect.DeepEqual(msg.Comment, "OK") {
		t.Errorf("results not match, want %v, have %v", "OK", msg.Comment)
		return
	}
}

// TestRefillAccountNotExistError проверяет, что если аккаунта не существует, он будет создан, и вернётся 200
func TestRefillAccountNotExist(t *testing.T) {
	var (
		userID             int64   = 1
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

	serverManager := manager.CreateNewManager(accountController, orderController, transactionController)

	accountHandler := CreateAccountHandler(serverManager)
	ts := httptest.NewServer(http.HandlerFunc(accountHandler.RefillBalance))
	defer ts.Close()

	bodyParams := messages.RefillParams{
		UserID:  userID,
		Sum:     sum,
		Comment: comment,
	}

	reqBody, _ := json.Marshal(bodyParams)

	searcherReq, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(reqBody))

	r, _ := ts.Client().Do(searcherReq)

	// Проверка
	if err != nil {
		t.Errorf("unexpected err: %v", err)
		return
	}

	msg := response.ShortResponseMessage{}
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

	if !reflect.DeepEqual(msg.Comment, "OK") {
		t.Errorf("results not match, want %v, have %v", "OK", msg.Comment)
		return
	}
}

// TestHandlerTransferSuccess проверяет, что перевод средств между двумя аккаунтами прошёл успешно и вернулось 200
func TestHandlerTransferSuccess(t *testing.T) {
	var (
		srcUserID          int64   = 1
		dstUserID          int64   = 2
		sum                float64 = 200
		balanceFirst       float64 = 400
		balanceSecond      float64 = 200
		comment                    = "Всё хорошо!"
		expectedStatusCode         = http.StatusOK
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

	accountMock.ExpectExec("UPDATE balanceApp.accounts SET amount = amoumt +").
		WithArgs(-sum, srcUserID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	accountMock.ExpectExec("UPDATE balanceApp.accounts SET amount = amoumt +").
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

	serverManager := manager.CreateNewManager(accountController, orderController, transactionController)

	accountHandler := CreateAccountHandler(serverManager)
	ts := httptest.NewServer(http.HandlerFunc(accountHandler.Transfer))
	defer ts.Close()

	bodyParams := messages.TransferMessage{
		SrcUserID: srcUserID,
		DstUserID: dstUserID,
		Sum:       sum,
		Comment:   comment,
	}

	reqBody, _ := json.Marshal(bodyParams)
	searcherReq, _ := http.NewRequest("POST", ts.URL, bytes.NewBuffer(reqBody))
	r, _ := ts.Client().Do(searcherReq)

	msg := response.ShortResponseMessage{}
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

// TestHandlerTransferAccNotExistError проверяет, что если аккаунта не существует, вернется 400
func TestHandlerTransferAccNotExistError(t *testing.T) {
	var (
		srcUserID          int64   = 1
		dstUserID          int64   = 2
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

	serverManager := manager.CreateNewManager(accountController, orderController, transactionController)

	accountHandler := CreateAccountHandler(serverManager)
	ts := httptest.NewServer(http.HandlerFunc(accountHandler.Transfer))
	defer ts.Close()

	bodyParams := messages.TransferMessage{
		SrcUserID: srcUserID,
		DstUserID: dstUserID,
		Sum:       sum,
		Comment:   comment,
	}

	reqBody, _ := json.Marshal(bodyParams)
	searcherReq, _ := http.NewRequest("POST", ts.URL, bytes.NewBuffer(reqBody))
	r, _ := ts.Client().Do(searcherReq)

	msg := response.ShortResponseMessage{}
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

// TestHandlerTransferNotEnoughMoneyError проверяет, что если не хватает денег для перевода, вернется 400
func TestHandlerTransferNotEnoughMoneyError(t *testing.T) {
	var (
		srcUserID          int64   = 1
		dstUserID          int64   = 2
		sum                float64 = 500
		balanceFirst       float64 = 400
		balanceSecond      float64 = 200
		comment                    = "Всё хорошо!"
		expectedStatusCode         = http.StatusBadRequest
	)
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

	serverManager := manager.CreateNewManager(accountController, orderController, transactionController)

	accountHandler := CreateAccountHandler(serverManager)
	ts := httptest.NewServer(http.HandlerFunc(accountHandler.Transfer))
	defer ts.Close()

	bodyParams := messages.TransferMessage{
		SrcUserID: srcUserID,
		DstUserID: dstUserID,
		Sum:       sum,
		Comment:   comment,
	}

	reqBody, _ := json.Marshal(bodyParams)
	searcherReq, _ := http.NewRequest("POST", ts.URL, bytes.NewBuffer(reqBody))
	r, _ := ts.Client().Do(searcherReq)

	msg := response.ShortResponseMessage{}
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
