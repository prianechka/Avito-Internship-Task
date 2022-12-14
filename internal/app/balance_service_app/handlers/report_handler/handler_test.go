package report_handler

import (
	ac "Avito-Internship-Task/internal/app/balance_service_app/account/account_controller"
	"Avito-Internship-Task/internal/app/balance_service_app/account/account_repo"
	"Avito-Internship-Task/internal/app/balance_service_app/handlers/models"
	"Avito-Internship-Task/internal/app/balance_service_app/managers/report_manager"
	"Avito-Internship-Task/internal/app/balance_service_app/order"
	oc "Avito-Internship-Task/internal/app/balance_service_app/order/order_controller"
	"Avito-Internship-Task/internal/app/balance_service_app/order/order_repo"
	"Avito-Internship-Task/internal/app/balance_service_app/report"
	"Avito-Internship-Task/internal/app/balance_service_app/transaction"
	tc "Avito-Internship-Task/internal/app/balance_service_app/transaction/transaction_controller"
	"Avito-Internship-Task/internal/app/balance_service_app/transaction/transaction_repo"
	"encoding/json"
	"fmt"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestHandlerGetFinanceReportSuccess показывает, что программа вернёт ссылку на отчёт при верных данных.
func TestHandlerGetFinanceReportSuccess(t *testing.T) {
	var (
		year               = 2022
		month              = 2
		expectedStatusCode = http.StatusOK
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

	testManager := report_manager.CreateNewReportManager(accountController, orderController, transactionController)
	handler := CreateReportHandler(testManager)

	ts := httptest.NewServer(http.HandlerFunc(handler.GetFinanceReport))
	defer ts.Close()

	searcherReq, err := http.NewRequest("GET", ts.URL+fmt.Sprintf("?month=%d&year=%d", month, year), nil)
	r, handleError := ts.Client().Do(searcherReq)

	// Проверка
	if handleError != nil {
		t.Errorf("unexpected err: %v", err)
		return
	}

	msg := models.FinanceReportResponseMessage{}
	body, _ := ioutil.ReadAll(r.Body)

	unmarshalError := json.Unmarshal(body, &msg)
	if unmarshalError != nil {
		t.Errorf("unexpected error: %v", unmarshalError)
		return
	}

	if r.StatusCode != expectedStatusCode {
		t.Errorf("unexpected status code: %d %v", r.StatusCode, expectedStatusCode)
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

// TestHandlerGetFinanceReportBadMonthError показывает, что программа вернёт ошибку, если месяц указан неправильно.
func TestHandlerGetFinanceReportBadMonthError(t *testing.T) {
	var (
		year               = 2022
		month              = 13
		expectedStatusCode = http.StatusBadRequest
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

	// Создание объектов
	accountRepo := account_repo.NewAccountRepo(accountDB)
	accountController := ac.CreateNewAccountController(accountRepo)

	orderRepo := order_repo.NewOrderRepo(orderDB)
	orderController := oc.CreateNewOrderController(orderRepo)

	transactionRepo := transaction_repo.NewTransactionRepo(transactionDB)
	transactionController := tc.CreateNewTransactionController(transactionRepo)

	testManager := report_manager.CreateNewReportManager(accountController, orderController, transactionController)
	handler := CreateReportHandler(testManager)

	ts := httptest.NewServer(http.HandlerFunc(handler.GetFinanceReport))
	defer ts.Close()

	searcherReq, err := http.NewRequest("GET", ts.URL+fmt.Sprintf("?month=%d&year=%d", month, year), nil)
	r, handleError := ts.Client().Do(searcherReq)

	// Проверка
	if handleError != nil {
		t.Errorf("unexpected err: %v", err)
		return
	}

	msg := models.FinanceReportResponseMessage{}
	body, _ := ioutil.ReadAll(r.Body)

	unmarshalError := json.Unmarshal(body, &msg)
	if unmarshalError != nil {
		t.Errorf("unexpected error: %v", unmarshalError)
		return
	}

	if r.StatusCode != expectedStatusCode {
		t.Errorf("unexpected status code: %d %v", r.StatusCode, expectedStatusCode)
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

// TestHandlerGetFinanceReportBadYearError показывает, что программа вернёт ошибку, если год указан неправильно.
func TestHandlerGetFinanceReportBadYearError(t *testing.T) {
	var (
		year               = 2025
		month              = 11
		expectedStatusCode = http.StatusBadRequest
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

	// Создание объектов
	accountRepo := account_repo.NewAccountRepo(accountDB)
	accountController := ac.CreateNewAccountController(accountRepo)

	orderRepo := order_repo.NewOrderRepo(orderDB)
	orderController := oc.CreateNewOrderController(orderRepo)

	transactionRepo := transaction_repo.NewTransactionRepo(transactionDB)
	transactionController := tc.CreateNewTransactionController(transactionRepo)

	testManager := report_manager.CreateNewReportManager(accountController, orderController, transactionController)
	handler := CreateReportHandler(testManager)

	ts := httptest.NewServer(http.HandlerFunc(handler.GetFinanceReport))
	defer ts.Close()

	searcherReq, err := http.NewRequest("GET", ts.URL+fmt.Sprintf("?month=%d&year=%d", month, year), nil)
	r, handleError := ts.Client().Do(searcherReq)

	// Проверка
	if handleError != nil {
		t.Errorf("unexpected err: %v", err)
		return
	}

	msg := models.FinanceReportResponseMessage{}
	body, _ := ioutil.ReadAll(r.Body)

	unmarshalError := json.Unmarshal(body, &msg)
	if unmarshalError != nil {
		t.Errorf("unexpected error: %v", unmarshalError)
		return
	}

	if r.StatusCode != expectedStatusCode {
		t.Errorf("unexpected status code: %d %v", r.StatusCode, expectedStatusCode)
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

// TestHandlerGetFinanceReportNoYearError показывает, что программа вернёт ошибку, если год не указан.
func TestHandlerGetFinanceReportNoYearError(t *testing.T) {
	var (
		month              = 11
		expectedStatusCode = http.StatusBadRequest
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

	// Создание объектов
	accountRepo := account_repo.NewAccountRepo(accountDB)
	accountController := ac.CreateNewAccountController(accountRepo)

	orderRepo := order_repo.NewOrderRepo(orderDB)
	orderController := oc.CreateNewOrderController(orderRepo)

	transactionRepo := transaction_repo.NewTransactionRepo(transactionDB)
	transactionController := tc.CreateNewTransactionController(transactionRepo)

	testManager := report_manager.CreateNewReportManager(accountController, orderController, transactionController)
	handler := CreateReportHandler(testManager)

	ts := httptest.NewServer(http.HandlerFunc(handler.GetFinanceReport))
	defer ts.Close()

	searcherReq, err := http.NewRequest("GET", ts.URL+fmt.Sprintf("?month=%d", month), nil)
	r, handleError := ts.Client().Do(searcherReq)

	// Проверка
	if handleError != nil {
		t.Errorf("unexpected err: %v", err)
		return
	}

	msg := models.FinanceReportResponseMessage{}
	body, _ := ioutil.ReadAll(r.Body)

	unmarshalError := json.Unmarshal(body, &msg)
	if unmarshalError != nil {
		t.Errorf("unexpected error: %v", unmarshalError)
		return
	}

	if r.StatusCode != expectedStatusCode {
		t.Errorf("unexpected status code: %d %v", r.StatusCode, expectedStatusCode)
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

// TestHandlerGetFinanceReportNoMonthError показывает, что программа вернёт ошибку, если месяц не указан.
func TestHandlerGetFinanceReportNoMonthError(t *testing.T) {
	var (
		year               = 2022
		expectedStatusCode = http.StatusBadRequest
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

	// Создание объектов
	accountRepo := account_repo.NewAccountRepo(accountDB)
	accountController := ac.CreateNewAccountController(accountRepo)

	orderRepo := order_repo.NewOrderRepo(orderDB)
	orderController := oc.CreateNewOrderController(orderRepo)

	transactionRepo := transaction_repo.NewTransactionRepo(transactionDB)
	transactionController := tc.CreateNewTransactionController(transactionRepo)

	testManager := report_manager.CreateNewReportManager(accountController, orderController, transactionController)
	handler := CreateReportHandler(testManager)

	ts := httptest.NewServer(http.HandlerFunc(handler.GetFinanceReport))
	defer ts.Close()

	searcherReq, err := http.NewRequest("GET", ts.URL+fmt.Sprintf("?year=%d", year), nil)
	r, handleError := ts.Client().Do(searcherReq)

	// Проверка
	if handleError != nil {
		t.Errorf("unexpected err: %v", err)
		return
	}

	msg := models.FinanceReportResponseMessage{}
	body, _ := ioutil.ReadAll(r.Body)

	unmarshalError := json.Unmarshal(body, &msg)
	if unmarshalError != nil {
		t.Errorf("unexpected error: %v", unmarshalError)
		return
	}

	if r.StatusCode != expectedStatusCode {
		t.Errorf("unexpected status code: %d %v", r.StatusCode, expectedStatusCode)
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

// TestHandlerGetUserReport показывает, что программа вернёт отчёт по пользователю при правильном вводе данных.
func TestHandlerGetUserReport(t *testing.T) {
	var (
		userID                     = 1
		serviceID                  = 1
		sum                float64 = 100
		balance            float64 = 100
		comment                    = "Хорошо"
		orderBy                    = "id"
		limit                      = 2
		offset                     = 0
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

	accountMock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(userID).
		WillReturnRows(accountFirstRows)

	// Подготовка БД для таблицы с транзакциями
	transactionDB, transactionMock, createTransactDBErr := sqlmock.New()
	if createTransactDBErr != nil {
		t.Fatalf("cant create mock: %s", createTransactDBErr)
	}
	defer transactionDB.Close()

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

	transactionMock.ExpectQuery("SELECT transactionID, userID, transactionType, sum, time," +
		" actionComments, addComments FROM balanceApp.transactions WHERE userID").
		WillReturnRows(rows).WillReturnError(nil)

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

	testManager := report_manager.CreateNewReportManager(accountController, orderController, transactionController)
	handler := CreateReportHandler(testManager)

	ts := httptest.NewServer(http.HandlerFunc(handler.GetUserReport))
	defer ts.Close()

	searcherReq, err := http.NewRequest("GET", ts.URL+fmt.Sprintf("?userID=%d&orderBy=%s&limit=%d&offset=%d",
		userID, orderBy, limit, offset), nil)
	r, handleError := ts.Client().Do(searcherReq)

	// Проверка
	if handleError != nil {
		t.Errorf("unexpected err: %v", err)
		return
	}

	msg := models.UserReportResponseMessage{}
	body, _ := ioutil.ReadAll(r.Body)

	unmarshalError := json.Unmarshal(body, &msg)
	if unmarshalError != nil {
		t.Errorf("unexpected error: %v", unmarshalError)
		return
	}

	if r.StatusCode != expectedStatusCode {
		t.Errorf("unexpected status code: %d %v", r.StatusCode, expectedStatusCode)
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

// TestHandlerGetUserReportWithBadOrderField показывает, что программа вернёт отчёт даже с неправильным полем order.
func TestHandlerGetUserReportWithBadOrderField(t *testing.T) {
	var (
		userID                     = 1
		serviceID                  = 1
		sum                float64 = 100
		balance            float64 = 100
		comment                    = "Хорошо"
		orderBy                    = "temp"
		limit                      = 2
		offset                     = 0
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

	accountMock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(userID).
		WillReturnRows(accountFirstRows)

	// Подготовка БД для таблицы с транзакциями
	transactionDB, transactionMock, createTransactDBErr := sqlmock.New()
	if createTransactDBErr != nil {
		t.Fatalf("cant create mock: %s", createTransactDBErr)
	}
	defer transactionDB.Close()

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

	transactionMock.ExpectQuery("SELECT transactionID, userID, transactionType, sum, time," +
		" actionComments, addComments FROM balanceApp.transactions WHERE userID").
		WillReturnRows(rows).WillReturnError(nil)

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

	testManager := report_manager.CreateNewReportManager(accountController, orderController, transactionController)
	handler := CreateReportHandler(testManager)

	ts := httptest.NewServer(http.HandlerFunc(handler.GetUserReport))
	defer ts.Close()

	searcherReq, err := http.NewRequest("GET", ts.URL+fmt.Sprintf("?userID=%d&orderBy=%s&limit=%d&offset=%d",
		userID, orderBy, limit, offset), nil)
	r, handleError := ts.Client().Do(searcherReq)

	// Проверка
	if handleError != nil {
		t.Errorf("unexpected err: %v", err)
		return
	}

	msg := models.UserReportResponseMessage{}
	body, _ := ioutil.ReadAll(r.Body)

	unmarshalError := json.Unmarshal(body, &msg)
	if unmarshalError != nil {
		t.Errorf("unexpected error: %v", unmarshalError)
		return
	}

	if r.StatusCode != expectedStatusCode {
		t.Errorf("unexpected status code: %d %v", r.StatusCode, expectedStatusCode)
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

// TestHandlerGetUserReportBadUserError показывает, что программа вернёт ошибку, если пользователь выбран неверно.
func TestHandlerGetUserReportBadUserError(t *testing.T) {
	var (
		userID             = 1
		orderBy            = "id"
		limit              = 2
		offset             = 0
		expectedStatusCode = http.StatusUnauthorized
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

	testManager := report_manager.CreateNewReportManager(accountController, orderController, transactionController)
	handler := CreateReportHandler(testManager)

	ts := httptest.NewServer(http.HandlerFunc(handler.GetUserReport))
	defer ts.Close()

	searcherReq, err := http.NewRequest("GET", ts.URL+fmt.Sprintf("?userID=%d&orderBy=%s&limit=%d&offset=%d",
		userID, orderBy, limit, offset), nil)
	r, handleError := ts.Client().Do(searcherReq)

	// Проверка
	if handleError != nil {
		t.Errorf("unexpected err: %v", err)
		return
	}

	msg := models.UserReportResponseMessage{}
	body, _ := ioutil.ReadAll(r.Body)

	unmarshalError := json.Unmarshal(body, &msg)
	if unmarshalError != nil {
		t.Errorf("unexpected error: %v", unmarshalError)
		return
	}

	if r.StatusCode != expectedStatusCode {
		t.Errorf("unexpected status code: %d %v", r.StatusCode, expectedStatusCode)
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
