package report_handler

import (
	"Avito-Internship-Task/internal/app/balance_service_app/manager"
	"net/http"
)

type ReportHandler struct {
	Manager manager.ManagerInterface
}

func (h *ReportHandler) GetFinanceReport(w http.ResponseWriter, r *http.Request) {}

func (h *ReportHandler) GetUserReport(w http.ResponseWriter, r *http.Request) {}
