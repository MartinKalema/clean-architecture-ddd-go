/**
 * Stress Test
 *
 * Pushes the system beyond normal capacity to find breaking points.
 * Tests all endpoints with realistic user flows.
 * Ramps up to 10,000 concurrent users.
 *
 * Usage:
 *   k6 run tests/load/stress.js
 *
 * With web dashboard:
 *   K6_WEB_DASHBOARD=true k6 run tests/load/stress.js
 */

import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Rate, Trend, Counter } from 'k6/metrics';

const errorRate = new Rate('errors');
const listBooksTrend = new Trend('list_books_duration');
const addBookTrend = new Trend('add_book_duration');
const getBookTrend = new Trend('get_book_duration');
const borrowBookTrend = new Trend('borrow_book_duration');
const returnBookTrend = new Trend('return_book_duration');

export const options = {
  stages: [
    { duration: '1m', target: 500 },
    { duration: '2m', target: 500 },
    { duration: '1m', target: 2000 },
    { duration: '2m', target: 2000 },
    { duration: '1m', target: 5000 },
    { duration: '2m', target: 5000 },
    { duration: '1m', target: 10000 },
    { duration: '3m', target: 10000 },
    { duration: '1m', target: 1000 },
    { duration: '2m', target: 1000 },
    { duration: '1m', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<5000'],
    http_req_failed: ['rate<0.20'],
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
  return {};
}

export default function () {
  let bookId;

  // List books
  group('List Books', function () {
    const res = http.get(`${BASE_URL}/api/v1/books`);
    listBooksTrend.add(res.timings.duration);

    const success = check(res, {
      'list: status 200': (r) => r.status === 200,
    });
    errorRate.add(!success);
  });

  sleep(0.1);

  // Add book
  group('Add Book', function () {
    const timestamp = Date.now();
    const vuId = __VU;
    const iter = __ITER;

    const res = http.post(
      `${BASE_URL}/api/v1/books`,
      JSON.stringify({
        title: `Stress Book ${vuId}-${iter}-${timestamp}`,
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

  sleep(0.1);

  // Get book
  if (bookId) {
    group('Get Book', function () {
      const res = http.get(`${BASE_URL}/api/v1/books/${bookId}`);
      getBookTrend.add(res.timings.duration);

      const success = check(res, {
        'get: status 200': (r) => r.status === 200,
        'get: correct ID': (r) => JSON.parse(r.body).ID === bookId,
      });
      errorRate.add(!success);
    });

    sleep(0.1);

    // Borrow book
    group('Borrow Book', function () {
      const res = http.post(
        `${BASE_URL}/api/v1/books/${bookId}/borrow`,
        JSON.stringify({
          borrower_email: `user${__VU}@example.com`,
        }),
        { headers: { 'Content-Type': 'application/json' } }
      );
      borrowBookTrend.add(res.timings.duration);

      const success = check(res, {
        'borrow: status 200': (r) => r.status === 200,
      });
      errorRate.add(!success);
    });

    sleep(0.1);

    // Return book
    group('Return Book', function () {
      const res = http.post(
        `${BASE_URL}/api/v1/books/${bookId}/return`,
        null,
        { headers: { 'Content-Type': 'application/json' } }
      );
      returnBookTrend.add(res.timings.duration);

      const success = check(res, {
        'return: status 200': (r) => r.status === 200,
      });
      errorRate.add(!success);
    });
  }

  sleep(0.1);
}

export function teardown() {
  console.log('Stress test completed.');
}
