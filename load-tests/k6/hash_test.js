import http from 'k6/http';
import { sleep } from 'k6';

export let options = {
  stages: [
    { duration: '30s', target: 50 },   // ramp to 50 VUs
    { duration: '30s', target: 100 },  // ramp to 100 VUs
    { duration: '30s', target: 500 },  // ramp to 500 VUs
    { duration: '60s', target: 1000 },  // sustain at 500 VUs
    { duration: '20s', target: 0 },    // ramp down
  ],
};

export default function () {
  const url = 'http://localhost:30080/hash';
  const payload = JSON.stringify({ input: "hello world" });

  http.post(url, payload, { headers: { 'Content-Type': 'application/json' }});

  sleep(0.05); // even smaller idle between requests to handle higher load
}