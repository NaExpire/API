package main

import (
	"database/sql"
	"io"
	"net/http"

	"github.com/NAExpire/API/src/util"
	"github.com/gorilla/mux"
)

type GetMealHandler struct {
	DB *sql.DB
}

type UpdateMealHandler struct {
	DB *sql.DB
}

type AddMealHandler struct {
	DB *sql.DB
}

type DeleteMealHandler struct {
	DB *sql.DB
}

type mealSchema struct {
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	RestaurantID int     `json:"restaurantId"`
	Price        float64 `json:"price"`
}

func (handler GetMealHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	x := &mealSchema{}

	rows, err := handler.DB.Query("SELECT `name`, `description`, `restaurantId`, `price` FROM menuitems WHERE id=?", vars["mealID"])

	defer rows.Close()

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	if !rows.Next() {
		writer.WriteHeader(http.StatusNotFound)
		util.WriteErrorJSON(writer, "Meal with ID "+vars["mealID"]+" could not be found")
		return
	}

	err = rows.Scan(&x.Name, &x.Description, &x.RestaurantID, &x.Price)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	util.EncodeJSON(writer, x)
}

func (handler UpdateMealHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	x := &mealSchema{}
	err := util.DecodeJSON(request.Body, x)

	if err != nil {
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	_, err = handler.DB.Exec("UPDATE menuitems SET `name` = ? , `description` = ? , `restaurantId` = ?, `price` = ? WHERE id = ?", x.Name, x.Description, x.RestaurantID, x.Price, vars["mealID"])

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		io.WriteString(writer, err.Error()+"\n")
		return
	}

	io.WriteString(writer, "{\"ok\": true}")
}

func (handler DeleteMealHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	_, err := handler.DB.Exec("DELETE FROM menuitems WHERE id = ?", vars["mealID"])

	if err != nil {
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	io.WriteString(writer, "{\"ok\": true}")
}

func (handler AddMealHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	x := &mealSchema{}
	err := util.DecodeJSON(request.Body, x)

	if err != nil {
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	_, err = handler.DB.Exec("INSERT INTO menuitems (`name` , `description`, `restaurantId`, `price`) VALUES (? , ? , ?, ?) ", x.Name, x.Description, x.RestaurantID, x.Price)

	if err != nil {
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	io.WriteString(writer, "{\"ok\": true}")
}
