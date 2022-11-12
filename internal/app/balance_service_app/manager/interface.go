package manager

//go:generate mockgen -source=interface.go -destination=mocks/manager_mock.go -package=mocks ManagerInterface
type ManagerInterface interface {
	RefillBalance(userID int64, sum float64, comments string) error
	GetUserBalance(userID int64) (float64, error)
	BuyService(userID, orderID, serviceID int64, sum float64, comment string) error
	AcceptBuy(userID, orderID, serviceID int64) error
	RefuseBuy(userID, orderID, serviceID int64, comment string) error
	Transfer(srcUserID, dstUserID int64, sum float64, comment string) error
	GetReport() error
	GetUserReport() error
}
