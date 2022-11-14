package order_repo

import (
	"Avito-Internship-Task/internal/app/balance_service_app/order"
	"Avito-Internship-Task/internal/app/balance_service_app/report"
	"reflect"
	"testing"
	"time"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

// go test -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html

func TestAddNewOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var newOrder order.Order

	mock.
		ExpectExec("INSERT INTO balanceApp.orders").
		WithArgs(newOrder.OrderID, newOrder.UserID, newOrder.ServiceID,
			newOrder.OrderCost, sqlmock.AnyArg(), newOrder.Comment, newOrder.OrderState).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewOrderRepo(db)

	execErr := repo.CreateOrder(newOrder)
	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}
}

func TestGetAllOrders(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	curTime := time.Now()

	rows := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})
	expect := []order.Order{{1, 1, 1, 100, curTime, "Good", 1},
		{2, 2, 2, 200, curTime, "Bad", 1}}
	for _, order := range expect {
		rows = rows.AddRow(order.OrderID, order.UserID, order.ServiceID, order.OrderCost,
			curTime, order.Comment, order.OrderState)
	}

	mock.
		ExpectQuery(MySQLGetAllOrders{}.GetString()).
		WillReturnRows(rows).WillReturnError(nil)

	repo := NewOrderRepo(db)

	allOrders, execErr := repo.GetAllOrders()
	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}

	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	for i := range allOrders {
		allOrders[i].CreatingTime = expect[i].CreatingTime
	}
	if !reflect.DeepEqual(allOrders, expect) {
		t.Errorf("results not match, want %v, have %v", expect, allOrders)
		return
	}
}

func TestGetUserOrders(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	curTime := time.Now()
	var userID int64 = 1

	rows := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})
	expect := []order.Order{{1, 1, 1, 100, curTime, "Good", 1},
		{2, 1, 2, 200, curTime, "Bad", 1}}
	for _, order := range expect {
		rows = rows.AddRow(order.OrderID, order.UserID, order.ServiceID, order.OrderCost,
			order.CreatingTime, order.Comment, order.OrderState)
	}

	mock.
		ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(rows).WillReturnError(nil)

	repo := NewOrderRepo(db)

	allOrders, execErr := repo.GetUserOrders(userID)
	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}

	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	for i := range allOrders {
		allOrders[i].CreatingTime = expect[i].CreatingTime
	}
	if !reflect.DeepEqual(allOrders, expect) {
		t.Errorf("results not match, want %v, have %v", expect, allOrders)
		return
	}
}

func TestGetServiceOrders(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	curTime := time.Now()
	var serviceType int64 = 2
	rows := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})
	expect := []order.Order{{1, 1, 2, 100, curTime, "Good", 1},
		{2, 1, 2, 200, curTime, "Bad", 1}}
	for _, order := range expect {
		rows = rows.AddRow(order.OrderID, order.UserID, order.ServiceID, order.OrderCost,
			order.CreatingTime, order.Comment, order.OrderState)
	}

	mock.
		ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(rows).WillReturnError(nil)

	repo := NewOrderRepo(db)

	allOrders, execErr := repo.GetServiceOrders(serviceType)
	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}

	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	for i := range allOrders {
		allOrders[i].CreatingTime = expect[i].CreatingTime
	}

	if !reflect.DeepEqual(allOrders, expect) {
		t.Errorf("results not match, want %v, have %v", expect, allOrders)
		return
	}
}

func TestGetOrderByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	curTime := time.Now()
	var orderID int64 = 1
	var userID int64 = 1
	var serviceType int64 = 1

	rows := sqlmock.NewRows([]string{"orderID", "userID", "serviceType", "orderCost",
		"creatingTime", "comment", "orderState"})
	expect := order.Order{orderID, userID, serviceType, 100, curTime, "Good", 1}
	rows.AddRow(expect.OrderID, expect.UserID, expect.ServiceID, expect.OrderCost,
		expect.CreatingTime, expect.Comment, expect.OrderState)

	mock.
		ExpectQuery("SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders").
		WillReturnRows(rows).WillReturnError(nil)

	repo := NewOrderRepo(db)

	allOrders, execErr := repo.GetOrderByID(orderID, userID, orderID)
	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}

	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	allOrders.CreatingTime = expect.CreatingTime

	if !reflect.DeepEqual(allOrders, expect) {
		t.Errorf("results not match, want %v, have %v", expect, allOrders)
		return
	}
}

func TestChangeState(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var orderID int64 = 1
	var userID int64 = 1
	var serviceType int64 = order.REGISTRATED
	var orderState int64 = 2

	mock.
		ExpectExec("UPDATE balanceApp.orders SET orderState = ").
		WithArgs(orderState, orderID, userID, serviceType).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewOrderRepo(db)

	execErr := repo.ChangeOrderState(orderID, userID, serviceType, orderState)
	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}
}

func TestGetFinanceReport(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var month int64 = 2
	var year int64 = 3

	rows := sqlmock.NewRows([]string{"serviceType", "sum"})
	allServicesReport := []report.FinanceReport{{1, 100}, {2, 150}}
	for _, service := range allServicesReport {
		rows.AddRow(service.ServiceType, service.Sum)
	}

	mock.ExpectQuery("SELECT serviceType ").
		WithArgs(month, year).WillReturnRows(rows)

	repo := NewOrderRepo(db)

	curReports, execErr := repo.GetSumOfFinishedServices(month, year)
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
