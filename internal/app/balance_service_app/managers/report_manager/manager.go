package report_manager

import (
	ac "Avito-Internship-Task/internal/app/balance_service_app/account/account_controller"
	oc "Avito-Internship-Task/internal/app/balance_service_app/order/order_controller"
	"Avito-Internship-Task/internal/app/balance_service_app/report/report_controller"
	"Avito-Internship-Task/internal/app/balance_service_app/transaction"
	tc "Avito-Internship-Task/internal/app/balance_service_app/transaction/transaction_controller"
	"Avito-Internship-Task/internal/pkg/utils"
)

type ReportManager struct {
	accountController     ac.AccountControllerInterface
	orderController       oc.OrderControllerInterface
	transactionController tc.TransactionControllerInterface
}

func CreateNewReportManager(accController ac.AccountControllerInterface, orderController oc.OrderControllerInterface,
	transactionController tc.TransactionControllerInterface) *ReportManager {
	return &ReportManager{
		accountController:     accController,
		orderController:       orderController,
		transactionController: transactionController,
	}
}

func (m *ReportManager) GetFinanceReport(month, year int, url string) error {
	dataToReport, err := m.orderController.GetFinanceReports(month, year)
	if err == nil {
		reportController := report_controller.CreateNewReportController()
		err = reportController.CreateFinancialReportCSV(dataToReport, url)
	}
	return err
}

func (m *ReportManager) GetUserReport(userID int, orderBy string, limit, offset int) ([]transaction.Transaction, error) {
	var allTransactions = make([]transaction.Transaction, utils.EMPTY)
	var err error

	if limit == utils.NotInQuery {
		limit = utils.DefaultLimit
	}
	if offset == utils.NotInQuery {
		offset = utils.DefaultOffset
	}
	if orderBy == utils.EmptyString {
		orderBy = utils.DefaultOrderBy
	}

	_, checkAccountError := m.accountController.CheckAccountIsExist(userID)
	if checkAccountError == nil {
		allTransactions, err = m.transactionController.GetUserTransactions(userID, orderBy, limit, offset)
	} else {
		err = checkAccountError
	}

	return allTransactions, err
}
