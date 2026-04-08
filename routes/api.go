package routes

import (
	"api-jwt-dua/controllers"
	"api-jwt-dua/middleware"

	"github.com/go-chi/chi/v5"
)

func MapRoutes(r chi.Router) {
	//public
	r.Post("/login", controllers.LoginHandler)
	//protected
	r.Route("/api", func(api chi.Router) {
		api.Use(middleware.JWTMiddleware)
		api.Get("/idbarang/{barkod}", controllers.GetBarangID) //nobarcode=65567675
		api.Get("/barang/{barkod}", controllers.GetBarang)     //nobarcode LIKE 009
		api.Get("/barang", controllers.GetAllBarang)
		api.Post("/newbarang", controllers.NewBarang)
		api.Put("/barang/{barkod}", controllers.UpdateBarang)
	})

}

/*
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
*/
