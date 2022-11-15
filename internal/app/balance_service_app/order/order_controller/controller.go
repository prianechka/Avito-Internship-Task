package order_controller

import (
	"Avito-Internship-Task/internal/app/balance_service_app/order"
	"Avito-Internship-Task/internal/app/balance_service_app/order/order_repo"
	"Avito-Internship-Task/internal/app/balance_service_app/report"
	"Avito-Internship-Task/internal/pkg/utils"
	"database/sql"
	"sync"
	"time"
)

type OrderController struct {
	mutex sync.RWMutex
	repo  order_repo.OrderRepoInterface
}

func CreateNewOrderController(repo order_repo.OrderRepoInterface) *OrderController {
	return &OrderController{mutex: sync.RWMutex{}, repo: repo}
}

func (c *OrderController) GetOrder(orderID, userID, serviceID int) (order.Order, error) {
	return c.repo.GetOrderByID(orderID, userID, serviceID)
}

func (c *OrderController) CreateNewOrder(orderID, userID, serviceID int, sum float64, comment string) error {
	isExist, err := c.CheckOrderIsExist(orderID, userID, serviceID)
	if err == nil {
		if !isExist {
			order := order.Order{
				OrderID:      orderID,
				UserID:       userID,
				ServiceID:    serviceID,
				OrderCost:    sum,
				CreatingTime: time.Now(),
				Comment:      comment,
				OrderState:   order.REGISTRATED,
			}

			c.mutex.Lock()
			err = c.repo.CreateOrder(order)
			c.mutex.Unlock()
		} else {
			err = OrderIsAlreadyExist
		}
	}
	return err
}

func (c *OrderController) CheckOrderIsExist(orderID, userID, serviceID int) (bool, error) {
	var result bool

	c.mutex.Lock()
	foundOrder, err := c.repo.GetOrderByID(orderID, userID, serviceID)
	c.mutex.Unlock()

	if err == nil {
		if foundOrder.OrderID == orderID && foundOrder.UserID == userID && foundOrder.ServiceID == serviceID {
			result = true
		} else {
			result = false
		}
	} else if err == sql.ErrNoRows {
		err = nil
		result = false
	}
	return result, err
}

func (c *OrderController) ReserveOrder(orderID, userID, serviceID int) error {
	isOrderExist, err := c.CheckOrderIsExist(orderID, userID, serviceID)

	if err == nil {
		if isOrderExist {
			curOrder, getOrderErr := c.GetOrder(orderID, userID, serviceID)
			if getOrderErr == nil && curOrder.OrderState == order.REGISTRATED {
				err = c.repo.ChangeOrderState(orderID, userID, serviceID, order.RESERVED)
			} else {
				if getOrderErr != nil {
					err = GetOrderError
				} else {
					err = WrongStateError
				}
			}
		}
	}
	return err
}

func (c *OrderController) FinishOrder(orderID, userID, serviceID int) error {
	isOrderExist, err := c.CheckOrderIsExist(orderID, userID, serviceID)

	if err == nil {
		if isOrderExist {
			curOrder, getOrderErr := c.GetOrder(orderID, userID, serviceID)
			if getOrderErr == nil && curOrder.OrderState == order.RESERVED {
				err = c.repo.ChangeOrderState(orderID, userID, serviceID, order.FINISHED)
			} else {
				if getOrderErr != nil {
					err = GetOrderError
				} else {
					err = WrongStateError
				}
			}
		}
	}
	return err
}

func (c *OrderController) ReturnOrder(orderID, userID, serviceID int) (float64, error) {
	var sum float64
	isOrderExist, err := c.CheckOrderIsExist(orderID, userID, serviceID)

	if err == nil {
		if isOrderExist {
			curOrder, getOrderErr := c.GetOrder(orderID, userID, serviceID)
			if getOrderErr == nil && curOrder.OrderState == order.RESERVED {
				err = c.repo.ChangeOrderState(orderID, userID, serviceID, order.RETURNED)
				if err == nil {
					sum = curOrder.OrderCost
				}
			} else {
				if getOrderErr != nil {
					err = GetOrderError
				} else {
					err = WrongStateError
				}
			}
		}
	}
	return sum, err
}

func (c *OrderController) GetFinanceReports(month, year int) ([]report.FinanceReport, error) {
	var report = make([]report.FinanceReport, utils.EMPTY)
	var err error = nil
	if month < utils.January || month > utils.December {
		err = BadMonthError
	} else if year < utils.FirstYear || year > utils.LastYear {
		err = BadYearError
	} else {
		report, err = c.repo.GetSumOfFinishedServices(month, year)
	}
	return report, err
}
