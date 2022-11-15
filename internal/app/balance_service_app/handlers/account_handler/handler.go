package account_handler

import (
	ac "Avito-Internship-Task/internal/app/balance_service_app/account/account_controller"
	"Avito-Internship-Task/internal/app/balance_service_app/handlers/account_handler/request_models"
	"Avito-Internship-Task/internal/app/balance_service_app/handlers/models"
	"Avito-Internship-Task/internal/app/balance_service_app/manager"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

type AccountHandler struct {
	logger  *logrus.Entry
	manager manager.ManagerInterface
}

func CreateAccountHandler(newManager manager.ManagerInterface) *AccountHandler {
	contextLogger := logrus.WithFields(logrus.Fields{})
	logrus.SetReportCaller(false)
	logrus.SetFormatter(&logrus.TextFormatter{PadLevelText: false, DisableLevelTruncation: false})
	return &AccountHandler{logger: contextLogger, manager: newManager}
}

// RefillBalance
// @Summary refill user balance
// @Description users refill balance in the app
// @Accept json
// @Produce json
// @Param data body request_models.RefillMessage true "body for transfer money"
// @Success 200 {object} models.ShortResponseMessage "OK"
// @Failure 400 {object} models.ShortResponseMessage "invalid body params"
// @Failure 401 {object} models.ShortResponseMessage "account is not exist"
// @Failure 422 {object} models.ShortResponseMessage "sum must be > 0"
// @Failure 500 {object} models.ShortResponseMessage "internal server error"
// @Router /api/v1/accounts/refill [POST]
func (h *AccountHandler) RefillBalance(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string

	var refillParams request_models.RefillMessage

	body, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	unmarshalError := json.Unmarshal(body, &refillParams)
	if unmarshalError != nil {
		http.Error(w, "invalid body params", http.StatusBadRequest)
		return
	}

	refillError := h.manager.RefillBalance(refillParams.UserID, refillParams.Sum, refillParams.Comment)

	switch refillError {
	case nil:
		statusCode = http.StatusOK
		handleMessage = "OK"
	case ac.AccountNotExistErr:
		statusCode = http.StatusUnauthorized
		handleMessage = fmt.Sprintf("%v", ac.AccountNotExistErr)
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = fmt.Sprintf("internal server error")
	}
	models.SendShortResponse(w, statusCode, handleMessage)
	h.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s, err = %v",
		r.Method, r.URL.Path, statusCode, handleMessage, refillError)
}

// Transfer
// @Summary transfer money from account to another account
// @Description money transfer between users
// @Accept json
// @Produce json
// @Param data body request_models.TransferMessage true "body for transfer money"
// @Success 200 {object} models.ShortResponseMessage "OK"
// @Failure 400 {object} models.ShortResponseMessage "invalid body params"
// @Failure 401 {object} models.ShortResponseMessage "account is not exist"
// @Failure 422 {object} models.ShortResponseMessage "not enough money"
// @Failure 500 {object} models.ShortResponseMessage "internal server error"
// @Router /api/v1/transfer [POST]
func (h *AccountHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string

	var transferParams request_models.TransferMessage

	body, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		http.Error(w, "server problems", http.StatusInternalServerError)
		return
	}

	unmarshalError := json.Unmarshal(body, &transferParams)
	if unmarshalError != nil {
		http.Error(w, "invalid body params", http.StatusBadRequest)
		return
	}

	transferError := h.manager.Transfer(transferParams.SrcUserID, transferParams.DstUserID,
		transferParams.Sum, transferParams.Comment)

	switch transferError {
	case nil:
		statusCode = http.StatusOK
		handleMessage = "OK"
	case ac.AccountNotExistErr:
		statusCode = http.StatusUnauthorized
		handleMessage = fmt.Sprintf("%v", ac.AccountNotExistErr)
	case ac.NotEnoughMoneyErr:
		statusCode = http.StatusUnprocessableEntity
		handleMessage = fmt.Sprintf("%v", ac.NotEnoughMoneyErr)
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = fmt.Sprintf("internal server error: %v", transferError)
	}
	models.SendShortResponse(w, statusCode, handleMessage)
	h.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s",
		r.Method, r.URL.Path, statusCode, handleMessage)
}

// GetBalance
// @Summary get user balance
// @Description get user balance
// @Produce json
// @Param userID query int true "user_id in balanceApp"
// @Success 200 {object} models.BalanceResponseMessage
// @Failure 400 {object} models.ShortResponseMessage "userID not found | userID isn't number"
// @Failure 401 {object} models.ShortResponseMessage "account is not exist"
// @Failure 500 {object} models.ShortResponseMessage "internal server error"
// @Router /api/v1/accounts [GET]
func (h *AccountHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string

	strUserID := r.URL.Query().Get("userID")

	if strUserID == "" {
		models.SendShortResponse(w, http.StatusBadRequest, "userID not found")
		return
	}

	userID, err := strconv.Atoi(strUserID)
	if err != nil {
		models.SendShortResponse(w, http.StatusBadRequest, "userID isn't number")
		return
	}

	balance, getBalanceErr := h.manager.GetUserBalance(int64(userID))
	switch getBalanceErr {
	case nil:
		models.BalanceResponse(w, balance, "OK")
		return
	case ac.AccountNotExistErr:
		statusCode = http.StatusUnauthorized
		handleMessage = fmt.Sprintf("%v", ac.AccountNotExistErr)
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = fmt.Sprintf("internal server error")
	}
	models.SendShortResponse(w, statusCode, handleMessage)
	h.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s",
		r.Method, r.URL.Path, statusCode, handleMessage)
}
