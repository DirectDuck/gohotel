package apiTest

import (
	"bytes"
	"encoding/json"
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
