package report_handler

import (
	ac "Avito-Internship-Task/internal/app/balance_service_app/account/account_controller"
	"Avito-Internship-Task/internal/app/balance_service_app/handlers/models"
	"Avito-Internship-Task/internal/app/balance_service_app/managers/report_manager"
	"Avito-Internship-Task/internal/app/balance_service_app/order/order_controller"
	"Avito-Internship-Task/internal/pkg/utils"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type ReportHandler struct {
	logger  *logrus.Entry
	manager report_manager.ReportManagerInterface
}

func CreateReportHandler(man report_manager.ReportManagerInterface) *ReportHandler {
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
// @Router /api/v1/reports/finance [GET]
func (h *ReportHandler) GetFinanceReport(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string

	monthFromQuery := r.URL.Query().Get("month")

	if monthFromQuery == utils.EmptyString {
		models.SendShortResponse(w, http.StatusBadRequest, "month not found")
		return
	}

	yearFromQuery := r.URL.Query().Get("year")

	if yearFromQuery == utils.EmptyString {
		models.SendShortResponse(w, http.StatusBadRequest, "year not found")
		return
	}

	month, castMonthErr := strconv.Atoi(monthFromQuery)
	if castMonthErr != nil {
		models.SendShortResponse(w, http.StatusBadRequest, "incorrect month")
		return
	}

	year, castYearErr := strconv.Atoi(yearFromQuery)
	if castYearErr != nil {
		models.SendShortResponse(w, http.StatusBadRequest, "incorrect year")
		return
	}

	fileURL := utils.FileURL

	getFinanceReportErr := h.manager.GetFinanceReport(month, year, fileURL)
	switch getFinanceReportErr {
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
	h.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s, err = %v",
		r.Method, r.URL.Path, statusCode, handleMessage, getFinanceReportErr)
}

// GetUserReport
// @Summary get report of user operations
// @Description get report of all user financial operations with pagination
// @Produce json
// @Param userID query int true "userID for report"
// @Param orderBy query string false "required field to sort"
// @Param limit query int false "limit to paginate query"
// @Param offset query int false "offset to paginate query"
// @Success 200 {object} models.UserReportResponseMessage
// @Failure 400 {object} models.ShortResponseMessage "userID not found | userID isn't number | limit isn't number | offset isn't number"
// @Failure 401 {object} models.ShortResponseMessage "account is not exist"
// @Failure 500 {object} models.ShortResponseMessage "internal server error"
// @Router /api/v1/reports/user [GET]
func (h *ReportHandler) GetUserReport(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string

	userIDFromQuery := r.URL.Query().Get("userID")

	if userIDFromQuery == utils.EmptyString {
		models.SendShortResponse(w, http.StatusBadRequest, "userID not found")
		return
	}

	userID, err := strconv.Atoi(userIDFromQuery)
	if err != nil {
		models.SendShortResponse(w, http.StatusBadRequest, "userID isn't number")
		return
	}

	orderBy := r.URL.Query().Get("orderBy")

	switch orderBy {
	case utils.DefaultOrderBy:
	case utils.FieldOrderTime:
	case utils.FieldOrderSum:
	default:
		orderBy = utils.DefaultOrderBy
	}

	limit, getLimitErr := utils.GetOptionalIntParam(r, "limit")
	if getLimitErr != nil {
		models.SendShortResponse(w, http.StatusBadRequest, "limit isn't number")
		return
	}

	offset, getOffsetErr := utils.GetOptionalIntParam(r, "offset")
	if getOffsetErr != nil {
		models.SendShortResponse(w, http.StatusBadRequest, "offset isn't number")
		return
	}

	userReport, getReportErr := h.manager.GetUserReport(userID, orderBy, limit, offset)
	switch getReportErr {
	case nil:
		models.UserReportResponse(w, userReport)
		return
	case ac.AccountNotExistErr:
		statusCode = http.StatusUnauthorized
		handleMessage = fmt.Sprintf("%v", ac.AccountNotExistErr)
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = fmt.Sprintf("internal server error")
	}
	models.SendShortResponse(w, statusCode, handleMessage)
	h.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s, error = %v",
		r.Method, r.URL.Path, statusCode, handleMessage, getReportErr)
}
