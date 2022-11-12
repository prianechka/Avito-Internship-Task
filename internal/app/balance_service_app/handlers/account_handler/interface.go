package account_handler

import "net/http"

//go:generate mockgen -source=interface.go -destination=mocks/handler_mock.go -package=mocks AccountHandlerInterface
type AccountHandlerInterface interface {
	GetBalance(w http.ResponseWriter, r *http.Request)
	RefillBalance(w http.ResponseWriter, r *http.Request)
	Transfer(w http.ResponseWriter, r *http.Request)
}
