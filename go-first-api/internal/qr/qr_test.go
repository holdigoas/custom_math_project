package qr

import (
	"math"
	"testing"
)

const tol = 1e-9

func TestDecompose(t *testing.T) {
	cases := []struct {
		name string
		a    [][]float64
	}{
		{"1x1", [][]float64{{5}}},
		{"cuadrada 2x2", [][]float64{{1, 2}, {3, 4}}},
		{"cuadrada 3x3", [][]float64{{12, -51, 4}, {6, 167, -68}, {-4, 24, -41}}},
		{"alta 3x2", [][]float64{{1, 2}, {3, 4}, {5, 6}}},
		{"identidad 3x3", [][]float64{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}}},
		{"columna cero", [][]float64{{0, 1}, {0, 2}}},
		{"rango deficiente", [][]float64{{1, 2}, {2, 4}, {3, 6}}},
		{"negativos y decimales", [][]float64{{-1.5, 2.25}, {0.5, -3}, {4, 1}}},
		{"todo ceros", [][]float64{{0, 0}, {0, 0}}},
		{"ancha 2x3", [][]float64{{1, 2, 3}, {4, 5, 6}}},
		{"ancha 2x5", [][]float64{{1, 2, 3, 4, 5}, {6, 7, 8, 9, 10}}},
		{"ancha 1x3", [][]float64{{2, 4, 6}}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			q, r, err := Decompose(tc.a)
			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}

			m, n := len(tc.a), len(tc.a[0])

			if len(q) != m || len(q[0]) != m {
				t.Fatalf("Q debería ser %dx%d", m, m)
			}
			if len(r) != m || len(r[0]) != n {
				t.Fatalf("R debería ser %dx%d", m, n)
			}

			assertClose(t, "Q*R == A", multiply(q, r), tc.a)
			assertClose(t, "QT*Q == I", multiply(transpose(q), q), identity(m))
			assertUpperTriangular(t, r)
		})
	}
}

func TestDecompose_Errors(t *testing.T) {
	cases := []struct {
		name string
		a    [][]float64
		want error
	}{
		{"nil", nil, ErrEmptyMatrix},
		{"sin filas", [][]float64{}, ErrEmptyMatrix},
		{"fila vacía", [][]float64{{}}, ErrEmptyMatrix},
		{"irregular", [][]float64{{1, 2}, {3}}, ErrIrregularMatrix},
		{"NaN", [][]float64{{1, math.NaN()}}, ErrInvalidValue},
		{"Inf", [][]float64{{math.Inf(1), 2}}, ErrInvalidValue},
		//{"ancha 2x3", [][]float64{{1, 2, 3}, {4, 5, 6}}, ErrWideMatrix},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := Decompose(tc.a)
			if err != tc.want {
				t.Fatalf("se esperaba %v, se obtuvo %v", tc.want, err)
			}
		})
	}
}

func TestDecompose_NoModificaLaEntrada(t *testing.T) {
	a := [][]float64{{1, 2}, {3, 4}}
	original := copyMatrix(a)

	if _, _, err := Decompose(a); err != nil {
		t.Fatal(err)
	}
	assertClose(t, "entrada intacta", a, original)
}

// Funciones auxiliares de los tests

func multiply(a, b [][]float64) [][]float64 {
	p, inner, s := len(a), len(b), len(b[0])
	res := make([][]float64, p)
	for i := 0; i < p; i++ {
		res[i] = make([]float64, s)
		for j := 0; j < s; j++ {
			for k := 0; k < inner; k++ {
				res[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	return res
}

func transpose(a [][]float64) [][]float64 {
	t := make([][]float64, len(a[0]))
	for j := range t {
		t[j] = make([]float64, len(a))
		for i := range a {
			t[j][i] = a[i][j]
		}
	}
	return t
}

func assertClose(t *testing.T, label string, got, want [][]float64) {
	t.Helper()
	for i := range want {
		for j := range want[i] {
			if math.Abs(got[i][j]-want[i][j]) > tol {
				t.Fatalf("%s: en [%d][%d] se obtuvo %v y se esperaba %v", label, i, j, got[i][j], want[i][j])
			}
		}
	}
}

func assertUpperTriangular(t *testing.T, r [][]float64) {
	t.Helper()
	for i := range r {
		for j := 0; j < i && j < len(r[i]); j++ {
			if math.Abs(r[i][j]) > tol {
				t.Fatalf("R no es triangular superior: [%d][%d] = %v", i, j, r[i][j])
			}
		}
	}
}

func identity(n int) [][]float64 {
	m := make([][]float64, n)
	for i := range m {
		m[i] = make([]float64, n)
		m[i][i] = 1
	}
	return m
}

func copyMatrix(a [][]float64) [][]float64 {
	c := make([][]float64, len(a))
	for i := range a {
		c[i] = make([]float64, len(a[i]))
		copy(c[i], a[i])
	}
	return c
}
