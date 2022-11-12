package account_handler

import (
	ac "Avito-Internship-Task/internal/app/balance_service_app/account/account_controller"
	"Avito-Internship-Task/internal/app/balance_service_app/handlers/account_handler/messages"
	"Avito-Internship-Task/internal/app/balance_service_app/handlers/response"
	"Avito-Internship-Task/internal/app/balance_service_app/manager"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

type AccountHandler struct {
	logger  logrus.Logger
	manager manager.ManagerInterface
}

func CreateAccountHandler(newManager manager.ManagerInterface) *AccountHandler {
	return &AccountHandler{manager: newManager}
}

func (h *AccountHandler) RefillBalance(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string

	var refillParams messages.RefillParams

	body, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		http.Error(w, "server problems", http.StatusInternalServerError)
		return
	}

	unmarshalError := json.Unmarshal(body, &refillParams)
	if unmarshalError != nil {
		http.Error(w, "unmarshal error", http.StatusInternalServerError)
		return
	}

	refillError := h.manager.RefillBalance(refillParams.UserID, refillParams.Sum, refillParams.Comment)

	switch refillError {
	case nil:
		statusCode = http.StatusOK
		handleMessage = "OK"
	case ac.AccountNotExistErr:
		statusCode = http.StatusBadRequest
		handleMessage = fmt.Sprintf("%v", ac.AccountNotExistErr)
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = fmt.Sprintf("internal server error: %v", refillError)
	}
	response.SendShortResponse(w, statusCode, handleMessage)
	h.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s",
		r.Method, r.URL.Path, statusCode, handleMessage)
}

func (h *AccountHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string

	var transferParams messages.TransferMessage

	body, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		http.Error(w, "server problems", http.StatusInternalServerError)
		return
	}

	unmarshalError := json.Unmarshal(body, &transferParams)
	if unmarshalError != nil {
		http.Error(w, "unmarshal error", http.StatusInternalServerError)
		return
	}

	transferError := h.manager.Transfer(transferParams.SrcUserID, transferParams.DstUserID,
		transferParams.Sum, transferParams.Comment)

	switch transferError {
	case nil:
		statusCode = http.StatusOK
		handleMessage = "OK"
	case ac.AccountNotExistErr:
		statusCode = http.StatusBadRequest
		handleMessage = fmt.Sprintf("%v", ac.AccountNotExistErr)
	case ac.NotEnoughMoneyErr:
		statusCode = http.StatusBadRequest
		handleMessage = fmt.Sprintf("%v", ac.NotEnoughMoneyErr)
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = fmt.Sprintf("internal server error: %v", transferError)
	}
	response.SendShortResponse(w, statusCode, handleMessage)
	h.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s",
		r.Method, r.URL.Path, statusCode, handleMessage)
}

func (h *AccountHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string

	strUserID := r.URL.Query().Get("id")

	if strUserID == "" {
		response.SendShortResponse(w, http.StatusBadRequest, "userID not found")
		return
	}

	userID, err := strconv.Atoi(strUserID)
	if err != nil {
		response.SendShortResponse(w, http.StatusBadRequest, "userID isn't number")
		return
	}

	balance, getBalanceErr := h.manager.GetUserBalance(int64(userID))
	switch getBalanceErr {
	case nil:
		response.BalanceResponse(w, balance, "OK")
		return
	case ac.AccountNotExistErr:
		statusCode = http.StatusBadRequest
		handleMessage = fmt.Sprintf("%v", ac.AccountNotExistErr)
	case ac.NegSumError:
		statusCode = http.StatusBadRequest
		handleMessage = fmt.Sprintf("%v", ac.NegSumError)
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = fmt.Sprintf("internal server error: %v", getBalanceErr)
	}
	response.SendShortResponse(w, statusCode, handleMessage)
	h.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s",
		r.Method, r.URL.Path, statusCode, handleMessage)
}
