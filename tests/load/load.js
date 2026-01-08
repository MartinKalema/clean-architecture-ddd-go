/**
 * Load Test
 *
 * Simulates normal and peak traffic patterns.
 * Ramps up to 5000 concurrent users.
 *
 * Stages:
 *   1. Warm up: 0 → 500 users
 *   2. Normal load: 500 → 2000 users
 *   3. Peak load: 2000 → 5000 users
 *   4. Sustain peak: hold at 5000 users
 *   5. Cool down: 5000 → 0 users
 *
 * Usage:
 *   k6 run tests/load/load.js
 *
 * With web dashboard:
 *   K6_WEB_DASHBOARD=true k6 run tests/load/load.js
 */

import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Rate, Trend } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');
const listBooksTrend = new Trend('list_books_duration');
const addBookTrend = new Trend('add_book_duration');
const getBookTrend = new Trend('get_book_duration');

export const options = {
  stages: [
    { duration: '1m', target: 500 },   // Warm up to 500 users
    { duration: '2m', target: 500 },   // Stay at 500 users
    { duration: '1m', target: 2000 },  // Ramp to 2000 users
    { duration: '3m', target: 2000 },  // Stay at 2000 users
    { duration: '1m', target: 5000 },  // Ramp to 5000 users
    { duration: '3m', target: 5000 },  // Stay at 5000 users (peak)
    { duration: '1m', target: 0 },     // Cool down
  ],
  thresholds: {
    http_req_duration: ['p(95)<2000', 'p(99)<3000'],  // 95% < 2s, 99% < 3s
    http_req_failed: ['rate<0.05'],                   // Less than 5% errors
    errors: ['rate<0.05'],                            // Custom error rate
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export function setup() {
  // Verify API is reachable before starting
  const res = http.get(`${BASE_URL}/api/v1/books`);
  if (res.status !== 200) {
    throw new Error(`API not reachable: ${res.status}`);
  }
  console.log('API is reachable. Starting load test...');
}

export default function() {
  let bookId;

  group('List Books', function() {
    const res = http.get(`${BASE_URL}/api/v1/books`);
    listBooksTrend.add(res.timings.duration);

    const success = check(res, {
      'list: status 200': (r) => r.status === 200,
    });
    errorRate.add(!success);
  });

  sleep(0.3);

  group('Add Book', function() {
    const timestamp = new Date().getTime();
    const vuId = __VU;
    const res = http.post(
      `${BASE_URL}/api/v1/books`,
      JSON.stringify({
        title: `Load Test Book ${vuId}-${timestamp}`,
        author: `Author ${vuId}`,
      }),
      { headers: { 'Content-Type': 'application/json' } }
    );
    addBookTrend.add(res.timings.duration);

    const success = check(res, {
      'add: status 201': (r) => r.status === 201,
    });
    errorRate.add(!success);

    if (res.status === 201) {
      bookId = JSON.parse(res.body).ID;
    }
  });

  sleep(0.3);

  if (bookId) {
    group('Get Book', function() {
      const res = http.get(`${BASE_URL}/api/v1/books/${bookId}`);
      getBookTrend.add(res.timings.duration);

      const success = check(res, {
        'get: status 200': (r) => r.status === 200,
        'get: correct ID': (r) => JSON.parse(r.body).ID === bookId,
      });
      errorRate.add(!success);
    });
  }

  sleep(0.4);
}

export function teardown(data) {
  console.log('Load test completed.');
}
