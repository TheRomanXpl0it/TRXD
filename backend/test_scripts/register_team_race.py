import random
import string
import requests
import threading

def rand_name():
	return ''.join(random.choices(string.ascii_lowercase + string.digits, k=10))

user = rand_name()
email = user + "@test.test"
print(user, email)

s = requests.Session()
r = s.post('http://localhost:1337/register', json={
	"username": user,
	"email": email,
	"password": "test1234",
})

if r.status_code == 409:
	print("User already exists, logging in with existing user.")
	r = s.post('http://localhost:1337/login', json={
		"email": email,
		"password": "test1234",
	})

print(r)
print(r.text)

def register_team(name):
	r = s.post('http://localhost:1337/register-team', json={
		"name": name,
		"password": "testpass",
	})
	print(r.text)

threads = []
for _ in range(50):
	thread = threading.Thread(target=register_team, args=(rand_name(),))
	threads.append(thread)
for thread in threads:
	thread.start()
for thread in threads:
	thread.join()

