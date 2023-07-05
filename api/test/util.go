package apiTest

import (
	"bytes"
	"context"
	"encoding/json"
	"hotel/controllers"
	"hotel/db"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func sendStructJSONRequest[T any](
	app *fiber.App, method string, path string, params T,
) (*http.Response, error) {

	b, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	return app.Test(req)
}

func setupDBStore() *db.DB {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("No .env file found")
	}
	return db.GetTestDatabase()
}

func setupCTStore() *controllers.Store {
	return controllers.NewStore(setupDBStore())
}

func teardown() {
	err := db.GetTestDatabase().Drop(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
