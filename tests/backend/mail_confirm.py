import requests
import os

server = os.getenv("TEST_EMAIL_SERVER", None)
port = os.getenv("TEST_EMAIL_PORT", None)
addr = os.getenv("TEST_EMAIL_ADDR", None)
passwd = os.getenv("TEST_EMAIL_PASSWD", None)
toAddr = os.getenv("TEST_TO_EMAIL_ADDR", None)

if server is None or port is None or addr is None or passwd is None or toAddr is None:
	print("Please set TEST_EMAIL_SERVER, TEST_EMAIL_PORT, TEST_EMAIL_ADDR, TEST_EMAIL_PASSWD and TEST_TO_EMAIL_ADDR environment variables to run this test.")
	exit(0)

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


admin = login('admin@email.com', 'testpass')

r = admin.patch(f'{url}/configs',
	json={'key': 'domain', 'value': 'localhost:1337'},
	headers={"X-CSRF-Token": admin.cookies.get('csrf_')})
assert r.status_code == 200, r.text
r = admin.patch(f'{url}/configs',
	json={'key': 'email-server', 'value': server},
	headers={"X-CSRF-Token": admin.cookies.get('csrf_')})
assert r.status_code == 200, r.text
r = admin.patch(f'{url}/configs',
	json={'key': 'email-port', 'value': str(port)},
	headers={"X-CSRF-Token": admin.cookies.get('csrf_')})
assert r.status_code == 200, r.text
r = admin.patch(f'{url}/configs',
	json={'key': 'email-addr', 'value': addr},
	headers={"X-CSRF-Token": admin.cookies.get('csrf_')})
assert r.status_code == 200, r.text
r = admin.patch(f'{url}/configs',
	json={'key': 'email-passwd', 'value': passwd},
	headers={"X-CSRF-Token": admin.cookies.get('csrf_')})
assert r.status_code == 200, r.text


name = "tester"
mail = toAddr
password = "testpass"

s = requests.Session()
r = s.get(f'{url}/info')
assert r.status_code == 200, r.text
r = s.post(f'{url}/register', json={
	"name": name,
	"email": mail,
	"password": password,
}, headers={"X-CSRF-Token": s.cookies.get('csrf_')})
assert r.status_code == 200, r.text


verify_url = input("Enter verification URL: ")
r = s.get(verify_url)
assert r.status_code == 200, r.text

r = s.get(f'{url}/info')
assert r.status_code == 200, r.text
print(r.text)

s2 = login(mail, password)
print("Login successful!")

