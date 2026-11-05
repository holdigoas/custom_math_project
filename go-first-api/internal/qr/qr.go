package qr

import (
	"errors"
	"math"

	"gonum.org/v1/gonum/mat"
)

var (
	ErrEmptyMatrix     = errors.New("la matriz no puede estar vacía")
	ErrIrregularMatrix = errors.New("todas las filas deben tener la misma cantidad de columnas")
	ErrInvalidValue    = errors.New("la matriz contiene valores no válidos (NaN o infinito)")
	//ErrWideMatrix      = errors.New("la matriz debe tener al menos tantas filas como columnas")
)

func Decompose(a [][]float64) (q, r [][]float64, err error) {
	if err = validate(a); err != nil {
		return nil, nil, err
	}

	rows, cols := len(a), len(a[0])

	k := cols
	if rows < cols {
		k = rows
	}

	left := mat.NewDense(rows, k, nil)
	for i := 0; i < rows; i++ {
		for j := 0; j < k; j++ {
			left.Set(i, j, a[i][j])
		}
	}

	var f mat.QR
	f.Factorize(left)

	var qm, rm mat.Dense
	f.QTo(&qm)
	f.RTo(&rm)

	// R completa (rows x cols). El bloque izquierdo viene de gonum.
	full := mat.NewDense(rows, cols, nil)
	for i := 0; i < rows; i++ {
		for j := 0; j < k; j++ {
			full.Set(i, j, rm.At(i, j))
		}
	}

	// Caso ancho: bloque derecho de R = Qᵀ * A2
	if cols > k {
		right := mat.NewDense(rows, cols-k, nil)
		for i := 0; i < rows; i++ {
			for j := 0; j < cols-k; j++ {
				right.Set(i, j, a[i][k+j])
			}
		}
		var c mat.Dense
		c.Mul(qm.T(), right)
		for i := 0; i < rows; i++ {
			for j := 0; j < cols-k; j++ {
				full.Set(i, k+j, c.At(i, j))
			}
		}
	}

	return toSlices(&qm), toSlices(full), nil
}

func validate(a [][]float64) error {
	if len(a) == 0 || len(a[0]) == 0 {
		return ErrEmptyMatrix
	}
	cols := len(a[0])
	for _, row := range a {
		if len(row) != cols {
			return ErrIrregularMatrix
		}
		for _, x := range row {
			if math.IsNaN(x) || math.IsInf(x, 0) {
				return ErrInvalidValue
			}
		}
	}
	return nil
}

func toSlices(m mat.Matrix) [][]float64 {
	rows, cols := m.Dims()
	out := make([][]float64, rows)
	for i := range out {
		out[i] = make([]float64, cols)
		for j := range out[i] {
			out[i][j] = m.At(i, j)
		}
	}
	return out
}
