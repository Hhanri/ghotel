package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/hhanri/ghotel/api/api_error"
	"github.com/hhanri/ghotel/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testDB struct {
	client *mongo.Client
	*db.Store
}

var testDbUri string = "mongodb://ghotel:secret@localhost:27017/"

func (tdb *testDB) teardown(t *testing.T) {
	if err := tdb.client.Database(db.TestDBNAME).Drop(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testDB {
	client, err := mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(testDbUri),
	)
	if err != nil {
		t.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client, db.TestDBNAME)
	return &testDB{
		client: client,
		Store: &db.Store{
			User:    db.NewMongoUserStore(client, db.TestDBNAME),
			Booking: db.NewMongoBookingStore(client, db.TestDBNAME),
			Hotel:   hotelStore,
			Room:    db.NewMongoRoomStore(client, db.TestDBNAME, hotelStore),
		},
	}
}

func newApp() *fiber.App {
	return fiber.New(fiber.Config{ErrorHandler: api_error.FiberErrorHandler})
}

func defaultStatusHandler(status int, err api_error.ErrorResponse) {}

func testRequest[T any](
	app *fiber.App,
	method string,
	path string,
	jwt string,
	body io.Reader,
	transform func(io.ReadCloser) T,
	handleStatus func(int, api_error.ErrorResponse),
) T {
	request := httptest.NewRequest(method, path, body)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-Api-Token", jwt)
	resp, _ := app.Test(request)

	if resp.StatusCode != http.StatusOK {
		var errorResp api_error.ErrorResponse
		_ = json.NewDecoder(resp.Body).Decode(&errorResp)
		handleStatus(resp.StatusCode, errorResp)
	}

	return transform(resp.Body)
}
