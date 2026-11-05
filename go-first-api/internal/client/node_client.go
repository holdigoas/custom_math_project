package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

// errores contemplados
var (
	ErrTimeout     = errors.New("la API de estadísticas no respondió a tiempo")
	ErrUnavailable = errors.New("no se pudo conectar con la API de estadísticas")
	ErrUpstream    = errors.New("la API de estadísticas devolvió una respuesta inválida")
)

// matrices diagonales
type DiagonalInfo struct {
	Q   bool `json:"q"`
	R   bool `json:"r"`
	Any bool `json:"any"`
}

// estructura de la respuesta node
type Statistics struct {
	Max        float64      `json:"max"`
	Min        float64      `json:"min"`
	Average    float64      `json:"average"`
	Sum        float64      `json:"sum"`
	IsDiagonal DiagonalInfo `json:"isDiagonal"`
}

// estructura del body
type statisticsRequest struct {
	Matrices struct {
		Q [][]float64 `json:"q"`
		R [][]float64 `json:"r"`
	} `json:"matrices"`
}

type StatisticsClient interface {
	Compute(ctx context.Context, q, r [][]float64) (*Statistics, error)
}

// estructura de peticion
type NodeClient struct {
	baseURL    string
	http       *http.Client
	maxRetries int
	backoff    time.Duration
}

func NewNodeClient(baseURL string, timeout time.Duration, maxRetries int, backoff time.Duration) *NodeClient {
	return &NodeClient{
		baseURL:    baseURL,
		http:       &http.Client{Timeout: timeout},
		maxRetries: maxRetries,
		backoff:    backoff,
	}
}

func (c *NodeClient) Compute(ctx context.Context, q, r [][]float64) (*Statistics, error) {
	var body statisticsRequest
	body.Matrices.Q = q
	body.Matrices.R = r

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("no se pudo armar la petición: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			// Espera creciente entre reintentos (backoff lineal), respetando el contexto.
			select {
			case <-time.After(c.backoff * time.Duration(attempt)):
			case <-ctx.Done():
				return nil, classify(ctx.Err())
			}
		}

		statistics, retryable, err := c.doOnce(ctx, payload)
		if err == nil {
			return statistics, nil
		}
		lastErr = err
		if !retryable {
			break
		}
	}
	return nil, lastErr
}

func (c *NodeClient) doOnce(ctx context.Context, payload []byte) (*Statistics, bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/v1/statistics", bytes.NewReader(payload))
	if err != nil {
		return nil, false, fmt.Errorf("%w: %v", ErrUnavailable, err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, true, classify(err) // fallo de red: reintentable
	}
	defer resp.Body.Close()

	switch {
	case resp.StatusCode >= 500:
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil, true, fmt.Errorf("%w: node respondió %d", ErrUpstream, resp.StatusCode)
	case resp.StatusCode != http.StatusOK:
		// 4xx: repetir la misma petición daría el mismo error, no se reintenta.
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil, false, fmt.Errorf("%w: node respondió %d", ErrUpstream, resp.StatusCode)
	}

	var statistics Statistics
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&statistics); err != nil {
		return nil, false, fmt.Errorf("%w: no se pudo leer la respuesta", ErrUpstream)
	}
	return &statistics, false, nil
}

func classify(err error) error {
	var netErr net.Error
	if errors.Is(err, context.DeadlineExceeded) || (errors.As(err, &netErr) && netErr.Timeout()) {
		return ErrTimeout
	}
	return fmt.Errorf("%w: %v", ErrUnavailable, err)
}
