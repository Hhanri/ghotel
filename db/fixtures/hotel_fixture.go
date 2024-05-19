package fixtures

import (
	"context"
	"log"

	"github.com/hhanri/ghotel/db"
	"github.com/hhanri/ghotel/types"
)

func AddHotel(
	store *db.Store,
	name string,
	location string,
	rating int,
) *types.Hotel {
	hotel := &types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []string{},
		Rating:   rating,
	}

	insertedHotel, err := store.Hotel.Insert(context.Background(), hotel)
	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel
}
