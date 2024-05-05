package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/hhanri/ghotel/db"
	"github.com/hhanri/ghotel/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const testDbUri string = "mongodb://localhost:27017"

var userParams = types.CreateUserParams{
	FirstName: "Foo",
	LastName:  "Bar",
	Email:     "some@email.com",
	Password:  "SomeRandomPassword",
}

type testDB struct {
	db.UserStore
}

func (db *testDB) teardown(t *testing.T) {
	if err := db.UserStore.Drop(context.Background()); err != nil {
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
		UserStore: db.NewMongoUserStore(client, db.TestDBNAME),
	}
}

func newApp() *fiber.App {
	return fiber.New()
}

func (h *UserHandler) seedUser(app *fiber.App, params *types.CreateUserParams) types.User {
	app.Post("/", h.HandlePostUser)
	b, _ := json.Marshal(params)
	request := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	request.Header.Add("Content-Type", "application/json")
	resp, _ := app.Test(request)

	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)
	return user
}

func testRequest[T any](app *fiber.App, method string, path string, body io.Reader, transform func(io.ReadCloser) T) T {
	request := httptest.NewRequest(method, path, body)
	request.Header.Add("Content-Type", "application/json")
	resp, _ := app.Test(request)
	return transform(resp.Body)
}

func TestPostUser(t *testing.T) {
	testDB := setup(t)
	defer testDB.teardown(t)

	userHandler := NewUserHandler(testDB)

	app := newApp()
	app.Post("/", userHandler.HandlePostUser)

	params := userParams
	b, _ := json.Marshal(params)

	user := testRequest[types.User](
		app,
		"POST",
		"/",
		bytes.NewReader(b),
		func(r io.ReadCloser) types.User {
			var user types.User
			json.NewDecoder(r).Decode(&user)
			return user
		},
	)

	if user.FirstName != params.FirstName {
		t.Errorf("expected first name (%s) but got %s", params.FirstName, user.FirstName)
	}

	if user.LastName != params.LastName {
		t.Errorf("expected last name (%s) but got %s", params.LastName, user.LastName)
	}

	if user.Email != params.Email {
		t.Errorf("expected email (%s) but got %s", params.Email, user.Email)
	}

	if user.ID == "" {
		t.Errorf("expected user id to bee set")
	}

	if user.EncryptedPassword != "" {
		t.Errorf("expected password not to be returned")
	}
}

func TestGetUserByID(t *testing.T) {
	testDB := setup(t)
	defer testDB.teardown(t)

	userHandler := NewUserHandler(testDB)

	app := newApp()
	app.Get("/:id", userHandler.HandleGetUser)

	dbUser := userHandler.seedUser(app, &userParams)

	user := testRequest[types.User](
		app,
		"GET",
		"/"+dbUser.ID,
		nil,
		func(r io.ReadCloser) types.User {
			var user types.User
			json.NewDecoder(r).Decode(&user)
			return user
		},
	)

	if user != dbUser {
		t.Errorf("Got %+v instead of %+v", user, dbUser)
	}
}
