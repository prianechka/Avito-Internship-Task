package order_controller

import (
	"Avito-Internship-Task/internal/app/balance_service_app/order"
	"Avito-Internship-Task/internal/app/balance_service_app/order/order_repo"
	"Avito-Internship-Task/internal/app/balance_service_app/report"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"reflect"
	"testing"
	"time"
)

// TestCheckOrderIsExist проверяет, что если заказ существует, то mocks ответит true
func TestCheckOrderIsExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	curTime := time.Now()
	var orderID int = 1
	var userID int = 1
	var serviceID int = 1
	var expectResult = true

	rows := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})
	expect := order.Order{orderID, userID, serviceID, 100, curTime, "Good", order.REGISTRATED}
	rows.AddRow(expect.OrderID, expect.UserID, expect.ServiceID, expect.OrderCost,
		expect.CreatingTime, expect.Comment, expect.OrderState)

	mock.
		ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(rows).WillReturnError(nil)

	repo := order_repo.NewOrderRepo(db)
	controller := CreateNewOrderController(repo)

	result, execErr := controller.CheckOrderIsExist(orderID, userID, serviceID)

	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	if !reflect.DeepEqual(result, expectResult) {
		t.Errorf("results not match, want %v, have %v", result, expectResult)
		return
	}
}

// TestCheckOrderIsNotExist проверяет, что если заказ не существует, то mocks ответит false
func TestCheckOrderIsNotExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var orderID int = 1
	var userID int = 1
	var serviceID int = 1
	var expectResult = false

	rows := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})

	mock.
		ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(rows).WillReturnError(nil)

	repo := order_repo.NewOrderRepo(db)
	controller := CreateNewOrderController(repo)

	result, execErr := controller.CheckOrderIsExist(orderID, userID, serviceID)

	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	if !reflect.DeepEqual(result, expectResult) {
		t.Errorf("results not match, want %v, have %v", result, expectResult)
		return
	}
}

// TestCreateOrderSuccess создаёт новый заказ, которого ещё не существует
func TestCreateOrderSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var orderID int = 1
	var userID int = 1
	var serviceID int = 1
	var sum float64 = 100
	var comment = ""

	rows := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})

	mock.ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(rows).WillReturnError(nil)

	mock.
		ExpectExec("INSERT INTO balanceApp.orders").
		WithArgs(orderID, userID, serviceID, sum, sqlmock.AnyArg(), comment, order.REGISTRATED).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := order_repo.NewOrderRepo(db)
	controller := CreateNewOrderController(repo)

	execErr := controller.CreateNewOrder(orderID, userID, serviceID, sum, comment)

	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}
}

// TestCreateOrderError не создаёт новый заказ, так как уже существует заказ с такими же параметрами
func TestCreateOrderError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	curTime := time.Now()
	var orderID int = 1
	var userID int = 1
	var serviceID int = 1
	var sum float64 = 100
	var comment = ""

	rows := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})
	expect := order.Order{orderID, userID, serviceID, 100, curTime, "Good", 1}
	rows.AddRow(expect.OrderID, expect.UserID, expect.ServiceID, expect.OrderCost,
		expect.CreatingTime, expect.Comment, expect.OrderState)

	mock.ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(rows).WillReturnError(nil)

	repo := order_repo.NewOrderRepo(db)
	controller := CreateNewOrderController(repo)

	execErr := controller.CreateNewOrder(orderID, userID, serviceID, sum, comment)

	if execErr != OrderIsAlreadyExist {
		t.Errorf("unexpected err: %v", execErr)
		return
	}
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}
}

// TestReserveOrderSuccess успешно меняет статус заказа на зарезервированный
func TestReserveOrderSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var orderID int = 1
	var userID int = 1
	var serviceID int = 1
	curTime := time.Now()

	rows := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})
	expect := order.Order{orderID, userID, serviceID, 100, curTime, "Good", order.REGISTRATED}
	rows.AddRow(expect.OrderID, expect.UserID, expect.ServiceID, expect.OrderCost,
		expect.CreatingTime, expect.Comment, expect.OrderState)

	rows2 := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})
	rows2.AddRow(expect.OrderID, expect.UserID, expect.ServiceID, expect.OrderCost,
		expect.CreatingTime, expect.Comment, expect.OrderState)

	mock.ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(rows).WillReturnError(nil)

	mock.ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(rows2).WillReturnError(nil)

	mock.ExpectExec("UPDATE balanceApp.orders SET orderState = ").
		WithArgs(order.RESERVED, orderID, userID, serviceID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := order_repo.NewOrderRepo(db)
	controller := CreateNewOrderController(repo)

	execErr := controller.ReserveOrder(orderID, userID, serviceID)

	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}
}

// TestReserveOrderWrongStateError выдаёт ошибку, так как у заказа не подходящий статус для резервирования
func TestReserveOrderWrongStateError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var orderID int = 1
	var userID int = 1
	var serviceID int = 1
	curTime := time.Now()

	rows := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})
	expect := order.Order{orderID, userID, serviceID, 100, curTime, "Good", order.FINISHED}
	rows.AddRow(expect.OrderID, expect.UserID, expect.ServiceID, expect.OrderCost,
		expect.CreatingTime, expect.Comment, expect.OrderState)

	rows2 := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})
	rows2.AddRow(expect.OrderID, expect.UserID, expect.ServiceID, expect.OrderCost,
		expect.CreatingTime, expect.Comment, expect.OrderState)

	mock.ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(rows).WillReturnError(nil)

	mock.ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(rows2).WillReturnError(nil)

	repo := order_repo.NewOrderRepo(db)
	controller := CreateNewOrderController(repo)

	execErr := controller.ReserveOrder(orderID, userID, serviceID)

	if execErr != WrongStateError {
		t.Errorf("unexpected err: %v", execErr)
		return
	}
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}
}

// TestFinishOrderSuccess успешно меняет статус заказа на завершенный
func TestFinishOrderSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var orderID int = 1
	var userID int = 1
	var serviceID int = 1
	curTime := time.Now()

	rows := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})
	expect := order.Order{orderID, userID, serviceID, 100, curTime, "Good", order.RESERVED}
	rows.AddRow(expect.OrderID, expect.UserID, expect.ServiceID, expect.OrderCost,
		expect.CreatingTime, expect.Comment, expect.OrderState)

	rows2 := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})
	rows2.AddRow(expect.OrderID, expect.UserID, expect.ServiceID, expect.OrderCost,
		expect.CreatingTime, expect.Comment, expect.OrderState)

	mock.ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(rows).WillReturnError(nil)

	mock.ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(rows2).WillReturnError(nil)

	mock.ExpectExec("UPDATE balanceApp.orders SET orderState = ").
		WithArgs(order.FINISHED, orderID, userID, serviceID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := order_repo.NewOrderRepo(db)
	controller := CreateNewOrderController(repo)

	execErr := controller.FinishOrder(orderID, userID, serviceID)

	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}
}

// TestFinishOrderWrongStateError выдаёт ошибку, так как у заказа не подходящий статус для завершения
func TestFinishOrderWrongStateError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var orderID int = 1
	var userID int = 1
	var serviceID int = 1
	curTime := time.Now()

	rows := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})
	expect := order.Order{orderID, userID, serviceID, 100, curTime, "Good", order.REGISTRATED}
	rows.AddRow(expect.OrderID, expect.UserID, expect.ServiceID, expect.OrderCost,
		expect.CreatingTime, expect.Comment, expect.OrderState)

	rows2 := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})
	rows2.AddRow(expect.OrderID, expect.UserID, expect.ServiceID, expect.OrderCost,
		expect.CreatingTime, expect.Comment, expect.OrderState)

	mock.ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(rows).WillReturnError(nil)

	mock.ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(rows2).WillReturnError(nil)

	repo := order_repo.NewOrderRepo(db)
	controller := CreateNewOrderController(repo)

	execErr := controller.FinishOrder(orderID, userID, serviceID)

	if execErr != WrongStateError {
		t.Errorf("unexpected err: %v", execErr)
		return
	}
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}
}

// TestReturnOrderSuccess успешно меняет статус заказа на возвращен
func TestReturnOrderSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var orderID int = 1
	var userID int = 1
	var serviceID int = 1
	var orderCost float64 = 100
	curTime := time.Now()

	rows := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})
	expect := order.Order{orderID, userID, serviceID, orderCost, curTime, "Good", order.RESERVED}
	rows.AddRow(expect.OrderID, expect.UserID, expect.ServiceID, expect.OrderCost,
		expect.CreatingTime, expect.Comment, expect.OrderState)

	rows2 := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})
	rows2.AddRow(expect.OrderID, expect.UserID, expect.ServiceID, expect.OrderCost,
		expect.CreatingTime, expect.Comment, expect.OrderState)

	mock.ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(rows).WillReturnError(nil)

	mock.ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(rows2).WillReturnError(nil)

	mock.ExpectExec("UPDATE balanceApp.orders SET orderState = ").
		WithArgs(order.RETURNED, orderID, userID, serviceID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := order_repo.NewOrderRepo(db)
	controller := CreateNewOrderController(repo)

	sum, execErr := controller.ReturnOrder(orderID, userID, serviceID)

	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	if !reflect.DeepEqual(sum, orderCost) {
		t.Errorf("results not match, want %v, have %v", orderCost, sum)
		return
	}
}

// TestReturnOrderWrongStateError выдаёт ошибку, так как у заказа не подходящий статус для возврата
func TestReturnOrderWrongStateError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var orderCost float64 = 100
	var orderID int = 1
	var userID int = 1
	var serviceID int = 1
	curTime := time.Now()

	rows := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})
	expect := order.Order{orderID, userID, serviceID, orderCost, curTime, "Good", order.REGISTRATED}
	rows.AddRow(expect.OrderID, expect.UserID, expect.ServiceID, expect.OrderCost,
		expect.CreatingTime, expect.Comment, expect.OrderState)

	rows2 := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})
	rows2.AddRow(expect.OrderID, expect.UserID, expect.ServiceID, expect.OrderCost,
		expect.CreatingTime, expect.Comment, expect.OrderState)

	mock.ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(rows).WillReturnError(nil)

	mock.ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(rows2).WillReturnError(nil)

	repo := order_repo.NewOrderRepo(db)
	controller := CreateNewOrderController(repo)

	_, execErr := controller.ReturnOrder(orderID, userID, serviceID)

	if execErr != WrongStateError {
		t.Errorf("unexpected err: %v", execErr)
		return
	}
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}
}

// TestGetFinancialReportSuccess проверяет, что из базы данные по услугам получаются корректно.
func TestGetFinancialReportSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var month int = 2
	var year int = 2022

	rows := sqlmock.NewRows([]string{"serviceType", "sum"})
	allServicesReport := []report.FinanceReport{{1, 100}, {2, 150}}
	for _, service := range allServicesReport {
		rows.AddRow(service.ServiceType, service.Sum)
	}

	mock.ExpectQuery("SELECT serviceType ").
		WithArgs(month, year).WillReturnRows(rows)

	repo := order_repo.NewOrderRepo(db)
	controller := CreateNewOrderController(repo)

	curReports, execErr := controller.GetFinanceReports(month, year)
	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	if !reflect.DeepEqual(curReports, allServicesReport) {
		t.Errorf("results not match, want %v, have %v", allServicesReport, curReports)
		return
	}
}

// TestGetFinancialReportBadMonth проверяет, что если месяц задан некорректно, то вернётся ошибка.
func TestGetFinancialReportBadMonth(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var month int = 13
	var year int = 2022

	repo := order_repo.NewOrderRepo(db)
	controller := CreateNewOrderController(repo)

	_, execErr := controller.GetFinanceReports(month, year)
	if execErr != BadMonthError {
		t.Errorf("unexpected err: %v", execErr)
		return
	}
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}
}

// TestGetFinancialReportBadYear проверяет, что если год задан некорректно, то вернётся ошибка.
func TestGetFinancialReportBadYear(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var month int = 13
	var year int = 2024

	repo := order_repo.NewOrderRepo(db)
	controller := CreateNewOrderController(repo)

	_, execErr := controller.GetFinanceReports(month, year)
	if execErr != BadMonthError {
		t.Errorf("unexpected err: %v", execErr)
		return
	}
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}
}
