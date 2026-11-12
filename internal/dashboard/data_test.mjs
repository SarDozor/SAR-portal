import test from 'node:test';
import assert from 'node:assert/strict';
import {filterRows, summarize, validMeasurement} from './static/data.mjs';

const measurement = (changes = {}) => ({id: 1, job_id: 1, reference_method: 'base', reference_date: '2018-01-01T00:00:00Z', secondary_date: '2018-01-13T00:00:00Z', displacement_mm: -4, dam_coherence_mean: .8, ...changes});

test('summary counts pairs once and weights each pair equally', () => {
  const rows = [measurement(), measurement({id: 2, reference_method: 'height_100m'}), measurement({id: 3, reference_method: 'large_window_161'}), measurement({id: 4, job_id: 2, dam_coherence_mean: .4, displacement_mm: 2})];
  const result = summarize(rows);
  assert.equal(result.pairs, 2);
  assert.equal(result.measurements, 4);
  assert.equal(result.methods, 3);
  assert.ok(Math.abs(result.coherence - .6) < 1e-12);
  assert.equal(result.maxDisplacement, 4);
});

test('filters use inclusive secondary dates and selected methods', () => {
  const rows = [measurement(), measurement({id: 2, reference_method: 'other'}), measurement({id: 3, secondary_date: '2018-02-01T00:00:00Z'})];
  assert.deepEqual(filterRows(rows, {from: '2018-01-13', to: '2018-01-13', minCoherence: .8, methods: new Set(['base'])}), [rows[0]]);
  assert.equal(filterRows(rows, {methods: new Set()}).length, 0);
  assert.equal(filterRows(rows, {minCoherence: .9, methods: new Set(['base'])}).length, 0);
});

test('empty selections do not invent zero quality or displacement', () => {
  assert.deepEqual(summarize([]), {pairs: 0, measurements: 0, methods: 0, coherence: null, maxDisplacement: null});
  assert.equal(summarize([measurement({displacement_mm: 0, dam_coherence_mean: 0})]).coherence, 0);
});

test('reject malformed measurements and preserve real zero values', () => {
  assert.ok(validMeasurement(measurement({displacement_mm: 0, dam_coherence_mean: 0})));
  for (const changes of [{displacement_mm: null}, {displacement_mm: NaN}, {dam_coherence_mean: 2}, {secondary_date: 'bad'}, {secondary_date: '2018-01-01T00:00:00Z'}]) {
    assert.equal(validMeasurement(measurement(changes)), false);
  }
});
