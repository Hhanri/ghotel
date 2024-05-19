package fixtures

import (
	"context"
	"log"
	"time"

	"github.com/hhanri/ghotel/db"
	"github.com/hhanri/ghotel/types"
)

func AddBooking(
	store *db.Store,
	userID string,
	roomID string,
	from time.Time,
	until time.Time,
) *types.Booking {
	booking := types.Booking{
		RoomdID:   roomID,
		UserID:    userID,
		FromDate:  from,
		UntilDate: until,
		NumPeople: 1,
	}
	insertedBooking, err := store.Booking.Insert(context.Background(), &booking)
	if err != nil {
		log.Fatal(err)
	}
	return insertedBooking
}
