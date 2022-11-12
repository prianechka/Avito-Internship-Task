package transaction_repo

import (
	"Avito-Internship-Task/internal/app/balance_service_app/transaction"
	"Avito-Internship-Task/internal/pkg/utils"
	"database/sql"
	"sync"
)

type TransactionRepo struct {
	mutex sync.RWMutex
	conn  *sql.DB
}

func NewTransactionRepo(conn *sql.DB) *TransactionRepo {
	return &TransactionRepo{conn: conn}
}

func (repo *TransactionRepo) AddNewTransaction(newTransaction transaction.Transaction) error {
	repo.mutex.Lock()
	query := MySQLAddNewTransaction{}.GetString()
	_, err := repo.conn.Exec(query, newTransaction.TransactionID, newTransaction.UserID,
		newTransaction.TransactionType, newTransaction.Sum, newTransaction.Time,
		newTransaction.ActionComments, newTransaction.AddComments)
	repo.mutex.Unlock()
	return err
}

func (repo *TransactionRepo) GetAllTransactions() ([]transaction.Transaction, error) {
	allTransactions := make([]transaction.Transaction, utils.EMPTY)

	repo.mutex.Lock()
	query := MySQLGetAllTransactions{}.GetString()
	rows, err := repo.conn.Query(query)
	repo.mutex.Unlock()

	if err == nil {
		for rows.Next() {
			newTransact := transaction.Transaction{}
			err = rows.Scan(&newTransact.TransactionID, &newTransact.UserID, &newTransact.TransactionType,
				&newTransact.Sum, &newTransact.Time, &newTransact.ActionComments, &newTransact.AddComments)
			if err != nil {
				break
			} else {
				allTransactions = append(allTransactions, newTransact)
			}
		}
	}
	return allTransactions, err
}

func (repo *TransactionRepo) GetUserTransactions(userID int64) ([]transaction.Transaction, error) {
	allTransactions := make([]transaction.Transaction, utils.EMPTY)

	repo.mutex.Lock()
	query := MySQLGetUserTransactions{}.GetString()
	rows, err := repo.conn.Query(query, userID)
	repo.mutex.Unlock()

	if err == nil {
		for rows.Next() {
			newTransact := transaction.Transaction{}
			err = rows.Scan(&newTransact.TransactionID, &newTransact.UserID, &newTransact.TransactionType,
				&newTransact.Sum, &newTransact.Time, &newTransact.ActionComments, &newTransact.AddComments)
			if err != nil {
				break
			} else {
				allTransactions = append(allTransactions, newTransact)
			}
		}
	}
	return allTransactions, err
}

func (repo *TransactionRepo) GetTransactionByID(transactionID int64) (transaction.Transaction, error) {
	newTransact := transaction.Transaction{}

	repo.mutex.Lock()
	query := MySQLGetTransactionByID{}.GetString()
	row := repo.conn.QueryRow(query, transactionID)
	repo.mutex.Unlock()

	err := row.Scan(&newTransact.TransactionID, &newTransact.UserID, &newTransact.TransactionType,
		&newTransact.Sum, &newTransact.Time, &newTransact.ActionComments, &newTransact.AddComments)

	return newTransact, err
}
