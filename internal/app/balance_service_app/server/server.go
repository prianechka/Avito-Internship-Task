package server

import (
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
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	logger logrus.Logger
}

func CreateDB(DBName string) *sql.DB {
	dsn := "root:love@tcp(localhost:3306)/" + DBName + "?&charset=utf8&interpolateParams=true"
	db, err := sql.Open("mysql", dsn)
	if err == nil {
		db.SetMaxOpenConns(10)
		err = db.Ping()
		if err != nil {
			panic(err)
		}
	} else {
		panic(err)
	}
	return db
}

func (s *Server) Start() error {
	r := mux.NewRouter()
	router := r.PathPrefix("/api/v1/").Subrouter()

	accountDB := CreateDB("accounts")
	orderDB := CreateDB("orders")
	transactionDB := CreateDB("transactions")

	accountRepo := account_repo.NewAccountRepo(accountDB)
	accountController := ac.CreateNewAccountController(accountRepo)

	orderRepo := order_repo.NewOrderRepo(orderDB)
	orderController := oc.CreateNewOrderController(orderRepo)

	transactionRepo := transaction_repo.NewTransactionRepo(transactionDB)
	transactionController := tc.CreateNewTransactionController(transactionRepo)

	serverManager := manager.CreateNewManager(accountController, orderController, transactionController)

	accountHandler := account_handler.CreateAccountHandler(serverManager)
	serviceHandler := service_handler.CreateServiceHandler(serverManager)
	reportHandler := report_handler.ReportHandler{Manager: serverManager}

	router.HandleFunc("/accounts/{id}", accountHandler.GetBalance).Methods("GET")
	router.HandleFunc("/accounts/refill", accountHandler.RefillBalance).Methods("POST")
	router.HandleFunc("/transfer", accountHandler.Transfer).Methods("POST")

	router.HandleFunc("/services/buy", serviceHandler.BuyService).Methods("POST")
	router.HandleFunc("/services/accept", serviceHandler.AcceptService).Methods("POST")
	router.HandleFunc("/services/refuse", serviceHandler.RefuseService).Methods("POST")

	router.HandleFunc("/reports/user/{id}", reportHandler.GetUserReport).Methods("GET")
	router.HandleFunc("/reports/{year}/{month}", reportHandler.GetFinanceReport).Methods("GET")

	upgradedRouter := middleware.Panic(router)

	addr := ":8080"
	return http.ListenAndServe(addr, upgradedRouter)
}
