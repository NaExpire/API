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

type restaurantSchema struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Address     string `json:"address"`
	City        string `json:"city"`
	State       string `json:"state"`
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
	}

	if rows.Next() {
		err := rows.Scan(&x.Name, &x.Description, &x.Address, &x.City, &x.State)
		if err != nil {
			io.WriteString(writer, err.Error()+"\n")
		}
		util.EncodeJSON(writer, x)
	} else {
		writer.WriteHeader(http.StatusInternalServerError)
		io.WriteString(writer, "{\"ok\": false}")
	}
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
