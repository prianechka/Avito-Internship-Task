package order_controller

import "Avito-Internship-Task/internal/app/balance_service_app/order"

//go:generate mockgen -source=interface.go -destination=mocks/order_controller_mock.go -package=mocks OrderControllerInterface
type OrderControllerInterface interface {
	GetOrder(orderID, userID, serviceID int64) (order.Order, error)
	CreateNewOrder(orderID, userID, serviceID int64, sum float64, comment string) error
	CheckOrderIsExist(orderID, userID, serviceID int64) (bool, error)
	ReserveOrder(orderID, userID, serviceID int64) error
	FinishOrder(orderID, userID, serviceID int64) error
	ReturnOrder(orderID, userID, serviceID int64) error
}
