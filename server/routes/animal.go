package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raissarib/crud-houses-of-got/server/controllers"
)

func Animal(app *fiber.App) {
	group := app.Group("animal/")
	controller := controllers.NewAnimal()

	group.Post("", controller.Create)
	group.Get("one/:id", controller.GetOne)
	group.Put(":id", controller.Update)
	group.Delete(":id", controller.Delete)
	group.Get("csv/", controller.GetCSV)
	group.Post("csv/", controller.BulkCreate)
}
