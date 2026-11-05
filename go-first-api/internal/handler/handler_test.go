package handler

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func newTestApp() *fiber.App {
	app := fiber.New()
	app.Post("/api/v1/qr", NewQRHandler().Decompose)
	return app
}

func post(t *testing.T, body string) (int, []byte) {
	t.Helper()
	req := httptest.NewRequest("POST", "/api/v1/qr", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := newTestApp().Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	buf := make([]byte, 1<<20)
	n, _ := resp.Body.Read(buf)
	return resp.StatusCode, buf[:n]
}

func TestDecompose_OK(t *testing.T) {
	status, body := post(t, `{"matrix": [[1,2],[3,4],[5,6]]}`)
	if status != 200 {
		t.Fatalf("status %d, cuerpo: %s", status, body)
	}

	var res QRResponse
	if err := json.Unmarshal(body, &res); err != nil {
		t.Fatal(err)
	}
	if len(res.Q) != 3 || len(res.Q[0]) != 3 || len(res.R) != 3 || len(res.R[0]) != 2 {
		t.Fatalf("dimensiones inesperadas de Q o R")
	}
}

func TestDecompose_Errores(t *testing.T) {
	cases := []struct {
		name string
		body string
		code string
	}{
		{"json inválido", `{"matrix": `, "INVALID_JSON"},
		{"matriz vacía", `{"matrix": []}`, "EMPTY_MATRIX"},
		{"irregular", `{"matrix": [[1,2],[3]]}`, "IRREGULAR_MATRIX"},
		//{"ancha", `{"matrix": [[1,2,3],[4,5,6]]}`, "WIDE_MATRIX"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := post(t, tc.body)
			if status != 400 {
				t.Fatalf("se esperaba 400, se obtuvo %d", status)
			}
			var res ErrorResponse
			if err := json.Unmarshal(body, &res); err != nil {
				t.Fatal(err)
			}
			if res.Error.Code != tc.code {
				t.Fatalf("se esperaba %s, se obtuvo %s", tc.code, res.Error.Code)
			}
		})
	}
}

func TestDecompose_MatrizAncha(t *testing.T) {
	status, body := post(t, `{"matrix": [[1,2,3],[4,5,6]]}`)
	if status != 200 {
		t.Fatalf("status %d, cuerpo: %s", status, body)
	}

	var res QRResponse
	if err := json.Unmarshal(body, &res); err != nil {
		t.Fatal(err)
	}
	if len(res.Q) != 2 || len(res.Q[0]) != 2 || len(res.R) != 2 || len(res.R[0]) != 3 {
		t.Fatalf("Q debería ser 2x2 y R debería ser 2x3")
	}
}
