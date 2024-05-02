package main

import (
	"flag"

	"github.com/gofiber/fiber/v2"
	"github.com/hhanri/ghotel/api"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(
			map[string]string{"error": err.Error()},
		)
	},
}

func main() {
	listenAddr := flag.String("listenAddr", "localhost:5000", "Listen address of the api server")
	flag.Parse()

	app := fiber.New(config)

	apiV1 := app.Group("/api/v1")

	apiV1.Get("/user", api.HandleGetUsers)
	apiV1.Get("/user/:id", api.HandleGetUser)

	app.Listen(*listenAddr)

}

func handleFoo(c *fiber.Ctx) error {
	return c.JSON(
		map[string]string{
			"msg": "Working just fine",
		},
	)
}
