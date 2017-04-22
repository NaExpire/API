package main

import (
	"database/sql"
	"net/http"

	"github.com/NAExpire/API/src/util"
)

type GetCartHandler struct {
	DB *sql.DB
}

type AddMealCartHandler struct {
	DB *sql.DB
}

type AddDealCartHandler struct {
	DB *sql.DB
}

type UpdateQuantityMealHandler struct {
	DB *sql.DB
}

type DeleteMealCartHandler struct {
	DB *sql.DB
}

type DeleteCartContentsHandler struct {
	DB *sql.DB
}

type DeleteDealCartHandler struct {
	DB *sql.DB
}

type cartSchema struct {
	MenuItems []mealSchema `json:"menuitems"`
	Deals     []dealSchema `json:"deals"`
}

func (handler GetCartHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	menuItems := make([]mealSchema, 0)
	deals := make([]dealSchema, 0)

	mealRows, err := handler.DB.Query("SELECT m.name, m.description, m.restaurantID, m.price, m.type FROM `carts-menuitems` AS c INNER JOIN `users` AS u ON c.id = u.`cart-id` INNER JOIN sessions AS s ON s.`user-id` = u.id INNER JOIN `menuitems` AS m ON c.`menuitem-id` = m.id WHERE s.`session-content` = ?", request.Header.Get("session"))

	defer mealRows.Close()

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	for mealRows.Next() {
		var meal mealSchema
		mealRows.Scan(&meal.Name, &meal.Description, &meal.RestaurantID, &meal.Price, &meal.Type)
		menuItems = append(menuItems, meal)
	}

	dealRows, err := handler.DB.Query("SELECT d.`meal-id`, d.`deal-price`, d.quantity FROM `carts-deals` AS c INNER JOIN `users` AS u ON c.id = u.`cart-id` INNER JOIN sessions AS s ON s.`user-id` = u.id INNER JOIN `deals` AS d ON c.`deal-id` = d.id WHERE s.`session-content` = ?", request.Header.Get("session"))

	defer dealRows.Close()

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	for dealRows.Next() {
		var deal dealSchema
		dealRows.Scan(&deal.MealID, &deal.DealPrice, &deal.Quantity)
		deals = append(deals, deal)
	}

	var cart cartSchema
	cart.MenuItems = menuItems
	cart.Deals = deals

	util.EncodeJSON(writer, cart)

}

func (handler AddMealCartHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

}

func (handler AddDealCartHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

}

func (handler UpdateQuantityMealHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

}

func (handler DeleteMealCartHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

}

func (handler DeleteCartContentsHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

}

func (handler DeleteDealCartHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

}
