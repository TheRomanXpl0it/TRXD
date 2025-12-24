import docker
import requests
from urllib.parse import urlparse
import socket
import json
import os


url = 'http://localhost:1337/api'

project_name = os.getenv('PROJECT_NAME', 'trxd')
client = docker.from_env()


def login(mail, password):
	s = requests.Session()
	r = s.get(f'{url}/info')
	assert r.status_code == 200, r.text
	r = s.post(f'{url}/login', json={
		"email": mail,
		"password": password,
	}, headers={"X-Csrf-Token": s.cookies.get('csrf_')})
	assert r.status_code == 200, r.text
	return s

admin = login('admin@email.com', 'testpass')
s1 = login('a@a.a', 'testpass')

r = s1.get(f'{url}/info')
assert r.status_code == 200
team_id = r.json()['team_id']

r = s1.get(f'{url}/challenges')
assert r.status_code == 200
for chall in r.json():
	if chall['name'] == "chall-3":
		chall_id_3 = chall['id']
	elif chall['name'] == "chall-4":
		chall_id_4 = chall['id']

def update_config(session, key, value):
	r = session.patch(f'{url}/configs', json={
		"key": key,
		"value": value,
	}, headers={"X-Csrf-Token": session.cookies.get('csrf_')})
	assert r.status_code == 200, r.text

def update_challenge(session, chall_id, hash_domain=None, image=None, compose=None):
	fields={"chall_id": chall_id}
	if hash_domain is not None:
		fields['hash_domain'] = hash_domain
	if image is not None:
		fields['image'] = image
	if compose is not None:
		fields['compose'] = compose
	r = session.patch(f'{url}/challenges', json=fields,
		headers={'X-Csrf-Token': session.cookies.get('csrf_')})
	assert r.status_code == 200, r.text

def spawn_instance(session, chall_id):
	r = session.post(f'{url}/instances', json={
		"chall_id": chall_id,
	}, headers={"X-Csrf-Token": session.cookies.get('csrf_')})
	return r

def spawn_good_instance(session, chall_id):
	r = spawn_instance(session, chall_id)
	assert r.status_code == 200, r.text
	i = r.json()
	print(i)
	return i

def kill_container_by_name(name):
	res = client.containers.list(filters={"name": name})
	assert len(res) == 1, res
	res[0].kill()

def remove_container_by_name(name):
	res = client.containers.list(filters={"name": name})
	assert len(res) == 1, res
	res[0].remove(force=True)

def fail_connection(url):
	try:
		r = requests.get(url)
		assert False, f'{r} -> {r.text}'
	except requests.exceptions.ConnectionError as e:
		pass

def kill_instance(session, chall_id):
	r = session.delete(f'{url}/instances', json={
		"chall_id": chall_id,
	}, headers={"X-Csrf-Token": session.cookies.get('csrf_')})
	return r

def kill_good_instance(session, chall_id):
	r = kill_instance(session, chall_id)
	assert r.status_code == 200, r.text

def format_request(r: requests.Response, hash_domain):
	req = r.request
	parsed = urlparse(r.url)
	raw = f"{req.method} {req.path_url} HTTP/1.{'0' if hash_domain else '1'}\r\n"
	raw += f"Host: {parsed.hostname}{':' + str(parsed.port) if parsed.port else ''}\r\n"
	raw += "\r\n".join(f"{k}: {v}" for k, v in req.headers.items() if k != 'Connection')
	return raw

def assert_request(r: requests.Response, hash_domain):
	req = format_request(r, hash_domain)
	resp = r.text.strip()
	for line in req.split('\r\n'):
		assert line in resp, f'{req}\n-----DIFF-----\n{resp}'

LOCALHOST = "127.0.0.1"
LOCAL_HOSTS = {}
original_getaddrinfo = socket.getaddrinfo
def custom_getaddrinfo(host, *args, **kwargs):
	if host in LOCAL_HOSTS:
		return original_getaddrinfo(LOCAL_HOSTS[host], *args, **kwargs)
	return original_getaddrinfo(host, *args, **kwargs)
socket.getaddrinfo = custom_getaddrinfo

def get_network_by_name(name):
	res = client.networks.list(filters={'name': name})
	assert len(res) == 1, res
	return res[0]

def net_disconnect(net, hash_name):
	res = client.containers.list(filters={'name': project_name+'-nginx-1'})
	assert len(res) == 1, res
	net.disconnect(res[0], force=True)
	res = client.containers.list(filters={'name': f'chall_{hash_name}'})
	assert len(res) == 1, res
	net.disconnect(res[0], force=True)

def hash_request(hash_name):
	host = hash_name + '.domain.com'
	LOCAL_HOSTS[host] = LOCALHOST
	r = requests.get(f'http://{host}')
	return r

def fail_hash_request(hash_name):
	host = hash_name + '.domain.com'
	LOCAL_HOSTS[host] = LOCALHOST
	try:
		r = requests.get(f'http://{host}')
		assert False, f'{r} -> {r.text}'
	except requests.exceptions.ConnectionError as e:
		pass


## DISABLE HASH DOMAIN
update_challenge(admin, chall_id_3, hash_domain=False)
update_challenge(admin, chall_id_4, hash_domain=False)

#! KILL CONTAINER
i1 = spawn_good_instance(s1, chall_id_3)
kill_container_by_name(f"chall_{chall_id_3}_{team_id}")
fail_connection(f'http://localhost:{i1["port"]}')
kill_good_instance(s1, chall_id_3)

#! REMOVE CONTAINER
i1 = spawn_good_instance(s1, chall_id_3)
remove_container_by_name(f"chall_{chall_id_3}_{team_id}")
fail_connection(f'http://localhost:{i1["port"]}')
kill_good_instance(s1, chall_id_3)

#! KILL COMPOSE
i1 = spawn_good_instance(s1, chall_id_4)
kill_container_by_name(f"chall_{chall_id_4}_{team_id}")
fail_connection(f'http://localhost:{i1["port"]}')
kill_good_instance(s1, chall_id_4)

#! REMOVE COMPOSE
i1 = spawn_good_instance(s1, chall_id_4)
remove_container_by_name(f"chall_{chall_id_4}_{team_id}")
fail_connection(f'http://localhost:{i1["port"]}')
kill_good_instance(s1, chall_id_4)

#! CONTAINER ALREADY EXISTS
update_config(admin, "min-port", "10000")
update_config(admin, "max-port", "10000")
container = client.containers.run(image="echo-server:latest", name=f"chall_{chall_id_3}_{team_id}", ports={'1337': '10000'}, detach=True)
print(container.id)
i1 = spawn_good_instance(s1, chall_id_3)
r = requests.get(f'http://localhost:{i1["port"]}')
assert_request(r, False)
kill_good_instance(s1, chall_id_3)
try:
	stats = list(container.stats())
	assert json.loads(stats[0].decode())['message'].startswith("No such container"), stats
except docker.errors.NotFound:
	pass
update_config(admin, "min-port", "10000")
update_config(admin, "max-port", "20000")

#! CONTAINER ALREADY EXISTS (COMPOSE)
container = client.containers.run(image="echo-server:latest", name=f"chall_{chall_id_4}_{team_id}", detach=True)
print(container.id)
r = spawn_instance(s1, chall_id_4)
assert r.status_code == 500 and r.text == '{"error":"Error creating instance"}', r.text
r = kill_instance(s1, chall_id_4)
assert r.status_code == 404, r.text
container.remove(force=True)

#! COMPOSE ALREADY EXISTS
compose = f'''
services:
  chall:
    image: echo-server:latest
    container_name: "chall_{chall_id_4}_{team_id}"
    ports:
      - "10000:1337"
'''
try:
	os.mkdir(f'/tmp/chall_{chall_id_4}_{team_id}')
except FileExistsError:
	pass
with open(f'/tmp/chall_{chall_id_4}_{team_id}/compose.yml', 'w') as f:
	f.write(compose)
res = os.system(f'cd /tmp/chall_{chall_id_4}_{team_id} && docker compose up -d')
assert res == 0, res
i1 = spawn_good_instance(s1, chall_id_4)
fail_connection('http://localhost:10000')
r = requests.get(f'http://localhost:{i1["port"]}')
assert_request(r, False)
kill_good_instance(s1, chall_id_4)

## ENABLE HASH DOMAIN
update_challenge(admin, chall_id_3, hash_domain=True)
update_challenge(admin, chall_id_4, hash_domain=True)

#! DISCONNECT NETWORK
i1 = spawn_good_instance(s1, chall_id_3)
hash_name = i1['host'].split('.')[0]
net = get_network_by_name(f'net_{chall_id_3}_{team_id}')
net_disconnect(net, hash_name)
r = hash_request(hash_name)
assert r.status_code == 502, r.text
kill_good_instance(s1, chall_id_3)

#! DISCONNECT NETWORK & KILL CONTAINER
i1 = spawn_good_instance(s1, chall_id_3)
hash_name = i1['host'].split('.')[0]
net = get_network_by_name(f'net_{chall_id_3}_{team_id}')
net_disconnect(net, hash_name)
kill_container_by_name(f'chall_{hash_name}')
r = hash_request(hash_name)
assert r.status_code == 502, r.text
kill_good_instance(s1, chall_id_3)

#! DISCONNECT NETWORK & REMOVE CONTAINER
i1 = spawn_good_instance(s1, chall_id_3)
hash_name = i1['host'].split('.')[0]
net = get_network_by_name(f'net_{chall_id_3}_{team_id}')
net_disconnect(net, hash_name)
remove_container_by_name(f'chall_{hash_name}')
r = hash_request(hash_name)
assert r.status_code == 502, r.text
kill_good_instance(s1, chall_id_3)

#! DISCONNECT NETWORK & KILL CONTAINER (COMPOSE)
i1 = spawn_good_instance(s1, chall_id_4)
hash_name = i1['host'].split('.')[0]
net = get_network_by_name(f'net_{chall_id_4}_{team_id}')
net_disconnect(net, hash_name)
kill_container_by_name(f'chall_{hash_name}')
r = hash_request(hash_name)
assert r.status_code == 502, r.text
kill_good_instance(s1, chall_id_4)

#! DISCONNECT NETWORK & REMOVE CONTAINER (COMPOSE)
i1 = spawn_good_instance(s1, chall_id_4)
hash_name = i1['host'].split('.')[0]
net = get_network_by_name(f'net_{chall_id_4}_{team_id}')
net_disconnect(net, hash_name)
remove_container_by_name(f'chall_{hash_name}')
r = hash_request(hash_name)
assert r.status_code == 502, r.text
kill_good_instance(s1, chall_id_4)

#! REMOVE NETWORK
i1 = spawn_good_instance(s1, chall_id_3)
hash_name = i1['host'].split('.')[0]
net = get_network_by_name(f'net_{chall_id_3}_{team_id}')
net_disconnect(net, hash_name)
net.remove()
r = hash_request(hash_name)
assert r.status_code == 502, r.text
kill_good_instance(s1, chall_id_3)

#! REMOVE NETWORK & KILL CONTAINER
i1 = spawn_good_instance(s1, chall_id_3)
hash_name = i1['host'].split('.')[0]
net = get_network_by_name(f'net_{chall_id_3}_{team_id}')
net_disconnect(net, hash_name)
net.remove()
kill_container_by_name(f'chall_{hash_name}')
r = hash_request(hash_name)
assert r.status_code == 502, r.text
kill_good_instance(s1, chall_id_3)

#! REMOVE NETWORK & REMOVE CONTAINER
i1 = spawn_good_instance(s1, chall_id_3)
hash_name = i1['host'].split('.')[0]
net = get_network_by_name(f'net_{chall_id_3}_{team_id}')
net_disconnect(net, hash_name)
net.remove()
remove_container_by_name(f'chall_{hash_name}')
r = hash_request(hash_name)
assert r.status_code == 502, r.text
kill_good_instance(s1, chall_id_3)

#! REMOVE NETWORK & KILL CONTAINER (COMPOSE)
i1 = spawn_good_instance(s1, chall_id_4)
hash_name = i1['host'].split('.')[0]
net = get_network_by_name(f'net_{chall_id_4}_{team_id}')
net_disconnect(net, hash_name)
net.remove()
kill_container_by_name(f'chall_{hash_name}')
r = hash_request(hash_name)
assert r.status_code == 502, r.text
kill_good_instance(s1, chall_id_4)

#! REMOVE NETWORK & REMOVE CONTAINER (COMPOSE)
i1 = spawn_good_instance(s1, chall_id_4)
hash_name = i1['host'].split('.')[0]
net = get_network_by_name(f'net_{chall_id_4}_{team_id}')
net_disconnect(net, hash_name)
net.remove()
remove_container_by_name(f'chall_{hash_name}')
r = hash_request(hash_name)
assert r.status_code == 502, r.text
kill_good_instance(s1, chall_id_4)

#! SAME NAME NETWORK
client.networks.create(name=f'net_{chall_id_3}_{team_id}')
i1 = spawn_good_instance(s1, chall_id_3)
r = hash_request(i1['host'].split('.')[0])
assert_request(r, True)
kill_good_instance(s1, chall_id_3)

#! SAME NAME NETWORK (COMPOSE)
client.networks.create(name=f'net_{chall_id_4}_{team_id}')
i1 = spawn_good_instance(s1, chall_id_4)
r = hash_request(i1['host'].split('.')[0])
assert_request(r, True)
kill_good_instance(s1, chall_id_4)

#! KILL NGINX
i1 = spawn_good_instance(s1, chall_id_3)
kill_container_by_name(project_name+'-nginx-1')
fail_hash_request(i1['host'].split('.')[0])
kill_good_instance(s1, chall_id_3)
res = os.system(f'cd .. && docker compose -p {project_name} up -d nginx')
assert res == 0, res

#! REMOVE NGINX
i1 = spawn_good_instance(s1, chall_id_3)
remove_container_by_name(project_name+'-nginx-1')
fail_hash_request(i1['host'].split('.')[0])
kill_good_instance(s1, chall_id_3)
res = os.system(f'cd .. && docker compose -p {project_name} up -d nginx')
assert res == 0, res

#! KILL NGINX (COMPOSE)
i1 = spawn_good_instance(s1, chall_id_4)
kill_container_by_name(project_name+'-nginx-1')
fail_hash_request(i1['host'].split('.')[0])
kill_good_instance(s1, chall_id_4)
res = os.system(f'cd .. && docker compose -p {project_name} up -d nginx')
assert res == 0, res

#! REMOVE NGINX (COMPOSE)
i1 = spawn_good_instance(s1, chall_id_4)
remove_container_by_name(project_name+'-nginx-1')
fail_hash_request(i1['host'].split('.')[0])
kill_good_instance(s1, chall_id_4)
res = os.system(f'cd .. && docker compose -p {project_name} up -d nginx')
assert res == 0, res

#! KILL NGINX (ROUND 2: CACHE RESET)
i1 = spawn_good_instance(s1, chall_id_3)
kill_container_by_name(project_name+'-nginx-1')
fail_hash_request(i1['host'].split('.')[0])
kill_good_instance(s1, chall_id_3)
res = os.system(f'cd .. && docker compose -p {project_name} up -d nginx')
assert res == 0, res

#! REMOVE NGINX (ROUND 2: CACHE RESET)
i1 = spawn_good_instance(s1, chall_id_3)
remove_container_by_name(project_name+'-nginx-1')
fail_hash_request(i1['host'].split('.')[0])
kill_good_instance(s1, chall_id_3)
res = os.system(f'cd .. && docker compose -p {project_name} up -d nginx')
assert res == 0, res

#! KILL NGINX (COMPOSE) (ROUND 2: CACHE RESET)
i1 = spawn_good_instance(s1, chall_id_4)
kill_container_by_name(project_name+'-nginx-1')
fail_hash_request(i1['host'].split('.')[0])
kill_good_instance(s1, chall_id_4)
res = os.system(f'cd .. && docker compose -p {project_name} up -d nginx')
assert res == 0, res

#! REMOVE NGINX (COMPOSE) (ROUND 2: CACHE RESET)
i1 = spawn_good_instance(s1, chall_id_4)
remove_container_by_name(project_name+'-nginx-1')
fail_hash_request(i1['host'].split('.')[0])
kill_good_instance(s1, chall_id_4)
res = os.system(f'cd .. && docker compose -p {project_name} up -d nginx')
assert res == 0, res

#! REMOVE NGINX BEFOREHAND
remove_container_by_name(project_name+'-nginx-1')
r = spawn_instance(s1, chall_id_3)
assert r.status_code == 500, r.text
r = kill_instance(s1, chall_id_3)
assert r.status_code == 404, r.text
res = os.system(f'cd .. && docker compose -p {project_name} up -d nginx')
assert res == 0, res

#! REMOVE NGINX BEFOREHAND (COMPOSE)
remove_container_by_name(project_name+'-nginx-1')
r = spawn_instance(s1, chall_id_4)
assert r.status_code == 500, r.text
r = kill_instance(s1, chall_id_4)
assert r.status_code == 404, r.text
res = os.system(f'cd .. && docker compose -p {project_name} up -d nginx')
assert res == 0, res

#! REMOVE NGINX BEFOREHAND (WITH NON EMPTY CACHE)
r = spawn_instance(s1, chall_id_3)
assert r.status_code == 200, r.text
r = kill_instance(s1, chall_id_3)
assert r.status_code == 200, r.text

remove_container_by_name(project_name+'-nginx-1')
res = os.system(f'cd .. && docker compose -p {project_name} up -d nginx')
assert res == 0, res

r = spawn_instance(s1, chall_id_3)
assert r.status_code == 200, r.text
r = kill_instance(s1, chall_id_3)
assert r.status_code == 200, r.text

r = spawn_instance(s1, chall_id_3)
assert r.status_code == 200, r.text
r = kill_instance(s1, chall_id_3)
assert r.status_code == 200, r.text

## BREAK IMAGE / COMPOSE

#! WRONG IMAGE
update_challenge(admin, chall_id_3, image='nonexistentrepo/nonexistentimage:latest')
r = spawn_instance(s1, chall_id_3)
assert r.status_code == 500 and r.text == '{"error":"Error creating instance"}', r.text
r = kill_instance(s1, chall_id_3)
assert r.status_code == 404, r.text

#! WRONG COMPOSE (YAML parse error) (invalid mapping)
update_challenge(admin, chall_id_4, compose='nonsense-content')
r = spawn_instance(s1, chall_id_4)
assert r.status_code == 500 and r.text == '{"error":"Error creating instance"}', r.text
r = kill_instance(s1, chall_id_4)
assert r.status_code == 404, r.text

#! WRONG COMPOSE (compose parse error) (additional properties not allowed)
update_challenge(admin, chall_id_4, compose='a: b')
r = spawn_instance(s1, chall_id_4)
assert r.status_code == 500 and r.text == '{"error":"Error creating instance"}', r.text
r = kill_instance(s1, chall_id_4)
assert r.status_code == 404, r.text

#! WRONG COMPOSE (run error) (invalid image) (pull access denied for nonexistentrepo/nonexistentimage)
compose = '''
services:
  chall:
    image: nonexistentrepo/nonexistentimage:latest
    container_name: ${CONTAINER_NAME}
    ports:
      - "${INSTANCE_PORT}:1337"
    environment:
      - ECHO_MESSAGE=Hello from app
      - INSTANCE_PORT=${INSTANCE_PORT}
      - INSTANCE_HOST=${INSTANCE_HOST}
'''
update_challenge(admin, chall_id_4, compose=compose)
r = spawn_instance(s1, chall_id_4)
assert r.status_code == 500 and r.text == '{"error":"Error creating instance"}', r.text
r = kill_instance(s1, chall_id_4)
assert r.status_code == 404, r.text
