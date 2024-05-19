package fixtures

import (
	"context"
	"log"

	"github.com/hhanri/ghotel/db"
	"github.com/hhanri/ghotel/types"
)

func AddRoom(store *db.Store, hotelID string, room *types.Room) *types.Room {
	room.HotelID = hotelID
	insertedRoom, err := store.Room.Insert(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}
