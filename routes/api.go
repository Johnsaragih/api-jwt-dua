package routes

import (
	"api-jwt-dua/controllers"
	"api-jwt-dua/middleware"
	"net/http"
)

func MapRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/login", controllers.LoginHandler)
	//protected route groups
	api := http.NewServeMux()
	api.Handle("/personal", middleware.AllowMethods(http.MethodGet)(controllers.GetPersonal))
	api.Handle("/barang/", middleware.AllowMethods(http.MethodGet)(controllers.GetBarang))
	api.Handle("/barang", middleware.AllowMethods(http.MethodGet)(controllers.GetAllBarang))
	// Pasang JWT
	mux.Handle("/api/", http.StripPrefix("/api", middleware.JWTMiddleware(api)))
}
