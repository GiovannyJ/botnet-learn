package configuration


import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

/*
*Parse the .env file and pass its attr
*/

/*
*TESTED WORKING
loads .env file
returns .env variable
*/
func EnvVar(key string) string {
	err := godotenv.Load(".env")
  
    if err != nil {
      log.Fatalf("Error loading .env file")
    } 
  
  	return os.Getenv(key)
}