package order_repo

import (
	"Avito-Internship-Task/internal/app/balance_service_app/order"
	"Avito-Internship-Task/internal/app/balance_service_app/report"
)

//go:generate mockgen -source=interface.go -destination=mocks/order_repo_mock.go -package=mocks OrderRepoInterface
type OrderRepoInterface interface {
	CreateOrder(order order.Order) error
	GetAllOrders() ([]order.Order, error)
	GetOrderByID(orderID, userID, serviceType int) (order.Order, error)
	GetUserOrders(userID int) ([]order.Order, error)
	GetServiceOrders(serviceType int) ([]order.Order, error)
	ChangeOrderState(orderID, userID, serviceType int, orderState int) error
	GetSumOfFinishedServices(month, year int) ([]report.FinanceReport, error)
}
