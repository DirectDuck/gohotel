package apiTest

import (
	"context"
	"encoding/json"
	"hotel/api"
	"hotel/types"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestCreateUser(t *testing.T) {
	store := setupStore()
	defer teardown()

	app := fiber.New()
	userHandler := api.NewUserHandler(store)
	app.Post("/", userHandler.HandleCreateUser)

	params := types.CreateUserParams{
		BaseUserParams: types.BaseUserParams{
			Email:     "hello@mail.ru",
			FirstName: "Alex",
			LastName:  "Xela",
		},
		Password: "12312321421421",
	}

	resp, err := sendStructJSONRequest(app, "POST", "/", params)
	if err != nil {
		t.Error(err)
	}

	var user *types.User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		t.Error(err)
	}

	if user.Email != params.Email {
		t.Errorf("Expected email %s but got %s", params.Email, user.Email)
	}
	if len(user.ID) == 0 {
		t.Errorf("Received empty user ID")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Errorf("API returned password when it shouldn't")
	}

	actualUser, err := store.Users.GetByID(context.TODO(), user.ID)
	if err != nil {
		t.Error(err)
	}
	if actualUser == nil {
		t.Errorf("Created user %+v not found", user)
	}
}
