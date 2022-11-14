package server

import (
	"Avito-Internship-Task/configs"
	_ "Avito-Internship-Task/docs"
	ac "Avito-Internship-Task/internal/app/balance_service_app/account/account_controller"
	"Avito-Internship-Task/internal/app/balance_service_app/account/account_repo"
	"Avito-Internship-Task/internal/app/balance_service_app/handlers/account_handler"
	"Avito-Internship-Task/internal/app/balance_service_app/handlers/report_handler"
	"Avito-Internship-Task/internal/app/balance_service_app/handlers/service_handler"
	"Avito-Internship-Task/internal/app/balance_service_app/manager"
	"Avito-Internship-Task/internal/app/balance_service_app/middleware"
	oc "Avito-Internship-Task/internal/app/balance_service_app/order/order_controller"
	"Avito-Internship-Task/internal/app/balance_service_app/order/order_repo"
	tc "Avito-Internship-Task/internal/app/balance_service_app/transaction/transaction_controller"
	"Avito-Internship-Task/internal/app/balance_service_app/transaction/transaction_repo"
	"Avito-Internship-Task/internal/pkg/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

type Server struct {
	config *configs.ServerConfig
	logger *logrus.Entry
}

func CreateServer(config *configs.ServerConfig, logger *logrus.Entry) *Server {
	return &Server{config: config, logger: logger}
}

// @title AvitoIntershipApp
// @description Task for Avito-Intership.

func (s *Server) Start() error {

	r := mux.NewRouter()
	router := r.PathPrefix("/api/v1/").Subrouter()
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	accountDB := utils.NewMySQLConnction(s.config.ConnParams)
	orderDB := utils.NewMySQLConnction(s.config.ConnParams)
	transactionDB := utils.NewMySQLConnction(s.config.ConnParams)

	accountRepo := account_repo.NewAccountRepo(accountDB)
	accountController := ac.CreateNewAccountController(accountRepo)

	orderRepo := order_repo.NewOrderRepo(orderDB)
	orderController := oc.CreateNewOrderController(orderRepo)

	transactionRepo := transaction_repo.NewTransactionRepo(transactionDB)
	transactionController := tc.CreateNewTransactionController(transactionRepo)

	serverManager := manager.CreateNewManager(accountController, orderController, transactionController)

	accountHandler := account_handler.CreateAccountHandler(serverManager)
	serviceHandler := service_handler.CreateServiceHandler(serverManager)
	reportHandler := report_handler.CreateReportHandler(serverManager)

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	router.HandleFunc("/accounts", accountHandler.GetBalance).Methods("GET")
	router.HandleFunc("/accounts/refill", accountHandler.RefillBalance).Methods("POST")
	router.HandleFunc("/transfer", accountHandler.Transfer).Methods("POST")

	router.HandleFunc("/services/buy", serviceHandler.BuyService).Methods("POST")
	router.HandleFunc("/services/accept", serviceHandler.AcceptService).Methods("POST")
	router.HandleFunc("/services/refuse", serviceHandler.RefuseService).Methods("POST")

	router.HandleFunc("/reports/user", reportHandler.GetUserReport).Methods("GET")
	router.HandleFunc("/reports/finance", reportHandler.GetFinanceReport).Methods("GET")

	withLogsRouter := middleware.Log(s.logger, router)
	upgradedRouter := middleware.Panic(withLogsRouter)

	s.logger.Infof("Server started to work!")
	return http.ListenAndServe(s.config.PortToStart, upgradedRouter)
}
