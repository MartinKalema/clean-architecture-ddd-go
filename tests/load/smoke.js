/**
 * Smoke Test
 *
 * Quick sanity check to verify the API is working.
 * Runs with minimal load to catch obvious issues.
 *
 * Usage:
 *   k6 run tests/load/smoke.js
 *
 * With web dashboard:
 *   K6_WEB_DASHBOARD=true k6 run tests/load/smoke.js
 */

import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  vus: 1,              // 1 virtual user
  duration: '10s',     // 10 seconds
  thresholds: {
    http_req_duration: ['p(95)<500'],  // 95% of requests < 500ms
    http_req_failed: ['rate<0.01'],    // Less than 1% errors
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export default function() {
  // List books
  const listRes = http.get(`${BASE_URL}/api/v1/books`);
  check(listRes, {
    'list: status 200': (r) => r.status === 200,
    'list: is array': (r) => JSON.parse(r.body).Books !== undefined,
  });

  sleep(0.5);

  // Add a book
  const timestamp = new Date().getTime();
  const addRes = http.post(
    `${BASE_URL}/api/v1/books`,
    JSON.stringify({
      title: `Smoke Test Book ${timestamp}`,
      author: 'Test Author',
    }),
    { headers: { 'Content-Type': 'application/json' } }
  );
  check(addRes, {
    'add: status 201': (r) => r.status === 201,
    'add: has ID': (r) => JSON.parse(r.body).ID !== undefined,
  });

  sleep(0.5);
}
