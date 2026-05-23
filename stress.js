import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 100 },   // разгон
    { duration: '1m', target: 500 },    // нагрузка
    { duration: '30s', target: 1000 },  // пик (stress)
    { duration: '1m', target: 500 },    // стабилизация
    { duration: '30s', target: 0 },     // спад
  ],
  thresholds: {
    http_req_duration: ['p(95)<200'], // SLA: p95 < 200ms
    http_req_failed: ['rate<0.01'],   // <1% ошибок
  },
};

export default function () {
  const n = 10;

  const res = http.get(`http://localhost:8080/top?limit=${n}`);

  check(res, {
    'status is 200': (r) => r.status === 200,
    'has body': (r) => r.body.length > 0,
  });

  sleep(0.1);
}