from locust import HttpUser, task, between
import random
import string

def random_string(length=10):
    return ''.join(random.choice(string.ascii_letters) for _ in range(length))

class HashUser(HttpUser):
    wait_time = between(0.1, 0.5)  # simulate user think time

    @task
    def send_hash_request(self):
        payload = {
            "input": random_string(12)
        }
        self.client.post("/hash", json=payload)