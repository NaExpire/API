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
	err := util.DecodeJSON(request.Body, x)
	fmt.Printf("Got %s request to LoginHandler\n", request.Method)
	if request.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
	} else if err != nil {
		io.WriteString(writer, err.Error()+"\n")
		return
	}

	rows, err := handler.DB.Query("SELECT password, confirmed FROM users WHERE email=?", x.Email)

	defer rows.Close()

	if err != nil {
		io.WriteString(writer, err.Error()+"\n")
		return
	}

	if rows.Next() {
		passwordHash := ""
		confirmed := 0
		err := rows.Scan(&passwordHash, &confirmed)
		if err != nil {
			io.WriteString(writer, err.Error()+"\n")
		}

		fmt.Printf("passwordHash: %s\n", passwordHash)

		if bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(x.Password)) == nil {
			if confirmed == 0 {
				io.WriteString(writer, "Email has not been confirmed yet\n")
			} else {
				myUniqueSessionID := seshin.GenerateSessionID()
				seshin.CreateSession(handler.DB, myUniqueSessionID)
				_, err = handler.DB.Exec("INSERT INTO users (`last-login`) VALUES (?)", time.Now())

				writer.WriteHeader(http.StatusOK)
				writer.Header().Add("session", myUniqueSessionID)
				io.WriteString(writer, "{\"ok\": true}")
			}
		} else {
			writer.WriteHeader(http.StatusConflict)
			io.WriteString(writer, "Email or password is incorrect!\n")
			return
		}
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
	}
	sessionID := request.Header.Get("session")

	isValidSession, err := seshin.ValidateSession(handler.DB, sessionID)
	if isValidSession {
		seshin.InvalidateSession(handler.DB, sessionID)
		writer.WriteHeader(http.StatusOK)
		io.WriteString(writer, "{\"ok\": true}")
	} else {
		if err != nil {
			io.WriteString(writer, err.Error()+"\n")
		} else {
			io.WriteString(writer, "{\"ok\": false}")
		}
	}
}
