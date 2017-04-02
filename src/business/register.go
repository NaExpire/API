package business

import (
	"database/sql"
	"io"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type RegistrationHandler struct {
	DB *sql.DB
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

	if rows.Next() {
		writer.WriteHeader(http.StatusConflict)
		io.WriteString(writer, "Email is already in use\n")
		return
	}

	registrationDate := time.Now()
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(x.Password), 14)

	result, err := handler.DB.Exec("INSERT INTO users (email, password, firstname, lastname, `registration-date`) VALUES (?, ?, ?, ?, ?)", x.Email, string(passwordHash), x.FirstName, x.LastName, registrationDate)

	if err != nil {
		io.WriteString(writer, err.Error()+"\n")
		return
	}

	lastID, _ := result.LastInsertId()

	_, err = handler.DB.Exec("INSERT INTO restaurants (ownerid, name, description, address, city, state, zip, `registration-date`) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", lastID, x.RestaurantName, x.Description, x.AddressLine1+"\n"+x.AddressLine2, x.City, x.State, x.Zip, registrationDate)

	if err != nil {
		io.WriteString(writer, err.Error()+"\n")
		return
	}

	io.WriteString(writer, "{\"ok\": true}")
}
