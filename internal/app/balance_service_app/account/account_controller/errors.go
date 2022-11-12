package account_controller

import "errors"

var (
	NegSumError        = errors.New("bad sum to update")
	NotEnoughMoneyErr  = errors.New("not enough money")
	AccountIsExistErr  = errors.New("account is already exist")
	AccountNotExistErr = errors.New("account is not exist")
)
