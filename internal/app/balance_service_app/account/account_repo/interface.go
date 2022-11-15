package account_repo

//go:generate mockgen -source=interface.go -destination=mocks/account_repo_mock.go -package=mocks AccountRepoInterface
type AccountRepoInterface interface {
	AddNewAccount(userID int) error
	GetCurrentAmount(userID int) (amount float64, err error)
	ChangeAmount(userID int, delta float64) error
}
