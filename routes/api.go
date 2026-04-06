package routes

import (
	"api-jwt-dua/controllers"
	"api-jwt-dua/middleware"
	"net/http"
)

func MapRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/login", controllers.LoginHandler)
	//protected routes group
	api := http.NewServeMux()
	api.Handle("/personal", http.HandlerFunc(controllers.GetPersonal))
	//	mux.Handle("/personal", middleware.JWTMiddleware(http.HandlerFunc(controllers.PersonalHandler)))
	// Pasang JWT
	mux.Handle("/api/", http.StripPrefix("/api", middleware.JWTMiddleware(api)))
}

//HandleFunc → gunakan untuk fungsi biasa (w,r)
//Handle → gunakan untuk interface http.Handler (misal middleware)
