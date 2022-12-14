package models

import (
	"Avito-Internship-Task/internal/app/balance_service_app/transaction"
	"encoding/json"
	"net/http"
)

type ShortResponseMessage struct {
	Comment string `json:"comment"`
}

type BalanceResponseMessage struct {
	Balance float64 `json:"balance"`
	Comment string  `json:"comment"`
}

type FinanceReportResponseMessage struct {
	FileURL string `json:"fileURL"`
}

type UserReportResponseMessage struct {
	AllTransactions []transaction.Transaction `json:"transactions"`
}

func SendShortResponse(w http.ResponseWriter, code int, comment string) {
	var msg = ShortResponseMessage{comment}
	result, err := json.Marshal(msg)
	if err == nil {
		w.WriteHeader(code)
		_, err = w.Write(result)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func BalanceResponse(w http.ResponseWriter, balance float64, comment string) {
	var msg = BalanceResponseMessage{balance, comment}
	result, err := json.Marshal(msg)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(result)
	} else {
		SendShortResponse(w, http.StatusInternalServerError, "internal server problems")
	}
}

func FinanceReportResponse(w http.ResponseWriter, fileURL string) {
	var msg = FinanceReportResponseMessage{fileURL}
	result, err := json.Marshal(msg)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(result)
	} else {
		SendShortResponse(w, http.StatusInternalServerError, "internal server problems")
	}
}

func UserReportResponse(w http.ResponseWriter, allTransactions []transaction.Transaction) {
	var msg = UserReportResponseMessage{allTransactions}
	result, err := json.Marshal(msg)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(result)
	} else {
		SendShortResponse(w, http.StatusInternalServerError, "internal server problems")
	}
}
