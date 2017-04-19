package main

import (
	"database/sql"
	"fmt"
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

type GetItemsHandler struct {
	DB *sql.DB
}

type UpdateItemsHandler struct {
	DB *sql.DB
}

type restaurantSchema struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Address     string `json:"address"`
	City        string `json:"city"`
	State       string `json:"state"`
}

type itemsSchema struct {
	Items           string `json:"items"`
	Price           string `json:"price"`
	PickupTime      string `json:"pickupTime"`
	PickupMax       string `json:"pickupMax"`
	PickupRemaining string `json:"pickupRemaining"`
}

func (handler GetRestaurantHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	x := &restaurantSchema{}
	err := util.DecodeJSON(request.Body, x)
	fmt.Printf("Got %s request to GetRestaurantHandler\n", request.Method)
	if err != nil {
		io.WriteString(writer, err.Error()+"\n")
		return
	}

	rows, err := handler.DB.Query("SELECT name, description, address, city, state FROM restaurants WHERE id=?", vars["id"])

	defer rows.Close()

	if err != nil {
		io.WriteString(writer, err.Error()+"\n")
		return
	} else if !rows.Next() {
		writer.WriteHeader(http.StatusInternalServerError)
		io.WriteString(writer, "{\"ok\": false}")
		return
	}
	err = rows.Scan(&x.Name, &x.Description, &x.Address, &x.City, &x.State)
	if err != nil {
		io.WriteString(writer, err.Error()+"\n")
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

	_, err = handler.DB.Exec("UPDATE restaurants SET name = ? , description = ? , address = ? , city = ? , state = ?  WHERE id = ?", x.Name, x.Description, x.Address, x.City, x.State, vars["id"])
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		io.WriteString(writer, err.Error()+"\n")
		return
	}
	io.WriteString(writer, "{\"ok\": true}")
}

func (handler GetItemsHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	x := &itemsSchema{}
	err := util.DecodeJSON(request.Body, x)
	fmt.Printf("Got %s request to GetItemHandler\n", request.Method)
	if err != nil {
		io.WriteString(writer, err.Error()+"\n")
		return
	}
	rows, err := handler.DB.Query("SELECT items, `pickup-tme`, price, `pickup-max`, `pickup-remaining` FROM restaurants WHERE id=?", vars["id"])

	defer rows.Close()

	if err != nil {
		io.WriteString(writer, err.Error()+"\n")
		return
	} else if !rows.Next() {
		writer.WriteHeader(http.StatusInternalServerError)
		io.WriteString(writer, "{\"ok\": false}")
		return
	}
	err = rows.Scan(&x.Items, &x.PickupTime, &x.Price, &x.PickupMax, &x.PickupRemaining)
	if err != nil {
		io.WriteString(writer, err.Error()+"\n")
		return
	}
	util.EncodeJSON(writer, x)
}

func (handler UpdateItemsHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	x := &itemsSchema{}
	err := util.DecodeJSON(request.Body, x)

	if err != nil {
		io.WriteString(writer, err.Error()+"\n")
		return
	}

	_, err = handler.DB.Exec("UPDATE restaurants SET items = ? , `pickup-tme` = ? , price = ? ,  `pickup-max` = ? , pickup-remaining` = ?  WHERE id = ?", x.Items, x.PickupTime, x.Price, x.PickupMax, x.PickupRemaining, vars["id"])
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		io.WriteString(writer, err.Error()+"\n")
		return
	}
	io.WriteString(writer, "{\"ok\": true}")
}
