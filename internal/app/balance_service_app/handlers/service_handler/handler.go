package service_handler

import (
	ac "Avito-Internship-Task/internal/app/balance_service_app/account/account_controller"
	"Avito-Internship-Task/internal/app/balance_service_app/handlers/response"
	"Avito-Internship-Task/internal/app/balance_service_app/handlers/service_handler/messages"
	"Avito-Internship-Task/internal/app/balance_service_app/manager"
	oc "Avito-Internship-Task/internal/app/balance_service_app/order/order_controller"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ServiceHandler struct {
	manager manager.ManagerInterface
}

func CreateServiceHandler(man manager.ManagerInterface) *ServiceHandler {
	return &ServiceHandler{manager: man}
}

func (h *ServiceHandler) BuyService(w http.ResponseWriter, r *http.Request) {
	var buyParams messages.BuyServiceMessage

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
		response.SendShortResponse(w, http.StatusOK, "OK")
	case ac.AccountNotExistErr:
		response.SendShortResponse(w, http.StatusBadRequest, fmt.Sprintf("%v", ac.AccountNotExistErr))
	case ac.NotEnoughMoneyErr:
		response.SendShortResponse(w, http.StatusBadRequest, fmt.Sprintf("%v", ac.NotEnoughMoneyErr))
	default:
		response.SendShortResponse(w, http.StatusInternalServerError, fmt.Sprintf("internal server error: %v", buyError))
	}
}

func (h *ServiceHandler) AcceptService(w http.ResponseWriter, r *http.Request) {
	var acceptParams messages.AcceptServiceMessage

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
		response.SendShortResponse(w, http.StatusOK, "OK")
	case ac.AccountNotExistErr:
		response.SendShortResponse(w, http.StatusBadRequest, fmt.Sprintf("%v", ac.AccountNotExistErr))
	case oc.OrderNotFound:
		response.SendShortResponse(w, http.StatusBadRequest, fmt.Sprintf("%v", oc.OrderNotFound))
	case oc.WrongStateError:
		response.SendShortResponse(w, http.StatusBadRequest, fmt.Sprintf("%v", oc.WrongStateError))
	default:
		response.SendShortResponse(w, http.StatusInternalServerError, fmt.Sprintf("internal server error: %v", acceptBuy))
	}
}

func (h *ServiceHandler) RefuseService(w http.ResponseWriter, r *http.Request) {
	var refuseParams messages.RefuseServiceMessage

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
		response.SendShortResponse(w, http.StatusOK, "OK")
	case ac.AccountNotExistErr:
		response.SendShortResponse(w, http.StatusBadRequest, fmt.Sprintf("%v", ac.AccountNotExistErr))
	case oc.OrderNotFound:
		response.SendShortResponse(w, http.StatusBadRequest, fmt.Sprintf("%v", oc.OrderNotFound))
	case oc.WrongStateError:
		response.SendShortResponse(w, http.StatusBadRequest, fmt.Sprintf("%v", oc.WrongStateError))
	default:
		response.SendShortResponse(w, http.StatusInternalServerError, fmt.Sprintf("internal server error: %v", acceptBuy))
	}
}
