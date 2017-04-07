package business

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB *sql.DB
}

type businessLoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (handler AuthHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	x := &businessLoginCredentials{}
	err := decodeJSON(request.Body, x)
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
	rows, err := handler.DB.Query("SELECT email FROM users WHERE email=? AND password=?", x.Email, string(passwordHash))

	defer rows.Close()

	if err != nil {
		io.WriteString(writer, err.Error()+"\n")
		return
	}

	if rows.Next() {
		// issue session
		io.WriteString(writer, "Congratulations, you have logged in.")
		_, err = handler.DB.Exec("INSERT INTO users (`last-login`) VALUES (?)", time.Now())
	} else {
		writer.WriteHeader(http.StatusConflict)
		io.WriteString(writer, "Email or password is incorrect\n")
		return
	}

}
