package transaction_controller

import (
	"Avito-Internship-Task/internal/app/balance_service_app/order"
	"Avito-Internship-Task/internal/app/balance_service_app/transaction"
	"Avito-Internship-Task/internal/app/balance_service_app/transaction/transaction_repo"
	"Avito-Internship-Task/internal/pkg/utils"
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

func (c *TransactionController) GetTransactionByID(transactionID int) (transaction.Transaction, error) {
	return c.repo.GetTransactionByID(transactionID)
}

func (c *TransactionController) AddNewRecordRefillBalance(userID int, sum float64, comments string) error {
	newTransact := transaction.Transaction{
		TransactionID:   int(c.transactCount),
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

func (c *TransactionController) AddNewRecordBuyService(userID int, sum float64, serviceID int, comments string) error {
	serviceName := order.Types[serviceID]
	if serviceName == utils.EmptyString {
		serviceName = fmt.Sprintf("Service with id %d", serviceID)
	}

	newTransact := transaction.Transaction{
		TransactionID:   int(c.transactCount),
		UserID:          userID,
		TransactionType: transaction.Buy,
		Sum:             sum,
		Time:            time.Now(),
		ActionComments:  "куплена услуга: " + serviceName,
		AddComments:     comments,
	}

	c.mutex.Lock()
	err := c.repo.AddNewTransaction(newTransact)
	atomic.AddInt64(&c.transactCount, 1)
	c.mutex.Unlock()

	return err
}

func (c *TransactionController) AddNewRecordReturnService(userID int, sum float64, serviceID int, comments string) error {
	serviceName := order.Types[serviceID]
	if serviceName == utils.EmptyString {
		serviceName = fmt.Sprintf("Service with id %d", serviceID)
	}

	newTransact := transaction.Transaction{
		TransactionID:   int(c.transactCount),
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

func (c *TransactionController) AddNewRecordTransferTo(srcUserID, dstUserID int, sum float64, comments string) error {
	newTransact := transaction.Transaction{
		TransactionID:   int(c.transactCount),
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

func (c *TransactionController) AddNewRecordTransferFrom(srcUserID, dstUserID int, sum float64, comments string) error {
	newTransact := transaction.Transaction{
		TransactionID:   int(c.transactCount),
		UserID:          srcUserID,
		TransactionType: transaction.Transfer,
		Sum:             sum,
		Time:            time.Now(),
		ActionComments:  "перевод от пользователя: " + fmt.Sprintf("%d", dstUserID),
		AddComments:     comments,
	}

	c.mutex.Lock()
	err := c.repo.AddNewTransaction(newTransact)
	atomic.AddInt64(&c.transactCount, 1)
	c.mutex.Unlock()
	return err
}

func (c *TransactionController) GetUserTransactions(userID int, orderBy string,
	limit, offset int) ([]transaction.Transaction, error) {
	return c.repo.GetUserTransactions(userID, orderBy, limit, offset)
}
