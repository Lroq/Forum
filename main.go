package main

import (
	"Forum/database"
	"Forum/routes"
	"Forum/handlers"
	"fmt"
	"log"
	"net/http"
)

func main() {
	database.InitDB("./users.db")
	defer database.CloseDB()

	go handlers.HubInstance.Run()

	r := routes.InitRoutes()

	r.HandleFunc("/create_admin", handlers.CreateAdminHandler)
	fmt.Println("Le serveur écoute sur le port 8080 http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
