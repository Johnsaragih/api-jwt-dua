package main

import (
	"fmt"
	"log"

	"api-jwt/configs"
	"api-jwt/routes"
	"net/http"
)

func main() {
	configs.LoadConfig()
	configs.LoadDB()
	defer configs.DB.Close()
	mux := http.NewServeMux()
	routes.MapRoutes(mux)
	fmt.Println("Server Aktif di Port :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
