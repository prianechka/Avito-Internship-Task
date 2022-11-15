package order_manager

import (
	ac "Avito-Internship-Task/internal/app/balance_service_app/account/account_controller"
	oc "Avito-Internship-Task/internal/app/balance_service_app/order/order_controller"
	tc "Avito-Internship-Task/internal/app/balance_service_app/transaction/transaction_controller"
)

type OrderManager struct {
	accountController     ac.AccountControllerInterface
	orderController       oc.OrderControllerInterface
	transactionController tc.TransactionControllerInterface
}

func CreateNewOrderManager(accController ac.AccountControllerInterface, orderController oc.OrderControllerInterface,
	transactionController tc.TransactionControllerInterface) *OrderManager {
	return &OrderManager{
		accountController:     accController,
		orderController:       orderController,
		transactionController: transactionController,
	}
}

func (m *OrderManager) BuyService(userID, orderID, serviceID int, sum float64, comment string) error {
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

func (m *OrderManager) checkUserCanBuyService(userID int, sum float64) (bool, error) {
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

func (m *OrderManager) createOrder(userID, orderID, serviceID int, sum float64, comment string) error {
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

func (m *OrderManager) AcceptBuy(userID, orderID, serviceID int) error {
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

func (*OrderManager) transferMoneyToCompanyAccount() {}

func (m *OrderManager) RefuseBuy(userID, orderID, serviceID int, comment string) error {
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
