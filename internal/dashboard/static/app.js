import {colors, methodLabel, day, finite, validMeasurement, filterRows, summarize} from './data.mjs';

const $ = id => document.getElementById(id);
const state = {objects: [], rows: [], filtered: [], methods: [], selected: new Set(), page: 0, request: 0, controller: null, loadedAt: null};
const pageSize = 12;
const number = (value, digits = 2) => finite(value) ? new Intl.NumberFormat('ru-RU', {maximumFractionDigits: digits}).format(value) : '—';
const date = value => {
  const parsed = new Date(value);
  return Number.isNaN(parsed.getTime()) ? '—' : parsed.toLocaleDateString('ru-RU', {timeZone: 'UTC'});
};
const color = method => colors[Math.max(0, state.methods.indexOf(method)) % colors.length];
const setText = (id, value) => { $(id).textContent = value; };

function message(text = '', error = false) {
  $('message').hidden = !text;
  $('message').classList.toggle('error', error);
  setText('message', text);
}

async function getJSON(path, signal) {
  const response = await fetch(path, {signal, headers: {Accept: 'application/json'}, cache: 'no-store'});
  if (!response.ok) throw new Error(`Сервер вернул HTTP ${response.status}.`);
  const value = await response.json();
  if (!Array.isArray(value)) throw new Error('Сервер вернул неожиданный формат данных.');
  return value;
}

function connection(text, status) {
  setText('connection', text);
  $('connection').className = `connection ${status}`;
}

async function loadObjects() {
  state.controller?.abort();
  const controller = new AbortController();
  state.controller = controller;
  const request = ++state.request;
  const previous = $('object-select').value;
  $('refresh').disabled = true;
  $('object-select').disabled = true;
  $('workspace').hidden = true;
  connection('Подключение…', '');
  message('Загружаем реестр объектов…');
  try {
    const objects = await getJSON('/objects', controller.signal);
    if (request !== state.request) return;
    state.objects = objects;
    $('object-select').replaceChildren();
    for (const obj of objects) {
      const option = document.createElement('option');
      option.value = String(obj.id);
      option.textContent = obj.name || `Объект №${obj.id}`;
      $('object-select').append(option);
    }
    setText('object-count', `В реестре: ${objects.length}`);
    if (!objects.length) {
      const option = document.createElement('option');
      option.textContent = 'Нет объектов';
      $('object-select').append(option);
      setText('object-name', 'Реестр пока пуст');
      message('Добавьте объект и импортируйте измерения, затем нажмите «Обновить».');
      connection('API доступен', 'online');
      return;
    }
    if (objects.some(obj => String(obj.id) === previous)) $('object-select').value = previous;
    $('object-select').disabled = false;
    await loadObject();
  } catch (error) {
    if (error.name !== 'AbortError' && request === state.request) {
      message(`Не удалось загрузить объекты. ${error.message} Проверьте сервер и нажмите «Обновить».`, true);
      connection('Нет соединения', 'offline');
    }
  } finally {
    if (state.controller === controller) $('refresh').disabled = false;
  }
}

async function loadObject() {
  state.controller?.abort();
  const controller = new AbortController();
  state.controller = controller;
  const request = ++state.request;
  const obj = state.objects.find(item => String(item.id) === $('object-select').value);
  if (!obj) return;
  state.rows = [];
  state.filtered = [];
  $('workspace').hidden = true;
  $('refresh').disabled = true;
  setText('object-name', obj.name || `Объект №${obj.id}`);
  setText('object-subtitle', [obj.country, 'Спутниковые измерения поверхности'].filter(Boolean).join(' · '));
  connection('Загрузка измерений…', '');
  message('Загружаем наблюдения выбранного объекта…');
  try {
    const rows = await getJSON(`/objects/${encodeURIComponent(obj.id)}/measurements`, controller.signal);
    if (request !== state.request) return;
    if (rows.some(row => !validMeasurement(row))) throw new Error('В ответе есть измерения с некорректными датами или значениями.');
    state.rows = rows;
    state.methods = [...new Set(rows.map(row => row.reference_method))].sort();
    state.selected = new Set(state.methods);
    state.loadedAt = new Date();
    renderMethods();
    renderLocation(obj);
    resetFilters();
    $('workspace').hidden = false;
    connection('Данные загружены', 'online');
    setText('updated-at', `Обновлено в ${state.loadedAt.toLocaleTimeString('ru-RU', {hour: '2-digit', minute: '2-digit'})}`);
  } catch (error) {
    if (error.name !== 'AbortError' && request === state.request) {
      message(`Не удалось загрузить измерения. ${error.message} Нажмите «Обновить», чтобы повторить.`, true);
      connection('Ошибка загрузки', 'offline');
    }
  } finally {
    if (request === state.request) $('refresh').disabled = false;
  }
}

function renderMethods() {
  $('methods').replaceChildren();
  for (const method of state.methods) {
    const label = document.createElement('label');
    label.className = 'method-chip';
    const input = document.createElement('input');
    input.type = 'checkbox';
    input.checked = true;
    input.value = method;
    input.addEventListener('change', () => {
      input.checked ? state.selected.add(method) : state.selected.delete(method);
      render();
    });
    label.append(input, document.createTextNode(methodLabel(method)));
    $('methods').append(label);
  }
  if (!state.methods.length) $('methods').textContent = 'Нет измерений';
}

function resetFilters() {
  const dates = state.rows.map(row => day(row.secondary_date)).sort();
  $('date-from').value = dates[0] || '';
  $('date-to').value = dates.at(-1) || '';
  $('coherence-min').value = '0';
  state.selected = new Set(state.methods);
  for (const input of $('methods').querySelectorAll('input')) input.checked = true;
  render();
}

function renderLocation(obj) {
  setText('country', obj.country || 'Не указана');
  setText('raise-type', obj.raise_type || 'Не указан');
  setText('height', finite(obj.max_height_m) ? `${number(obj.max_height_m)} м` : 'Не указана');
  setText('object-status', obj.status || 'Не указан');
  const {latitude: lat, longitude: lon} = obj;
  $('map').replaceChildren();
  $('map-link').hidden = true;
  setText('coordinates', 'Не указаны');
  if (!finite(lat) || !finite(lon) || Math.abs(lat) > 90 || Math.abs(lon) > 180) {
    const p = document.createElement('p');
    p.textContent = 'У объекта нет корректных координат';
    $('map').append(p);
    return;
  }
  const iframe = document.createElement('iframe');
  const bbox = [Math.max(-180, lon - .025), Math.max(-90, lat - .015), Math.min(180, lon + .025), Math.min(90, lat + .015)].join(',');
  iframe.src = `https://www.openstreetmap.org/export/embed.html?bbox=${encodeURIComponent(bbox)}&layer=mapnik&marker=${lat},${lon}`;
  iframe.title = `Карта расположения: ${obj.name}`;
  iframe.loading = 'lazy';
  iframe.referrerPolicy = 'no-referrer';
  $('map').append(iframe);
  setText('coordinates', `${lat.toFixed(6)}, ${lon.toFixed(6)}`);
  $('map-link').href = `https://www.openstreetmap.org/?mlat=${lat}&mlon=${lon}#map=14/${lat}/${lon}`;
  $('map-link').hidden = false;
}

function render() {
  const from = $('date-from').value;
  const to = $('date-to').value;
  const threshold = $('coherence-min').valueAsNumber;
  const invalid = (from && to && from > to) || !finite(threshold) || threshold < 0 || threshold > 1;
  message(invalid ? 'Проверьте фильтры: начало периода не должно быть позже конца, а когерентность должна быть от 0 до 1.' : (!state.rows.length ? 'Для этого объекта пока нет измерений. После импорта нажмите «Обновить».' : ''), !!invalid);
  state.filtered = invalid ? [] : filterRows(state.rows, {from, to, minCoherence: threshold, methods: state.selected});
  const summary = summarize(state.filtered);
  setText('pair-count', number(summary.pairs, 0));
  setText('measurement-count', number(summary.measurements, 0));
  setText('method-count', `Методов опоры: ${summary.methods}`);
  setText('mean-coherence', number(summary.coherence, 3));
  setText('max-displacement', number(summary.maxDisplacement));
  $('export').disabled = !state.filtered.length;
  drawChart('displacement-chart', 'displacement_mm', false);
  drawChart('coherence-chart', 'dam_coherence_mean', true);
  $('chart-legend').replaceChildren();
  for (const method of state.methods.filter(method => state.selected.has(method))) {
    const item = document.createElement('span');
    item.append(swatch(method), document.createTextNode(methodLabel(method)));
    $('chart-legend').append(item);
  }
  state.page = 0;
  renderTable();
}

function swatch(method) {
  const node = document.createElement('i');
  node.className = 'swatch';
  node.style.background = color(method);
  return node;
}

function svgElement(tag, attrs = {}, text) {
  const element = document.createElementNS('http://www.w3.org/2000/svg', tag);
  for (const [key, value] of Object.entries(attrs)) element.setAttribute(key, value);
  if (text !== undefined) element.textContent = text;
  return element;
}

function drawChart(id, field, coherence) {
  const container = $(id);
  container.replaceChildren();
  const rows = state.filtered;
  if (!rows.length) {
    const empty = document.createElement('div');
    empty.className = 'chart-empty';
    empty.textContent = 'Нет измерений для выбранных фильтров';
    container.append(empty);
    return;
  }
  const width = Math.max(300, container.clientWidth ? container.clientWidth - 26 : 760);
  const height = coherence ? 170 : 270;
  const left = 51, right = width - 22, top = 20, bottom = height - 35;
  const times = rows.map(row => Date.parse(row.secondary_date));
  let start = Math.min(...times), end = Math.max(...times);
  if (start === end) { start -= 86400000; end += 86400000; }
  const values = rows.map(row => row[field]);
  let min = coherence ? 0 : Math.min(0, ...values);
  let max = coherence ? 1 : Math.max(0, ...values);
  if (!coherence) { const pad = Math.max((max - min) * .12, .5); min -= pad; max += pad; }
  const x = value => left + (value - start) / (end - start) * (right - left);
  const y = value => bottom - (value - min) / (max - min) * (bottom - top);
  const svg = svgElement('svg', {viewBox: `0 0 ${width} ${height}`, role: 'group', 'aria-label': coherence ? 'Когерентность по датам второго снимка. Значения доступны в таблице ниже.' : 'Смещение за интервал по датам второго снимка. Значения доступны в таблице ниже.'});
  for (let i = 0; i <= 4; i++) {
    const value = min + (max - min) * i / 4;
    svg.append(svgElement('line', {x1: left, x2: right, y1: y(value), y2: y(value), class: 'chart-grid'}));
    svg.append(svgElement('text', {x: left - 10, y: y(value) + 3, 'text-anchor': 'end', class: 'chart-label'}, number(value, coherence ? 2 : 1)));
  }
  if (!coherence) svg.append(svgElement('line', {x1: left, x2: right, y1: y(0), y2: y(0), class: 'chart-zero'}));
  const dateTicks = width < 480 ? 2 : 4;
  for (let i = 0; i <= dateTicks; i++) {
    const value = start + (end - start) * i / dateTicks;
    svg.append(svgElement('text', {x: x(value), y: height - 8, 'text-anchor': i === 0 ? 'start' : i === dateTicks ? 'end' : 'middle', class: 'chart-label'}, date(value)));
  }
  const tooltip = document.createElement('div');
  tooltip.className = 'chart-tooltip';
  tooltip.hidden = true;
  tooltip.setAttribute('role', 'status');
  for (const method of state.methods) {
    const points = rows.filter(row => row.reference_method === method);
    if (!points.length) continue;
    const path = points.map((row, i) => `${i ? 'L' : 'M'}${x(Date.parse(row.secondary_date))},${y(row[field])}`).join(' ');
    svg.append(svgElement('path', {d: path, stroke: color(method), class: 'chart-line'}));
    for (const row of points) {
      const description = `${methodLabel(method)}\n${date(row.reference_date)} → ${date(row.secondary_date)}\n${coherence ? 'Когерентность' : 'Смещение'}: ${number(row[field], 3)}${coherence ? '' : ' мм'}`;
      const point = svgElement('circle', {cx: x(Date.parse(row.secondary_date)), cy: y(row[field]), r: 3.3, fill: color(method), class: 'chart-dot', tabindex: '0', 'aria-label': description});
      point.append(svgElement('title', {}, description));
      const show = () => {
        tooltip.textContent = description;
        tooltip.hidden = false;
        const bounds = point.getBoundingClientRect(), outer = container.getBoundingClientRect();
        tooltip.style.left = `${Math.max(4, Math.min(bounds.left - outer.left + 10, outer.width - tooltip.offsetWidth - 8))}px`;
        tooltip.style.top = `${Math.max(0, bounds.top - outer.top - tooltip.offsetHeight - 8)}px`;
      };
      point.addEventListener('mouseenter', show);
      point.addEventListener('focus', show);
      point.addEventListener('click', show);
      point.addEventListener('mouseleave', () => { tooltip.hidden = true; });
      point.addEventListener('blur', () => { tooltip.hidden = true; });
      point.addEventListener('keydown', event => { if (event.key === 'Escape') tooltip.hidden = true; });
      svg.append(point);
    }
  }
  container.append(svg, tooltip);
}

function renderTable() {
  const rows = state.filtered;
  const start = state.page * pageSize;
  $('measurement-rows').replaceChildren();
  setText('table-count', number(rows.length, 0));
  for (const row of rows.slice(start, start + pageSize)) {
    const tr = document.createElement('tr');
    const values = [`${date(row.reference_date)} → ${date(row.secondary_date)}`, methodLabel(row.reference_method), number(row.displacement_mm, 3), number(row.dam_coherence_mean, 3), number(row.dam_coherence_valid_pixels, 0)];
    values.forEach((value, i) => {
      const td = document.createElement('td');
      if (i >= 2) td.className = 'numeric';
      if (i === 1) {
        const span = document.createElement('span');
        span.className = 'method-cell';
        span.append(swatch(row.reference_method), document.createTextNode(value));
        td.append(span);
      } else if (i === 3) {
        const span = document.createElement('span'); span.className = 'coherence-value';
        const bar = document.createElement('span'); bar.className = 'mini-bar';
        const fill = document.createElement('i'); fill.style.width = `${row.dam_coherence_mean * 100}%`;
        bar.append(fill); span.append(bar, document.createTextNode(value)); td.append(span);
      } else td.textContent = value;
      tr.append(td);
    });
    $('measurement-rows').append(tr);
  }
  if (!rows.length) {
    const tr = document.createElement('tr'), td = document.createElement('td');
    td.colSpan = 5; td.textContent = 'Нет измерений для выбранных фильтров'; tr.append(td); $('measurement-rows').append(tr);
  }
  setText('page-info', rows.length ? `${start + 1}–${Math.min(start + pageSize, rows.length)} из ${rows.length}` : '0 измерений');
  $('previous').disabled = state.page === 0;
  $('next').disabled = start + pageSize >= rows.length;
}

$('refresh').addEventListener('click', loadObjects);
$('object-select').addEventListener('change', loadObject);
$('reset').addEventListener('click', resetFilters);
for (const id of ['date-from', 'date-to', 'coherence-min']) $(id).addEventListener('change', render);
$('previous').addEventListener('click', () => { if (state.page > 0) { state.page--; renderTable(); } });
$('next').addEventListener('click', () => { if ((state.page + 1) * pageSize < state.filtered.length) { state.page++; renderTable(); } });
$('export').addEventListener('click', () => {
  const blob = new Blob([JSON.stringify(state.filtered, null, 2)], {type: 'application/json'});
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url; link.download = `sar-dozor-object-${$('object-select').value}.json`;
  document.body.append(link); link.click(); link.remove();
  setTimeout(() => URL.revokeObjectURL(url), 1000);
});
const chartWidths = new Map();
const resizeCharts = new ResizeObserver(entries => {
  for (const entry of entries) {
    const width = Math.round(entry.contentRect.width);
    if (width > 0 && chartWidths.get(entry.target.id) !== width) {
      chartWidths.set(entry.target.id, width);
      const coherence = entry.target.id === 'coherence-chart';
      drawChart(entry.target.id, coherence ? 'dam_coherence_mean' : 'displacement_mm', coherence);
    }
  }
});
resizeCharts.observe($('displacement-chart'));
resizeCharts.observe($('coherence-chart'));
loadObjects();
