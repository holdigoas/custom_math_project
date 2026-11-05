"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const express_1 = __importDefault(require("express"));
const statistics_controller_1 = require("./controllers/statistics.controller");
const matrices_validator_1 = require("./validators/matrices.validator");
const app = (0, express_1.default)();
app.disable('x-powered-by');
app.use(express_1.default.json({ limit: '1mb' }));
function errorBody(code, message) {
    return { error: { code, message } };
}
// Ruta principal
app.get('/health', (_req, res) => {
    res.json({ status: 'ok' });
});
// Rutas versionadas.
app.post('/api/v1/stats', statistics_controller_1.getStatistics);
// Ruta inexistente.
app.use((_req, res) => {
    res.status(404).json(errorBody('NOT_FOUND', 'ruta no encontrada'));
});
// Manejo centralizado de errores (mismo formato que la API de Go).
app.use((err, _req, res, _next) => {
    if (err instanceof matrices_validator_1.ValidationError) {
        res.status(400).json(errorBody(err.code, err.message));
        return;
    }
    // Errores del lector de JSON de Express traen un "type" identificable.
    const type = typeof err === 'object' && err !== null && 'type' in err
        ? err.type
        : undefined;
    if (type === 'entity.parse.failed') {
        res.status(400).json(errorBody('INVALID_JSON', 'el cuerpo no es un JSON válido'));
        return;
    }
    if (type === 'entity.too.large') {
        res
            .status(413)
            .json(errorBody('PAYLOAD_TOO_LARGE', 'el cuerpo supera el tamaño permitido'));
        return;
    }
    console.error(err);
    res.status(500).json(errorBody('INTERNAL_ERROR', 'error interno del servidor'));
});
exports.default = app;
