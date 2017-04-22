package main

import (
	"database/sql"
	"io"
	"math/rand"
	"net/http"
	"net/mail"
	"time"

	"errors"

	"github.com/NAExpire/API/src/util"
	"golang.org/x/crypto/bcrypt"
)

const senderEmail = "fill me in"
const senderPassword = "fill me in"
const smtpHost = "smtp.gmail.com"
const smtpPort = 587

type BusinessRegistrationHandler struct {
	DB *sql.DB
}

type ConsumerRegistrationHandler struct {
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

type consumerRegistrationCredentials struct {
	FirstName           string `json:"firstName"`
	LastName            string `json:"lastName"`
	Email               string `json:"email"`
	Password            string `json:"password"`
	PersonalPhoneNumber string `json:"personalPhoneNumber"`
}

func (handler ConsumerRegistrationHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	x := &consumerRegistrationCredentials{}
	err := util.DecodeJSON(request.Body, x)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		util.WriteErrorJSON(writer, "Malformed request syntax")
		return
	}

	responseCode, err := validateCredentials(handler.DB, x.Email, x.Password, x.PersonalPhoneNumber)
	if err != nil {
		writer.WriteHeader(responseCode)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	confirmationCode := randomConfirmationCode()
	_, responseCode, err = registerUser(handler.DB, x.Email, x.Password, x.FirstName, x.LastName, "customer", confirmationCode)
	if err != nil {
		writer.WriteHeader(responseCode)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	writer.WriteHeader(http.StatusCreated)
	io.WriteString(writer, "{\"ok\": true}")
}

func (handler BusinessRegistrationHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	x := &businessRegistrationCredentials{}
	err := util.DecodeJSON(request.Body, x)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		util.WriteErrorJSON(writer, "Malformed request syntax")
		return
	}

	responseCode, err := validateCredentials(handler.DB, x.Email, x.Password, x.PersonalPhoneNumber)
	if err != nil {
		writer.WriteHeader(responseCode)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	confirmationCode := randomConfirmationCode()
	insertedID, responseCode, err := registerUser(handler.DB, x.Email, x.Password, x.FirstName, x.LastName, "restaurant", confirmationCode)
	if err != nil {
		writer.WriteHeader(responseCode)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	_, err = handler.DB.Exec("INSERT INTO `restaurants` (`ownerid`, `name`, `description`, `address`, `city`, `state`, `zip`, `phone-number`, `registration-date`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", insertedID, x.RestaurantName, x.Description, x.AddressLine1+"\n"+x.AddressLine2, x.City, x.State, x.Zip, x.BusinessPhoneNumber, time.Now())
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.WriteErrorJSON(writer, err.Error())
		return
	}

	writer.WriteHeader(http.StatusCreated)
	io.WriteString(writer, "{\"ok\": true}")
}

func validateCredentials(db *sql.DB, email string, password string, phoneNumber string) (int, error) {
	emailUsed, err := userWithEmailExists(db, email)
	if err != nil {
		return http.StatusInternalServerError, err
	} else if emailUsed {
		return http.StatusConflict, errors.New("Email is already used by another user")
	} else if !validatePassword(password) {
		return http.StatusUnprocessableEntity, errors.New("Password does not pass validation")
	} else if !validatePhoneNumber(phoneNumber) {
		return http.StatusUnprocessableEntity, errors.New("Phone number does not pass validation")
	} else if !validateEmail(email) {
		return http.StatusUnprocessableEntity, errors.New("Email does not pass valdiation")
	}
	return 0, nil
}

func registerUser(db *sql.DB, email, password, firstName, lastName, userType string, confirmationCode int) (int64, int, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	result, err := db.Exec("INSERT INTO `carts` VALUES ()")
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	cartID, err := result.LastInsertId()
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	registrationDate := time.Now()
	result, err = db.Exec("INSERT INTO `users` (`email`, `password`, `firstname`, `lastname`, `type`, `confirmation-code`, `registration-date`, `cart-id`) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", email, string(passwordHash), firstName, lastName, userType, confirmationCode, registrationDate, cartID)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	return userID, 0, nil
}

func userWithEmailExists(db *sql.DB, email string) (bool, error) {
	rows, err := db.Query("SELECT `id` FROM `users` WHERE email=?", email)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}

func randomConfirmationCode() int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(89999999) + 10000000
}

func validatePassword(password string) bool {
	return len(password) > 7
}

func validatePhoneNumber(phoneNumber string) bool {
	return true
}

func validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}
	// auth := smtp.PlainAuth("", senderEmail, senderPassword, "smtp.gmail.com")
	// msg := "From: " + senderEmail + "\n" +
	// 	"To: " + email + "\n" +
	// 	"Subject: NAExpire Registration\n\n" +
	// 	"Hello! Your confirmation code is " + strconv.Itoa(confirmationCode) + "."
	// err = smtp.SendMail(smtpHost+":"+strconv.Itoa(smtpPort), auth, senderEmail, []string{email}, []byte(msg))
	// if err != nil {
	// 	return true, err
	// }
	return true
}
