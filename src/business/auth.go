package business

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/NAExpire/API/src/seshin"
	"github.com/NAExpire/API/src/util"
)

type LogoutHandler struct {
	DB *sql.DB
}

type LoginHandler struct {
	DB *sql.DB
}

type businessLoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (handler LoginHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	x := &businessLoginCredentials{}
	err := util.decodeJSON(request.Body, x)
	fmt.Printf("Got %s request to LoginHandler\n", request.Method)
	if request.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
	} else if err != nil {
		io.WriteString(writer, err.Error()+"\n")
		return
	} else {
		io.WriteString(writer, x.Email+"\n")
		io.WriteString(writer, x.Password+"\n")
	}

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(x.Password), 14)
	rows, err := handler.DB.Query("SELECT email FROM users WHERE email=? AND password=? AND confirmed=?", x.Email, string(passwordHash), 1)

	defer rows.Close()

	if err != nil {
		io.WriteString(writer, err.Error()+"\n")
		return
	}

	if rows.Next() {
		myUniqueSessionID := seshin.GenerateSessionID()
		seshin.CreateSession(myUniqueSessionID)
		_, err = handler.DB.Exec("INSERT INTO users (`last-login`) VALUES (?)", time.Now())

		io.WriteString(writer, "{  }")
	} else {
		writer.WriteHeader(http.StatusConflict)
		io.WriteString(writer, "Email or password is incorrect\n")
		return
	}

}

func (handler LogoutHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	fmt.Printf("Got %s request to LogoutHandler\n", request.Method)
	if request.Method != "GET" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
	sessionID := request.Header.Get("session")

	isValidSession, err := seshin.ValidateSession(handler.DB, sessionID)
	if isValidSession {
		sessin.InvalidateSession(handler.DB, sessionID)
		io.WriteString(writer, "{\"ok\": true}")
	} else {
		if err != nil {
			io.WriteString(writer, err.Error()+"\n")
		} else {
			io.WriteString(writer, "{\"ok\": false}")
		}
	}
}
