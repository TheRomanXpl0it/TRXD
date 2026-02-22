import requests
import yaml
import os

url = "https://cyberchallenge.diag.uniroma1.it/api"
email = ""
password = ""

category = ""
sorted_challs = sorted(os.listdir(category))
challenges = {
	category: sorted_challs,
}

print(challenges)
input("Press Enter to continue...")


s = requests.Session()
r = s.get(url + "/info")
assert r.status_code == 200, r.text

r = s.post(url + "/login",
	json={"email": email, "password": password},
	headers={"X-CSRF-Token": s.cookies.get('csrf_')},
)
assert r.status_code == 200, r.text

for category, challs in challenges.items():
	category = category.capitalize()
	print(f"Category: {category}")

	r = s.post(url + "/categories",
		json={"name": category},
		headers={"X-CSRF-Token": s.cookies.get('csrf_')}
	)
	print(r, r.text)

	for chall in challs:
		with open(f'{category.lower()}/{chall}/chall.yml', 'r') as f:
			chall_info = yaml.safe_load(f)
		print(chall_info)

		# Create Challenge
		r = s.post(url + f"/challenges", 
			json={
				"name": chall_info['name'],
		 		"category": chall_info['category'].capitalize(),
				"description": chall_info['description'],
				"max_points": chall_info['max_points'],
				"score_type": chall_info['score_type'],
				"type": chall_info['type'],
			},
			headers={"X-CSRF-Token": s.cookies.get('csrf_')}
		)
		print(chall, r, r.text)

		# Fetch Challenge ID
		r = s.get(url + '/challenges')
		assert r.status_code == 200, r.text
		challs_list = r.json()
		chall_id = None
		for c in challs_list:
			if c['name'] == chall_info['name']:
				chall_id = c['id']
				break
		assert chall_id is not None, f"Challenge {chall_info['name']} not found"

		# Update Challenge
		data = {'chall_id': chall_id}
		for key in ['authors', 'tags', 'host', 'port', 'conn_type']:
			if key in chall_info:
				data[key] = chall_info[key]
		deployment = chall_info.get("deployment", None)
		if deployment is not None:
			for key in ['image', 'compose', 'hash_domain', 'lifetime', 'envs', 'max_memory', 'max_cpu']:
				if key in deployment:
					data[key] = deployment[key]
		print(data)

		r = s.patch(url + f"/challenges", 
			json=data,
			headers={"X-CSRF-Token": s.cookies.get('csrf_')}
		)
		print(chall, r, r.text)

		# Load Challenge Flags
		for flag in chall_info.get("flags", []):
			data = {'chall_id': chall_id}
			data['flag'] = flag['flag']
			data['regex'] = flag.get('regex', False)
			r = s.post(url + f"/flags", 
				json=data,
				headers={"X-CSRF-Token": s.cookies.get('csrf_')}
			)
			print(chall, r, r.text)

		# Load Challenge Attachments
		data = {'chall_id': chall_id}
		for attachment in chall_info.get("attachments", []):
			path = os.path.join(f'{category.lower()}/{chall}/', attachment)
			print(path)
			with open(path, 'rb') as f:
				r = s.post(f'{url}/attachments',
					files={'attachments': (attachment, f, "application/octet-stream")},
					data={"chall_id": chall_id},
					headers={"X-Csrf-Token": s.cookies.get('csrf_')},
				)
			print(chall, r, r.text)
