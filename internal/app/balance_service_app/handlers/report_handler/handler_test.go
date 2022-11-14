package report_handler

import (
	ac "Avito-Internship-Task/internal/app/balance_service_app/account/account_controller"
	"Avito-Internship-Task/internal/app/balance_service_app/account/account_repo"
	"Avito-Internship-Task/internal/app/balance_service_app/handlers/models"
	"Avito-Internship-Task/internal/app/balance_service_app/manager"
	oc "Avito-Internship-Task/internal/app/balance_service_app/order/order_controller"
	"Avito-Internship-Task/internal/app/balance_service_app/order/order_repo"
	"Avito-Internship-Task/internal/app/balance_service_app/report"
	tc "Avito-Internship-Task/internal/app/balance_service_app/transaction/transaction_controller"
	"Avito-Internship-Task/internal/app/balance_service_app/transaction/transaction_repo"
	"encoding/json"
	"fmt"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHandlerGetFinanceReportSuccess показывает, что программа вернёт ссылку на отчёт при верных данных.
func TestHandlerGetFinanceReportSuccess(t *testing.T) {
	var (
		year               int64 = 2022
		month              int64 = 2
		expectedStatusCode       = http.StatusOK
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

	testManager := manager.CreateNewManager(accountController, orderController, transactionController)
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
