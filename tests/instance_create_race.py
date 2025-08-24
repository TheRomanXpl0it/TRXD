import os
import sys
import requests
import threading

N = int(os.getenv("TEST_WORKERS", 50))
print(f"Running with {N} threads")

COUNTER = 0
COUNTER_LOCK = threading.Lock()
def gen_name():
	global COUNTER
	with COUNTER_LOCK:
		num = COUNTER
		COUNTER += 1
	return f"test-user-{num}"


user = gen_name()
email = user + "@test.test"
s = requests.Session()
r = s.post('http://localhost:1337/api/register', json={
	"name": user,
	"email": email,
	"password": "test1234",
})

r = s.post('http://localhost:1337/api/teams/register', json={
	"name": "test-team"+user,
	"password": "test1234",
})

r = s.get('http://localhost:1337/api/challenges')
challs = r.json()
for c in challs:
	if c['name'] == "chall-3":
		chall_id = c['id']
		break

counter = {
	"instanced": 0,
	"already_active": 0,
	"invalid": 0,
}
lock = threading.Lock()

def instance(user):
	s = requests.Session()

	r = s.post('http://localhost:1337/api/login', json={
		"email": email,
		"password": "test1234",
	})

	r = s.post('http://localhost:1337/api/instances', json={
		"chall_id": chall_id,
	})

	resp = r.json()
	with lock:
		if "timeout" in resp:
			counter["instanced"] += 1
		elif resp['error'] == "Already an active instance":
			counter["already_active"] += 1
		else:
			print(f"Unexpected response: {resp}")
			counter["invalid"] += 1

threads = []
for _ in range(N):
	thread = threading.Thread(target=instance, args=(gen_name(),))
	threads.append(thread)
for thread in threads:
	thread.start()
for thread in threads:
	thread.join()

for key, value in counter.items():
	print(f"{key}: {value}")

if counter["instanced"] != 1:
	print("Test failed: Expected exactly one valid instance.")
	sys.exit(1)
else:
	print("Test passed: Exactly one valid instance.")
