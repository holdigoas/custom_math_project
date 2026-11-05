export type Matrix = number[][];

export interface MatricesInput {
    q: Matrix;
    r: Matrix;
}

export interface DiagonalInfo {
    q: boolean;
    r: boolean;
    any: boolean;
}

export interface Statistics {
    max: number;
    min: number;
    average: number;
    sum: number;
    isDiagonal: DiagonalInfo;
}

export interface ErrorResponse {
    error: {
        code: string;
        message: string;
    };
}