package transaction_controller

import "Avito-Internship-Task/internal/app/balance_service_app/transaction"

//go:generate mockgen -source=interface.go -destination=mocks/transaction_controller_mock.go -package=mocks TransactionControllerInterface
type TransactionControllerInterface interface {
	GetTransactionByID(transactionID int) (transaction.Transaction, error)
	AddNewRecordRefillBalance(userID int, sum float64, comments string) error
	AddNewRecordBuyService(userID int, sum float64, serviceID int, comments string) error
	AddNewRecordReturnService(userID int, sum float64, serviceID int, comments string) error
	AddNewRecordTransferTo(srcUserID, dstUserID int, sum float64, comments string) error
	AddNewRecordTransferFrom(srcUserID, dstUserID int, sum float64, comments string) error
	GetUserTransactions(userID int, orderBy string, limit, offset int) ([]transaction.Transaction, error)
}
