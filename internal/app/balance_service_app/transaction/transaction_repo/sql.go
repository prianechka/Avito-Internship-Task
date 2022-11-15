package transaction_repo

type MySQLAddNewTransaction struct{}
type MySQLGetAllTransactions struct{}
type MySQLGetUserTransactions struct{}
type MySQLGetTransactionByID struct{}

func (sql MySQLAddNewTransaction) GetString() string {
	return "INSERT INTO balanceApp.transactions(`transactionID`, `userID`, `transactionType`, `sum`, " +
		"`time`, `actionComments`, `addComments`) VALUES (?, ?, ?, ?, ?, ?, ?);"
}

func (sql MySQLGetAllTransactions) GetString() string {
	return "SELECT transactionID, userID, transactionType, sum, time," +
		" actionComments, addComments FROM balanceApp.transactions"
}

func (sql MySQLGetUserTransactions) GetString() string {
	return "SELECT transactionID, userID, transactionType, sum, time," +
		" actionComments, addComments FROM balanceApp.transactions WHERE userID = ? ORDER BY ? DESC LIMIT ? OFFSET ?"
}

func (sql MySQLGetTransactionByID) GetString() string {
	return "SELECT transactionID, userID, transactionType, sum, time," +
		" actionComments, addComments FROM balanceApp.transactions WHERE transactionID = ?"
}
