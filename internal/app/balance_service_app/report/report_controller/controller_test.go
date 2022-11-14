package report_controller

import (
	"Avito-Internship-Task/internal/app/balance_service_app/report"
	"testing"
)

// TestCreateNewReport проверяет, что контроллер корректно создаёт отчёт.
func TestCreateNewReport(t *testing.T) {
	allServicesReport := []report.FinanceReport{{1, 100}, {2, 150}}

	controller := CreateNewReportController()

	execErr := controller.CreateFinancialReportCSV(allServicesReport, "report.csv")
	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}
}
