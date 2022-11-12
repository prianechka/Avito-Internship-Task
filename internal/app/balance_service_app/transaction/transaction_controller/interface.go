package transaction_controller

import "Avito-Internship-Task/internal/app/balance_service_app/transaction"

//go:generate mockgen -source=interface.go -destination=mocks/transaction_controller_mock.go -package=mocks TransactionControllerInterface
type TransactionControllerInterface interface {
	GetTransactionByID(transactionID int64) (transaction.Transaction, error)
	AddNewRecordRefillBalance(userID int64, sum float64, comments string) error
	AddNewRecordBuyService(userID int64, sum float64, serviceID int64, comments string) error
	AddNewRecordReturnService(userID int64, sum float64, serviceID int64, comments string) error
	AddNewRecordTransferTo(srcUserID, dstUserID int64, sum float64, comments string) error
	GetUserTransactions(userID int64) ([]transaction.Transaction, error)
}
