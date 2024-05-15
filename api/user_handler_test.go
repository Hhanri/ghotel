package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/hhanri/ghotel/types"
)

var userParams = types.CreateUserParams{
	FirstName: "Foo",
	LastName:  "Bar",
	Email:     "some@email.com",
	Password:  "SomeRandomPassword",
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

func (h *UserHandler) getUser(app *fiber.App, id string) types.User {
	app.Get("/:id", h.HandleGetUser)
	user := testRequest[types.User](
		app,
		"GET",
		"/"+id,
		nil,
		func(r io.ReadCloser) types.User {
			var user types.User
			json.NewDecoder(r).Decode(&user)
			return user
		},
		defaultStatusHandler,
	)
	return user
}

func TestPostUser(t *testing.T) {
	testDB := setup(t)
	defer testDB.teardown(t)

	userHandler := NewUserHandler(testDB.Store)

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
		defaultStatusHandler,
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

	userHandler := NewUserHandler(testDB.Store)

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
		defaultStatusHandler,
	)

	if user != dbUser {
		t.Errorf("Got %+v instead of %+v", user, dbUser)
	}
}

func TestGetUsers(t *testing.T) {
	testDB := setup(t)
	defer testDB.teardown(t)

	userHandler := NewUserHandler(testDB.Store)

	app := newApp()
	app.Get("/", userHandler.HandleGetUsers)

	dbUser1 := userHandler.seedUser(app, &userParams)
	dbUser2 := userHandler.seedUser(app, &userParams)

	users := testRequest[[]types.User](
		app,
		"GET",
		"/",
		nil,
		func(r io.ReadCloser) []types.User {
			var users []types.User
			json.NewDecoder(r).Decode(&users)
			return users
		},
		defaultStatusHandler,
	)

	if users[0] != dbUser1 && users[1] != dbUser2 {
		t.Errorf("expected %+v but got %+v", []types.User{dbUser1, dbUser2}, users)
	}
}

func TestUpdateUser(t *testing.T) {
	testDB := setup(t)
	defer testDB.teardown(t)

	userHandler := NewUserHandler(testDB.Store)

	app := newApp()
	app.Put("/:id", userHandler.HandleUpdateUser)

	dbUser := userHandler.seedUser(app, &userParams)

	update := types.UpdateUserParams{
		FirstName: "FooFoo",
		LastName:  "BarBar",
	}
	b, _ := json.Marshal(update)

	id := testRequest[string](
		app,
		"PUT",
		"/"+dbUser.ID,
		bytes.NewReader(b),
		func(r io.ReadCloser) string {
			b, _ := io.ReadAll(r)
			m := make(map[string]string)
			_ = json.Unmarshal(b, &m)
			return m["data"]
		},
		defaultStatusHandler,
	)

	if id != dbUser.ID {
		t.Errorf("Expected ID %s but got %s", dbUser.ID, id)
	}

	updatedUser := userHandler.getUser(app, id)

	if updatedUser.FirstName != update.FirstName {
		t.Errorf("expected first name (%s) but got %s", update.FirstName, updatedUser.FirstName)
	}

	if updatedUser.LastName != update.LastName {
		t.Errorf("expected last name (%s) but got %s", update.FirstName, updatedUser.FirstName)
	}
}

func TestDeleteUser(t *testing.T) {
	testDB := setup(t)
	defer testDB.teardown(t)

	userHandler := NewUserHandler(testDB.Store)

	app := newApp()
	app.Delete("/:id", userHandler.HandleDeleteUser)

	dbUser := userHandler.seedUser(app, &userParams)

	id := testRequest[string](
		app,
		"DELETE",
		"/"+dbUser.ID,
		nil,
		func(r io.ReadCloser) string {
			b, _ := io.ReadAll(r)
			m := make(map[string]string)
			_ = json.Unmarshal(b, &m)
			return m["data"]
		},
		defaultStatusHandler,
	)

	if id != dbUser.ID {
		t.Errorf("Expected ID %s but got %s", dbUser.ID, id)
	}

	deletedUser := userHandler.getUser(app, id)

	if deletedUser != (types.User{}) {
		t.Errorf("expected empty user but got %+v", deletedUser)
	}
}
