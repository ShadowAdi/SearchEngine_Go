package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/joho/godotenv"
)

func main() {
	env := godotenv.Load()
	if env != nil {
		panic("Cannot Find Environment Variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":4000"
	}

	app := fiber.New(fiber.Config{
		IdleTimeout: 5 * time.Second,
	})

	app.Use(compress.New())

	go func() {
		if err := app.Listen(port); err != nil {
			log.Panic("ERROR Happend ", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	app.Shutdown()
	fmt.Println("Shutting Down The Server")

}
