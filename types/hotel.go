package types

type Hotel struct {
	ID       string   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string   `bson:"name" json:"name"`
	Location string   `bson:"location" json:"location"`
	Rooms    []string `bson:"rooms" json:"rooms"`
	Rating   int      `bson:"rating" json:"rating"`
}

type Room struct {
	ID      string  `bson:"_id,omitempty" json:"id,omitempty"`
	Price   float64 `bson:"price" json:"price"`
	HotelID string  `bson:"hotelId" json:"hotelId"`

	// small, normal, king
	Size    string `bson:"size" json:"size"`
	Seaside bool   `bson:"seaside" json:"seaside"`
}
