package client

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

const okBody = `{"max":4,"min":0,"average":1.375,"sum":11,"isDiagonal":{"q":true,"r":false,"any":true}}`

var (
	q = [][]float64{{1, 0}, {0, 1}}
	r = [][]float64{{2, 3}, {0, 4}}
)

func TestCompute_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(okBody))
	}))
	defer srv.Close()

	c := NewNodeClient(srv.URL, time.Second, 0, time.Millisecond)
	statistics, err := c.Compute(context.Background(), q, r)
	if err != nil {
		t.Fatal(err)
	}
	if statistics.Max != 4 || statistics.Sum != 11 || !statistics.IsDiagonal.Any {
		t.Fatalf("estadísticas inesperadas: %+v", statistics)
	}
}

func TestCompute_ReintentaAnteError5xx(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if atomic.AddInt32(&calls, 1) == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = w.Write([]byte(okBody))
	}))
	defer srv.Close()

	c := NewNodeClient(srv.URL, time.Second, 2, time.Millisecond)
	if _, err := c.Compute(context.Background(), q, r); err != nil {
		t.Fatalf("debió recuperarse con el reintento: %v", err)
	}
	if got := atomic.LoadInt32(&calls); got != 2 {
		t.Fatalf("se esperaban 2 llamadas, hubo %d", got)
	}
}

func TestCompute_NoReintentaAnteError4xx(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer srv.Close()

	c := NewNodeClient(srv.URL, time.Second, 2, time.Millisecond)
	_, err := c.Compute(context.Background(), q, r)
	if !errors.Is(err, ErrUpstream) {
		t.Fatalf("se esperaba ErrUpstream, se obtuvo %v", err)
	}
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("un 4xx no debe reintentarse; hubo %d llamadas", got)
	}
}

func TestCompute_ServidorCaido(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	url := srv.URL
	srv.Close() // simulacion de node no disponible (apagando)

	c := NewNodeClient(url, time.Second, 1, time.Millisecond)
	_, err := c.Compute(context.Background(), q, r)
	if !errors.Is(err, ErrUnavailable) {
		t.Fatalf("se esperaba ErrUnavailable, se obtuvo %v", err)
	}
}

func TestCompute_Timeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(300 * time.Millisecond)
		_, _ = w.Write([]byte(okBody))
	}))
	defer srv.Close()

	c := NewNodeClient(srv.URL, 50*time.Millisecond, 0, time.Millisecond)
	_, err := c.Compute(context.Background(), q, r)
	if !errors.Is(err, ErrTimeout) {
		t.Fatalf("se esperaba ErrTimeout, se obtuvo %v", err)
	}
}
