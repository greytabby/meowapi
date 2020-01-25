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
	dbAccessor.Db.AddTableWithName(model.Wash{}, "wash")
	dbAccessor.Db.AddTableWithName(model.User{}, "user")

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
	e.Use(middleware.Recover())
	e.Debug = true

	// Create api request handler
	catHandler := handler.CatHandler{Db: dbAccessor}
	toiletHandler := handler.ToiletHandler{Db: dbAccessor}
	useToiletHandler := handler.UseToiletHandler{Db: dbAccessor}
	washHandler := handler.WashHandler{Db: dbAccessor}
	authHandler := handler.AuthHandler{Db: dbAccessor}

	// Routing
	// Cat Endpoint
	// Use JWT authentication
	// TODO: apply JWT auth to other endpoint
	r := e.Group("/api")
	r.Use(middleware.JWTWithConfig(handler.JWTConfig))
	r.GET("/cat", catHandler.GetAllCats)
	r.POST("/cat", catHandler.AddCat)
	r.PUT("/cat", catHandler.UpdateCat)
	r.DELETE("/cat", catHandler.DeleteCat)

	// Toilet Endpoint
	r.GET("/toilet", toiletHandler.GetAllToilets)
	r.POST("/toilet", toiletHandler.AddToilet)
	r.PUT("/toilet", toiletHandler.UpdateToilet)
	r.DELETE("/toilet", toiletHandler.DeleteToilet)

	// UseToilet Endpoint
	r.GET("/usetoilet", useToiletHandler.GetAllUseToilets)
	r.POST("/usetoilet", useToiletHandler.AddUseToilet)
	r.PUT("/usetoilet", useToiletHandler.UpdateUseToilet)
	r.DELETE("/usetoilet", useToiletHandler.DeleteUseToilet)

	// Wash Endpoint
	e.GET("/api/wash", washHandler.GetAllWashes)
	e.GET("/api/wash/:toiletid", washHandler.GetWashesByToiletId)
	e.POST("/api/wash", washHandler.AddWash)
	e.PUT("/api/wash", washHandler.UpdateWash)
	e.DELETE("/api/wash", washHandler.DeleteWash)

	// Auth Endpiont
	e.POST("/signup", authHandler.Signup)
	e.POST("/login", authHandler.Login)

	// Service Start
	port := os.Getenv("BIND_PORT")
	err = e.Start(":" + port)
	if err != nil {
		e.Logger.Fatal(err)
		return 1
	}
	return 0
}
