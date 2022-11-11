package account_repo

//go:generate mockgen -source=interface.go -destination=mocks/account_repo_mock.go -package=mocks AccountRepoInterface
type AccountRepoInterface interface {
	AddNewAccount(userID int64) error
	GetCurrentAmount(userID int64) (amount float64, err error)
	ChangeAmount(userID int64, delta float64) error
	DeleteAccount(userID int64) error
}
