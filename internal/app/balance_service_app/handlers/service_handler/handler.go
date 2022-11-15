package service_handler

import (
	ac "Avito-Internship-Task/internal/app/balance_service_app/account/account_controller"
	"Avito-Internship-Task/internal/app/balance_service_app/handlers/models"
	"Avito-Internship-Task/internal/app/balance_service_app/handlers/service_handler/request_models"
	"Avito-Internship-Task/internal/app/balance_service_app/managers/order_manager"
	oc "Avito-Internship-Task/internal/app/balance_service_app/order/order_controller"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type ServiceHandler struct {
	logger  *logrus.Entry
	manager order_manager.OrderManagerInterface
}

func CreateServiceHandler(man order_manager.OrderManagerInterface) *ServiceHandler {
	contextLogger := logrus.WithFields(logrus.Fields{})
	logrus.SetReportCaller(false)
	logrus.SetFormatter(&logrus.TextFormatter{PadLevelText: false, DisableLevelTruncation: false})
	return &ServiceHandler{logger: contextLogger, manager: man}
}

// BuyService
// @Summary user buy service
// @Description user buy service
// @Accept json
// @Produce json
// @Param data body request_models.BuyServiceMessage true "body for buy service"
// @Success 200 {object} models.ShortResponseMessage "OK"
// @Failure 400 {object} models.ShortResponseMessage "invalid body params"
// @Failure 401 {object} models.ShortResponseMessage "account is not exist"
// @Failure 422 {object} models.ShortResponseMessage "not enough money | sum must be > 0"
// @Failure 500 {object} models.ShortResponseMessage "internal server error"
// @Router /api/v1/services/buy [POST]
func (h *ServiceHandler) BuyService(w http.ResponseWriter, r *http.Request) {
	var buyParams request_models.BuyServiceMessage
	var statusCode int
	var handleMessage string

	body, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	unmarshalError := json.Unmarshal(body, &buyParams)
	if unmarshalError != nil {
		http.Error(w, "unmarshal error", http.StatusBadRequest)
		return
	}

	buyError := h.manager.BuyService(buyParams.UserID, buyParams.OrderID, buyParams.ServiceID,
		buyParams.Sum, buyParams.Comment)

	switch buyError {
	case nil:
		statusCode = http.StatusOK
		handleMessage = "OK"
	case ac.AccountNotExistErr:
		statusCode = http.StatusUnauthorized
		handleMessage = fmt.Sprintf("%v", ac.AccountNotExistErr)
	case ac.NotEnoughMoneyErr:
		statusCode = http.StatusUnprocessableEntity
		handleMessage = fmt.Sprintf("%v", ac.NotEnoughMoneyErr)
	case ac.NegSumError:
		statusCode = http.StatusUnprocessableEntity
		handleMessage = fmt.Sprintf("%v", ac.NegSumError)
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = fmt.Sprintf("internal server error")
	}
	models.SendShortResponse(w, statusCode, handleMessage)
	h.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s, err = %v",
		r.Method, r.URL.Path, statusCode, handleMessage, buyError)
}

// AcceptService
// @Summary service accepted
// @Description service bought by user is accepted
// @Accept json
// @Produce json
// @Param data body request_models.AcceptServiceMessage true "body for accept service"
// @Success 200 {object} models.ShortResponseMessage "OK"
// @Failure 400 {object} models.ShortResponseMessage "invalid body params"
// @Failure 401 {object} models.ShortResponseMessage "account is not exist"
// @Failure 403 {object} models.ShortResponseMessage "state isn't right to change order state"
// @Failure 404 {object} models.ShortResponseMessage "order not found"
// @Failure 500 {object} models.ShortResponseMessage "internal server error"
// @Router /api/v1/services/accept [POST]
func (h *ServiceHandler) AcceptService(w http.ResponseWriter, r *http.Request) {
	var acceptParams request_models.AcceptServiceMessage
	var statusCode int
	var handleMessage string

	body, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	unmarshalError := json.Unmarshal(body, &acceptParams)
	if unmarshalError != nil {
		http.Error(w, "invalid body params", http.StatusBadRequest)
		return
	}

	acceptBuyErr := h.manager.AcceptBuy(acceptParams.UserID, acceptParams.OrderID, acceptParams.ServiceID)

	switch acceptBuyErr {
	case nil:
		statusCode = http.StatusOK
		handleMessage = "OK"
	case ac.AccountNotExistErr:
		statusCode = http.StatusUnauthorized
		handleMessage = fmt.Sprintf("%v", ac.AccountNotExistErr)
	case oc.OrderNotFound:
		statusCode = http.StatusNotFound
		handleMessage = fmt.Sprintf("%v", oc.OrderNotFound)
	case oc.WrongStateError:
		statusCode = http.StatusForbidden
		handleMessage = fmt.Sprintf("%v", oc.WrongStateError)
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = fmt.Sprintf("internal server error")
	}
	models.SendShortResponse(w, statusCode, handleMessage)
	h.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s, err = %v",
		r.Method, r.URL.Path, statusCode, handleMessage, acceptBuyErr)
}

// RefuseService
// @Summary service refused
// @Description service bought by user is refused and money returned to user
// @Accept json
// @Produce json
// @Param data body request_models.RefuseServiceMessage true "body for refuse service"
// @Success 200 {object} models.ShortResponseMessage "OK"
// @Failure 400 {object} models.ShortResponseMessage "invalid body params"
// @Failure 401 {object} models.ShortResponseMessage "account is not exist"
// @Failure 403 {object} models.ShortResponseMessage "state isn't right to change order state"
// @Failure 404 {object} models.ShortResponseMessage "order not found"
// @Failure 500 {object} models.ShortResponseMessage "internal server error"
// @Router /api/v1/services/refuse [POST]
func (h *ServiceHandler) RefuseService(w http.ResponseWriter, r *http.Request) {
	var refuseParams request_models.RefuseServiceMessage
	var statusCode int
	var handleMessage string

	body, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	unmarshalError := json.Unmarshal(body, &refuseParams)
	if unmarshalError != nil {
		http.Error(w, "invalid body params", http.StatusBadRequest)
		return
	}

	refuseBuyErr := h.manager.RefuseBuy(refuseParams.UserID, refuseParams.OrderID, refuseParams.ServiceID,
		refuseParams.Comment)

	switch refuseBuyErr {
	case nil:
		statusCode = http.StatusOK
		handleMessage = "OK"
	case ac.AccountNotExistErr:
		statusCode = http.StatusUnauthorized
		handleMessage = fmt.Sprintf("%v", ac.AccountNotExistErr)
	case oc.OrderNotFound:
		statusCode = http.StatusNotFound
		handleMessage = fmt.Sprintf("%v", oc.OrderNotFound)
	case oc.WrongStateError:
		statusCode = http.StatusForbidden
		handleMessage = fmt.Sprintf("%v", oc.WrongStateError)
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = fmt.Sprintf("internal server error")
	}
	models.SendShortResponse(w, statusCode, handleMessage)
	h.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s, err = %v",
		r.Method, r.URL.Path, statusCode, handleMessage, refuseBuyErr)
}
