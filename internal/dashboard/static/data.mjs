export const colors = ['#238a78', '#c58b3d', '#728cae', '#9b79a6', '#b76556'];
export const methodLabel = method => ({base: 'Базовая опора', height_100m: 'Высота ±100 м', large_window_161: 'Окно 161 × 161'}[method] || method);
export const day = value => typeof value === 'string' ? value.slice(0, 10) : '';
export const finite = value => typeof value === 'number' && Number.isFinite(value);

export function validMeasurement(row) {
  return row && finite(row.id) && finite(row.job_id) &&
    typeof row.reference_method === 'string' && row.reference_method.length > 0 &&
    finite(row.displacement_mm) && finite(row.dam_coherence_mean) &&
    row.dam_coherence_mean >= 0 && row.dam_coherence_mean <= 1 &&
    Number.isFinite(Date.parse(row.reference_date)) && Number.isFinite(Date.parse(row.secondary_date)) &&
    day(row.secondary_date) > day(row.reference_date);
}

export function filterRows(rows, {from = '', to = '', minCoherence = 0, methods}) {
  return rows.filter(row => (!from || day(row.secondary_date) >= from) &&
    (!to || day(row.secondary_date) <= to) && row.dam_coherence_mean >= minCoherence &&
    methods.has(row.reference_method)).sort((a, b) =>
    day(a.secondary_date).localeCompare(day(b.secondary_date)) ||
    day(a.reference_date).localeCompare(day(b.reference_date)) ||
    a.reference_method.localeCompare(b.reference_method) || a.id - b.id);
}

export function summarize(rows) {
  const pairs = new Map();
  for (const row of rows) {
    const values = pairs.get(row.job_id) || [];
    values.push(row.dam_coherence_mean);
    pairs.set(row.job_id, values);
  }
  // Each pair has equal weight even when some methods are absent or filtered.
  const means = [...pairs.values()].map(values => values.reduce((a, b) => a + b, 0) / values.length);
  return {
    pairs: pairs.size,
    measurements: rows.length,
    methods: new Set(rows.map(row => row.reference_method)).size,
    coherence: means.length ? means.reduce((a, b) => a + b, 0) / means.length : null,
    maxDisplacement: rows.length ? rows.reduce((max, row) => Math.max(max, Math.abs(row.displacement_mm)), 0) : null,
  };
}
