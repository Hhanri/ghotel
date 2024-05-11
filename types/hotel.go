package types

type Hotel struct {
	ID       string   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string   `bson:"name" json:"name"`
	Location string   `bson:"location" json:"location"`
	Rooms    []string `bson:"rooms" json:"rooms"`
}

type RoomType int

const (
	_ RoomType = iota
	SinglePersonRoomType
	DoubleRoomType
	SeeSideRoomType
	DeluxeRoomType
)

type Room struct {
	ID        string   `bson:"_id,omitempty" json:"id,omitempty"`
	Type      RoomType `bson:"type" json:"type"`
	BasePrice float64  `bson:"basePrice" json:"basePrice"`
	Price     float64  `bson:"price" json:"price"`
	HotelID   string   `bson:"hotelId" json:"hotelId"`
}