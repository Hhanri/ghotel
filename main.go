package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/hhanri/ghotel/api"
	"github.com/hhanri/ghotel/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbUri string = "mongodb://localhost:27017"

func main() {
	listenAddr := flag.String("listenAddr", "localhost:5000", "Listen address of the api server")
	flag.Parse()

	app := fiber.New()
	apiV1 := app.Group("/api/v1")

	client, err := mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(dbUri),
	)
	if err != nil {
		log.Fatal(err)
	}

	// stores initialization
	userStore := db.NewMongoUserStore(client)

	// handlers initialization
	userHandler := api.NewUserHandler(userStore)

	apiV1.Get("/user", userHandler.HandleGetUser)
	apiV1.Get("/user/:id", userHandler.HandleGetUsers)

	app.Listen(*listenAddr)

}

func handleFoo(c *fiber.Ctx) error {
	return c.JSON(
		map[string]string{
			"msg": "Working just fine",
		},
	)
}
