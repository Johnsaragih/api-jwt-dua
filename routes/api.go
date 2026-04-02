package routes

import (
	"api-jwt/controllers"
	"net/http"
)

func MapRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/login", controllers.LoginHandler)
	mux.HandleFunc("/personal", controllers.PersonalHandler)

}
