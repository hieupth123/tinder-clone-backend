package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/phamtrunghieu/tinder-clone-backend/config"
	"github.com/phamtrunghieu/tinder-clone-backend/database"
	"github.com/phamtrunghieu/tinder-clone-backend/routes"
	"os"
)
var engine *gin.Engine
func init() {
	engine = gin.New()
	engine.Use(gin.Logger())
}
func main() {
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}
	config.Init()
	cfg := config.GetConfig()

	_, err := database.Init()
	if err == nil {
		fmt.Println("\nDatabase connected!")
	} else {
		fmt.Errorf("Fatal error database connection: %s \n", err)
	}
	port := cfg.GetString("server.port")
	fmt.Println("port: ", port)
	StartRest(port)
}

func StartRest(port string) {
	routes.RouteInit(engine)
	engine.Run(":" + port)
}