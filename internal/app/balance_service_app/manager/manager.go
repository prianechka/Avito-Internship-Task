package manager

import (
	ac "Avito-Internship-Task/internal/app/balance_service_app/account/account_controller"
	oc "Avito-Internship-Task/internal/app/balance_service_app/order/order_controller"
	"Avito-Internship-Task/internal/app/balance_service_app/report/report_controller"
	"Avito-Internship-Task/internal/app/balance_service_app/transaction"
	tc "Avito-Internship-Task/internal/app/balance_service_app/transaction/transaction_controller"
	"Avito-Internship-Task/internal/pkg/utils"
	"fmt"
)

type Manager struct {
	accountController     ac.AccountControllerInterface
	orderController       oc.OrderControllerInterface
	transactionController tc.TransactionControllerInterface
}

func CreateNewManager(accController ac.AccountControllerInterface, orderController oc.OrderControllerInterface,
	transactionController tc.TransactionControllerInterface) *Manager {
	return &Manager{
		accountController:     accController,
		orderController:       orderController,
		transactionController: transactionController,
	}
}

func (m *Manager) RefillBalance(userID int64, sum float64, comments string) error {
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

func (m *Manager) GetUserBalance(userID int64) (float64, error) {
	isAccExist, err := m.accountController.CheckAccountIsExist(userID)
	if err == nil {
		if isAccExist {
			return m.accountController.CheckBalance(userID)
		} else {
			return 0, ac.AccountNotExistErr
		}
	} else {
		return 0, fmt.Errorf("manager call: %w", err)
	}
}

func (m *Manager) BuyService(userID, orderID, serviceID int64, sum float64, comment string) error {
	canUserBuy, err := m.checkUserCanBuyService(userID, sum)
	if err == nil {
		if canUserBuy {
			err = m.createOrder(userID, orderID, serviceID, sum, comment)
			if err == nil {
				err = m.orderController.ReserveOrder(orderID, userID, serviceID)
				if err == nil {
					err = m.accountController.SpendMoney(userID, sum)
					if err == nil {
						err = m.transactionController.AddNewRecordBuyService(userID, sum, serviceID, comment)
					}
				}
			}
		} else {
			err = ac.NotEnoughMoneyErr
		}
	}
	return err
}

func (m *Manager) checkUserCanBuyService(userID int64, sum float64) (bool, error) {
	var canBuy bool
	isAccExist, err := m.accountController.CheckAccountIsExist(userID)
	if err == nil {
		if isAccExist {
			canBuy, err = m.accountController.CheckAbleToBuyService(userID, sum)
		} else {
			err = ac.AccountNotExistErr
		}
	}
	return canBuy, err
}

func (m *Manager) createOrder(userID, orderID, serviceID int64, sum float64, comment string) error {
	isOrderExist, err := m.orderController.CheckOrderIsExist(orderID, userID, serviceID)
	if err == nil {
		if !isOrderExist {
			err = m.orderController.CreateNewOrder(orderID, userID, serviceID, sum, comment)
		} else {
			err = oc.OrderIsAlreadyExist
		}
	}
	return err
}

func (m *Manager) AcceptBuy(userID, orderID, serviceID int64) error {
	isOrderExist, err := m.orderController.CheckOrderIsExist(orderID, userID, serviceID)
	if err == nil {
		if isOrderExist {
			err = m.orderController.FinishOrder(orderID, userID, serviceID)
			if err == nil {
				m.transferMoneyToCompanyAccount()
			}
		} else {
			err = oc.OrderNotFound
		}
	}
	return err
}

func (*Manager) transferMoneyToCompanyAccount() {}

func (m *Manager) RefuseBuy(userID, orderID, serviceID int64, comment string) error {
	isOrderExist, err := m.orderController.CheckOrderIsExist(orderID, userID, serviceID)
	if err == nil {
		if isOrderExist {
			sum, returnOrderErr := m.orderController.ReturnOrder(orderID, userID, serviceID)
			if returnOrderErr == nil {
				err = m.accountController.DonateMoney(userID, sum)
				if err == nil {
					err = m.transactionController.AddNewRecordReturnService(userID, sum, serviceID, comment)
				}
			} else {
				err = returnOrderErr
			}
		} else {
			err = oc.OrderNotFound
		}
	}
	return err
}

func (m *Manager) Transfer(srcUserID, dstUserID int64, sum float64, comment string) error {
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

func (m *Manager) checkAllAccountsAreExists(srcUserID, dstUserID int64) (bool, error) {
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

func (m *Manager) makeReportsForAllUsers(srcUserID, dstUserID int64, sum float64, comment string) error {
	err := m.transactionController.AddNewRecordTransferTo(srcUserID, dstUserID, sum, comment)
	if err == nil {
		err = m.transactionController.AddNewRecordTransferFrom(dstUserID, srcUserID, sum, comment)
	}
	return err
}

func (m *Manager) GetFinanceReport(month, year int64, url string) error {
	dataToReport, err := m.orderController.GetFinanceReports(month, year)
	if err == nil {
		reportController := report_controller.CreateNewReportController()
		err = reportController.CreateFinancialReportCSV(dataToReport, url)
	}
	return err
}

func (m *Manager) GetUserReport(userID int64, orderBy string, limit, offset int) ([]transaction.Transaction, error) {
	var allTransactions = make([]transaction.Transaction, utils.EMPTY)
	var err error

	if limit == utils.NotInQuery {
		limit = utils.DefaultLimit
	}
	if offset == utils.NotInQuery {
		offset = utils.DefaultOffset
	}
	if orderBy == utils.EmptyString {
		orderBy = utils.DefaultOrderBy
	}

	_, checkAccountError := m.accountController.CheckAccountIsExist(userID)
	if checkAccountError == nil {
		allTransactions, err = m.transactionController.GetUserTransactions(userID, orderBy, limit, offset)
	} else {
		err = checkAccountError
	}

	return allTransactions, err
}
