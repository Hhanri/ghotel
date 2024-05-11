package main

import (
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/hhanri/ghotel/api"
	"github.com/hhanri/ghotel/db"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(
			map[string]string{"error": err.Error()},
		)
	},
}

func main() {

	dbUri := flag.String("dbUri", "mongodb://localhost:27017", "DB Uri")
	listenAddr := flag.String("listenAddr", ":9090", "Listen address of the api server")
	flag.Parse()

	app := fiber.New(config)

	apiV1 := app.Group("/api/v1")

	client, err := db.NewMongoClient(*dbUri)
	if err != nil {
		log.Fatal(err)
	}

	// stores initialization
	userStore := db.NewMongoUserStore(client, db.DBNAME)
	hotelStore := db.NewMongoHotelStore(client, db.DBNAME)
	roomStore := db.NewMongoRoomStore(client, db.DBNAME, hotelStore)

	// handlers initialization
	userHandler := api.NewUserHandler(userStore)
	hotelHandler := api.NewHotelHandler(hotelStore, roomStore)

	// user handlers
	apiV1.Get("/user", userHandler.HandleGetUsers)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)
	apiV1.Post("/user", userHandler.HandlePostUser)
	apiV1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiV1.Put("/user/:id", userHandler.HandleUpdateUser)

	// hotel handlers
	apiV1.Get("/hotel", hotelHandler.HandleGetHotels)

	app.Listen(*listenAddr)

}

func handleFoo(c *fiber.Ctx) error {
	return c.JSON(
		map[string]string{
			"msg": "Working just fine",
		},
	)
}
