package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/hhanri/ghotel/db"
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

	if err := roomStore.Drop(context.Background()); err != nil {
		log.Fatal(err)
	}

	if err := hotelStore.Drop(context.Background()); err != nil {
		log.Fatal(err)
	}

	seedHotel(
		hotelStore,
		roomStore,
		"Belluga",
		"France",
		0,
	)

	seedHotel(
		hotelStore,
		roomStore,
		"Bellutwo",
		"France",
		5,
	)

	seedHotel(
		hotelStore,
		roomStore,
		"Cheese",
		"United Kingdom",
		3,
	)

	seedUser(
		userStore,
		"f1",
		"l1",
		"email1@email.com",
		"password",
	)
	seedUser(
		userStore,
		"f2",
		"l2",
		"email2@email.com",
		"password",
	)

}

func insertHotelAndRooms(
	hotelStore db.HotelStore,
	roomStore db.RoomStore,
	hotel *types.Hotel,
	rooms []*types.Room,
) error {
	ctx := context.Background()

	insertedHotel, err := hotelStore.Insert(ctx, hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		_, err := roomStore.Insert(ctx, room)
		if err != nil {
			return err
		}
		fmt.Printf("Inserted %+v\n", room)
	}

	return nil
}

func seedHotel(
	hotelStore db.HotelStore,
	roomStore db.RoomStore,
	name string,
	location string,
	rating int,
) error {
	hotel := &types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []string{},
		Rating:   rating,
	}
	roomA := &types.Room{
		Size:    "small",
		Seaside: false,
		Price:   99.99,
	}
	roomB := &types.Room{
		Size:    "king",
		Seaside: true,
		Price:   179.99,
	}

	return insertHotelAndRooms(
		hotelStore,
		roomStore,
		hotel,
		[]*types.Room{
			roomA,
			roomB,
		},
	)
}

func seedUser(userStore db.UserStore, firstName, lastName, email, password string) error {
	params := types.CreateUserParams{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
	}
	user, _ := types.NewUserFromParams(params)
	_, err := userStore.InsertUser(context.Background(), user)
	return err
}
