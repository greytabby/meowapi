package main

import (
	"log"
	"os"
	"time"

	"github.com/greytabby/meowapi/lib/db"
	"github.com/greytabby/meowapi/lib/handler"
	"github.com/greytabby/meowapi/lib/model"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	exitCode := run()
	os.Exit(exitCode)
}

func run() int {
	// Initialize database connection
	dsn := os.Getenv("DATA_SOURCE_NAME")
	dbAccessor, err := db.NewMysqlDbAccessor(dsn)
	if err != nil {
		log.Fatalf("Can not create db accessor. %v\n", err)
		return 1
	}
	defer dbAccessor.Db.Db.Close()

	// Prepare database
	dbAccessor.Db.AddTableWithName(model.Cat{}, "cat")
	dbAccessor.Db.AddTableWithName(model.Toilet{}, "toilet")
	dbAccessor.Db.AddTableWithName(model.UseToilet{}, "usetoilet")

	for i := 0; i < 10; i++ {
		err = dbAccessor.Db.CreateTablesIfNotExists()
		if err == nil {
			break
		}
		log.Printf("Can not connect database %d times. %v\n", i+1, err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Printf("Can not create table. %v\n", err)
	}

	// prepare middleware
	e := echo.New()
	e.Use(middleware.Logger())
	e.Debug = true

	// Create api request handler
	catHandler := handler.CatHandler{Db: dbAccessor}
	toiletHandler := handler.ToiletHandler{Db: dbAccessor}
	useToiletHandler := handler.UseToiletHandler{Db: dbAccessor}

	// Routing
	// Cat Endpoint
	e.GET("/api/cat", catHandler.GetAllCats)
	e.POST("/api/cat", catHandler.AddCat)
	e.PUT("/api/cat", catHandler.UpdateCat)
	e.DELETE("/api/cat", catHandler.DeleteCat)

	// Toilet Endpoint
	e.GET("/api/toilet", toiletHandler.GetAllToilets)
	e.POST("/api/toilet", toiletHandler.AddToilet)
	e.PUT("/api/toilet", toiletHandler.UpdateToilet)
	e.DELETE("/api/toilet", toiletHandler.DeleteToilet)

	// UseToilet Endpoint
	e.GET("/api/usetoilet", useToiletHandler.GetAllUseToilets)
	e.POST("/api/usetoilet", useToiletHandler.AddUseToilet)
	e.PUT("/api/usetoilet", useToiletHandler.UpdateUseToilet)
	e.DELETE("/api/usetoilet", useToiletHandler.DeleteUseToilet)

	// Service Start
	port := os.Getenv("BIND_PORT")
	err = e.Start(":" + port)
	if err != nil {
		e.Logger.Fatal(err)
		return 1
	}
	return 0
}
