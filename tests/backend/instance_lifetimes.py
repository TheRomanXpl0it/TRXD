import requests
from urllib.parse import urlparse
import socket
import time

url = 'http://localhost:1337/api'


def login(mail, password):
	s = requests.Session()
	r = s.get(f'{url}/info')
	assert r.status_code == 200, r.text
	r = s.post(f'{url}/login', json={
		"email": mail,
		"password": password,
	}, headers={"X-CSRF-Token": s.cookies.get('csrf_')})
	assert r.status_code == 200, r.text
	return s

admin = login('admin@email.com', 'testpass')
s1 = login('a@a.a', 'testpass')
s2 = login('b@b.b', 'testpass')
s3 = login('c@c.c', 'testpass')

r = s1.get(f'{url}/challenges')
assert r.status_code == 200
for chall in r.json():
	if chall['name'] == "chall-3":
		chall_id_3 = chall['id']
	elif chall['name'] == "chall-4":
		chall_id_4 = chall['id']


def update_challenge(session, chall_id, hash_domain=None, lifetime=None):
	data = {"chall_id": chall_id}
	if hash_domain is not None:
		data["hash_domain"] = hash_domain
	if lifetime is not None:
		data["lifetime"] = lifetime
	r = session.patch(f'{url}/challenges', json=data,
		headers={'X-CSRF-Token': session.cookies.get('csrf_'),})
	assert r.status_code == 200, r.text

update_challenge(admin, chall_id_3, False, 5)
update_challenge(admin, chall_id_4, False, 5)


def spawn_instance(session, chall_id):
	r = session.post(f'{url}/instances', json={
		"chall_id": chall_id,
	}, headers={"X-CSRF-Token": session.cookies.get('csrf_')})
	return r

def kill_instance(session, chall_id):
	r = session.delete(f'{url}/instances', json={
		"chall_id": chall_id,
	}, headers={"X-CSRF-Token": session.cookies.get('csrf_')})
	assert r.status_code == 200, r.text

def format_request(r: requests.Response, hash_domain):
	req = r.request
	parsed = urlparse(r.url)
	# raw = f"{req.method} {req.path_url} HTTP/1.{'0' if hash_domain else '1'}\r\n"
	raw = f"Host: {parsed.hostname}{':' + str(parsed.port) if parsed.port else ''}\r\n"
	raw += "\r\n".join(f"{k}: {v}" for k, v in req.headers.items() if k != 'Connection')
	return raw

def assert_request(r: requests.Response, hash_domain):
	req = format_request(r, hash_domain)
	resp = r.text.strip()
	for line in req.split('\r\n'):
		assert line in resp, f'{req}\n-----DIFF-----\n{resp}'

def connection_refused(port):
	try:
		r = requests.get(f'http://localhost:{port}')
		assert False, "Instance should be down: " + str(r.status_code) + "\n" + r.text
	except requests.ConnectionError:
		pass


#! CONTAINER WORKS

r = spawn_instance(s1, chall_id_3)
assert r.status_code == 200, r.text
i1 = r.json()
print(i1)
r = spawn_instance(s2, chall_id_3)
assert r.status_code == 409, r.text
r = spawn_instance(s3, chall_id_3)
assert r.status_code == 200, r.text
i3 = r.json()
print(i3)

r = requests.get(f'http://localhost:{i1["port"]}')
assert_request(r, False)
r = requests.get(f'http://localhost:{i3["port"]}')
assert_request(r, False)

kill_instance(s1, chall_id_3)
kill_instance(s3, chall_id_3)


#! COMPOSE WORKS

r = spawn_instance(s1, chall_id_4)
assert r.status_code == 200, r.text
i1 = r.json()
print(i1)
r = spawn_instance(s3, chall_id_4)
assert r.status_code == 200, r.text
i3 = r.json()
print(i3)

r = requests.get(f'http://localhost:{i1["port"]}')
assert_request(r, False)
r = requests.get(f'http://localhost:{i3["port"]}')
assert_request(r, False)

kill_instance(s1, chall_id_4)
kill_instance(s3, chall_id_4)


#! CONTAINER EXPIRES

r = spawn_instance(s1, chall_id_3)
assert r.status_code == 200, r.text
i1 = r.json()
print(i1)
r = spawn_instance(s2, chall_id_3)
assert r.status_code == 409, r.text
r = spawn_instance(s3, chall_id_3)
assert r.status_code == 200, r.text
i3 = r.json()
print(i3)

time.sleep(10)

connection_refused(i1["port"])
connection_refused(i3["port"])


#! COMPOSE EXPIRES

r = spawn_instance(s1, chall_id_4)
assert r.status_code == 200, r.text
i1 = r.json()
print(i1)
r = spawn_instance(s3, chall_id_4)
assert r.status_code == 200, r.text
i3 = r.json()
print(i3)

time.sleep(10)

connection_refused(i1["port"])
connection_refused(i3["port"])


#! CONTAINER EXPIRES, SO RECREATED

r = spawn_instance(s1, chall_id_3)
assert r.status_code == 200, r.text
i1 = r.json()
print(i1)

r = requests.get(f'http://localhost:{i1["port"]}')
assert_request(r, False)

time.sleep(10)

r = spawn_instance(s2, chall_id_3)
assert r.status_code == 200, r.text
i2 = r.json()
print(i2)

r = requests.get(f'http://localhost:{i2["port"]}')
assert_request(r, False)

kill_instance(s2, chall_id_3)


#! CONTAINER EXPIRES, SO RECREATED

r = spawn_instance(s1, chall_id_4)
assert r.status_code == 200, r.text
i1 = r.json()
print(i1)

r = requests.get(f'http://localhost:{i1["port"]}')
assert_request(r, False)

time.sleep(10)

r = spawn_instance(s2, chall_id_4)
assert r.status_code == 200, r.text
i2 = r.json()
print(i2)

r = requests.get(f'http://localhost:{i2["port"]}')
assert_request(r, False)

kill_instance(s2, chall_id_4)



# Set Hash Domain and test again

LOCALHOST = "127.0.0.1"
LOCAL_HOSTS = {}

original_getaddrinfo = socket.getaddrinfo
def custom_getaddrinfo(host, *args, **kwargs):
	if host in LOCAL_HOSTS:
		return original_getaddrinfo(LOCAL_HOSTS[host], *args, **kwargs)
	return original_getaddrinfo(host, *args, **kwargs)
socket.getaddrinfo = custom_getaddrinfo

def hash_request(resp):
	host = resp['host'].split('.')[0] + '.domain.com'
	LOCAL_HOSTS[host] = LOCALHOST
	r = requests.get(f'http://{host}')
	return r

def hash_connection_refused(resp):
	try:
		r = hash_request(resp)
		assert r.status_code == 502, "Instance should be down: " + str(r.status_code) + "\n" + r.text
	except requests.ConnectionError:
		pass


update_challenge(admin, chall_id_3, True)
update_challenge(admin, chall_id_4, True)


#! (HASH DOMAIN) CONTAINER WORKS

r = spawn_instance(s1, chall_id_3)
assert r.status_code == 200, r.text
i1 = r.json()
print(i1)
r = spawn_instance(s2, chall_id_3)
assert r.status_code == 409, r.text
r = spawn_instance(s3, chall_id_3)
assert r.status_code == 200, r.text
i3 = r.json()
print(i3)

r = hash_request(i1)
assert_request(r, True)
r = hash_request(i3)
assert_request(r, False)

kill_instance(s1, chall_id_3)
kill_instance(s3, chall_id_3)


#! (HASH DOMAIN) COMPOSE WORKS

r = spawn_instance(s1, chall_id_4)
assert r.status_code == 200, r.text
i1 = r.json()
print(i1)
r = spawn_instance(s3, chall_id_4)
assert r.status_code == 200, r.text
i3 = r.json()
print(i3)

r = hash_request(i1)
assert_request(r, True)
r = hash_request(i3)
assert_request(r, True)

kill_instance(s1, chall_id_4)
kill_instance(s3, chall_id_4)


#! (HASH DOMAIN) CONTAINER EXPIRES

r = spawn_instance(s1, chall_id_3)
assert r.status_code == 200, r.text
i1 = r.json()
print(i1)
r = spawn_instance(s2, chall_id_3)
assert r.status_code == 409, r.text
r = spawn_instance(s3, chall_id_3)
assert r.status_code == 200, r.text
i3 = r.json()
print(i3)

time.sleep(10)

hash_connection_refused(i1)
hash_connection_refused(i3)


#! (HASH DOMAIN) COMPOSE EXPIRES

r = spawn_instance(s1, chall_id_4)
assert r.status_code == 200, r.text
i1 = r.json()
print(i1)
r = spawn_instance(s3, chall_id_4)
assert r.status_code == 200, r.text
i3 = r.json()
print(i3)

time.sleep(10)

hash_connection_refused(i1)
hash_connection_refused(i3)


#! (HASH DOMAIN) CONTAINER EXPIRES, SO RECREATED

r = spawn_instance(s1, chall_id_3)
assert r.status_code == 200, r.text
i1 = r.json()
print(i1)

r = hash_request(i1)
assert_request(r, True)

time.sleep(10)

r = spawn_instance(s2, chall_id_3)
assert r.status_code == 200, r.text
i2 = r.json()
print(i2)

r = hash_request(i2)
assert_request(r, True)

kill_instance(s2, chall_id_3)


#! (HASH DOMAIN) CONTAINER EXPIRES, SO RECREATED

r = spawn_instance(s1, chall_id_4)
assert r.status_code == 200, r.text
i1 = r.json()
print(i1)

r = hash_request(i1)
assert_request(r, True)

time.sleep(10)

r = spawn_instance(s2, chall_id_4)
assert r.status_code == 200, r.text
i2 = r.json()
print(i2)

r = hash_request(i2)
assert_request(r, True)

kill_instance(s2, chall_id_4)



socket.getaddrinfo = original_getaddrinfo
