package account_handler

import (
	ac "Avito-Internship-Task/internal/app/balance_service_app/account/account_controller"
	"Avito-Internship-Task/internal/app/balance_service_app/handlers/account_handler/messages"
	"Avito-Internship-Task/internal/app/balance_service_app/handlers/response"
	"Avito-Internship-Task/internal/app/balance_service_app/manager"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type AccountHandler struct {
	manager manager.ManagerInterface
}

func CreateAccountHandler(newManager manager.ManagerInterface) *AccountHandler {
	return &AccountHandler{manager: newManager}
}

func (h *AccountHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)

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
	case ac.AccountNotExistErr:
		response.SendShortResponse(w, http.StatusBadRequest, fmt.Sprintf("%v", ac.AccountNotExistErr))
	case ac.NegSumError:
		response.SendShortResponse(w, http.StatusBadRequest, fmt.Sprintf("%v", ac.NegSumError))
	default:
		response.SendShortResponse(w, http.StatusInternalServerError, fmt.Sprintf("internal server error: %v", getBalanceErr))
	}
}

func (h *AccountHandler) RefillBalance(w http.ResponseWriter, r *http.Request) {
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
		response.SendShortResponse(w, http.StatusOK, "OK")
	case ac.AccountNotExistErr:
		response.SendShortResponse(w, http.StatusBadRequest, fmt.Sprintf("%v", ac.AccountNotExistErr))
	default:
		response.SendShortResponse(w, http.StatusInternalServerError, fmt.Sprintf("internal server error: %v", refillError))
	}
}

func (h *AccountHandler) Transfer(w http.ResponseWriter, r *http.Request) {
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
		response.SendShortResponse(w, http.StatusOK, "OK")
	case ac.AccountNotExistErr:
		response.SendShortResponse(w, http.StatusBadRequest, fmt.Sprintf("%v", ac.AccountNotExistErr))
	case ac.NotEnoughMoneyErr:
		response.SendShortResponse(w, http.StatusBadRequest, fmt.Sprintf("%v", ac.NotEnoughMoneyErr))
	default:
		response.SendShortResponse(w, http.StatusInternalServerError, fmt.Sprintf("internal server error: %v", transferError))
	}
}
