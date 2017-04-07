package business

import (
	"database/sql"
	"io"
	"math/rand"
	"net/http"
	"net/mail"
	"net/smtp"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const senderEmail = "fill me in"
const senderPassword = "fill me in"
const smtpHost = "smtp.gmail.com"
const smtpPort = 587

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

	// validate phone and address
	confirmationCode := randomConfirmationCode()
	if !validatePassword(x.Password) {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		io.WriteString(writer, "Password does not pass validation")
	}
	if !validatePhoneNumber(x.BusinessPhoneNumber) {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		io.WriteString(writer, "Phone number does not pass validation")
	}
	emailValidated, err := validateEmail(x.Email, confirmationCode)
	if !emailValidated {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		io.WriteString(writer, "Email is not valid")
	} else if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		io.WriteString(writer, "Could not send confirmation email")
	}

	registrationDate := time.Now()
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(x.Password), 14)

	result, err := handler.DB.Exec("INSERT INTO users (email, password, firstname, lastname, `registration-date`, `confirmation-code`) VALUES (?, ?, ?, ?, ?, ?)", x.Email, string(passwordHash), x.FirstName, x.LastName, registrationDate, confirmationCode)

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

func randomConfirmationCode() int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(89999999) + 10000000
}

func validatePassword(password string) bool {
	return true
}

func validatePhoneNumber(phoneNumber string) bool {
	return true
}

func validateEmail(email string, confirmationCode int) (bool, error) {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false, nil
	}
	auth := smtp.PlainAuth("", senderEmail, senderPassword, "smtp.gmail.com")
	msg := "From: " + senderEmail + "\n" +
		"To: " + email + "\n" +
		"Subject: NAExpire Registration\n\n" +
		"Hello! Your confirmation code is " + strconv.Itoa(confirmationCode) + "."
	err = smtp.SendMail(smtpHost+":"+strconv.Itoa(smtpPort), auth, senderEmail, []string{email}, []byte(msg))
	if err != nil {
		return true, err
	}
	return true, nil
}
