import { describe, it, expect } from 'vitest';
import request from 'supertest';
import app from '../src/app';

describe('POST /api/v1/stats', () => {
    it('devuelve las estadísticas', async () => {
        const res = await request(app)
            .post('/api/v1/stats')
            .send({ matrices: { q: [[1, 0], [0, 1]], r: [[2, 3], [0, 4]] } });

        expect(res.status).toBe(200);
        expect(res.body.max).toBe(4);
        expect(res.body.min).toBe(0);
        expect(res.body.sum).toBe(11);
        expect(res.body.average).toBe(1.375);
        expect(res.body.isDiagonal).toEqual({ q: true, r: false, any: true });
    });

    const errorCases: Array<[string, string, string]> = [
        ['json inválido', '{"matrices":', 'INVALID_JSON'],
        ['falta r', JSON.stringify({ matrices: { q: [[1]] } }), 'INVALID_BODY'],
        ['sin matrices', JSON.stringify({}), 'INVALID_BODY'],
        ['q vacía', JSON.stringify({ matrices: { q: [], r: [[1]] } }), 'INVALID_BODY'],
        [
            'irregular',
            JSON.stringify({ matrices: { q: [[1, 2], [3]], r: [[1]] } }),
            'INVALID_BODY',
        ],
        [
            'valor no numérico',
            JSON.stringify({ matrices: { q: [[1, 'a']], r: [[1]] } }),
            'INVALID_BODY',
        ],
        [
            'valor demasiado grande',
            JSON.stringify({ matrices: { q: [[1e200]], r: [[1]] } }),
            'INVALID_BODY',
        ],
    ];

    it.each(errorCases)('error 400 (%s) -> %s', async (_name, body, code) => {
        const res = await request(app)
            .post('/api/v1/stats')
            .set('Content-Type', 'application/json')
            .send(body);

        expect(res.status).toBe(400);
        expect(res.body.error.code).toBe(code);
    });
});

describe('GET /health', () => {
    it('responde ok', async () => {
        const res = await request(app).get('/health');
        expect(res.status).toBe(200);
        expect(res.body).toEqual({ status: 'ok' });
    });
});