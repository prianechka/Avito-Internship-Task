package account_controller

import (
	"Avito-Internship-Task/internal/app/balance_service_app/account/account_repo"
	"sync"
)

type AccountManager struct {
	mutex sync.RWMutex
	repo  account_repo.AccountRepoInterface
}

func CreateNewAccountManager(repo account_repo.AccountRepoInterface) *AccountManager {
	return &AccountManager{mutex: sync.RWMutex{}, repo: repo}
}

func (m *AccountManager) CheckAccountIsExist(userID int64) (result bool, err error) {
	_, isAccExistErr := m.repo.GetCurrentAmount(userID)
	switch isAccExistErr {
	case account_repo.AccountNotExist:
		result = false
	case nil:
		result = true
	default:
		err = isAccExistErr
	}
	return result, err
}

func (m *AccountManager) CreateNewAccount(userID int64) error {
	m.mutex.Lock()
	isAccExist, err := m.CheckAccountIsExist(userID)
	if err == nil {
		if !isAccExist {
			err = m.repo.AddNewAccount(userID)
		} else {
			err = AccountIsExistErr
		}
	}
	m.mutex.Unlock()
	return err
}

func (m *AccountManager) CheckBalance(userID int64) (float64, error) {
	return m.repo.GetCurrentAmount(userID)
}

func (m *AccountManager) CheckAbleToBuyService(userID int64, servicePrice float64) (bool, error) {
	var result bool
	balance, err := m.repo.GetCurrentAmount(userID)
	if err == nil {
		if servicePrice <= balance {
			result = true
		}
	}
	return result, err
}

func (m *AccountManager) DonateMoney(userID int64, sum float64) (err error) {
	m.mutex.Lock()
	if sum >= 0 {
		err = m.repo.ChangeAmount(userID, sum)
	} else {
		err = NegSumError
	}
	m.mutex.Unlock()
	return err
}

func (m *AccountManager) SpendMoney(userID int64, sum float64) error {
	m.mutex.Lock()
	canSpendMoney, err := m.CheckAbleToBuyService(userID, sum)
	if err == nil && canSpendMoney {
		err = m.repo.ChangeAmount(userID, -sum)
	} else {
		err = NotEnoughMoneyErr
	}
	m.mutex.Unlock()
	return err
}
