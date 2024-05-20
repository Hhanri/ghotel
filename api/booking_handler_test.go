package api

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/hhanri/ghotel/api/api_error"
	"github.com/hhanri/ghotel/api/middleware"
	"github.com/hhanri/ghotel/db/fixtures"
	"github.com/hhanri/ghotel/types"
)

func TestUserGetBooking(t *testing.T) {

	db := setup(t)
	defer db.teardown(t)

	bookingHandler := NewBookingHandler(db.Store)
	jwt := middleware.NewJWTMiddleware(db.Store)
	app := newApp()
	jwtApp := app.Group("/", jwt.JWTAuthentication)
	jwtApp.Get("/:id", bookingHandler.HandleGetBooking)

	//admin := fixtures.AddUser(db.Store, "ad", "min", "admin@ghotel.com", "password", true)
	ownerUser := fixtures.AddUser(db.Store, "us", "er", "owner@ghotel.com", "password", false)
	randomUser := fixtures.AddUser(db.Store, "ran", "dom", "random@ghotel.com", "password", false)
	hotel := fixtures.AddHotel(db.Store, "hotel", "somewhere", 5)

	room := fixtures.AddRoom(
		db.Store,
		hotel.ID,
		&types.Room{
			Size:    "normal",
			Price:   99,
			Seaside: true,
		},
	)

	booking := fixtures.AddBooking(
		db.Store,
		ownerUser.ID,
		room.ID,
		time.Now(),
		time.Now().Add(time.Hour*24),
	)

	// -- Testing Owner call
	res := testRequest[*types.Booking](
		app,
		"GET",
		"/"+booking.ID,
		ownerUser.CreateJWT(),
		nil,
		func(r io.ReadCloser) *types.Booking {
			var booking *types.Booking
			json.NewDecoder(r).Decode(&booking)
			return booking
		},
		defaultStatusHandler,
	)

	// setting dates because formatting and parsing doesn't return exactly the same dates
	res.FromDate = booking.FromDate
	res.UntilDate = booking.UntilDate
	if !reflect.DeepEqual(res, booking) {
		t.Fatal("expected bookings to be the same")
	}

	// -- Testing non-Owner call
	_ = testRequest[*types.Booking](
		app,
		"GET",
		"/"+booking.ID,
		randomUser.CreateJWT(),
		nil,
		func(r io.ReadCloser) *types.Booking {
			var booking *types.Booking
			json.NewDecoder(r).Decode(&booking)
			return booking
		},
		func(code int, err api_error.ErrorResponse) {
			if err != api_error.UnauthorizedErrorResponse {
				t.Fatal("expected unauthorized")
			}
		},
	)

	// setting dates because formatting and parsing doesn't return exactly the same dates
	res.FromDate = booking.FromDate
	res.UntilDate = booking.UntilDate
	if !reflect.DeepEqual(res, booking) {
		t.Fatal("expected bookings to be the same")
	}
}

func TestAdminGetBookings(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	bookingHandler := NewBookingHandler(db.Store)
	jwt := middleware.NewJWTMiddleware(db.Store)
	app := newApp()
	adminApp := app.Group("/", jwt.JWTAuthentication, middleware.AdminAuthentication)
	adminApp.Get("/", bookingHandler.HandleGetBookings)

	admin := fixtures.AddUser(db.Store, "ad", "min", "admin@ghotel.com", "password", true)
	user := fixtures.AddUser(db.Store, "us", "er", "user@ghotel.com", "password", false)
	hotel := fixtures.AddHotel(db.Store, "hotel", "somewhere", 5)

	room := fixtures.AddRoom(
		db.Store,
		hotel.ID,
		&types.Room{
			Size:    "normal",
			Price:   99,
			Seaside: true,
		},
	)

	booking := fixtures.AddBooking(
		db.Store,
		admin.ID,
		room.ID,
		time.Now(),
		time.Now().Add(time.Hour*24),
	)

	// -- Testing Admin call
	bookings := testRequest[[]*types.Booking](
		app,
		"GET",
		"/",
		admin.CreateJWT(),
		nil,
		func(r io.ReadCloser) []*types.Booking {
			var bookings []*types.Booking
			json.NewDecoder(r).Decode(&bookings)
			return bookings
		},
		func(code int, err api_error.ErrorResponse) {
			if code != http.StatusOK {
				t.Fatal(err.Error)
			}
		},
	)

	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking but got %d", len(bookings))
	}

	// setting dates because formatting and parsing doesn't return exactly the same dates
	bookings[0].FromDate = booking.FromDate
	bookings[0].UntilDate = booking.UntilDate
	if !reflect.DeepEqual(booking, bookings[0]) {
		t.Fatal("expected bookings to be equal")
	}

	// -- Testing non-Admin call
	_ = testRequest[[]*types.Booking](
		app,
		"GET",
		"/",
		user.CreateJWT(),
		nil,
		func(r io.ReadCloser) []*types.Booking {
			var bookings []*types.Booking
			json.NewDecoder(r).Decode(&bookings)
			return bookings
		},
		func(code int, err api_error.ErrorResponse) {
			if err != api_error.UnauthorizedErrorResponse {
				t.Fatal("expected unauthorized")
			}
		},
	)
}
