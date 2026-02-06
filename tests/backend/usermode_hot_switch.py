import requests
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

r = admin.get(f'{url}/users')
assert r.status_code == 200, r.text

r = admin.patch(f'{url}/configs',
	json={'key': 'user-mode', 'value': 'true'},
	headers={"X-CSRF-Token": admin.cookies.get('csrf_')})
assert r.status_code == 200, r.text

time.sleep(0.1)

r = admin.get(f'{url}/users')
assert r.status_code == 404, r.text

r = admin.patch(f'{url}/configs',
	json={'key': 'user-mode', 'value': 'false'},
	headers={"X-CSRF-Token": admin.cookies.get('csrf_')})
assert r.status_code == 200, r.text

time.sleep(0.1)

r = admin.get(f'{url}/users')
assert r.status_code == 200, r.text

r = admin.patch(f'{url}/configs',
	json={'key': 'user-mode', 'value': 'true'},
	headers={"X-CSRF-Token": admin.cookies.get('csrf_')})
assert r.status_code == 200, r.text

time.sleep(0.1)

r = admin.get(f'{url}/users')
assert r.status_code == 404, r.text

r = admin.patch(f'{url}/configs',
	json={'key': 'user-mode', 'value': 'false'},
	headers={"X-CSRF-Token": admin.cookies.get('csrf_')})
assert r.status_code == 200, r.text

time.sleep(0.1)

r = admin.get(f'{url}/users')
assert r.status_code == 200, r.text
