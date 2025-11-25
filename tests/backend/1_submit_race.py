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
s.get('http://localhost:1337/api/info')

r = s.post('http://localhost:1337/api/register', json={
	"name": user,
	"email": email,
	"password": "test1234",
}, headers={'X-Csrf-Token': s.cookies.get('csrf_')})

r = s.post('http://localhost:1337/api/teams/register', json={
	"name": "test-team",
	"password": "test1234",
}, headers={'X-Csrf-Token': s.cookies.get('csrf_')})

r = s.get('http://localhost:1337/api/challenges')
challs = r.json()
for c in challs:
	if c['name'] == "chall-1":
		chall_id = c['id']
		break

counter = {
	"correct": 0,
	"repeated": 0,
	"invalid": 0,
}
lock = threading.Lock()

def submit(user):
	email = user + "@test.test"
	s = requests.Session()
	s.get('http://localhost:1337/api/info')

	r = s.post('http://localhost:1337/api/register', json={
		"name": user,
		"email": email,
		"password": "test1234",
	}, headers={'X-Csrf-Token': s.cookies.get('csrf_')})

	r = s.post('http://localhost:1337/api/teams/join', json={
		"name": "test-team",
		"password": "test1234",
	}, headers={'X-Csrf-Token': s.cookies.get('csrf_')})

	r = s.post('http://localhost:1337/api/submissions', json={
		"chall_id": chall_id,
		"flag": "flag{test-1}",
	}, headers={'X-Csrf-Token': s.cookies.get('csrf_')})

	resp = r.json()
	with lock:
		if "status" in resp:
			if resp['status'] == "Correct":
				counter["correct"] += 1
			elif resp['status'] == "Repeated":
				counter["repeated"] += 1
			else:
				print(f"Unexpected status: {resp['status']}")
				counter["invalid"] += 1
		else:
			print(f"Unexpected response: {resp}")
			counter["invalid"] += 1

threads = []
for _ in range(N):
	thread = threading.Thread(target=submit, args=(gen_name(),))
	threads.append(thread)
for thread in threads:
	thread.start()
for thread in threads:
	thread.join()

for key, value in counter.items():
	print(f"{key}: {value}")

if counter["correct"] != 1:
	print("Test failed: Expected exactly one valid submission.")
	sys.exit(1)
else:
	print("Test passed: Exactly one valid submission.")
