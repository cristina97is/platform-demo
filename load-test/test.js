import http from 'k6/http';
import { sleep } from 'k6';

export const options = {
  vus: 20,        // виртуальные пользователи
  duration: '30s' // длительность теста
};

export default function () {

  const payload = JSON.stringify({
    user_id: 1,
    type: "bet",
    amount: 100
  });

  const params = {
    headers: {
      'Content-Type': 'application/json'
    }
  };

  http.post('http://localhost:8081/events', payload, params);

  http.get('http://localhost:8081/events');

  sleep(1);
}
