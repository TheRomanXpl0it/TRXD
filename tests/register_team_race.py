import sys
import random
import string
import requests
import threading

try:
	N = int(sys.argv[1])
except:
	N = 50
print(f"Running with {N} threads")

def rand_name():
	return ''.join(random.choices(string.ascii_lowercase + string.digits, k=10))

user = rand_name()
email = user + "@test.test"
# print(user, email)

s = requests.Session()
r = s.post('http://localhost:1337/api/register', json={
	"username": user,
	"email": email,
	"password": "test1234",
})

if r.status_code == 409:
	print("User already exists, logging in with existing user.")
	r = s.post('http://localhost:1337/api/login', json={
		"email": email,
		"password": "test1234",
	})

# print(r)
# print(r.text)

counter = {
	"valid": 0,
	"Already in a team": 0,
	"Error registering team": 0,
	"invalid": 0,
}
lock = threading.Lock()

def register_team(name):
	r = s.post('http://localhost:1337/api/teams', json={
		"name": name,
		"password": "testpass",
	})
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
			elif res["error"] == "Error registering team":
				counter["Error registering team"] += 1
			else:
				print(f"Unexpected error: {res['error']}")
				counter["invalid"] += 1
		else:
			print(f"Unexpected result: {res}")
			counter["invalid"] += 1

threads = []
for _ in range(N):
	thread = threading.Thread(target=register_team, args=(rand_name(),))
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
