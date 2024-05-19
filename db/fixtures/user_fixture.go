package fixtures

import (
	"context"
	"log"

	"github.com/hhanri/ghotel/db"
	"github.com/hhanri/ghotel/types"
)

func AddUser(store *db.Store, firstName, lastName, email, password string, isAdmin bool) *types.User {
	params := types.CreateUserParams{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin

	_, err = store.User.InsertUser(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

	return user
}
