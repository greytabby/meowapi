package main

import (
	"log"
	"os"
	"time"

	"github.com/greytabby/meowapi/lib/model"
	"github.com/greytabby/meowapi/lib/db"
	"github.com/greytabby/meowapi/lib/handler"

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
	dbAccessor.Db.AddTableWithName(model.Item{}, "items")

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

	// Create api request handler
	handler := handler.Handler{Db: dbAccessor}

	// prepare middleware
	e := echo.New()
	e.Use(middleware.Logger())
	e.Debug = true

	// Routing
	e.GET("/item", handler.GetItems)
	e.POST("/item", handler.InsertItem)
	e.PATCH("/item", handler.UpdateItem)
	e.DELETE("/item", handler.DeleteItem)

	// Service Start
	port := os.Getenv("BIND_PORT")
	err = e.Start(":" + port)
	if err != nil {
		e.Logger.Fatal(err)
		return 1
	}
	return 0
}
