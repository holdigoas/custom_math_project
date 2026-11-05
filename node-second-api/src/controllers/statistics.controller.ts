import type { Request, Response } from 'express';
import { validateBody } from '../validators/matrices.validator';
import { computeStatistics } from '../services/statistics.service';

export function getStatistics(req: Request, res: Response): void {
    const matrices = validateBody(req.body);
    res.json(computeStatistics(matrices));
}