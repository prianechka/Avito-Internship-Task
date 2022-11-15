package report_manager

import "Avito-Internship-Task/internal/app/balance_service_app/transaction"

//go:generate mockgen -source=interface.go -destination=mocks/manager_mock.go -package=mocks ReportManagerInterface
type ReportManagerInterface interface {
	GetFinanceReport(month, year int, url string) error
	GetUserReport(userID int, orderBy string, limit, offset int) ([]transaction.Transaction, error)
}
