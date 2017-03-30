package business

import (
	"io"
	"net/http"
	"time"

	"database/sql"
)

type RegistrationHandler struct {
	DB *sql.DB
}

type restaurantDetails struct {
}

type businessRegistrationCredentials struct {
	FirstName           string `json:"firstName"`
	LastName            string `json:"lastName"`
	Email               string `json:"email"`
	Password            string `json:"password"`
	PersonalPhoneNumber string `json:"personalPhoneNumber"`
	RestaurantName      string `json:"restaurantName"`
	AddressLine1        string `json:"addressLine1"`
	AddressLine2        string `json:"addressLine2"`
	City                string `json:"city"`
	State               string `json:"state"`
	Zip                 string `json:"zip"`
	BusinessPhoneNumber string `json:"businessPhoneNumber"`
	Description         string `json:"description"`
}

func (handler RegistrationHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	x := &businessRegistrationCredentials{}
	err := decodeJSON(request.Body, x)
	if err != nil {
		io.WriteString(writer, err.Error()+"\n")
		return
	}

	rows, err := handler.DB.Query("SELECT email FROM users WHERE email=?", x.Email)

	defer rows.Close()

	if err != nil {
		io.WriteString(writer, err.Error()+"\n")
		return
	}

	for rows.Next() {
		writer.WriteHeader(http.StatusConflict)
		io.WriteString(writer, "Email is already in use\n")
		return
	}

	_, err = handler.DB.Exec("INSERT INTO users (email, password, firstname, lastname, `registration-date`) VALUES (?, ?, ?, ?, ?)", x.Email, x.Password, x.FirstName, x.LastName, time.Now())
	if err != nil {
		io.WriteString(writer, err.Error()+"\n")
		return
	}
}
