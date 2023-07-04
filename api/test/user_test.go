package apiTest

import (
	"context"
	"encoding/json"
	"hotel/api"
	"hotel/controllers"
	"hotel/types"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestCreateUser(t *testing.T) {
	store := setupCTStore()
	defer teardown()

	app := fiber.New()
	userHandler := api.NewUserHandler(
		&controllers.UserController{Store: store},
	)
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
		t.Fatal(err)
	}

	var user *types.User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		t.Fatal(err)
	}

	if user.Email != params.Email {
		t.Fatalf("Expected email %s but got %s", params.Email, user.Email)
	}
	if len(user.ID) == 0 {
		t.Fatalf("Received empty user ID")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Fatalf("API returned password when it shouldn't")
	}

	actualUser, err := store.CT.Users.GetByID(context.Background(), user.ID)
	if err != nil {
		t.Fatal(err)
	}
	if actualUser == nil {
		t.Fatalf("Created user %+v not found", user)
	}

	// Since api won't return EncryptedPassword we need to conceal it
	actualUser.EncryptedPassword = ""
	if !reflect.DeepEqual(user, actualUser) {
		t.Fatalf("User's aren't equal")
	}
}

func TestLoginUser(t *testing.T) {
	store := setupCTStore()
	defer teardown()

	userEmail := "helloworld@gmail.com"
	userPassword := "12345678"

	userUnsaved, err := types.NewUserFromCreateParams(
		types.CreateUserParams{
			BaseUserParams: types.BaseUserParams{
				Email:     userEmail,
				FirstName: "Hello",
				LastName:  "World",
			},
			Password: userPassword,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	userController := &controllers.UserController{Store: store}
	_, err = userController.Create(context.Background(), userUnsaved)
	if err != nil {
		t.Fatal(err)
	}

	app := fiber.New()
	userHandler := api.NewUserHandler(userController)
	app.Post("/", userHandler.HandleLogin)

	// User not found
	resp, err := sendStructJSONRequest(
		app, "POST", "/",
		types.LoginUserParams{
			Email:    "somethingwrong",
			Password: "123",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("Incorrect login status")
	}

	// Incorrect password
	resp, err = sendStructJSONRequest(
		app, "POST", "/",
		types.LoginUserParams{
			Email:    userEmail,
			Password: "totallyincorrectpassword",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("Incorrect login status")
	}

	// Everything good
	resp, err = sendStructJSONRequest(
		app, "POST", "/",
		types.LoginUserParams{
			Email:    userEmail,
			Password: userPassword,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("Incorrect login status")
	}
	respMap := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&respMap)
	if err != nil {
		t.Fatal(err)
	}
	if len(respMap["token"].(string)) == 0 {
		t.Fatalf("Invalid token length")
	}

	if respMap["user"].(map[string]interface{})["email"].(string) != userEmail {
		t.Fatalf("Invalid user email")
	}
}
