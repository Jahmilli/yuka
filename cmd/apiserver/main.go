package main

import (
	"context"
	"os"
	"yuka/internal/database"
	"yuka/internal/routers"
	"yuka/pkg/utils"

	"gorm.io/gorm"
)

// @title          Yuka API
// @version        1.0
// @description	This is the Yuka API Server.
// @securityDefinitions.basic BasicAuth

// @BasePath  		/
func main() {
	ctx := context.Background()
	logger, err := utils.GetLogger()
	if err != nil {
		logger.Fatal(err.Error())
	}

	// Local, Prod etc
	environment := os.Getenv("ENVIRONMENT")
	databaseHostname := os.Getenv("DATABASE_HOSTNAME")
	databaseUsername := os.Getenv("DATABASE_USERNAME")
	databasePassword := os.Getenv("DATABASE_PASSWORD")
	databaseName := os.Getenv("DATABASE_NAME")
	databasePort := os.Getenv("DATABASE_PORT")

	var db *gorm.DB
	if environment == "local" {
		// db, err = database.ConnectTestDatabase(logger)
		db, err = database.ConnectDatabase(context.TODO(), logger.Sugar(), databaseHostname, databaseUsername, databasePassword, databaseName, databasePort, "disable")
	} else {
		db, err = database.ConnectDatabase(context.TODO(), logger.Sugar(), databaseHostname, databaseUsername, databasePassword, databaseName, databasePort, "disable")
	}
	if err != nil {
		logger.Fatal(err.Error())
	}

	routerOptions := routers.NewRouterOptions(logger, db)

	if err := routers.Run(ctx, &routerOptions); err != nil {
		logger.Fatal(err.Error())

	}
}
