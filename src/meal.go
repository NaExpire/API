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
	MealID       int     `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	RestaurantID int     `json:"restaurantID"`
	Price        float64 `json:"price"`
	Type         string  `json:"type"`
}

func (handler GetMealHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	x := &mealSchema{}

	rows, err := handler.DB.Query("SELECT `id`, `name`, `description`, `restaurantid`, `price`, `type` FROM menuitems WHERE id=?", vars["mealID"])

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

	err = rows.Scan(&x.MealID, &x.Name, &x.Description, &x.RestaurantID, &x.Price, &x.Type)
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

	_, err = handler.DB.Exec("UPDATE menuitems SET `name` = ?, `description` = ?,  `type` = ?, `price` = ? WHERE id = ?", x.Name, x.Description, x.Type, x.Price, vars["mealID"])

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
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

	_, err = handler.DB.Exec("INSERT INTO menuitems (`name`, `description`, `restaurantid`, `price`, `type`) VALUES (?, ?, ?, ?, ?)", x.Name, x.Description, x.RestaurantID, x.Price, x.Type)

	if err != nil {
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	io.WriteString(writer, "{\"ok\": true}")
}
