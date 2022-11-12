package order_repo

import (
	"Avito-Internship-Task/internal/app/balance_service_app/order"
	"Avito-Internship-Task/internal/pkg/utils"
	"database/sql"
	"sync"
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
	_, err := repo.conn.Exec(query, order.OrderID, order.UserID, order.ServiceID,
		order.OrderCost, order.CreatingTime, order.Comment, order.OrderState)
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
			err = rows.Scan(&newOrder.OrderID, &newOrder.UserID, &newOrder.ServiceID,
				&newOrder.OrderCost, &newOrder.CreatingTime, &newOrder.Comment, &newOrder.OrderState)
			if err != nil {
				break
			} else {
				allOrders = append(allOrders, *newOrder)
			}
		}
	}
	return allOrders, err
}

func (repo *OrderRepo) GetOrderByID(orderID, userID, serviceType int64) (order.Order, error) {
	order := order.Order{}

	repo.mutex.Lock()
	query := MySQLGetOrderByID{}.GetString()
	row := repo.conn.QueryRow(query, orderID, userID, serviceType)
	repo.mutex.Unlock()

	err := row.Scan(&order.OrderID, &order.UserID, &order.ServiceID,
		&order.OrderCost, &order.CreatingTime, &order.Comment, &order.OrderState)

	return order, err
}

func (repo *OrderRepo) GetUserOrders(userID int64) ([]order.Order, error) {
	allOrders := make([]order.Order, utils.EMPTY)

	repo.mutex.Lock()
	query := MySQLGetUserOrders{}.GetString()
	rows, err := repo.conn.Query(query, userID)
	repo.mutex.Unlock()

	if err == nil {
		for rows.Next() {
			newOrder := &order.Order{}
			err = rows.Scan(&newOrder.OrderID, &newOrder.UserID, &newOrder.ServiceID,
				&newOrder.OrderCost, &newOrder.CreatingTime, &newOrder.Comment, &newOrder.OrderState)
			if err != nil {
				break
			} else {
				allOrders = append(allOrders, *newOrder)
			}
		}
	}
	return allOrders, err
}

func (repo *OrderRepo) GetServiceOrders(serviceType int64) ([]order.Order, error) {
	allOrders := make([]order.Order, utils.EMPTY)

	repo.mutex.Lock()
	query := MySQLGetServiceOrders{}.GetString()
	rows, err := repo.conn.Query(query, serviceType)
	repo.mutex.Unlock()

	if err == nil {
		for rows.Next() {
			newOrder := &order.Order{}
			err = rows.Scan(&newOrder.OrderID, &newOrder.UserID, &newOrder.ServiceID,
				&newOrder.OrderCost, &newOrder.CreatingTime, &newOrder.Comment, &newOrder.OrderState)
			if err != nil {
				break
			} else {
				allOrders = append(allOrders, *newOrder)
			}
		}
	}
	return allOrders, err
}

func (repo *OrderRepo) ChangeOrderState(orderID, userID, serviceType int64, orderState int64) error {
	repo.mutex.Lock()
	query := MySQLChangeOrderState{}.GetString()
	_, err := repo.conn.Exec(query, orderState, orderID, userID, serviceType)
	repo.mutex.Unlock()
	return err
}
