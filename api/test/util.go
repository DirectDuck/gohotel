package apiTest

import (
	"bytes"
	"context"
	"encoding/json"
	"hotel/db"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/gofiber/fiber/v2"
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

func setupStore() *db.Store {
	return db.GetTestDatabase().Store
}

func teardown() {
	err := db.GetTestDatabase().Drop(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}
