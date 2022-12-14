package transaction_repo

import "Avito-Internship-Task/internal/app/balance_service_app/transaction"

//go:generate mockgen -source=interface.go -destination=mocks/transaction_repo_mock.go -package=mocks TransactionRepoInterface
type TransactionRepoInterface interface {
	AddNewTransaction(newTransaction transaction.Transaction) error
	GetAllTransactions() ([]transaction.Transaction, error)
	GetUserTransactions(userID int, orderBy string, limit int, offset int) ([]transaction.Transaction, error)
	GetTransactionByID(transactionID int) (transaction.Transaction, error)
}
