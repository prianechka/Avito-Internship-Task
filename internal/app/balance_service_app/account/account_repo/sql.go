package account_repo

type MySQLAddNewAccount struct{}
type MySQLGetCurrentAmount struct{}
type MySQLChangeAmount struct{}

func (sql MySQLAddNewAccount) GetString() string {
	return "INSERT INTO balanceApp.accounts(`userID`, `amount`) VALUES (?, 0);"
}

func (sql MySQLGetCurrentAmount) GetString() string {
	return "SELECT amount FROM balanceApp.accounts WHERE userID = ?;"
}

func (sql MySQLChangeAmount) GetString() string {
	return "UPDATE balanceApp.accounts SET amount = amoumt + ? WHERE userID = ?;"
}
