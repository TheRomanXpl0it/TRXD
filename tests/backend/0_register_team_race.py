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
# print(user, email)

s = requests.Session()
s.get('http://localhost:1337/api/info')

r = s.post('http://localhost:1337/api/register', json={
	"name": user,
	"email": email,
	"password": "test1234",
}, headers={'X-Csrf-Token': s.cookies.get('csrf_')})

if r.status_code == 409:
	print("User already exists, logging in with existing user.")
	r = s.post('http://localhost:1337/api/login', json={
		"email": email,
		"password": "test1234",
	}, headers={'X-Csrf-Token': s.cookies.get('csrf_')})


counter = {
	"valid": 0,
	"Already in a team": 0,
	"invalid": 0,
}
lock = threading.Lock()

def register_team(name):
	r = s.post('http://localhost:1337/api/teams/register', json={
		"name": name,
		"password": "testpass",
	}, headers={'X-Csrf-Token': s.cookies.get('csrf_')})
	if r.status_code != 200:
		res = r.json()
	else:
		res = None
	with lock:
		if res is None:
			counter["valid"] += 1
		elif "error" in res:
			if res["error"] == "Already in a team":
				counter["Already in a team"] += 1
			else:
				print(f"Unexpected error: {res['error']}")
				counter["invalid"] += 1
		else:
			print(f"Unexpected result: {res}")
			counter["invalid"] += 1

threads = []
for _ in range(N):
	thread = threading.Thread(target=register_team, args=(gen_name(),))
	threads.append(thread)
for thread in threads:
	thread.start()
for thread in threads:
	thread.join()

for key, value in counter.items():
	print(f"{key}: {value}")

if counter["valid"] != 1:
	print("Test failed: Expected exactly one valid team registration.")
	sys.exit(1)
else:
	print("Test passed: Exactly one valid team registration.")
