package main

import (
	"database/sql"
	"io"
	"net/http"

	"github.com/NAExpire/API/src/util"
	"github.com/gorilla/mux"
)

type IssueTransactionHandler struct {
	DB *sql.DB
}

type CancelTransactionHandler struct {
	DB *sql.DB
}

type RejectTransactionHandler struct {
	DB *sql.DB
}

type AcceptTransactionHandler struct {
	DB *sql.DB
}

type FulfillTransactionHandler struct {
	DB *sql.DB
}

func (handler IssueTransactionHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	io.WriteString(writer, "connected to issue transaction endpoint")
}

func (handler CancelTransactionHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	rows, err := handler.DB.Query("SELECT `status` FROM `transactions` INNER JOIN `sessions` ON sessions.`user-id` = transactions.`user-id` AND sessions.`session-content` = ? AND transactions.`id` = ?", request.Header.Get("session"), vars["transactionID"])

	if err != nil {
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	defer rows.Close()

	if !rows.Next() {
		writer.WriteHeader(http.StatusForbidden)
		return
	}

	var transactionStatus string
	err = rows.Scan(&transactionStatus)

	if err != nil {
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	if transactionStatus != "open" && transactionStatus != "accepted" {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	_, err = handler.DB.Exec("UPDATE transactions SET status = ? WHERE id = ?", "cancelled", vars["transactionID"])

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	io.WriteString(writer, "{\"ok\": true}")
}

func (handler RejectTransactionHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	rows, err := handler.DB.Query("SELECT `status` FROM `transactions` INNER JOIN `sessions` ON sessions.`user-id` = transactions.`restaurant-id` AND sessions.`session-content` = ? AND transactions.`id` = ?", request.Header.Get("session"), vars["transactionID"])

	if err != nil {
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	defer rows.Close()

	if !rows.Next() {
		writer.WriteHeader(http.StatusForbidden)
		return
	}

	var transactionStatus string
	err = rows.Scan(&transactionStatus)

	if err != nil {
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	if transactionStatus != "open" {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	_, err = handler.DB.Exec("UPDATE transactions SET status = ? WHERE id = ?", "rejected", vars["transactionID"])

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		io.WriteString(writer, err.Error()+"\n")
		return
	}

	io.WriteString(writer, "{\"ok\": true}")
}

func (handler AcceptTransactionHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	rows, err := handler.DB.Query("SELECT `status` FROM `transactions` INNER JOIN `sessions` ON sessions.`user-id` = transactions.`restaurant-id` AND sessions.`session-content` = ? AND transactions.`id` = ?", request.Header.Get("session"), vars["transactionID"])

	if err != nil {
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	defer rows.Close()

	if !rows.Next() {
		writer.WriteHeader(http.StatusForbidden)
		return
	}

	var transactionStatus string
	err = rows.Scan(&transactionStatus)

	if err != nil {
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	if transactionStatus != "open" {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	_, err = handler.DB.Exec("UPDATE transactions SET status = ? WHERE id = ?", "accepted", vars["transactionID"])

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		io.WriteString(writer, err.Error()+"\n")
		return
	}

	io.WriteString(writer, "{\"ok\": true}")
}

func (handler FulfillTransactionHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	rows, err := handler.DB.Query("SELECT `status` FROM `transactions` INNER JOIN `sessions` ON sessions.`user-id` = transactions.`user-id` AND sessions.`session-content` = ? AND transactions.`id` = ?", request.Header.Get("session"), vars["transactionID"])

	if err != nil {
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	defer rows.Close()

	if !rows.Next() {
		writer.WriteHeader(http.StatusForbidden)
		return
	}

	var transactionStatus string
	err = rows.Scan(&transactionStatus)

	if err != nil {
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	if transactionStatus != "accepted" {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	_, err = handler.DB.Exec("UPDATE transactions SET status = ? WHERE id = ?", "fulfilled", vars["transactionID"])

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		io.WriteString(writer, err.Error()+"\n")
		return
	}

	io.WriteString(writer, "{\"ok\": true}")
}
