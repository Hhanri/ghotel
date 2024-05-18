package main

import (
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/hhanri/ghotel/api"
	"github.com/hhanri/ghotel/api/middleware"
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

	client, err := db.NewMongoClient(*dbUri)
	if err != nil {
		log.Fatal(err)
	}

	// stores initialization
	userStore := db.NewMongoUserStore(client, db.DBNAME)
	hotelStore := db.NewMongoHotelStore(client, db.DBNAME)
	roomStore := db.NewMongoRoomStore(client, db.DBNAME, hotelStore)
	bookingStore := db.NewMongoBookingStore(client, db.DBNAME)

	store := &db.Store{
		User:    userStore,
		Hotel:   hotelStore,
		Room:    roomStore,
		Booking: bookingStore,
	}

	// middlewares
	jwtMiddleware := middleware.NewJWTMiddleware(store)

	// handlers initialization
	authHandler := api.NewAuthHandler(store)
	userHandler := api.NewUserHandler(store)
	hotelHandler := api.NewHotelHandler(store)
	roomHandler := api.NewRoomHandler(store)
	bookingHandler := api.NewBookingHandler(store)

	// app
	app := fiber.New(config)
	apiRoot := app.Group("/api")
	apiV1 := app.Group("/api/v1", jwtMiddleware.JWTAuthentication)
	adminRoute := apiV1.Group("/admin", middleware.AdminAuthentication)

	// auth
	apiRoot.Post("/auth", authHandler.HandleAuthenticate)

	// user handlers
	apiV1.Get("/user", userHandler.HandleGetUsers)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)
	apiV1.Post("/user", userHandler.HandlePostUser)
	apiV1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiV1.Put("/user/:id", userHandler.HandleUpdateUser)

	// hotel handlers
	apiV1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiV1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	// room handlers
	apiV1.Post("/room/:id/book", roomHandler.HandleBookRoom)
	apiV1.Get("/room", roomHandler.HandleGetAllRooms)

	// booking handlers
	adminRoute.Get("/booking", bookingHandler.HandleGetBookings)
	apiV1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiV1.Post("/booking/:id/cancel", bookingHandler.HandleCancelBooking)

	app.Listen(*listenAddr)

}

func handleFoo(c *fiber.Ctx) error {
	return c.JSON(
		map[string]string{
			"msg": "Working just fine",
		},
	)
}
