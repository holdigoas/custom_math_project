package handler

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"go-first-api/internal/qr"
)

const maxDimension = 200

type QRRequest struct {
	Matrix [][]float64 `json:"matrix"`
}

type QRResponse struct {
	Q [][]float64 `json:"q"`
	R [][]float64 `json:"r"`
}

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type QRHandler struct{}

func NewQRHandler() *QRHandler {
	return &QRHandler{}
}

func (h *QRHandler) Decompose(c fiber.Ctx) error {
	var req QRRequest
	if err := c.Bind().Body(&req); err != nil {
		return sendError(c, fiber.StatusBadRequest, "INVALID_JSON",
			"el cuerpo debe ser un JSON válido con la forma {\"matrix\": [[...]]}")
	}

	if len(req.Matrix) > maxDimension || (len(req.Matrix) > 0 && len(req.Matrix[0]) > maxDimension) {
		return sendError(c, fiber.StatusBadRequest, "MATRIX_TOO_LARGE",
			"la matriz supera el tamaño máximo permitido de 200x200")
	}

	q, r, err := qr.Decompose(req.Matrix)
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, errorCode(err), err.Error())
	}

	return c.JSON(QRResponse{Q: q, R: r})
}

func errorCode(err error) string {
	switch {
	case errors.Is(err, qr.ErrEmptyMatrix):
		return "EMPTY_MATRIX"
	case errors.Is(err, qr.ErrIrregularMatrix):
		return "IRREGULAR_MATRIX"
	case errors.Is(err, qr.ErrInvalidValue):
		return "INVALID_VALUE"
	case errors.Is(err, qr.ErrWideMatrix):
		return "WIDE_MATRIX"
	default:
		return "INVALID_MATRIX"
	}
}

func sendError(c fiber.Ctx, status int, code, message string) error {
	return c.Status(status).JSON(ErrorResponse{Error: ErrorBody{Code: code, Message: message}})
}
