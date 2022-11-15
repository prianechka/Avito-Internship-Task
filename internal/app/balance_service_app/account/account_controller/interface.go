package account_controller

//go:generate mockgen -source=interface.go -destination=mocks/controller_mocks.go -package=mocks AccountControllerInterface
type AccountControllerInterface interface {
	CheckAccountIsExist(userID int) (result bool, err error)
	CreateNewAccount(userID int) error
	CheckBalance(userID int) (float64, error)
	CheckAbleToBuyService(userID int, servicePrice float64) (bool, error)
	DonateMoney(userID int, sum float64) (err error)
	SpendMoney(userID int, sum float64) error
}
