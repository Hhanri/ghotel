package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/hhanri/ghotel/db"
	"github.com/hhanri/ghotel/db/fixtures"
	"github.com/hhanri/ghotel/types"
)

func main() {

	dbUri := flag.String("dbUri", "mongodb://localhost:27017", "DB Uri")
	flag.Parse()

	client, err := db.NewMongoClient(*dbUri)
	if err != nil {
		log.Fatal(err)
	}

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

	if err := roomStore.Drop(context.Background()); err != nil {
		log.Fatal(err)
	}

	if err := hotelStore.Drop(context.Background()); err != nil {
		log.Fatal(err)
	}

	if err := userStore.Drop(context.Background()); err != nil {
		log.Fatal(err)
	}

	if err := bookingStore.Drop(context.Background()); err != nil {
		log.Fatal(err)
	}

	hotel := fixtures.AddHotel(
		store,
		"Belluga",
		"France",
		0,
	)

	room1 := fixtures.AddRoom(
		store,
		hotel.ID,
		&types.Room{
			Size:    "small",
			Seaside: false,
			Price:   99.99,
		},
	)

	room2 := fixtures.AddRoom(
		store,
		hotel.ID,
		&types.Room{
			Size:    "king",
			Seaside: true,
			Price:   179.99,
		},
	)

	fixtures.AddUser(
		store,
		"ad",
		"min",
		"admin@ghotel.com",
		"password",
		true,
	)
	user1 := fixtures.AddUser(
		store,
		"f1",
		"l1",
		"email1@email.com",
		"password",
		false,
	)
	user2 := fixtures.AddUser(
		store,
		"f2",
		"l2",
		"email2@email.com",
		"password",
		false,
	)

	fixtures.AddBooking(
		store,
		user1.ID,
		room1.ID,
		time.Now(),
		time.Now().Add(time.Hour*24),
	)

	fixtures.AddBooking(
		store,
		user2.ID,
		room2.ID,
		time.Now(),
		time.Now().Add(time.Hour*24),
	)

	for i := range 100 {
		fixtures.AddHotel(
			store,
			fmt.Sprintf("Hotel_%d", i),
			"France",
			i%5+1,
		)
	}
}
