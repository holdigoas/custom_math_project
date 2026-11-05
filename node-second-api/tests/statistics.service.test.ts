import { describe, it, expect } from 'vitest';
import { computeStatistics, isDiagonal } from '../src/services/statistics.service';

describe('computeStatistics', () => {
  it('calcula max, min, suma y promedio sobre Q y R combinadas', () => {
    const statistics = computeStatistics({
      q: [[1, 2], [3, 4]],
      r: [[5, 6], [7, 8]],
    });
    expect(statistics.max).toBe(8);
    expect(statistics.min).toBe(1);
    expect(statistics.sum).toBe(36);
    expect(statistics.average).toBe(4.5);
  });

  it('la suma compensada evita el error clásico de sumar 0.1 diez veces', () => {
    const statistics = computeStatistics({
      q: [[0.1, 0.1, 0.1, 0.1, 0.1]],
      r: [[0.1, 0.1, 0.1, 0.1, 0.1]],
    });
    expect(statistics.sum).toBe(1); // 0.9999999999999999
  });

  it('maneja valores negativos', () => {
    const statistics = computeStatistics({ q: [[-5, 2]], r: [[-1, 0]] });
    expect(statistics.max).toBe(2);
    expect(statistics.min).toBe(-5);
  });

  it('isDiagonal.any es verdadero si al menos una es diagonal', () => {
    const statistics = computeStatistics({ q: [[1, 2], [3, 4]], r: [[2, 0], [0, 3]] });
    expect(statistics.isDiagonal).toEqual({ q: false, r: true, any: true });
  });
});

describe('isDiagonal', () => {
  it('detecta matrices diagonales y no diagonales', () => {
    expect(isDiagonal([[1, 0], [0, 1]])).toBe(true);
    expect(isDiagonal([[1, 2], [0, 3]])).toBe(false);
  });

  it('funciona con matrices rectangulares', () => {
    expect(isDiagonal([[1, 0, 0], [0, 2, 0]])).toBe(true);
  });

  it('ignora residuos de punto flotante', () => {
    expect(isDiagonal([[1, 1e-12], [0, 1]])).toBe(true);
  });
});