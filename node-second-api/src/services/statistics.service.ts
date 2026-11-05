import { max, mean, min, sum } from 'simple-statistics';
import type { Matrix, MatricesInput, Statistics } from '../types';

export const EPSILON = 1e-9;

export function isDiagonal(matrix: Matrix, epsilon: number = EPSILON): boolean {
    for (let i = 0; i < matrix.length; i++) {
        for (let j = 0; j < matrix[i].length; j++) {
            if (i !== j && Math.abs(matrix[i][j]) > epsilon) {
                return false;
            }
        }
    }
    return true;
}


export function computeStatistics({ q, r }: MatricesInput): Statistics {
    const values = q.flat().concat(r.flat());
    const diagonalQ = isDiagonal(q);
    const diagonalR = isDiagonal(r);

    return {
        max: max(values),
        min: min(values),
        average: mean(values),
        sum: sum(values),
        isDiagonal: { q: diagonalQ, r: diagonalR, any: diagonalQ || diagonalR },
    };
}