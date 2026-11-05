import { z } from 'zod';
import type { MatricesInput } from '../types';

// limite de la dimencion de la matriz
export const MAX_DIMENSION = 200;
// Límite de magnitud por valor o cada elemento.
export const MAX_ABS_VALUE = 1e150;

export class ValidationError extends Error {
    readonly code: string;

    constructor(code: string, message: string) {
        super(message);
        this.name = 'ValidationError';
        this.code = code;
    }
}

// definicion del schema zod
const matrixSchema = z
    .array(
        z
            .array(z.number().min(-MAX_ABS_VALUE).max(MAX_ABS_VALUE))
            .min(1)
            .max(MAX_DIMENSION)
    )
    .min(1)
    .max(MAX_DIMENSION)
    .refine((rows) => rows.every((row) => row.length === rows[0]?.length), {
        message: 'todas las filas deben tener la misma cantidad de columnas',
    });

// EStructura del cuerpo
const bodySchema = z.object({
    matrices: z.object({ q: matrixSchema, r: matrixSchema }),
});

// Valida el cuerpo { matrices: { q, r } } y devuelve { q, r } ya tipado.
export function validateBody(body: unknown): MatricesInput {
    const result = bodySchema.safeParse(body);
    if (!result.success) {
        const issue = result.error.issues[0];
        const where = issue.path.length > 0 ? issue.path.join('.') : 'body';
        throw new ValidationError('INVALID_BODY', `${where}: ${issue.message}`);
    }
    return result.data.matrices;
}