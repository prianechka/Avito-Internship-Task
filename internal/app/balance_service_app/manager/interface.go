package manager

import "Avito-Internship-Task/internal/app/balance_service_app/transaction"

//go:generate mockgen -source=interface.go -destination=mocks/manager_mock.go -package=mocks ManagerInterface
type ManagerInterface interface {
	RefillBalance(userID int, sum float64, comments string) error
	GetUserBalance(userID int) (float64, error)
	BuyService(userID, orderID, serviceID int, sum float64, comment string) error
	AcceptBuy(userID, orderID, serviceID int) error
	RefuseBuy(userID, orderID, serviceID int, comment string) error
	Transfer(srcUserID, dstUserID int, sum float64, comment string) error
	GetFinanceReport(month, year int, url string) error
	GetUserReport(userID int, orderBy string, limit, offset int) ([]transaction.Transaction, error)
}
