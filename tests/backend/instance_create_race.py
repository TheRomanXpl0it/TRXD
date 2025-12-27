import os
import sys
import requests
import threading
from random import randint

url = 'http://localhost:1337/api'

N = int(os.getenv("TEST_WORKERS", 25))
if N > 25:
	print("Limiting to 25 threads")
	N = 25
print(f"Running with {N} threads")

id = randint(1, 100)

COUNTER = 0
COUNTER_LOCK = threading.Lock()
def gen_name():
	global COUNTER
	with COUNTER_LOCK:
		num = COUNTER
		COUNTER += 1
	return f"test-user-{id}-{num}"


s = requests.Session()
s.get(f'{url}/info') # TODO: use url and fmt string

r = s.post(f'{url}/login', json={
	"email": 'admin@email.com',
	"password": "testpass",
}, headers={'X-CSRF-Token': s.cookies.get('csrf_')})

r = s.get(f'{url}/challenges')
challs = r.json()
for c in challs:
	if c['name'] == "chall-3":
		chall_id = c['id']
		break

r = s.patch(f'{url}/challenges',
	json={
		"chall_id": chall_id,
		"hash_domain": False,
	},
	headers={'X-CSRF-Token': s.cookies.get('csrf_')},
)


sessions = [None] * N
for i in range(N):
	user = gen_name()
	email = user + "@test.test"
	s = requests.Session()
	s.get(f'{url}/info')

	r = s.post(f'{url}/register', json={
		"name": user,
		"email": email,
		"password": "test1234",
	}, headers={'X-CSRF-Token': s.cookies.get('csrf_')})

	if r.status_code == 200:
		r = s.post(f'{url}/teams/register', json={
			"name": "test-team"+user,
			"password": "test1234",
		}, headers={'X-CSRF-Token': s.cookies.get('csrf_')})
	else:
		r = s.post(f'{url}/login', json={
			"email": email,
			"password": "test1234",
		}, headers={'X-CSRF-Token': s.cookies.get('csrf_')})

	sessions[i] = s


counter = {
	"instanced": 0,
	"invalid": 0,
}
lock = threading.Lock()

def instance(i):
	s = sessions[i]
	r = s.post(f'{url}/instances', json={
		"chall_id": chall_id,
	}, headers={'X-CSRF-Token': s.cookies.get('csrf_')})
	resp = r.json()

	r = s.delete(f'{url}/instances', json={
		"chall_id": chall_id,
	}, headers={'X-CSRF-Token': s.cookies.get('csrf_')})
	success = r.status_code == 200
	del_resp = r.text

	with lock:
		if "timeout" in resp:
			if success:
				counter["instanced"] += 1
			else:
				print(f"Failed to delete instance: {del_resp}")
				counter["invalid"] += 1
		else:
			print(f"Unexpected response: {resp}")
			counter["invalid"] += 1

threads = []
for i in range(N):
	thread = threading.Thread(target=instance, args=(i,))
	threads.append(thread)
for thread in threads:
	thread.start()
for thread in threads:
	thread.join()

for key, value in counter.items():
	print(f"{key}: {value}")

if counter["instanced"] != N:
	print(f"Test failed: Expected exactly {N} valid instance.")
	sys.exit(1)
else:
	print(f"Test passed: Exactly {N} valid instance.")
