package transaction_controller

import (
	"Avito-Internship-Task/internal/app/balance_service_app/order"
	"Avito-Internship-Task/internal/app/balance_service_app/transaction"
	"Avito-Internship-Task/internal/app/balance_service_app/transaction/transaction_repo"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type TransactionController struct {
	mutex         sync.RWMutex
	repo          transaction_repo.TransactionRepoInterface
	transactCount int64
}

func CreateNewTransactionController(repo transaction_repo.TransactionRepoInterface) *TransactionController {
	return &TransactionController{
		mutex: sync.RWMutex{},
		repo:  repo,
	}
}

func (c *TransactionController) GetTransactionByID(transactionID int64) (transaction.Transaction, error) {
	return c.repo.GetTransactionByID(transactionID)
}

func (c *TransactionController) AddNewRecordRefillBalance(userID int64, sum float64, comments string) error {
	newTransact := transaction.Transaction{
		TransactionID:   c.transactCount,
		UserID:          userID,
		TransactionType: transaction.Refill,
		Sum:             sum,
		Time:            time.Now(),
		ActionComments:  "зачислены средства на баланс",
		AddComments:     comments,
	}

	c.mutex.Lock()
	err := c.repo.AddNewTransaction(newTransact)
	c.mutex.Unlock()
	return err
}

func (c *TransactionController) AddNewRecordBuyService(userID int64, sum float64, serviceID int64, comments string) error {
	newTransact := transaction.Transaction{
		TransactionID:   c.transactCount,
		UserID:          userID,
		TransactionType: transaction.Buy,
		Sum:             sum,
		Time:            time.Now(),
		ActionComments:  "куплена услуга: " + order.Types[serviceID],
		AddComments:     comments,
	}

	c.mutex.Lock()
	err := c.repo.AddNewTransaction(newTransact)
	atomic.AddInt64(&c.transactCount, 1)
	c.mutex.Unlock()

	return err
}

func (c *TransactionController) AddNewRecordReturnService(userID int64, sum float64, serviceID int64, comments string) error {
	newTransact := transaction.Transaction{
		TransactionID:   c.transactCount,
		UserID:          userID,
		TransactionType: transaction.Return,
		Sum:             sum,
		Time:            time.Now(),
		ActionComments:  "возврат за услугу: " + order.Types[serviceID],
		AddComments:     comments,
	}

	c.mutex.Lock()
	err := c.repo.AddNewTransaction(newTransact)
	atomic.AddInt64(&c.transactCount, 1)
	c.mutex.Unlock()
	return err
}

func (c *TransactionController) AddNewRecordTransferTo(srcUserID, dstUserID int64, sum float64, comments string) error {
	newTransact := transaction.Transaction{
		TransactionID:   c.transactCount,
		UserID:          srcUserID,
		TransactionType: transaction.Transfer,
		Sum:             sum,
		Time:            time.Now(),
		ActionComments:  "перевод пользователю: " + fmt.Sprintf("%d", dstUserID),
		AddComments:     comments,
	}

	c.mutex.Lock()
	err := c.repo.AddNewTransaction(newTransact)
	atomic.AddInt64(&c.transactCount, 1)
	c.mutex.Unlock()
	return err
}

func (c *TransactionController) GetUserTransactions(userID int64) ([]transaction.Transaction, error) {
	return c.repo.GetUserTransactions(userID)
}
