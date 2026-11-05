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
	ErrWideMatrix      = errors.New("la matriz debe tener al menos tantas filas como columnas")
)

func Decompose(a [][]float64) (q, r [][]float64, err error) {
	if err = validate(a); err != nil {
		return nil, nil, err
	}

	rows, cols := len(a), len(a[0])
	if rows < cols {
		return nil, nil, ErrWideMatrix
	}

	data := make([]float64, 0, rows*cols)
	for _, row := range a {
		data = append(data, row...)
	}

	var f mat.QR
	f.Factorize(mat.NewDense(rows, cols, data))

	var qm, rm mat.Dense
	f.QTo(&qm)
	f.RTo(&rm)

	return toSlices(&qm), toSlices(&rm), nil
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
