package account_manager

import (
	ac "Avito-Internship-Task/internal/app/balance_service_app/account/account_controller"
	tc "Avito-Internship-Task/internal/app/balance_service_app/transaction/transaction_controller"
	"fmt"
)

type AccountManager struct {
	accountController     ac.AccountControllerInterface
	transactionController tc.TransactionControllerInterface
}

func CreateNewAccountManager(accController ac.AccountControllerInterface,
	transactionController tc.TransactionControllerInterface) *AccountManager {
	return &AccountManager{
		accountController:     accController,
		transactionController: transactionController,
	}
}

func (m *AccountManager) RefillBalance(userID int, sum float64, comments string) error {
	isAccExist, err := m.accountController.CheckAccountIsExist(userID)
	if err == nil {
		if !isAccExist {
			err = m.accountController.CreateNewAccount(userID)
		}
		if err == nil {
			err = m.accountController.DonateMoney(userID, sum)
			if err == nil {
				err = m.transactionController.AddNewRecordRefillBalance(userID, sum, comments)
			}
		}
	}
	return err
}

func (m *AccountManager) GetUserBalance(userID int) (float64, error) {
	isAccExist, err := m.accountController.CheckAccountIsExist(userID)
	if err == nil {
		if isAccExist {
			return m.accountController.CheckBalance(userID)
		} else {
			return 0, ac.AccountNotExistErr
		}
	} else {
		return 0, fmt.Errorf("managers call: %w", err)
	}
}

func (m *AccountManager) Transfer(srcUserID, dstUserID int, sum float64, comment string) error {
	isAccExists, err := m.checkAllAccountsAreExists(srcUserID, dstUserID)
	if err == nil {
		if isAccExists {
			canBuy, checkBuyErr := m.checkUserCanBuyService(srcUserID, sum)
			if checkBuyErr == nil {
				if canBuy {
					err = m.accountController.SpendMoney(srcUserID, sum)
					if err == nil {
						err = m.accountController.DonateMoney(dstUserID, sum)
						if err == nil {
							err = m.makeReportsForAllUsers(srcUserID, dstUserID, sum, comment)
						}
					}
				} else {
					err = ac.NotEnoughMoneyErr
				}
			}
		} else {
			err = ac.AccountNotExistErr
		}
	}
	return err
}

func (m *AccountManager) checkAllAccountsAreExists(srcUserID, dstUserID int) (bool, error) {
	var result bool
	var err error
	isAccExist, firstCheck := m.accountController.CheckAccountIsExist(srcUserID)
	if firstCheck == nil {
		if isAccExist {
			isSecAccExist, secCheckErr := m.accountController.CheckAccountIsExist(dstUserID)
			if secCheckErr == nil {
				if isSecAccExist {
					result = true
				} else {
					err = ac.AccountNotExistErr
				}
			}
		} else {
			err = ac.AccountNotExistErr
		}
	}
	return result, err
}

func (m *AccountManager) checkUserCanBuyService(userID int, sum float64) (bool, error) {
	var canBuy bool
	isAccExist, err := m.accountController.CheckAccountIsExist(userID)
	if err == nil {
		if sum <= 0 {
			err = ac.NegSumError
		} else if !isAccExist {
			err = ac.AccountNotExistErr
		} else {
			canBuy, err = m.accountController.CheckAbleToBuyService(userID, sum)
		}
	}
	return canBuy, err
}

func (m *AccountManager) makeReportsForAllUsers(srcUserID, dstUserID int, sum float64, comment string) error {
	err := m.transactionController.AddNewRecordTransferTo(srcUserID, dstUserID, sum, comment)
	if err == nil {
		err = m.transactionController.AddNewRecordTransferFrom(dstUserID, srcUserID, sum, comment)
	}
	return err
}
