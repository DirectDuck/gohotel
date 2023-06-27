package apiTest

import (
	"context"
	"encoding/json"
	"hotel/api"
	"hotel/db"
	"hotel/types"
	"log"
	"testing"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const testdbUri = "mongodb://admin:admin@localhost:27017"
const testdbName = "hotel-reservation-test"

type testUserStore struct {
	*db.MongoUserStore
}

func (userStore *testUserStore) teardown() {
	err := userStore.Drop(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}

func setup() *testUserStore {
	dbClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testdbUri))
	if err != nil {
		log.Fatal(err)
	}
	return &testUserStore{
		MongoUserStore: db.NewMongoUserStore(dbClient.Database(testdbName)),
	}
}

func TestCreateUser(t *testing.T) {
	userStore := setup()
	defer userStore.teardown()

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
