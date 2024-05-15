package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/hhanri/ghotel/types"
)

func insertTestUser(testDB *testDB, t *testing.T) *types.User {
	params := types.CreateUserParams{
		FirstName: "first",
		LastName:  "last",
		Email:     "some@email.com",
		Password:  "password",
	}
	user, _ := types.NewUserFromParams(params)
	insertedUser, err := testDB.User.InsertUser(context.Background(), user)
	if err != nil {
		t.Fatal(err)
	}
	return insertedUser
}

func TestAutheticateFailureWrongPassword(t *testing.T) {
	testDB := setup(t)
	defer testDB.teardown(t)

	authHandler := NewAuthHandler(testDB.Store)

	app := newApp()
	app.Post("/auth", authHandler.HandleAuthenticate)
	insertTestUser(testDB, t)

	authParams := AuthParams{
		Email:    "some@email.com",
		Password: "someWrongPassword",
	}
	b, _ := json.Marshal(authParams)

	_ = testRequest[AuthResponse](
		app,
		"POST",
		"/auth",
		bytes.NewReader(b),
		func(r io.ReadCloser) AuthResponse {
			b, _ := io.ReadAll(r)
			var resp AuthResponse
			_ = json.Unmarshal(b, &resp)
			return resp
		},
		func(status int, body string) {
			if status == http.StatusOK {
				t.Fatal("expected authentication failure")
			}
			if body != "invalid credentials" {
				t.Fatal("expected to return invalid credentials")
			}
		},
	)

}

func TestAutheticateFailureNonExistingUser(t *testing.T) {
	testDB := setup(t)
	defer testDB.teardown(t)

	authHandler := NewAuthHandler(testDB.Store)

	app := newApp()
	app.Post("/auth", authHandler.HandleAuthenticate)
	insertTestUser(testDB, t)

	authParams := AuthParams{
		Email:    "some123@email.com",
		Password: "password",
	}
	b, _ := json.Marshal(authParams)

	_ = testRequest[AuthResponse](
		app,
		"POST",
		"/auth",
		bytes.NewReader(b),
		func(r io.ReadCloser) AuthResponse {
			b, _ := io.ReadAll(r)
			var resp AuthResponse
			_ = json.Unmarshal(b, &resp)
			return resp
		},
		func(status int, body string) {

			fmt.Println(body)
			if status == http.StatusOK {
				t.Fatal("expected authentication failure")
			}
			if body != "not found" {
				t.Fatal("expected to return not found/")
			}
		},
	)

}

func TestAutheticateSuccess(t *testing.T) {
	testDB := setup(t)
	defer testDB.teardown(t)

	authHandler := NewAuthHandler(testDB.Store)

	app := newApp()
	app.Post("/auth", authHandler.HandleAuthenticate)
	user := insertTestUser(testDB, t)

	authParams := AuthParams{
		Email:    user.Email,
		Password: "password",
	}
	b, _ := json.Marshal(authParams)

	authResp := testRequest[AuthResponse](
		app,
		"POST",
		"/auth",
		bytes.NewReader(b),
		func(r io.ReadCloser) AuthResponse {
			b, _ := io.ReadAll(r)
			var authResp AuthResponse
			_ = json.Unmarshal(b, &authResp)
			return authResp
		},
		func(status int, _ string) {
			if status != http.StatusOK {
				t.Fatalf("expected authentication success but got %d", status)
			}
		},
	)

	if authResp.Token == "" {
		t.Fatal("Expected token, got empty string")
	}

	user.EncryptedPassword = "" // json doesn't return EncryptedPassword
	if !reflect.DeepEqual(user, authResp.User) {
		t.Fatalf("Expected %+v but got %+v", user, authResp.User)
	}

}
