import http from 'k6/http';
import { check } from 'k6';

export let options = {
  vus: 100,         // number of virtual users
  iterations: 1000000, // total number of POST requests
  teardownTimeout: '3m',  // Increase teardown timeout to 2 minutes
};

export default function () {
  let timestamp = Date.now(); // Current Unix timestamp in milliseconds

  let payload = JSON.stringify({ timestamp: timestamp });

  let res = http.post('http://localhost:8080/test', payload, {
    headers: { 'Content-Type': 'application/json' },
  });

  check(res, {
    'status is 200': (r) => r.status === 200,
  });
}

// teardown is called once after all iterations are done
export function teardown() {
  let res = http.get('http://localhost:8080/shutdown', {
    timeout: '3m', 
  }); 
  check(res, {
    'shutdown status is 200': (r) => r.status === 200,
  });
}