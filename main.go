package main

import (
	"log"	

	"cwl/database"
	"cwl/api"
	"cwl/config"
)





func main() {

	listenAddr := ":3000"

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}


	db, err := database.ConnectDb(&cfg.DB)
	if err != nil {
        log.Fatalf("Database initialization failed: %v", err)
    }
	

	store := database.NewStore(db)
	defer db.Close()


	clashClient := api.NewClashClient(&cfg.Clash)
	defer clashClient.Close()

	
	log.Println("Running on port: ", listenAddr)
	server := api.NewServer(listenAddr, cfg, store, clashClient)
	log.Fatal(server.Start())
}