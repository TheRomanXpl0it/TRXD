import requests
import os
import time

webhook_url = os.getenv('DISCORD_WEBHOOK_URL', None)
if webhook_url is None:
	raise RuntimeError("DISCORD_WEBHOOK_URL environment variable not set")

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

def register(name, mail, password):
	s = requests.Session()
	r = s.get(f'{url}/info')
	assert r.status_code == 200, r.text
	r = s.post(f'{url}/register', json={
		"name": name,
		"email": mail,
		"password": password,
	}, headers={"X-CSRF-Token": s.cookies.get('csrf_')})
	assert r.status_code == 200, r.text
	return s

def update_challenge(session, chall_id, hash_domain):
	r = session.patch(f'{url}/challenges',
		json={
			"chall_id": chall_id,
			"hash_domain": hash_domain,
		}, headers={'X-CSRF-Token': session.cookies.get('csrf_'),})
	assert r.status_code == 200, r.text


admin = login('admin@email.com', 'testpass')
player = register('player', 'player@test.com', 'testpass')

r = player.post(f'{url}/teams/register',
	json={
		"name": "@everyone team-webhook",
		"password": "testpass",
	},
	headers={"X-CSRF-Token": player.cookies.get('csrf_')})
assert r.status_code == 200, r.text

r = admin.patch(f'{url}/configs',
	json={'key': 'discord-webhook', 'value': webhook_url},
	headers={"X-CSRF-Token": admin.cookies.get('csrf_')})
assert r.status_code == 200, r.text

r = admin.post(f'{url}/challenges',
	json={
		"name": "chall-webhook",
		"category": "cat-1",
		"description": "test challenge for webhook",
		"type": "Normal",
		"max_points": 500,
		"score_type": "Dynamic",
	},
	headers={"X-CSRF-Token": admin.cookies.get('csrf_')},
)
assert r.status_code == 200, r.text

r = admin.get(f'{url}/challenges')
assert r.status_code == 200, r.text
for chall in r.json():
	if chall['name'] == "chall-webhook":
		chall_id = chall['id']

r = admin.post(f'{url}/flags', 
	json={
		"chall_id": chall_id,
		"flag": "flag{webhook-test-flag}",
		"regex": False,
	},
	headers={"X-CSRF-Token": admin.cookies.get('csrf_')},
)
assert r.status_code == 200, r.text

r = admin.patch(f'{url}/challenges',
	json={"chall_id": chall_id, "hidden": False},
	headers={'X-CSRF-Token': admin.cookies.get('csrf_'),})
assert r.status_code == 200, r.text


r = player.post(f'{url}/submissions',
	json={
		"chall_id": chall_id,
		"flag": "flag{webhook-test-flag}",
	}, 
	headers={"X-CSRF-Token": player.cookies.get('csrf_')})
assert r.status_code == 200, r.text

time.sleep(6) # Wait for the webhook message to be sent

print("Discord webhook test passed.")
