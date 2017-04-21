package main

import (
	"database/sql"
	"io"
	"net/http"

	"github.com/NAExpire/API/src/util"
	"github.com/gorilla/mux"
)

type GetDealHandler struct {
	DB *sql.DB
}

type UpdateDealHandler struct {
	DB *sql.DB
}

type AddDealHandler struct {
	DB *sql.DB
}

type DeleteDealHandler struct {
	DB *sql.DB
}

type dealSchema struct {
	MealID    int     `json:"meal-id"`
	DealPrice float64 `json:"deal-price"`
	Quantity  int     `json:"quantity"`
}

func (handler GetDealHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	x := &dealSchema{}

	rows, err := handler.DB.Query("SELECT `meal-id`, `deal-price`, `quantity` FROM deals WHERE id=?", vars["dealID"])

	defer rows.Close()

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	if !rows.Next() {
		writer.WriteHeader(http.StatusNotFound)
		util.WriteErrorJSON(writer, "Deal with ID "+vars["dealID"]+" could not be found")
		return
	}

	err = rows.Scan(&x.MealID, &x.DealPrice, &x.Quantity)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	util.EncodeJSON(writer, x)
}

func (handler UpdateDealHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	x := &dealSchema{}
	err := util.DecodeJSON(request.Body, x)

	if err != nil {
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	_, err = handler.DB.Exec("UPDATE deals SET `meal-id` = ? , `deal-price` = ? , `quantity` = ? WHERE id = ?", x.MealID, x.DealPrice, x.Quantity, vars["dealID"])

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		io.WriteString(writer, err.Error()+"\n")
		return
	}

	io.WriteString(writer, "{\"ok\": true}")
}

func (handler DeleteDealHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	_, err := handler.DB.Exec("DELETE FROM deals WHERE id = ?", vars["dealID"])

	if err != nil {
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	io.WriteString(writer, "{\"ok\": true}")
}

func (handler AddDealHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	x := &dealSchema{}
	err := util.DecodeJSON(request.Body, x)

	if err != nil {
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	_, err = handler.DB.Exec("INSERT INTO deals (`meal-id` , `deal-price`, `quantity`) VALUES (? , ? , ?) ", x.MealID, x.DealPrice, x.Quantity)

	if err != nil {
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	io.WriteString(writer, "{\"ok\": true}")
}
