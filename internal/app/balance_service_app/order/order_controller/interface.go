package order_controller

import (
	"Avito-Internship-Task/internal/app/balance_service_app/order"
	"Avito-Internship-Task/internal/app/balance_service_app/report"
)

//go:generate mockgen -source=interface.go -destination=mocks/order_controller_mock.go -package=mocks OrderControllerInterface
type OrderControllerInterface interface {
	GetOrder(orderID, userID, serviceID int) (order.Order, error)
	CreateNewOrder(orderID, userID, serviceID int, sum float64, comment string) error
	CheckOrderIsExist(orderID, userID, serviceID int) (bool, error)
	ReserveOrder(orderID, userID, serviceID int) error
	FinishOrder(orderID, userID, serviceID int) error
	ReturnOrder(orderID, userID, serviceID int) (float64, error)
	GetFinanceReports(month, year int) ([]report.FinanceReport, error)
}
