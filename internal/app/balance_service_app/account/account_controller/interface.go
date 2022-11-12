package account_controller

//go:generate mockgen -source=interface.go -destination=mocks/controller_mocks.go -package=mocks AccountControllerInterface
type AccountControllerInterface interface {
	CheckAccountIsExist(userID int64) (result bool, err error)
	CreateNewAccount(userID int64) error
	CheckBalance(userID int64) (float64, error)
	CheckAbleToBuyService(userID int64, servicePrice float64) (bool, error)
	DonateMoney(userID int64, sum float64) (err error)
	SpendMoney(userID int64, sum float64) error
}
