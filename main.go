package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"search_engine/db"
	"search_engine/routes"
	"search_engine/utils"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/joho/godotenv"
)

func main() {
	env := godotenv.Load()
	if env != nil {
		panic("cannot find environment variables")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = ":4000"
	} else {
		port = ":" + port
	}
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		panic("DATABASE_URL is not set")
	}

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		panic("SECRET_KEY is not set")
	}

	app := fiber.New(fiber.Config{
		IdleTimeout: 5 * time.Second,
	})

	app.Use(compress.New())
	db.InitDB()
	routes.SetRoutes(app)
	utils.StartCronJobs()

	go func() {
		if err := app.Listen(port); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c // Block the main thread until interupted
	app.Shutdown()
	fmt.Println("shutting down server")
}
