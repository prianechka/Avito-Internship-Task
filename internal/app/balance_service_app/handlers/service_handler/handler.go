package service_handler

import (
	ac "Avito-Internship-Task/internal/app/balance_service_app/account/account_controller"
	"Avito-Internship-Task/internal/app/balance_service_app/handlers/response"
	"Avito-Internship-Task/internal/app/balance_service_app/handlers/service_handler/messages"
	"Avito-Internship-Task/internal/app/balance_service_app/manager"
	oc "Avito-Internship-Task/internal/app/balance_service_app/order/order_controller"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type ServiceHandler struct {
	logger  logrus.Logger
	manager manager.ManagerInterface
}

func CreateServiceHandler(man manager.ManagerInterface) *ServiceHandler {
	return &ServiceHandler{manager: man}
}

func (h *ServiceHandler) BuyService(w http.ResponseWriter, r *http.Request) {
	var buyParams messages.BuyServiceMessage
	var statusCode int
	var handleMessage string

	body, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		http.Error(w, "server problems", http.StatusInternalServerError)
		return
	}

	unmarshalError := json.Unmarshal(body, &buyParams)
	if unmarshalError != nil {
		http.Error(w, "unmarshal error", http.StatusInternalServerError)
		return
	}

	buyError := h.manager.BuyService(buyParams.UserID, buyParams.OrderID, buyParams.ServiceID,
		buyParams.Sum, buyParams.Comment)

	switch buyError {
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
		handleMessage = fmt.Sprintf("internal server error: %v", buyError)
	}
	response.SendShortResponse(w, statusCode, handleMessage)
	h.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s",
		r.Method, r.URL.Path, statusCode, handleMessage)
}

func (h *ServiceHandler) AcceptService(w http.ResponseWriter, r *http.Request) {
	var acceptParams messages.AcceptServiceMessage
	var statusCode int
	var handleMessage string

	body, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		http.Error(w, "server problems", http.StatusInternalServerError)
		return
	}

	unmarshalError := json.Unmarshal(body, &acceptParams)
	if unmarshalError != nil {
		http.Error(w, "unmarshal error", http.StatusInternalServerError)
		return
	}

	acceptBuy := h.manager.AcceptBuy(acceptParams.UserID, acceptParams.OrderID, acceptParams.ServiceID)

	switch acceptBuy {
	case nil:
		statusCode = http.StatusOK
		handleMessage = "OK"
	case ac.AccountNotExistErr:
		statusCode = http.StatusBadRequest
		handleMessage = fmt.Sprintf("%v", ac.AccountNotExistErr)
	case oc.OrderNotFound:
		statusCode = http.StatusBadRequest
		handleMessage = fmt.Sprintf("%v", oc.OrderNotFound)
	case oc.WrongStateError:
		statusCode = http.StatusBadRequest
		handleMessage = fmt.Sprintf("%v", oc.OrderNotFound)
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = fmt.Sprintf("internal server error: %v", acceptBuy)
	}
	response.SendShortResponse(w, statusCode, handleMessage)
	h.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s",
		r.Method, r.URL.Path, statusCode, handleMessage)
}

func (h *ServiceHandler) RefuseService(w http.ResponseWriter, r *http.Request) {
	var refuseParams messages.RefuseServiceMessage
	var statusCode int
	var handleMessage string

	body, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		http.Error(w, "server problems", http.StatusInternalServerError)
		return
	}

	unmarshalError := json.Unmarshal(body, &refuseParams)
	if unmarshalError != nil {
		http.Error(w, "unmarshal error", http.StatusInternalServerError)
		return
	}

	acceptBuy := h.manager.RefuseBuy(refuseParams.UserID, refuseParams.OrderID, refuseParams.ServiceID,
		refuseParams.Comment)

	switch acceptBuy {
	case nil:
		statusCode = http.StatusOK
		handleMessage = "OK"
	case ac.AccountNotExistErr:
		statusCode = http.StatusBadRequest
		handleMessage = fmt.Sprintf("%v", ac.AccountNotExistErr)
	case oc.OrderNotFound:
		statusCode = http.StatusBadRequest
		handleMessage = fmt.Sprintf("%v", oc.OrderNotFound)
	case oc.WrongStateError:
		statusCode = http.StatusBadRequest
		handleMessage = fmt.Sprintf("%v", oc.WrongStateError)
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = fmt.Sprintf("internal server error: %v", acceptBuy)
	}
	response.SendShortResponse(w, statusCode, handleMessage)
	h.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s",
		r.Method, r.URL.Path, statusCode, handleMessage)
}
