package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"

	"go-first-api/internal/client"
)

type fakeStatistics struct {
	statistics *client.Statistics
	err        error
}

func (f fakeStatistics) Compute(_ context.Context, _, _ [][]float64) (*client.Statistics, error) {
	return f.statistics, f.err
}

var okStatistics = fakeStatistics{statistics: &client.Statistics{Max: 4, Min: 0, Average: 1.5, Sum: 12}}

func post(t *testing.T, statistics client.StatisticsClient, body string) (int, []byte) {
	t.Helper()

	app := fiber.New()
	app.Post("/api/v1/qr", NewQRHandler(statistics).Decompose)

	req := httptest.NewRequest("POST", "/api/v1/qr", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return resp.StatusCode, data
}

func TestDecompose_OK(t *testing.T) {
	status, body := post(t, okStatistics, `{"matrix": [[1,2],[3,4],[5,6]]}`)
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
	if res.Statistics == nil || res.Statistics.Max != 4 {
		t.Fatalf("faltan las estadísticas en la respuesta: %s", body)
	}
}

func TestDecompose_MatrizAncha(t *testing.T) {
	status, body := post(t, okStatistics, `{"matrix": [[1,2,3],[4,5,6]]}`)
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

func TestDecompose_Errores(t *testing.T) {
	cases := []struct {
		name string
		body string
		code string
	}{
		{"json inválido", `{"matrix": `, "INVALID_JSON"},
		{"matriz vacía", `{"matrix": []}`, "EMPTY_MATRIX"},
		{"irregular", `{"matrix": [[1,2],[3]]}`, "IRREGULAR_MATRIX"},
		{"valor enorme", `{"matrix": [[1e101, 1]]}`, "VALUE_TOO_LARGE"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := post(t, okStatistics, tc.body)
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

func TestDecompose_FallosDeNode(t *testing.T) {
	cases := []struct {
		name   string
		err    error
		status int
		code   string
	}{
		{"timeout", client.ErrTimeout, 504, "STATISTICS_TIMEOUT"},
		{"no disponible", client.ErrUnavailable, 502, "STATISTICS_UNAVAILABLE"},
		{"respuesta inválida", client.ErrUpstream, 502, "STATISTICS_ERROR"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := post(t, fakeStatistics{err: tc.err}, `{"matrix": [[1,2],[3,4]]}`)
			if status != tc.status {
				t.Fatalf("se esperaba %d, se obtuvo %d", tc.status, status)
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
