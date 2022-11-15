package report_handler

import "net/http"

//go:generate mockgen -source=interface.go -destination=mocks/handler_mock.go -package=mocks ReportHandlerInterface
type ReportHandlerInterface interface {
	GetFinanceReport(w http.ResponseWriter, r *http.Request)
	GetUserReport(w http.ResponseWriter, r *http.Request)
}
