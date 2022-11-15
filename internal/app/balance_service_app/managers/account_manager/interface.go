package account_manager

//go:generate mockgen -source=interface.go -destination=mocks/manager_mock.go -package=mocks AccountManagerInterface
type AccountManagerInterface interface {
	RefillBalance(userID int, sum float64, comments string) error
	GetUserBalance(userID int) (float64, error)
	Transfer(srcUserID, dstUserID int, sum float64, comment string) error
}
