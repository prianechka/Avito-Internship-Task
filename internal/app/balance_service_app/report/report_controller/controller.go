package report_controller

import (
	"Avito-Internship-Task/internal/app/balance_service_app/order"
	"Avito-Internship-Task/internal/app/balance_service_app/report"
	"encoding/csv"
	"fmt"
	"os"
)

type ReportController struct{}

func CreateNewReportController() *ReportController {
	return &ReportController{}
}

func (c *ReportController) CreateFinancialReportCSV(serviceReport []report.FinanceReport, fileURL string) error {
	csvFile, err := os.Create(fileURL)
	if err == nil {
		defer csvFile.Close()
		writer := csv.NewWriter(csvFile)
		for _, record := range serviceReport {
			err = writer.Write([]string{order.Types[record.ServiceType], fmt.Sprintf("%f", record.Sum)})
			if err != nil {
				break
			}
		}
		writer.Flush()
	} else {
		err = BadFilePath
	}
	return err
}
