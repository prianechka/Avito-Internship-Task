package account_repo

//go:generate mockgen -source=interface.go -destination=mocks/account_repo_mock.go -package=mocks AccountRepoInterface
type AccountRepoInterface interface {
	AddNewAccount(accountID int64) error
	GetCurrentAmount(accountID int64) (amount float64, err error)
	ChangeAmount(accountID int64, delta float64) error
	DeleteAccount(accountID int64) error
}
