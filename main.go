package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raissarib/crud-houses-of-got/server/routes"
)

func main() {

	app := fiber.New()
	routes.Animal(app)

	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}

}
