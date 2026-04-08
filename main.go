package main

import (
	"fmt"
	"log"

	"api-jwt-dua/configs"
	"api-jwt-dua/routes"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {

	configs.LoadConfig()
	configs.LoadDB()
	defer configs.DB.Close()
	r := chi.NewRouter()
	routes.MapRoutes(r)
	fmt.Println("Server Aktif di Port :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

/*
func main() {
	configs.LoadConfig()
	configs.LoadDB()
	defer configs.DB.Close()
	mux := http.NewServeMux()
	routes.MapRoutes(mux)
	fmt.Println("Server Aktif di Port :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
*/
