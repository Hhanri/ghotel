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

	hotelStore := db.NewMongoHotelStore(client, db.DBNAME)
	roomStore := db.NewMongoRoomStore(client, db.DBNAME, hotelStore)

	if err := roomStore.Drop(context.Background()); err != nil {
		log.Fatal(err)
	}

	if err := hotelStore.Drop(context.Background()); err != nil {
		log.Fatal(err)
	}

	hotel := &types.Hotel{
		Name:     "Belluga",
		Location: "France",
		Rooms:    []string{},
	}

	roomA := &types.Room{
		Type:      types.SinglePersonRoomType,
		BasePrice: 99.99,
	}
	roomB := &types.Room{
		Type:      types.DeluxeRoomType,
		BasePrice: 179.99,
	}

	insertHotelAndRooms(
		hotelStore,
		roomStore,
		hotel,
		[]*types.Room{
			roomA,
			roomB,
		},
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
