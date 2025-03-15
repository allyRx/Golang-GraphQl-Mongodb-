//Helper function for load the environement variable

package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)



func LoadEnvUrl() string{
	err := godotenv.Load()

	if err != nil{
		log.Fatal("Error loading .env file ")
	}
	return os.Getenv("MONGOURL")	
}