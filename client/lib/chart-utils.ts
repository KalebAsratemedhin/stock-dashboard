/**
 * Exponential moving average for a series of numbers.
 */
export function ema(data: number[], period: number): number[] {
  const out: number[] = [];
  const k = 2 / (period + 1);
  for (let i = 0; i < data.length; i++) {
    if (i === 0) out.push(data[0]);
    else out.push(data[i] * k + out[i - 1] * (1 - k));
  }
  return out;
}
