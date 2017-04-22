package main

import (
	"database/sql"
	"io"
	"net/http"

	"github.com/NAExpire/API/src/util"
	"github.com/gorilla/mux"
)

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

type reviewSchema struct {
	Score      int    `json:"score"`
	ReviewBody string `json:"review-body"`
}

func (handler GetReviewHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	x := &reviewSchema{}

	rows, err := handler.DB.Query("SELECT `score` ,`review-body`, `quantity` FROM reviews WHERE id=?", vars["reviewID"])

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

	err = rows.Scan(&x.Score, &x.ReviewBody)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	util.EncodeJSON(writer, x)
}

func (handler AddReviewHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	x := &reviewSchema{}

	err := util.DecodeJSON(request.Body, x)

	if err != nil {
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	_, err = handler.DB.Exec("INSERT INTO reviews (`score`,`review-body`) VALUES (?, ?) ", x.Score, x.ReviewBody)

	if err != nil {
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
