package handler

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/gofiber/fiber/v3"

	"go-first-api/internal/qr"

	"go-first-api/internal/client"
)

// dimencion maxima
const maxDimension = 200

// tiempo de tolerancia de espera de la API node
const upstreamTimeout = 20 * time.Second

// estructura del body
type QRRequest struct {
	Matrix [][]float64 `json:"matrix"`
}

// estructura error de la API go
type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// estructura de la respuesta de error API go
type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

// estructura de respuesta de la API go
type QRResponse struct {
	Q          [][]float64        `json:"q"`
	R          [][]float64        `json:"r"`
	Statistics *client.Statistics `json:"Statistics"`
}

// estructura respuesta API node
type QRHandler struct {
	statistics client.StatisticsClient
}

// estructura cleinte node
func NewQRHandler(statistics client.StatisticsClient) *QRHandler {
	return &QRHandler{statistics: statistics}
}

// funcion principal del archivo handler
func (h *QRHandler) Decompose(c fiber.Ctx) error {
	var req QRRequest

	// validando el formato del body
	if err := c.Bind().Body(&req); err != nil {
		return sendError(c, fiber.StatusBadRequest, "INVALID_JSON",
			"el cuerpo debe ser un JSON válido con la forma {\"matrix\": [[...]]}")
	}

	// validando dimenciones
	if len(req.Matrix) > maxDimension || (len(req.Matrix) > 0 && len(req.Matrix[0]) > maxDimension) {
		return sendError(c, fiber.StatusBadRequest, "MATRIX_TOO_LARGE",
			"la matriz supera el tamaño máximo permitido de 200x200")
	}

	// Factorizacion QR
	q, r, err := qr.Decompose(req.Matrix)
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, errorCode(err), err.Error())
	}

	// preparando la consulta al cliente
	ctx, cancel := context.WithTimeout(context.Background(), upstreamTimeout)
	defer cancel()

	// obteniendo las estadisticas de node
	statistics, err := h.statistics.Compute(ctx, q, r)
	if err != nil {
		log.Printf("error al llamar a la API de Node: %v", err)
		return sendUpstreamError(c, err)
	}

	return c.JSON(QRResponse{Q: q, R: r, Statistics: statistics})
}

// customizando el mensaje de error API go
func errorCode(err error) string {
	switch {
	case errors.Is(err, qr.ErrEmptyMatrix):
		return "EMPTY_MATRIX"
	case errors.Is(err, qr.ErrIrregularMatrix):
		return "IRREGULAR_MATRIX"
	case errors.Is(err, qr.ErrInvalidValue):
		return "INVALID_VALUE"
	case errors.Is(err, qr.ErrValueTooLarge):
		return "VALUE_TOO_LARGE"
	default:
		return "INVALID_MATRIX"
	}
}

// customizando mensaje error de la API Node
func sendUpstreamError(c fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, client.ErrTimeout):
		return sendError(c, fiber.StatusGatewayTimeout, "STATISTICS_TIMEOUT",
			"el servicio de estadísticas no respondió a tiempo")
	case errors.Is(err, client.ErrUnavailable):
		return sendError(c, fiber.StatusBadGateway, "STATISTICS_UNAVAILABLE",
			"el servicio de estadísticas no está disponible")
	default:
		return sendError(c, fiber.StatusBadGateway, "STATISTICS_ERROR",
			"el servicio de estadísticas devolvió un error")
	}
}

func sendError(c fiber.Ctx, status int, code, message string) error {
	return c.Status(status).JSON(ErrorResponse{Error: ErrorBody{Code: code, Message: message}})
}
