package main

import (
	"database/sql"
	"io"
	"net/http"
	"time"

	"strconv"

	"github.com/NAExpire/API/src/seshin"
	"github.com/NAExpire/API/src/util"
)

type LogoutHandler struct {
	DB *sql.DB
}

type ConsumerLoginHandler struct {
	DB *sql.DB
}

type BusinessLoginHandler struct {
	DB *sql.DB
}

type loginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (handler BusinessLoginHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	x := &loginCredentials{}
	err := util.DecodeJSON(request.Body, x)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		util.WriteErrorJSON(writer, "Malformed request syntax")
		return
	}

	rows, err := handler.DB.Query("SELECT u.`id`, u.`password`, u.`confirmed`, r.`id` FROM `users` AS u INNER JOIN `restaurants` AS r ON u.`id` = r.`ownerid` WHERE email=?", x.Email)
	defer rows.Close()
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, "Malformed request syntax")
		return
	}
	if !rows.Next() {
		writer.WriteHeader(http.StatusForbidden)
		util.WriteErrorJSON(writer, "Invalid username or password")
		return
	}

	var userID int
	var passwordHash string
	var confirmed int
	var restaurantID int
	err = rows.Scan(&userID, &passwordHash, &confirmed, &restaurantID)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	} else if confirmed == 0 {
		writer.WriteHeader(http.StatusForbidden)
		util.WriteErrorJSON(writer, "Account not confirmed.")
		return
	}

	myUniqueSessionID := seshin.GenerateSessionID()
	err = seshin.CreateSession(handler.DB, myUniqueSessionID, userID)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}
	_, err = handler.DB.Exec("UPDATE `users` SET `last-login`=? WHERE `email`=?", time.Now(), x.Email)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	writer.WriteHeader(http.StatusOK)
	responseBody := "{\"ok\": true,\"sessionID\": \"" + myUniqueSessionID + "\",\"restaurantID\":" + strconv.FormatInt(int64(restaurantID), 10) + "}"
	io.WriteString(writer, responseBody)
}

func (handler ConsumerLoginHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	x := &loginCredentials{}
	err := util.DecodeJSON(request.Body, x)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		util.WriteErrorJSON(writer, "Malformed request syntax")
		return
	}

	rows, err := handler.DB.Query("SELECT `id`, `password`, `confirmed` FROM `users` WHERE email=?", x.Email)
	defer rows.Close()
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, "Malformed request syntax")
		return
	}
	if !rows.Next() {
		writer.WriteHeader(http.StatusUnauthorized)
		util.WriteErrorJSON(writer, "Invalid username or password")
		return
	}

	var userID int
	var passwordHash string
	var confirmed int
	err = rows.Scan(&userID, &passwordHash, &confirmed)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	} else if confirmed == 0 {
		writer.WriteHeader(http.StatusForbidden)
		util.WriteErrorJSON(writer, "Account not confirmed.")
		return
	}

	myUniqueSessionID := seshin.GenerateSessionID()
	err = seshin.CreateSession(handler.DB, myUniqueSessionID, userID)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}
	_, err = handler.DB.Exec("UPDATE `users` SET `last-login`=? WHERE `email`=?", time.Now(), x.Email)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	writer.WriteHeader(http.StatusOK)
	responseBody := "{\"ok\": true,\"sessionID\": \"" + myUniqueSessionID + "\"}"
	io.WriteString(writer, responseBody)
}

func (handler LogoutHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	sessionID := request.Header.Get("session")

	isValidSession, err := seshin.ValidateSession(handler.DB, sessionID)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	} else if !isValidSession {
		writer.WriteHeader(http.StatusUnauthorized)
		util.WriteErrorJSON(writer, "Invalid session")
		return
	}

	err = seshin.InvalidateSession(handler.DB, sessionID)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	writer.WriteHeader(http.StatusOK)
	io.WriteString(writer, "{\"ok\": true}")
}
