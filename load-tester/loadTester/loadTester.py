import time
from locust import HttpUser, task, between
from locust import events
import json
import random

userID = 1

class testUserAdd(HttpUser):
    wait_time = between(0.1, 2)


    @task(3)
    def getUser(self):
        self.client.get("/")
        pass

    @task
    def addUser(self):    
        with self.client.post(f"/insert?name=testUser{self.userID}&email=testEmail{self.userID}&password=secret", {}, catch_response=True) as response:
            if response.status_code != 200:
                response.failure(f"Got wrong response: {response.text}")
            elif response.elapsed.total_seconds() > 5.0:
                response.failure("Request took too long")
        self.userID = random.random()

    def on_start(self):
        global userID
        self.userID = 1

        
    