package report_controller

import (
	"Avito-Internship-Task/internal/app/balance_service_app/order"
	"Avito-Internship-Task/internal/app/balance_service_app/report"
	"Avito-Internship-Task/internal/pkg/utils"
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
		writer.Comma = utils.DefaultSeparator
		for _, record := range serviceReport {
			serviceName := order.Types[record.ServiceType]
			if serviceName == utils.EmptyString {
				serviceName = fmt.Sprintf("Service with id %d", record.ServiceType)
			}
			err = writer.Write([]string{serviceName, fmt.Sprintf("%f", record.Sum)})
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
