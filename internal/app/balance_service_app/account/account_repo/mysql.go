package account_repo

type MySQLAddNewAccount struct{}
type MySQLGetCurrentAmount struct{}
type MySQLChangeAmount struct{}
type MySQLDeleteAccount struct{}

func (sql MySQLAddNewAccount) GetString() string {
	return "INSERT INTO balanceApp.accounts(`accountID`, `amount`) VALUES (?, 0);"
}

func (sql MySQLGetCurrentAmount) GetString() string {
	return "SELECT amount FROM balanceApp.accounts WHERE accountID = ?;"
}

func (sql MySQLChangeAmount) GetString() string {
	return "UPDATE balanceApp.accounts SET amount = amoumt + ? WHERE accountID = ?;"
}

func (sql MySQLDeleteAccount) GetString() string {
	return "DELETE FROM balanceApp.accounts WHERE accountID = ?;"
}
