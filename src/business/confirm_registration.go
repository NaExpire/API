package business

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

type confirmBusinessRegistrationCredentials struct {
	ConfirmationCode string `json:"confirmationCode"`
	EmailAddress     string `json:"emailAddress"`
}

func (handler ConfirmRegistrationHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	x := &confirmBusinessRegistrationCredentials{}
	util.DecodeJSON(request.Body, x)
	rows, err := handler.DB.Query("SELECT `confirmation-code` FROM `users` WHERE `email` = ?", x.EmailAddress)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		io.WriteString(writer, err.Error()+"\n")
		return
	}
	defer rows.Close()
	var confirmationCode int
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		io.WriteString(writer, err.Error()+"\n")
		return
	}
	if !rows.Next() {
		io.WriteString(writer, "{\"ok\": false}")
		return
	}
	err = rows.Scan(&confirmationCode)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		io.WriteString(writer, err.Error()+"\n")
		return
	}
	if strconv.Itoa(confirmationCode) == x.ConfirmationCode {
		io.WriteString(writer, "{\"ok\": true}")
		_, err := handler.DB.Exec("UPDATE `users` SET `confirmed` = ? WHERE `email` = ?", 1, x.EmailAddress)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			io.WriteString(writer, err.Error()+"\n")
			return
		}
	} else {
		io.WriteString(writer, "{\"ok\": false}")
	}
}
