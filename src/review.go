package main

import (
	"database/sql"
	"io"
	"net/http"

	"github.com/NAExpire/API/src/util"
	"github.com/gorilla/mux"
)

type GetAllReviewsHandler struct {
	DB *sql.DB
}

type GetReviewHandler struct {
	DB *sql.DB
}

type AddReviewHandler struct {
	DB *sql.DB
}

type UpdateReviewHandler struct {
	DB *sql.DB
}

type DeleteReviewHandler struct {
	DB *sql.DB
}

type createReviewSchema struct {
	RestaurantID  int    `json:"restaurantID"`
	Score         int    `json:"score"`
	ReviewBody    string `json:"reviewBody"`
	TransactionID int    `json:"transactionID"`
}

type reviewSchema struct {
	RestaurantID int    `json:"restaurantID"`
	Score        int    `json:"score"`
	ReviewBody   string `json:"reviewBody"`
}

func (handler GetAllReviewsHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	x := &reviewSchema{}

	rows, err := handler.DB.Query("SELECT `score` ,`review-body`, `restaurant-id` FROM `reviews` WHERE `restaurant-id` = ?", vars["restaurantID"])

	defer rows.Close()

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	if !rows.Next() {
		writer.WriteHeader(http.StatusNotFound)
		util.WriteErrorJSON(writer, "Reviews with restaurantID"+vars["restaurantID"]+"could not be found")
		return
	}

	err = rows.Scan(&x.Score, &x.ReviewBody, &x.RestaurantID)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	util.EncodeJSON(writer, x)
}

func (handler GetReviewHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	x := &reviewSchema{}

	rows, err := handler.DB.Query("SELECT `score`, `review-body`, `restaurant-id` FROM `reviews` WHERE id=?", vars["reviewID"])

	defer rows.Close()

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	if !rows.Next() {
		writer.WriteHeader(http.StatusNotFound)
		util.WriteErrorJSON(writer, "Review with ID "+vars["reviewID"]+" could not be found")
		return
	}

	err = rows.Scan(&x.Score, &x.ReviewBody, &x.RestaurantID)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	util.EncodeJSON(writer, x)
}

func (handler AddReviewHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	sessionID := request.Header.Get("session")
	rows, err := handler.DB.Query("SELECT `user-id` FROM `sessions` WHERE `session-content` = ?", sessionID)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	if !rows.Next() {
		writer.WriteHeader(http.StatusUnauthorized)
		util.WriteErrorJSON(writer, "Unauthorized")
		return
	}

	var userID int
	err = rows.Scan(&userID)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	x := &createReviewSchema{}

	err = util.DecodeJSON(request.Body, x)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	_, err = handler.DB.Exec("INSERT INTO reviews (`score`,`review-body`, `restaurant-id`, `transaction-id`, `user-id`) VALUES (?, ?, ?, ?, ?) ", x.Score, x.ReviewBody, x.RestaurantID, x.TransactionID, userID)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	writer.WriteHeader(http.StatusCreated)
	io.WriteString(writer, "{\"ok\": true}")
}

func (handler UpdateReviewHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	x := &reviewSchema{}
	err := util.DecodeJSON(request.Body, x)

	if err != nil {
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	_, err = handler.DB.Exec("UPDATE reviews SET `score` = ? ,`review-body` = ? WHERE id = ?", x.Score, x.ReviewBody, vars["reviewID"])

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	io.WriteString(writer, "{\"ok\": true}")
}

func (handler DeleteReviewHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	_, err := handler.DB.Exec("DELETE FROM reviews WHERE id = ?", vars["reviewID"])

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	io.WriteString(writer, "{\"ok\": true}")
}
