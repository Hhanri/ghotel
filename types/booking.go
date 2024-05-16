package types

import "time"

type Booking struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	RoomdID   string    `bson:"roomId" json:"roomId"`
	UserID    string    `bson:"userId" json:"userId"`
	NumPeople int       `bson:"numPeople" json:"numPeople"`
	FromDate  time.Time `bson:"fromDate" json:"fromDate"`
	UntilDate time.Time `bson:"untilDate" json:"untilDate"`
}
