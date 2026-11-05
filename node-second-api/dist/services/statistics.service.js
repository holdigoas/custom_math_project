"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.EPSILON = void 0;
exports.isDiagonal = isDiagonal;
exports.computeStatistics = computeStatistics;
const simple_statistics_1 = require("simple-statistics");
exports.EPSILON = 1e-9;
function isDiagonal(matrix, epsilon = exports.EPSILON) {
    for (let i = 0; i < matrix.length; i++) {
        for (let j = 0; j < matrix[i].length; j++) {
            if (i !== j && Math.abs(matrix[i][j]) > epsilon) {
                return false;
            }
        }
    }
    return true;
}
function computeStatistics({ q, r }) {
    const values = q.flat().concat(r.flat());
    const diagonalQ = isDiagonal(q);
    const diagonalR = isDiagonal(r);
    return {
        max: (0, simple_statistics_1.max)(values),
        min: (0, simple_statistics_1.min)(values),
        average: (0, simple_statistics_1.mean)(values),
        sum: (0, simple_statistics_1.sum)(values),
        isDiagonal: { q: diagonalQ, r: diagonalR, any: diagonalQ || diagonalR },
    };
}
