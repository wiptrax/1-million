import http from 'k6/http';
import { check } from 'k6';

export let options = {
  vus: 100,         // number of virtual users
  iterations: 1000000, // total number of POST requests
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

