package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"

	"go-first-api/internal/client"
	"go-first-api/internal/handler"
)

// getEnv lee una variable de entorno o devuelve un valor por defecto.
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit: 1 * 1024 * 1024, // 1 MB máximo por petición
	})

	// Ruta principal (comprobar que el servidor está activo)
	app.Get("/home", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Dfiniendo el cliente API de Node:
	// - 5 s por intento
	// - 2 reintentos
	// - espera base de 200 ms.
	nodeClient := client.NewNodeClient(
		getEnv("NODE_API_URL", "http://localhost:3001"),
		5*time.Second,
		2,
		200*time.Millisecond,
	)

	// Definiendo rutas de la API versionada
	qrHandler := handler.NewQRHandler(nodeClient)
	api := app.Group("/api/v1")
	api.Post("/qr", qrHandler.Decompose)

	// corriendo servidor
	port := getEnv("PORT", "3000")
	log.Fatal(app.Listen(":" + port))
}
