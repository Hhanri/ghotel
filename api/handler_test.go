package api

import (
	"context"
	"io"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/hhanri/ghotel/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testDB struct {
	*db.Store
}

var testDbUri string = "mongodb://ghotel:secret@localhost:27017/"

func (db *testDB) teardown(t *testing.T) {
	if err := db.User.Drop(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testDB {

	client, err := mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(testDbUri),
	)
	if err != nil {
		log.Fatal(err)
	}
	return &testDB{
		Store: &db.Store{
			User: db.NewMongoUserStore(client, db.TestDBNAME),
		},
	}
}

func newApp() *fiber.App {
	return fiber.New()
}

func testRequest[T any](app *fiber.App, method string, path string, body io.Reader, transform func(io.ReadCloser) T) T {
	request := httptest.NewRequest(method, path, body)
	request.Header.Add("Content-Type", "application/json")
	resp, _ := app.Test(request)
	return transform(resp.Body)
}
