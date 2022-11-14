package report_controller

import "Avito-Internship-Task/internal/app/balance_service_app/report"

//go:generate mockgen -source=interface.go -destination=mocks/report_controller_mock.go -package=mocks ReportControllerInterface
type ReportControllerInterface interface {
	CreateFinancialReportCSV([]report.FinanceReport, string) error
}
