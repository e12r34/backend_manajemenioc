package routes

import (
	"fiberioc/controllers"

	"github.com/gofiber/fiber/v2"
)

func Router(app *fiber.App) {
	app.Post("/user", controllers.CreateUser) //add this
	app.Get("/user/:userId", controllers.GetAUser)
	app.Put("/user/:userId", controllers.EditAUser)
	app.Delete("/user/:userId", controllers.DeleteAUser)
	app.Get("/users", controllers.GetAllUsers)

	app.Post("/ioc", controllers.CreateIoc) //add this
	app.Put("/ioc/:id", controllers.EditAIoc)
	app.Delete("/ioc/:id", controllers.DeleteAIoc)
	app.Get("/iocs", controllers.GetAllIocs)
	app.Get("/ioc/csv", controllers.GetAllIocsCSV)
	app.Get("/ioc/some", controllers.GetSomeIoc)
	app.Post("/ioc/some/csv", controllers.GetSomeIocCSV)
	app.Get("/ioc/:id", controllers.GetAIoc)
	app.Post("/ioczip/csv", controllers.DuaCSV)
	app.Get("/ioczip/csv/all", controllers.DuaCSV)
	app.Post("/ioc/suricata", controllers.ExportSuricataSome)
	app.Get("/ioc/suricata/all", controllers.ExportSuricataAll)
	app.Post("/ioc/thor", controllers.ExportThorSome)
	app.Get("/ioc/thor/all", controllers.ExportThorAll)
	// app.Post("/ioczip/binalyze", controllers.DuaCSV)
	// app.Get("/ioczip/binalyze/all", controllers.DuaCSV)
	app.Post("/ioc/mgnx", controllers.ExportMgSome)
	app.Get("/ioc/mgnx/all", controllers.ExportMgAll)
	app.Post("/ioc/bae", controllers.ExportBAESome)
	app.Get("/ioc/bae/all", controllers.ExportBAEAll)
}
