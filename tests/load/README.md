# Load Tests

Performance and load testing using [k6](https://k6.io/).

## Prerequisites

Install k6:

```bash
# macOS
brew install k6

# Linux
sudo gpg -k
sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt-get update
sudo apt-get install k6

# Windows
choco install k6
```

## Test Types

| Test | Purpose | Duration | Users |
|------|---------|----------|-------|
| `smoke.js` | Sanity check | 10s | 1 |
| `load.js` | Normal traffic | 12m | 100-1000 |
| `stress.js` | Find breaking point | 17m | 500-10000 |

## Running Tests

### Basic Run

```bash
# Smoke test (quick sanity check)
k6 run tests/load/smoke.js

# Load test (normal traffic simulation)
k6 run tests/load/load.js

# Stress test (find breaking points)
k6 run tests/load/stress.js
```

### With Web Dashboard

Real-time visualization in browser at http://localhost:5665:

```bash
K6_WEB_DASHBOARD=true k6 run tests/load/smoke.js
```

### Custom Base URL

```bash
BASE_URL=http://production.example.com k6 run tests/load/load.js
```

### Output to JSON

```bash
k6 run --out json=results.json tests/load/load.js
```

## Understanding Results

### Key Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| `http_req_duration` | Request latency | p95 < 500ms |
| `http_req_failed` | Error rate | < 1% |
| `http_reqs` | Requests per second | Higher is better |
| `vus` | Virtual users | Configured |

### Percentiles

- **p50** (median): 50% of requests faster than this
- **p95**: 95% of requests faster than this
- **p99**: 99% of requests faster than this

### Example Output

```
     checks.........................: 100.00% ✓ 234      ✗ 0
     data_received..................: 45 kB   4.5 kB/s
     data_sent......................: 23 kB   2.3 kB/s
     http_req_duration..............: avg=12.3ms min=5ms med=10ms max=45ms p(95)=25ms
     http_req_failed................: 0.00%   ✓ 0        ✗ 234
     http_reqs......................: 234     23.4/s
     vus............................: 1       min=1      max=1
```

## Thresholds

Tests will **fail** if thresholds are not met:

| Test | Threshold |
|------|-----------|
| Smoke | p95 < 500ms, errors < 1% |
| Load | p95 < 1000ms, errors < 5% |
| Stress | p95 < 3000ms, errors < 15% |

## CI/CD Integration

```yaml
# GitHub Actions example
- name: Run load tests
  run: |
    k6 run tests/load/smoke.js
    # Fail pipeline if thresholds not met
```

## Grafana Integration

For persistent dashboards:

1. Install InfluxDB
2. Run k6 with InfluxDB output:
   ```bash
   k6 run --out influxdb=http://localhost:8086/k6 tests/load/load.js
   ```
3. Import Grafana dashboard: https://grafana.com/grafana/dashboards/2587
