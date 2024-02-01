package main

import (
	"log"
	"net/http"
	"os"

	"github.com/daneofmanythings/diceroni/internal/config"
	"github.com/daneofmanythings/diceroni/internal/handlers"
	"github.com/daneofmanythings/diceroni/internal/routes"
	"github.com/joho/godotenv"
)

var app config.Config

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("error loading .env: %s", err.Error())
	}

	portNumber := os.Getenv("PORT")
	// dbURL := os.Getenv("DBURL")
	//
	// db, err := sql.Open("postgres", dbURL)
	// if err != nil {
	// 	log.Fatalf("error creating database: %s", err)
	// }
	//
	// app.DB = database.New(db)
	//
	repo := handlers.NewRepo(&app)
	handlers.LinkRepository(repo)

	// routers := []routes.RouterPath{
	// 	routes.V1ROUTER,
	// }

	// debugMode := flag.Bool("debug", false, "Enable debug mode")
	// flag.Parse()
	//
	// if *debugMode {
	// 	log.Println("Starting in DEBUG mode")
	// 	routers = append(routers, routes.TESTROUTER)
	// }
	//
	handler := routes.InitRouter()

	server := http.Server{
		Addr:    portNumber,
		Handler: handler,
	}

	log.Printf("Starting server on port: %s", portNumber)
	log.Fatal(server.ListenAndServe())
}
