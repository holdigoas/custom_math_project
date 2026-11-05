import express from 'express';
import type { NextFunction, Request, Response } from 'express';
import { getStatistics } from './controllers/statistics.controller';
import { ValidationError } from './validators/matrices.validator';
import type { ErrorResponse } from './types';

const app = express();
app.disable('x-powered-by');
app.use(express.json({ limit: '1mb' }));

function errorBody(code: string, message: string): ErrorResponse {
    return { error: { code, message } };
}

// Ruta principal
app.get('/home', (_req: Request, res: Response) => {
    res.json({ status: 'ok' });
});

// Rutas versionadas
app.post('/api/v1/statistics', getStatistics);

// Ruta inexistente
app.use((_req: Request, res: Response) => {
    res.status(404).json(errorBody('NOT_FOUND', 'ruta no encontrada'));
});

app.use((err: unknown, _req: Request, res: Response, _next: NextFunction) => {
    if (err instanceof ValidationError) {
        res.status(400).json(errorBody(err.code, err.message));
        return;
    }

    const type =
        typeof err === 'object' && err !== null && 'type' in err
            ? (err as { type?: unknown }).type
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

export default app;