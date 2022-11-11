package order_repo

import "Avito-Internship-Task/internal/app/balance_service_app/order"

//go:generate mockgen -source=interface.go -destination=mocks/order_repo_mock.go -package=mocks OrderRepoInterface
type OrderRepoInterface interface {
	CreateOrder(order order.Order) error
	GetAllOrders() ([]order.Order, error)
	GetOrderByID(orderID int64) (order.Order, error)
	GetUserOrders(userID int64) ([]order.Order, error)
	GetServiceOrders(serviceType int64) ([]order.Order, error)
	ChangeOrderState(orderID int64, orderState int64) error
}
