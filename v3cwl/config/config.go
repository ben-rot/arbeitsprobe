package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)





type Config struct {

	DB   	PostgresConfig
	Auth  	AuthConfig
	Clash 	ClashConfig
}



type PostgresConfig struct {

	Host     	string
	Port     	int
	User     	string
	Password 	string
	Database 	string
}


type AuthConfig struct {

	ClientId     	string
	ClientSecret 	string
	CallbackUrl  	string
}


type ClashConfig struct {
	ApiKey 		string
	BaseURL 	string
}


func LoadConfig() (*Config, error) {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: No .env file found")
	}


	portStr := os.Getenv("DB_PORT")
	port, err := strconv.Atoi(portStr)
    if err != nil {
        // Fallback to a default port if the env is missing or invalid
        port = 5432 
    }





	cfg := &Config {
		
		Auth: AuthConfig {
			ClientId:     os.Getenv("DISCORD_CLIENT_ID"),
			ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
			CallbackUrl:  os.Getenv("DISCORD_CALLBACK_URL"),
		},
		
		DB: PostgresConfig {
			Host:     os.Getenv("DB_HOST"),
			Port:     port,
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Database: os.Getenv("DB_NAME"),
		},

		Clash: ClashConfig {
			ApiKey: os.Getenv("CLASH_OF_CLANS_API_KEY"),
			BaseURL: "https://api.clashofclans.com/v1",
		},
	}

	return cfg, nil
}