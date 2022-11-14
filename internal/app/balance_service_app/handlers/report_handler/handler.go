package report_handler

import (
	"Avito-Internship-Task/internal/app/balance_service_app/handlers/models"
	"Avito-Internship-Task/internal/app/balance_service_app/manager"
	"Avito-Internship-Task/internal/app/balance_service_app/order/order_controller"
	"Avito-Internship-Task/internal/pkg/utils"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type ReportHandler struct {
	logger  *logrus.Entry
	manager manager.ManagerInterface
}

func CreateReportHandler(man manager.ManagerInterface) *ReportHandler {
	contextLogger := logrus.WithFields(logrus.Fields{})
	logrus.SetReportCaller(false)
	logrus.SetFormatter(&logrus.TextFormatter{PadLevelText: false, DisableLevelTruncation: false})
	return &ReportHandler{manager: man, logger: contextLogger}
}

// GetFinanceReport
// @Summary get finance report of company
// @Description generate report of company's revenue with csv format
// @Produce json
// @Param month query int true "required month"
// @Param year query int true "required year"
// @Success 200 {object} models.FinanceReportResponseMessage
// @Failure 400 {object} models.ShortResponseMessage "month not found | year not found | incorrect month | incorrect year"
// @Failure 500 {object} models.ShortResponseMessage "internal server error"
// @Router /api/v1//reports/finance [GET]
func (h *ReportHandler) GetFinanceReport(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string

	strMonth := r.URL.Query().Get("month")

	if strMonth == "" {
		models.SendShortResponse(w, http.StatusBadRequest, "month not found")
		return
	}

	strYear := r.URL.Query().Get("year")

	if strYear == "" {
		models.SendShortResponse(w, http.StatusBadRequest, "year not found")
		return
	}

	month, castMonthErr := strconv.Atoi(strMonth)
	if castMonthErr != nil {
		models.SendShortResponse(w, http.StatusBadRequest, "incorrect month")
		return
	}

	year, castYearErr := strconv.Atoi(strYear)
	if castYearErr != nil {
		models.SendShortResponse(w, http.StatusBadRequest, "incorrect year")
		return
	}

	fileURL := utils.FileURL

	getFinanceReportError := h.manager.GetFinanceReport(int64(month), int64(year), fileURL)
	switch getFinanceReportError {
	case nil:
		models.FinanceReportResponse(w, utils.FileURL)
		return
	case order_controller.BadYearError:
		statusCode = http.StatusBadRequest
		handleMessage = fmt.Sprintf("%v", order_controller.BadYearError)
	case order_controller.BadMonthError:
		statusCode = http.StatusBadRequest
		handleMessage = fmt.Sprintf("%v", order_controller.BadMonthError)
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = fmt.Sprintf("internal server error")
	}
	models.SendShortResponse(w, statusCode, handleMessage)
	h.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s",
		r.Method, r.URL.Path, statusCode, handleMessage)
}

func (h *ReportHandler) GetUserReport(w http.ResponseWriter, r *http.Request) {}
