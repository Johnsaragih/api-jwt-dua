package routes

import (
	"api-jwt/controllers"
	"api-jwt/middleware"
	"net/http"
)

func MapRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/login", controllers.LoginHandler)

	mux.Handle("/personal", middleware.JWTMiddleware(http.HandlerFunc(controllers.PersonalHandler)))

}

//HandleFunc → gunakan untuk fungsi biasa (w,r)
//Handle → gunakan untuk interface http.Handler (misal middleware)
