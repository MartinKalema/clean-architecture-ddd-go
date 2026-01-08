/**
 * Stress Test
 *
 * Pushes the system beyond normal capacity to find breaking points.
 * Ramps up to 10,000 concurrent users.
 *
 * Helps identify:
 *   - Maximum capacity
 *   - Recovery behavior
 *   - Error handling under load
 *
 * Stages:
 *   1. Warm up: 0 → 500 users
 *   2. Ramp: 500 → 2000 users
 *   3. Push: 2000 → 5000 users
 *   4. Peak: 5000 → 10000 users
 *   5. Sustain: hold at 10000 users
 *   6. Recovery: 10000 → 1000 users
 *   7. Cool down: 1000 → 0 users
 *
 * Usage:
 *   k6 run tests/load/stress.js
 *
 * With web dashboard:
 *   K6_WEB_DASHBOARD=true k6 run tests/load/stress.js
 *
 * WARNING: This test generates EXTREME load (10k users). Use with caution.
 */

import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Rate, Trend, Counter } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');
const requestCount = new Counter('total_requests');
const listBooksTrend = new Trend('list_books_duration');
const addBookTrend = new Trend('add_book_duration');

export const options = {
  stages: [
    { duration: '1m', target: 500 },    // Warm up to 500
    { duration: '2m', target: 500 },    // Stay at 500
    { duration: '1m', target: 2000 },   // Ramp to 2000
    { duration: '2m', target: 2000 },   // Stay at 2000
    { duration: '1m', target: 5000 },   // Push to 5000
    { duration: '2m', target: 5000 },   // Stay at 5000
    { duration: '1m', target: 10000 },  // Peak at 10000
    { duration: '3m', target: 10000 },  // Sustain 10000 (stress point)
    { duration: '1m', target: 1000 },   // Recovery
    { duration: '2m', target: 1000 },   // Verify recovery
    { duration: '1m', target: 0 },      // Cool down
  ],
  thresholds: {
    http_req_duration: ['p(95)<5000'],   // 95% < 5s (relaxed for extreme load)
    http_req_failed: ['rate<0.20'],      // Less than 20% errors (relaxed)
    errors: ['rate<0.20'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export function setup() {
  const res = http.get(`${BASE_URL}/api/v1/books`);
  if (res.status !== 200) {
    throw new Error(`API not reachable: ${res.status}`);
  }
  console.log('Starting stress test...');
  console.log('WARNING: This will generate significant load!');
}

export default function() {
  requestCount.add(1);

  // 70% read operations
  if (Math.random() < 0.7) {
    group('Read: List Books', function() {
      const res = http.get(`${BASE_URL}/api/v1/books`);
      listBooksTrend.add(res.timings.duration);

      const success = check(res, {
        'list: status 200': (r) => r.status === 200,
      });
      errorRate.add(!success);
    });
  }
  // 30% write operations
  else {
    group('Write: Add Book', function() {
      const timestamp = new Date().getTime();
      const vuId = __VU;
      const iteration = __ITER;

      const res = http.post(
        `${BASE_URL}/api/v1/books`,
        JSON.stringify({
          title: `Stress Book ${vuId}-${iteration}-${timestamp}`,
          author: `Stress Author ${vuId}`,
        }),
        { headers: { 'Content-Type': 'application/json' } }
      );
      addBookTrend.add(res.timings.duration);

      const success = check(res, {
        'add: status 201': (r) => r.status === 201,
      });
      errorRate.add(!success);
    });
  }

  // Minimal sleep to maximize load
  sleep(0.1);
}

export function teardown(data) {
  console.log('Stress test completed.');
  console.log('Check metrics for breaking points and recovery behavior.');
}
