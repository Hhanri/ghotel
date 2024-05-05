package api

import (
	"bytes"
	"context"
	"encoding/json"
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

func TestPostUser(t *testing.T) {
	testDB := setup(t)
	defer testDB.teardown(t)

	userHandler := NewUserHandler(testDB)

	app := fiber.New()
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		FirstName: "Foo",
		LastName:  "Bar",
		Email:     "some@email.com",
		Password:  "SomeRandomPassword",
	}
	b, _ := json.Marshal(params)
	request := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	request.Header.Add("Content-Type", "application/json")
	resp, _ := app.Test(request)

	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)

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

	if user.EncryptedPassword == "" {
		t.Errorf("expected password not to be returned")
	}
}
