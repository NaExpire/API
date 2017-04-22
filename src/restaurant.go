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

type GetItemsHandler struct {
	DB *sql.DB
}

type UpdateItemsHandler struct {
	DB *sql.DB
}

type getRestaurantSchema struct {
	Name          string       `json:"name"`
	BusinessPhone string       `json:"phone-number"`
	PickupTime    string       `json:"pickup-time"`
	Description   string       `json:"description"`
	Address       string       `json:"address"`
	City          string       `json:"city"`
	State         string       `json:"state"`
	Meals         []mealSchema `json:"meals"`
	Deals         []dealSchema `json:"deals"`
	Items         string       `json:"items"`
}

type updateRestaurantSchema struct {
	Name          string       `json:"name"`
	BusinessPhone string       `json:"phone-number"`
	PickupTime    string       `json:"pickup-time"`
	Description   string       `json:"description"`
	Address       string       `json:"address"`
	City          string       `json:"city"`
	State         string       `json:"state"`
	Meals         []mealSchema `json:"meals"`
	Deals         []dealSchema `json:"deals"`
	Items         string       `json:"items"`
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
	x := &getRestaurantSchema{}

	rows, err := handler.DB.Query("SELECT `name`, `phone-number`, `description`, `address`, `city`, `state`, `items`, `pickup-time` FROM restaurants WHERE id=?", vars["restaurantID"])

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

	err = rows.Scan(&x.Name, &x.BusinessPhone, &x.Description, &x.Address, &x.City, &x.State, &x.Items, &x.PickupTime)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	meals := make([]mealSchema, 0)
	deals := make([]dealSchema, 0)

	mealRows, mealErr := handler.DB.Query("SELECT `name`, `description`, `restaurantid`, `price`, `type` FROM `menuitems` WHERE `menuitems`.`restaurantid` = ?", vars["restaurantID"])
	defer mealRows.Close()

	dealRows, dealErr := handler.DB.Query("SELECT `meal-id`, `deal-price`, `quantity` FROM `deals` WHERE `deals`.`restaurant-id` = ?", vars["restaurantID"])
	defer dealRows.Close()

	if mealErr != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, mealErr.Error())
		return
	}

	if dealErr != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, dealErr.Error())
		return
	}

	for mealRows.Next() {
		var meal mealSchema
		mealRows.Scan(&meal.Name, &meal.Description, &meal.RestaurantID, &meal.Price, &meal.Type)
		meals = append(meals, meal)
	}

	for dealRows.Next() {
		var deal dealSchema
		dealRows.Scan(&deal.MealID, &deal.DealPrice, &deal.Quantity)
		deals = append(deals, deal)
	}

	x.Meals = meals
	x.Deals = deals

	util.EncodeJSON(writer, x)
}

func (handler UpdateRestaurantHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	x := &updateRestaurantSchema{}
	err := util.DecodeJSON(request.Body, x)

	if err != nil {
		io.WriteString(writer, err.Error()+"\n")
		return
	}

	_, err = handler.DB.Exec("UPDATE restaurants SET name = ? , `phone-number` = ? , description = ? , address = ? , city = ? , state = ?, `pickup-time` = ?, `items` = ? WHERE id = ?", x.Name, x.BusinessPhone, x.Description, x.Address, x.City, x.State, x.PickupTime, x.Items, vars["restaurantID"])

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		io.WriteString(writer, err.Error()+"\n")
		return
	}

	io.WriteString(writer, "{\"ok\": true}")
}

func canSessionOwnerModifyRestaurant(db *sql.DB, sessionID, restaurantID string) (bool, error) {
	rows, err := db.Query("SELECT u.`id` FROM `users` AS u INNER JOIN `sessions` AS s ON s.`user-id` = u.`id` INNER JOIN `restaurants` AS r ON r.`ownerid` = u.`id` WHERE s.`session-content` = ? AND r.`id` = ?", sessionID, restaurantID)
	defer rows.Close()
	if err != nil {
		return false, err
	}
	return rows.Next(), nil
}
