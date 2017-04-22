package main

import (
	"database/sql"
	"net/http"

	"io"

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

type DeleteMealCartHandler struct {
	DB *sql.DB
}

type DeleteDealCartHandler struct {
	DB *sql.DB
}

type DeleteCartContentsHandler struct {
	DB *sql.DB
}

type cartSchema struct {
	MenuItems []mealSchema `json:"menuitems"`
	Deals     []dealSchema `json:"deals"`
}

type addMealToCartSchema struct {
	MealID   int `json:"mealID"`
	Quantity int `json:"quantity"`
}

type addDealToCartSchema struct {
	DealID int `json:"dealID"`
}

type deleteMealSchema struct {
	MealID int `json: "mealID"`
}

type deleteDealSchema struct {
	DealID int `json: "dealID"`
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
	var x addMealToCartSchema
	util.DecodeJSON(request.Body, x)

	sessionID := request.Header.Get("session")

	var cartID int
	cartIDRows, err := handler.DB.Query("SELECT `cart-id` FROM `users` INNER JOIN `sessions` AS s ON s.`user-id` = u.`id` WHERE s.`session-content` = ?", sessionID)
	defer cartIDRows.Close()

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	if !cartIDRows.Next() {
		writer.WriteHeader(http.StatusUnauthorized)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	cartIDRows.Scan(&cartID)

	rows, err := handler.DB.Query("SELECT `quantity` FROM `carts-menuitems` WHERE `cart-id` = ? AND `menuitems-id` = ?", cartID, x.MealID)
	defer rows.Close()

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	// Item not in the cart already
	if !rows.Next() {
		_, err = handler.DB.Exec("INSERT INTO `carts-menuitems` (`cart-id`, `menuitem-id`, `quantity`) VALUES (?, ?, ?)", cartID, x.MealID, x.Quantity)
	} else {
		_, err = handler.DB.Exec("UPDATE `carts-menuitems` SET `quantity` = ? WHERE `menuitem-id` = ? AND `cart-id` = ?", x.MealID, cartID)
	}

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	io.WriteString(writer, "{\"ok\": true}")
}

func (handler AddDealCartHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var x addDealToCartSchema
	util.DecodeJSON(request.Body, x)

	sessionID := request.Header.Get("session")

	var cartID int
	cartIDRows, err := handler.DB.Query("SELECT `cart-id` FROM `users` INNER JOIN `sessions` AS s ON s.`user-id` = u.`id` WHERE s.`session-content` = ?", sessionID)
	defer cartIDRows.Close()

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	if !cartIDRows.Next() {
		writer.WriteHeader(http.StatusUnauthorized)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	cartIDRows.Scan(&cartID)

	rows, err := handler.DB.Query("SELECT * FROM `carts-deals` WHERE `cart-id` = ? AND `deal-id` = ?", cartID, x.DealID)
	defer rows.Close()

	if !rows.Next() {
		_, err = handler.DB.Exec("INSERT INTO `carts-deals` (`cart-id`, `deal-id`) VALUES (?, ?)", cartID, x.DealID)
	} else {
		io.WriteString(writer, "Duplicate Deal already in Cart")
	}

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	io.WriteString(writer, "{\"ok\": true}")
}

func (handler DeleteMealCartHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	x := &deleteMealSchema{}
	err := util.DecodeJSON(request.Body, x)

	sessionID := request.Header.Get("session")

	var cartID int
	cartIDRows, err := handler.DB.Query("SELECT `cart-id` FROM `users` INNER JOIN `sessions` AS s ON s.`user-id` = u.`id` WHERE s.`session-content` = ?", sessionID)
	defer cartIDRows.Close()

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	if !cartIDRows.Next() {
		writer.WriteHeader(http.StatusUnauthorized)
		util.WriteErrorJSON(writer, "You do not have permission to make changes to this cart item.")
		return
	}

	cartIDRows.Scan(&cartID)

	_, err = handler.DB.Exec("DELETE FROM `carts-menuitems` WHERE `menuitem-id` = ? AND `cart-id` = ?", x.MealID, cartID)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	io.WriteString(writer, "{\"ok\": true}")
}

func (handler DeleteDealCartHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	x := &deleteDealSchema{}
	err := util.DecodeJSON(request.Body, x)

	sessionID := request.Header.Get("session")

	var cartID int
	cartIDRows, err := handler.DB.Query("SELECT `cart-id` FROM `users` INNER JOIN `sessions` AS s ON s.`user-id` = u.`id` WHERE s.`session-content` = ?", sessionID)
	defer cartIDRows.Close()

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	if !cartIDRows.Next() {
		writer.WriteHeader(http.StatusUnauthorized)
		util.WriteErrorJSON(writer, "You do not have permission to make changes to this cart item.")
		return
	}

	cartIDRows.Scan(&cartID)

	_, err = handler.DB.Exec("DELETE FROM `carts-deals` WHERE `deal-id` = ? AND `cart-id` = ?", x.DealID, cartID)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	io.WriteString(writer, "{\"ok\": true}")
}

func (handler DeleteCartContentsHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	sessionID := request.Header.Get("session")

	mealRows, err := handler.DB.Query("DELETE c FROM `carts-menuitems` AS c INNER JOIN `users` AS u ON u.`cart-id` = c.`cart-id` INNER JOIN `sessions` AS s ON s.`user-id` = u.id WHERE s.`session-content` = ?", sessionID)
	defer mealRows.Close()

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	dealRows, err := handler.DB.Query("DELETE c FROM `carts-deals` AS c INNER JOIN `users` AS u ON u.`cart-id` = c.`cart-id` INNER JOIN `sessions` AS s ON s.`user-id` = u.id WHERE s.`session-content` = ?", sessionID)
	defer dealRows.Close()

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	io.WriteString(writer, "{\"ok\": true}")
}
