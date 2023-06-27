package apiTest

import (
	"context"
	"encoding/json"
	"hotel/api"
	"hotel/db"
	"hotel/types"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func setupUserStore() *db.MongoUserStore {
	return db.NewMongoUserStore(db.GetTestDatabase())
}

func TestCreateUser(t *testing.T) {
	userStore := setupUserStore()
	defer teardown()

	app := fiber.New()
	userHandler := api.NewUserHandler(userStore)
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

	actualUser, err := userStore.GetUserByID(context.TODO(), user.ID)
	if err != nil {
		t.Error(err)
	}
	if actualUser == nil {
		t.Errorf("Created user %+v not found", user)
	}
}
