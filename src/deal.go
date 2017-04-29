package main

import (
	"database/sql"
	"io"
	"net/http"
	"strconv"

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

type getDealSchema struct {
	DealID            int     `json:"id"`
	MealID            int     `json:"mealID"`
	DealPrice         float64 `json:"dealPrice"`
	Quantity          int     `json:"quantity"`
	RestaurantID      int     `json:"restaurantID"`
	ItemName          string  `json:"itemName"`
	RestaurantName    string  `json:"restaurantName"`
	RestaurantAddress string  `json:"restaurantAddress"`
}

type dealSchema struct {
	DealID       int     `json:"id"`
	MealID       int     `json:"mealID"`
	DealPrice    float64 `json:"dealPrice"`
	Quantity     int     `json:"quantity"`
	RestaurantID int     `json:"restaurantID"`
}

func (handler GetDealHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	x := &getDealSchema{}

	rows, err := handler.DB.Query("SELECT d.`id`, d.`meal-id`, d.`deal-price`, d.`quantity`, d.`restaurant-id`, m.`name`, r.`name`, r.`address` FROM deals as d INNER JOIN `restaurants` as r ON r.`id` = d.`restaurant-id` INNER JOIN `menuitems` as m ON m.`id` = d.`meal-id` WHERE d.`id`=?", vars["dealID"])

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

	err = rows.Scan(&x.DealID, &x.MealID, &x.DealPrice, &x.Quantity, &x.RestaurantID, &x.ItemName, &x.RestaurantName, &x.RestaurantAddress)
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

	_, err = handler.DB.Exec("UPDATE deals SET `deal-price` = ?, `quantity` = ? WHERE id = ?", x.DealPrice, x.Quantity, vars["dealID"])

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
		writer.WriteHeader(http.StatusBadRequest)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	result, err := handler.DB.Exec("INSERT INTO deals (`meal-id`, `deal-price`, `quantity`, `restaurant-id`) VALUES (?, ?, ?, ?) ", x.MealID, x.DealPrice, x.Quantity, x.RestaurantID)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	io.WriteString(writer, "{\"ok\": true, \"id\": "+strconv.FormatInt(insertedID, 10)+"}")
}
