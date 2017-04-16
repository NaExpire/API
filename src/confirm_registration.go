package main

import (
	"database/sql"
	"io"
	"net/http"
	"strconv"

	"github.com/NAExpire/API/src/util"
)

type ConfirmRegistrationHandler struct {
	DB *sql.DB
}

type confirmRegistrationCredentials struct {
	ConfirmationCode string `json:"confirmationCode"`
	EmailAddress     string `json:"emailAddress"`
}

func (handler ConfirmRegistrationHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	x := &confirmRegistrationCredentials{}
	err := util.DecodeJSON(request.Body, x)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		util.WriteErrorJSON(writer, "Malformed request syntax")
		return
	}

	rows, err := handler.DB.Query("SELECT `confirmation-code` FROM `users` WHERE `email` = ?", x.EmailAddress)
	defer rows.Close()
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	} else if !rows.Next() {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		util.WriteErrorJSON(writer, "Email address is not registered to any user")
		return
	}

	var confirmationCode int
	err = rows.Scan(&confirmationCode)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	} else if strconv.Itoa(confirmationCode) != x.ConfirmationCode {
		writer.WriteHeader(http.StatusUnauthorized)
		util.WriteErrorJSON(writer, "Invalid email or confirmation code")
		return
	}

	_, err = handler.DB.Exec("UPDATE `users` SET `confirmed` = ? WHERE `email` = ?", 1, x.EmailAddress)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	io.WriteString(writer, "{\"ok\": true")
}
