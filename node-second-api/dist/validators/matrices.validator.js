"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ValidationError = exports.MAX_ABS_VALUE = exports.MAX_DIMENSION = void 0;
exports.validateBody = validateBody;
const zod_1 = require("zod");
exports.MAX_DIMENSION = 20;
exports.MAX_ABS_VALUE = 1e150;
class ValidationError extends Error {
    code;
    constructor(code, message) {
        super(message);
        this.name = 'ValidationError';
        this.code = code;
    }
}
exports.ValidationError = ValidationError;
const matrixSchema = zod_1.z
    .array(zod_1.z
    .array(zod_1.z.number().min(-exports.MAX_ABS_VALUE).max(exports.MAX_ABS_VALUE))
    .min(1)
    .max(exports.MAX_DIMENSION))
    .min(1)
    .max(exports.MAX_DIMENSION)
    .refine((rows) => rows.every((row) => row.length === rows[0]?.length), {
    message: 'todas las filas deben tener la misma cantidad de columnas',
});
const bodySchema = zod_1.z.object({
    matrices: zod_1.z.object({ q: matrixSchema, r: matrixSchema }),
});
// Valida el cuerpo { matrices: { q, r } } y devuelve { q, r } ya tipado.
function validateBody(body) {
    const result = bodySchema.safeParse(body);
    if (!result.success) {
        const issue = result.error.issues[0];
        const where = issue.path.length > 0 ? issue.path.join('.') : 'body';
        throw new ValidationError('INVALID_BODY', `${where}: ${issue.message}`);
    }
    return result.data.matrices;
}
