package main

import (
	"log"

	"os"

	"github.com/gofiber/fiber/v3"

	"go-first-api/internal/handler"
)

func main() {

	app := fiber.New(fiber.Config{
		BodyLimit: 1 * 1024 * 1024, // 1 MB máximo por solicitud
	})

	// Ruta principal
	app.Get("/home", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Rutas de calculo versionada
	qrHandler := handler.NewQRHandler()
	api := app.Group("/api/v1")
	api.Post("/qr", qrHandler.Decompose)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}
