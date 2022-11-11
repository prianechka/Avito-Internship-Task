package order_repo

type MySQLAddNewOrder struct{}
type MySQLGetAllOrders struct{}
type MySQLGetOrderByID struct{}
type MySQLGetUserOrders struct{}
type MySQLGetServiceOrders struct{}
type MySQLChangeOrderState struct{}

func (sql MySQLAddNewOrder) GetString() string {
	return "INSERT INTO balanceApp.orders(`orderID`, `userID`, `serviceType`, `orderCost`, " +
		"`creatingTime`, `comments`, `orderState`) VALUES (?, ?, ?, ?, ?, ?, ?);"
}

func (sql MySQLGetAllOrders) GetString() string {
	return "SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders;"
}

func (sql MySQLGetOrderByID) GetString() string {
	return "SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders WHERE orderID = ?;"
}

func (sql MySQLGetUserOrders) GetString() string {
	return "SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders WHERE userID = ?;"
}

func (sql MySQLGetServiceOrders) GetString() string {
	return "SELECT orderID, userID, serviceType, orderCost, creatingTime, comments, orderState FROM balanceApp.orders WHERE orderType = ?;"
}

func (sql MySQLChangeOrderState) GetString() string {
	return "UPDATE balanceApp.orders SET orderState = ? WHERE orderID = ?"
}
