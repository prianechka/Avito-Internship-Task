package service_handler

import "net/http"

//go:generate mockgen -source=interface.go -destination=mocks/handler_mock.go -package=mocks ServiceHandlerInterface
type ServiceHandlerInterface interface {
	BuyService(w http.ResponseWriter, r *http.Request)
	AcceptService(w http.ResponseWriter, r *http.Request)
	RefuseService(w http.ResponseWriter, r *http.Request)
}
