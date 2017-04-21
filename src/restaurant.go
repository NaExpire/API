package main

import (
	"database/sql"
	"io"
	"net/http"

	"github.com/NAExpire/API/src/util"
	"github.com/gorilla/mux"
)

type GetRestaurantHandler struct {
	DB *sql.DB
}

type UpdateRestaurantHandler struct {
	DB *sql.DB
}

type restaurantSchema struct {
	Name          string `json:"name"`
	BusinessPhone string `json:"phone-number"`
	//PickupTime    time.Time `json:"pickup-time"`
	Description string `json:"description"`
	Address     string `json:"address"`
	City        string `json:"city"`
	State       string `json:"state"`
}

func (handler GetRestaurantHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	x := &restaurantSchema{}

	rows, err := handler.DB.Query("SELECT `name`, `phone-number`, `description`, `address`, `city`, `state` FROM restaurants WHERE id=?", vars["restaurantID"])

	defer rows.Close()

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	if !rows.Next() {
		writer.WriteHeader(http.StatusNotFound)
		util.WriteErrorJSON(writer, "Restaurant with ID "+vars["restaurantID"]+" could not be found")
		return
	}

	err = rows.Scan(&x.Name, &x.BusinessPhone, &x.Description, &x.Address, &x.City, &x.State)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	util.EncodeJSON(writer, x)
}

func (handler UpdateRestaurantHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	x := &restaurantSchema{}
	err := util.DecodeJSON(request.Body, x)

	if err != nil {
		io.WriteString(writer, err.Error()+"\n")
		return
	}

	_, err = handler.DB.Exec("UPDATE restaurants SET name = ? , `phone-number` = ? , description = ? , address = ? , city = ? , state = ?  WHERE id = ?", x.Name, x.BusinessPhone, x.Description, x.Address, x.City, x.State, vars["restaurantID"])

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		io.WriteString(writer, err.Error()+"\n")
		return
	}

	io.WriteString(writer, "{\"ok\": true}")
}
