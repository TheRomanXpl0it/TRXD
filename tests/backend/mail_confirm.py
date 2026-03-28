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

def register(name=None, mail=None, password=None, token=None):
	data = {}
	if name is not None:
		data["name"] = name
	if mail is not None:
		data["email"] = mail
	if password is not None:
		data["password"] = password
	if token is not None:
		data["token"] = token

	s = requests.Session()
	r = s.get(f'{url}/info')
	assert r.status_code == 200, r.text
	r = s.post(f'{url}/register', json=data,
		headers={"X-CSRF-Token": s.cookies.get('csrf_')})
	assert r.status_code == 200, r.text
	return s

def change_conf(s, key, value):
	r = s.patch(f'{url}/configs',
		json={'key': key, 'value': value},
		headers={"X-CSRF-Token": s.cookies.get('csrf_')})
	assert r.status_code == 200, r.text

admin = login('admin@email.com', 'testpass')

change_conf(admin, 'domain', 'localhost:1337')
change_conf(admin, 'email-server', server)
change_conf(admin, 'email-port', str(port))
change_conf(admin, 'email-addr', addr)
change_conf(admin, 'email-passwd', passwd)


name = "tester"
mail = toAddr
password = "testpass"

s = register(mail=mail)


token = input("Enter verification TOKEN: ").strip()
s = register(name=name, password=password, token=token)

r = s.get(f'{url}/info')
assert r.status_code == 200, r.text
print(r.text)

s2 = login(mail, password)
print("Login successful!")

