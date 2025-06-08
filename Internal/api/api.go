package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"micro-service/Internal/api/middleware"
	"micro-service/Internal/service"
)

type Routers struct {
	Service service.Service
}

func NewRouters(r *Routers, token string) *fiber.App {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowMethods:  "GET, POST, PUT, PATCH, DELETE",
		AllowHeaders:  "Accept, Authorization, Content-Type, X-CSRF-Token, X-REQUEST-ID",
		ExposeHeaders: "Link",
		MaxAge:        300,
	}))
	apiGroup := app.Group("/v1", middleware.Authorization(token))

	apiGroup.Post("/tasks", r.Service.CreateTask)
	apiGroup.Get("/tasks/:id", r.Service.GetTask)
	apiGroup.Put("/tasks/:id", r.Service.PutTask)
	apiGroup.Delete("/tasks/:id", r.Service.DeleteTask)
	apiGroup.Get("/tasks", r.Service.GetTasks)
	return app
}
