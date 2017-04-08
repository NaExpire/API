package business

import (
	"database/sql"
	"io"
	"net/http"
)

type ConfirmRegistrationHandler struct {
	DB *sql.DB
}

type confirmBusinessRegistrationCredentials struct {
	confirmationCode int    `json:"confirmationCode"`
	emailAddress     string `json:"emailAddress"`
}

func (handler ConfirmRegistrationHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	x := &confirmBusinessRegistrationCredentials{}
	rows, err := handler.DB.Query("SELECT `confirmation-code` FROM `users` WHERE `email` = ?", x.emailAddress)
	defer rows.Close()
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		io.WriteString(writer, err.Error()+"\n")
		return
	}
	var confirmationCode int
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		io.WriteString(writer, err.Error()+"\n")
		return
	}
	rows.Scan(&confirmationCode)
	if confirmationCode == x.confirmationCode {
		io.WriteString(writer, "{\"ok\": true}")
	} else {
		io.WriteString(writer, "{\"ok\": false}")
	}
}
