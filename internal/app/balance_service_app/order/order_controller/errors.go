package order_controller

import "errors"

var (
	OrderNotFound       = errors.New("order not found")
	OrderIsAlreadyExist = errors.New("order is already exist")
	GetOrderError       = errors.New("bad order get")
	WrongStateError     = errors.New("state isn't right to change order state")
	BadMonthError       = errors.New("bad month")
	BadYearError        = errors.New("bad year")
)
