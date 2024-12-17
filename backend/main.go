package main

import (
	"backend/HTTP"
	"backend/Mongo"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("/root/.env")

	var port string
	if err == nil {
		port = os.Getenv("BACKEND_PORT")
	} else {
		port = "8080"
	}

	//var endpointRouter = HTTP.Routes{} // Inicializacija router-jev za endpoint-e
	Mongo.ConnectToMongoDB() // Vzpostavitev povezave s podatkovno bazo MongoDB
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	HTTP.Router(router)

	if err := router.Run("0.0.0.0:" + port); err != nil {
		panic(err)
	}

}
