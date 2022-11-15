package order_manager

//go:generate mockgen -source=interface.go -destination=mocks/manager_mock.go -package=mocks OrderManagerInterface
type OrderManagerInterface interface {
	BuyService(userID, orderID, serviceID int, sum float64, comment string) error
	AcceptBuy(userID, orderID, serviceID int) error
	RefuseBuy(userID, orderID, serviceID int, comment string) error
}
