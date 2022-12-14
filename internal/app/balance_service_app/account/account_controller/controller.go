package account_controller

import (
	"Avito-Internship-Task/internal/app/balance_service_app/account/account_repo"
	"sync"
)

type AccountController struct {
	mutex sync.RWMutex
	repo  account_repo.AccountRepoInterface
}

func CreateNewAccountController(repo account_repo.AccountRepoInterface) *AccountController {
	return &AccountController{mutex: sync.RWMutex{}, repo: repo}
}

func (m *AccountController) CheckAccountIsExist(userID int) (result bool, err error) {
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

func (m *AccountController) CreateNewAccount(userID int) error {
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

func (m *AccountController) CheckBalance(userID int) (float64, error) {
	return m.repo.GetCurrentAmount(userID)
}

func (m *AccountController) CheckAbleToBuyService(userID int, servicePrice float64) (bool, error) {
	var result bool
	balance, err := m.repo.GetCurrentAmount(userID)
	if err == nil {
		if servicePrice <= balance {
			result = true
		}
	}
	return result, err
}

func (m *AccountController) DonateMoney(userID int, sum float64) (err error) {
	m.mutex.Lock()
	if sum >= 0 {
		err = m.repo.ChangeAmount(userID, sum)
	} else {
		err = NegSumError
	}
	m.mutex.Unlock()
	return err
}

func (m *AccountController) SpendMoney(userID int, sum float64) error {
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
