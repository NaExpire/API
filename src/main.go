package main

import (
	"net/http"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"fmt"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize connection to DB
	db, err := sql.Open("mysql", "root:root@tcp(138.197.33.88:3306)/naexpire?parseTime=true&charset=utf8")
	panicOnErr(err)
	err = db.Ping()
	panicOnErr(err)

	// Initialize request routing
	apiRouter := mux.NewRouter().
		StrictSlash(false)
	initBusinessRouter(apiRouter, db)
	initConsumerRotuer(apiRouter, db)

	// Listen for incoming connections on port 8000
	http.ListenAndServe(":8000", apiRouter)
}

func panicOnErr(err error) {
	if err != nil {
		fmt.Printf("%s", err.Error())
		panic(err)
	}
}

func initBusinessRouter(parent *mux.Router, db *sql.DB) {
	// All subrouted requests will be suffixes of the URL pattern /api/business
	businessRouter := parent.PathPrefix("/api/business").
		Subrouter()

	// e.g. /api/business/login/
	businessRouter.Handle("/login/", Chain(LoginHandler{DB: db}, AllowCORS())).
		Methods("POST")
	businessRouter.Handle("/logout/", Chain(LogoutHandler{DB: db}, AllowCORS())).
		Methods("POST")
	businessRouter.Handle("/register/", Chain(BusinessRegistrationHandler{DB: db}, AllowCORS())).
		Methods("POST")
	businessRouter.Handle("/register/confirm/", Chain(ConfirmRegistrationHandler{DB: db}, AllowCORS())).
		Methods("POST")
	// businessRouter.HandleFunc("/restaurant/{restaurantID}/menu/{menuItemID}", MenuGetHandler)
	// businessRouter.HandleFunc("/restaurant/{restaurantID}/menu/{menuItemID}/update/", MenuUpdateHandler).
	// 	Methods("POST")
	businessRouter.Handle("/restaurant/{restaurantID}/", Chain(GetRestaurantHandler{DB: db}, AllowCORS())).
		Methods("GET")
	businessRouter.Handle("/restaurant/{restaurantID}/update/", Chain(UpdateRestaurantHandler{DB: db}, AllowCORS())).
		Methods("POST")
	// businessRouter.HandleFunc("/discount/create/{restaurantID}/{menuItemID}", DiscountCreateHandler).
	// 	Methods("POST")
}

func initConsumerRotuer(parent *mux.Router, db *sql.DB) {
	consumerRouter := parent.PathPrefix("/api/consumer").
		Subrouter()

	consumerRouter.Handle("/login/", Chain(LoginHandler{DB: db}, AllowCORS())).
		Methods("POST")
	consumerRouter.Handle("/logout/", Chain(LogoutHandler{DB: db}, AllowCORS())).
		Methods("POST")
	consumerRouter.Handle("/register/", Chain(ConsumerRegistrationHandler{DB: db}, AllowCORS())).
		Methods("POST")
	consumerRouter.Handle("/register/confirm/", Chain(ConfirmRegistrationHandler{DB: db}, AllowCORS())).
		Methods("POST")
	consumerRouter.Handle("/restaurant/{restaurantID}/update/", Chain(UpdateRestaurantHandler{DB: db}, AllowCORS())).
		Methods("POST")
	consumerRouter.Handle("/restaurant/{restaurantID}/", Chain(GetRestaurantHandler{DB: db}, AllowCORS())).
		Methods("GET")
}
