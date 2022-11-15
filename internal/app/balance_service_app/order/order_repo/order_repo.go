package order_repo

import (
	"Avito-Internship-Task/internal/app/balance_service_app/order"
	"Avito-Internship-Task/internal/app/balance_service_app/report"
	"Avito-Internship-Task/internal/pkg/utils"
	"database/sql"
	"sync"
	"time"
)

type OrderRepo struct {
	mutex sync.RWMutex
	conn  *sql.DB
}

func NewOrderRepo(conn *sql.DB) *OrderRepo {
	return &OrderRepo{conn: conn}
}

func (repo *OrderRepo) CreateOrder(order order.Order) error {
	repo.mutex.Lock()
	query := MySQLAddNewOrder{}.GetString()
	time := order.CreatingTime.Format(utils.TimeLayout)
	_, err := repo.conn.Exec(query, order.OrderID, order.UserID, order.ServiceID,
		order.OrderCost, time, order.Comment, order.OrderState)
	repo.mutex.Unlock()
	return err
}

func (repo *OrderRepo) GetAllOrders() ([]order.Order, error) {
	allOrders := make([]order.Order, utils.EMPTY)

	repo.mutex.Lock()
	query := MySQLGetAllOrders{}.GetString()
	rows, err := repo.conn.Query(query)
	repo.mutex.Unlock()

	if err == nil {
		for rows.Next() {
			newOrder := &order.Order{}
			var orderTime string
			err = rows.Scan(&newOrder.OrderID, &newOrder.UserID, &newOrder.ServiceID,
				&newOrder.OrderCost, &orderTime, &newOrder.Comment, &newOrder.OrderState)
			newOrder.CreatingTime, _ = time.Parse(utils.TimeLayout, orderTime)
			if err != nil {
				break
			} else {
				allOrders = append(allOrders, *newOrder)
			}
		}
	}
	return allOrders, err
}

func (repo *OrderRepo) GetOrderByID(orderID, userID, serviceType int) (order.Order, error) {
	order := order.Order{}

	repo.mutex.Lock()
	query := MySQLGetOrderByID{}.GetString()
	row := repo.conn.QueryRow(query, orderID, userID, serviceType)
	repo.mutex.Unlock()

	var orderTime string

	err := row.Scan(&order.OrderID, &order.UserID, &order.ServiceID,
		&order.OrderCost, &orderTime, &order.Comment, &order.OrderState)

	order.CreatingTime, _ = time.Parse(utils.TimeLayout, orderTime)

	return order, err
}

func (repo *OrderRepo) GetUserOrders(userID int) ([]order.Order, error) {
	allOrders := make([]order.Order, utils.EMPTY)

	repo.mutex.Lock()
	query := MySQLGetUserOrders{}.GetString()
	rows, err := repo.conn.Query(query, userID)
	repo.mutex.Unlock()

	if err == nil {
		for rows.Next() {
			newOrder := &order.Order{}

			var orderTime string
			err = rows.Scan(&newOrder.OrderID, &newOrder.UserID, &newOrder.ServiceID,
				&newOrder.OrderCost, &orderTime, &newOrder.Comment, &newOrder.OrderState)
			newOrder.CreatingTime, _ = time.Parse(utils.TimeLayout, orderTime)
			if err != nil {
				break
			} else {
				allOrders = append(allOrders, *newOrder)
			}
		}
	}
	return allOrders, err
}

func (repo *OrderRepo) GetServiceOrders(serviceType int) ([]order.Order, error) {
	allOrders := make([]order.Order, utils.EMPTY)

	repo.mutex.Lock()
	query := MySQLGetServiceOrders{}.GetString()
	rows, err := repo.conn.Query(query, serviceType)
	repo.mutex.Unlock()

	if err == nil {
		for rows.Next() {
			newOrder := &order.Order{}
			var orderTime string
			err = rows.Scan(&newOrder.OrderID, &newOrder.UserID, &newOrder.ServiceID,
				&newOrder.OrderCost, &orderTime, &newOrder.Comment, &newOrder.OrderState)
			newOrder.CreatingTime, _ = time.Parse(utils.TimeLayout, orderTime)
			if err != nil {
				break
			} else {
				allOrders = append(allOrders, *newOrder)
			}
		}
	}
	return allOrders, err
}

func (repo *OrderRepo) ChangeOrderState(orderID, userID, serviceType int, orderState int) error {
	repo.mutex.Lock()
	query := MySQLChangeOrderState{}.GetString()
	_, err := repo.conn.Exec(query, orderState, orderID, userID, serviceType)
	repo.mutex.Unlock()
	return err
}

func (repo *OrderRepo) GetSumOfFinishedServices(month, year int) ([]report.FinanceReport, error) {
	allServices := make([]report.FinanceReport, utils.EMPTY)

	repo.mutex.Lock()
	query := MySQLGetAllOrdersStat{}.GetString()
	rows, err := repo.conn.Query(query, month, year)
	repo.mutex.Unlock()

	if err == nil {
		for rows.Next() {
			newServiceReport := report.FinanceReport{}
			err = rows.Scan(&newServiceReport.ServiceType, &newServiceReport.Sum)
			if err != nil {
				break
			} else {
				allServices = append(allServices, newServiceReport)
			}
		}
	}
	return allServices, err
}
