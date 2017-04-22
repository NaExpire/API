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
	businessRouter.Handle("/restaurant/{restaurantID}/", Chain(GetRestaurantHandler{DB: db}, AllowCORS())).
		Methods("GET")
	businessRouter.Handle("/restaurant/{restaurantID}/update/", Chain(UpdateRestaurantHandler{DB: db}, AllowCORS())).
		Methods("POST")
	businessRouter.Handle("/meal/{mealID}/", Chain(GetMealHandler{DB: db}, AllowCORS())).
		Methods("GET")
	businessRouter.Handle("/meal/{mealID}/update/", Chain(UpdateMealHandler{DB: db}, AllowCORS())).
		Methods("PUT")
	businessRouter.Handle("/meal/create/", Chain(AddMealHandler{DB: db}, AllowCORS())).
		Methods("POST")
	businessRouter.Handle("/meal/{mealID}/delete/", Chain(DeleteMealHandler{DB: db}, AllowCORS())).
		Methods("DELETE")
	businessRouter.Handle("/deal/{dealID}/", Chain(GetDealHandler{DB: db}, AllowCORS())).
		Methods("GET")
	businessRouter.Handle("/deal/{dealID}/update/", Chain(UpdateDealHandler{DB: db}, AllowCORS())).
		Methods("PUT")
	businessRouter.Handle("/deal/create/", Chain(AddDealHandler{DB: db}, AllowCORS())).
		Methods("POST")
	businessRouter.Handle("/deal/{dealID}/delete/", Chain(DeleteDealHandler{DB: db}, AllowCORS())).
		Methods("DELETE")
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
	consumerRouter.Handle("/restaurant/{restaurantID}/", Chain(GetRestaurantHandler{DB: db}, AllowCORS())).
		Methods("GET")
	consumerRouter.Handle("/meal/{mealID}/", Chain(GetMealHandler{DB: db}, AllowCORS())).
		Methods("GET")
	consumerRouter.Handle("/deal/{dealID}/", Chain(GetDealHandler{DB: db}, AllowCORS())).
		Methods("GET")
	consumerRouter.Handle("/cart/", Chain(GetCartHandler{DB: db}, AllowCORS())).
		Methods("GET")
	consumerRouter.Handle("/cart/empty/", Chain(DeleteCartContentsHandler{DB: db}, AllowCORS())).
		Methods("DELETE")
	consumerRouter.Handle("/cart/add/meal/", Chain(AddMealCartHandler{DB: db}, AllowCORS())).
		Methods("POST")
	consumerRouter.Handle("/cart/add/deal/", Chain(AddDealCartHandler{DB: db}, AllowCORS())).
		Methods("POST")
	consumerRouter.Handle("/cart/delete/meal/", Chain(DeleteMealCartHandler{DB: db}, AllowCORS())).
		Methods("DELETE")
	consumerRouter.Handle("/cart/delete/deal/", Chain(DeleteDealCartHandler{DB: db}, AllowCORS())).
		Methods("DELETE")
}
