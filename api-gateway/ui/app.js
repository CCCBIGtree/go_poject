const requestsEl = document.getElementById('requests');
const errorsEl = document.getElementById('errors');
const latencyEl = document.getElementById('latency');
const qpsEl = document.getElementById('qps');
const canvas = document.getElementById('qpsChart');
const ctx = canvas.getContext('2d');

let lastReq = 0;
let lastTs = Date.now();
const series = [];

function drawChart() {
  const w = canvas.width;
  const h = canvas.height;
  ctx.clearRect(0, 0, w, h);

  ctx.strokeStyle = '#334155';
  ctx.lineWidth = 1;
  for (let i = 1; i <= 5; i++) {
    const y = (h / 6) * i;
    ctx.beginPath();
    ctx.moveTo(0, y);
    ctx.lineTo(w, y);
    ctx.stroke();
  }

  if (!series.length) return;
  const max = Math.max(...series, 1);
  const step = w / Math.max(series.length - 1, 1);

  ctx.strokeStyle = '#38bdf8';
  ctx.lineWidth = 2;
  ctx.beginPath();
  series.forEach((v, i) => {
    const x = i * step;
    const y = h - (v / max) * (h - 20) - 10;
    if (i === 0) ctx.moveTo(x, y);
    else ctx.lineTo(x, y);
  });
  ctx.stroke();
}

async function refresh() {
  try {
    const resp = await fetch('/metrics');
    const data = await resp.json();
    const now = Date.now();

    const req = data.requests_total || 0;
    const err = data.errors_total || 0;
    const avgNs = data.avg_latency_ns || 0;

    const dt = (now - lastTs) / 1000;
    const qps = dt > 0 ? (req - lastReq) / dt : 0;

    requestsEl.textContent = req;
    errorsEl.textContent = err;
    latencyEl.textContent = (avgNs / 1e6).toFixed(2);
    qpsEl.textContent = qps.toFixed(2);

    series.push(Number(qps.toFixed(2)));
    if (series.length > 30) series.shift();
    drawChart();

    lastReq = req;
    lastTs = now;
  } catch (e) {
    console.error('refresh failed', e);
  }
}

refresh();
setInterval(refresh, 1000);
