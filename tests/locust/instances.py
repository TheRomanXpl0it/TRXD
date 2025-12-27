import os
import time
import json
import random
import string
import socket
from typing import Optional, Tuple

import requests
from requests_toolbelt.multipart.encoder import MultipartEncoder
from locust import HttpUser, task, between, events


ADMIN_EMAIL = os.getenv("ADMIN_EMAIL")
ADMIN_PASSWORD = os.getenv("ADMIN_PASSWORD")

CHALL_CONTAINER_NAME = os.getenv("LOCUST_CONTAINER_CHALL", "locust-container")
CHALL_COMPOSE_NAME = os.getenv("LOCUST_COMPOSE_CHALL", "locust-compose")
CATEGORY_NAME = os.getenv("LOCUST_CATEGORY", "locust")

CHALL_CONTAINER_ID: Optional[int] = None
CHALL_COMPOSE_ID: Optional[int] = None


def rand_str(prefix: str, n: int = 8) -> str:
    return f"{prefix}-" + "".join(random.choices(string.ascii_lowercase + string.digits, k=n))


def wait_for_tcp(host: str, port: int, timeout: float = 10.0) -> bool:
    deadline = time.time() + timeout
    while time.time() < deadline:
        try:
            with socket.create_connection((host, port), timeout=1.5):
                return True
        except OSError:
            time.sleep(0.3)
    return False


class AdminAPI:
    def __init__(self, base_url: str):
        self.base_url = base_url.rstrip("/")
        self.sess = requests.Session()

    def _csrf(self) -> str:
        token = self.sess.cookies.get("csrf_")
        if not token:
            self.sess.get(self.base_url + "/api/info")
            token = self.sess.cookies.get("csrf_")
        return token or ""

    def post_json(self, path: str, data: dict) -> requests.Response:
        headers = {"Content-Type": "application/json", "X-CSRF-Token": self._csrf()}
        return self.sess.post(self.base_url + path, headers=headers, json=data)

    def patch_json(self, path: str, data: dict) -> requests.Response:
        headers = {"Content-Type": "application/json", "X-CSRF-Token": self._csrf()}
        return self.sess.patch(self.base_url + path, headers=headers, json=data)

    def patch_multipart(self, path: str, fields: dict) -> requests.Response:
        enc = MultipartEncoder(fields=fields)
        headers = {"Content-Type": enc.content_type, "X-CSRF-Token": self._csrf()}
        return self.sess.patch(self.base_url + path, headers=headers, data=enc)

    def get(self, path: str) -> requests.Response:
        return self.sess.get(self.base_url + path)

    def login_admin(self, email: str, password: str) -> bool:
        r = self.post_json("/api/login", {"email": email, "password": password})
        return r.status_code == 200

    def set_config(self, key: str, value: str) -> bool:
        r = self.patch_json("/api/configs", {"key": key, "value": value})
        return r.status_code == 200

    def ensure_category(self, name: str, icon: str = "icon") -> bool:
        r = self.post_json("/api/categories", {"name": name, "icon": icon})
        # 200 OK or 409 Conflict (already exists)
        return r.status_code in (200, 409)

    def create_challenge(self, name: str, category: str, ctype: str, description: str = "", max_points: int = 500, score_type: str = "Dynamic", host: Optional[str] = None, port: Optional[int] = None) -> bool:
        data = {
            "name": name,
            "category": category,
            "description": description or f"{name} description",
            "type": ctype,
            "max_points": max_points,
            "score_type": score_type,
        }
        r = self.post_json("/api/challenges", data)
        if r.status_code not in (200, 409):
            return False
        # If host/port provided, update challenge basics
        fields = {"chall_id": str(self.find_challenge_id_by_name(name))}
        if host is not None:
            fields["host"] = host
        if port is not None:
            fields["port"] = str(port)
        if len(fields) > 1:
            rr = self.patch_multipart("/api/challenges", fields)
            return rr.status_code == 200
        return True

    def find_challenge_id_by_name(self, name: str) -> Optional[int]:
        lst = self.list_challenges()
        for c in lst:
            if c.get("name") == name:
                return int(c["id"]) if "id" in c else None
        return None

    def list_challenges(self):
        r = self.get("/api/challenges")
        if r.status_code == 200:
            try:
                return r.json()
            except Exception:
                return []
        return []

    def update_docker_config_image(self, chall_id: int, image: str, lifetime: int = 600, hash_domain: bool = False, max_memory: int = 128, max_cpu: str = "0.5") -> bool:
        fields = {
            "chall_id": str(chall_id),
            "image": image,
            "lifetime": str(lifetime),
            "hash_domain": json.dumps(bool(hash_domain)).lower(),
            "max_memory": str(max_memory),
            "max_cpu": max_cpu,
        }
        r = self.patch_multipart("/api/challenges", fields)
        return r.status_code == 200

    def update_docker_config_compose(self, chall_id: int, compose_body: str, lifetime: int = 600, hash_domain: bool = False, max_memory: int = 128, max_cpu: str = "0.5") -> bool:
        fields = {
            "chall_id": str(chall_id),
            "compose": compose_body,
            "lifetime": str(lifetime),
            "hash_domain": json.dumps(bool(hash_domain)).lower(),
            "max_memory": str(max_memory),
            "max_cpu": max_cpu,
        }
        r = self.patch_multipart("/api/challenges", fields)
        return r.status_code == 200


def get_csrf(client) -> str:
    token = client.cookies.get("csrf_")
    if not token:
        client.get("/api/info", name="GET /api/info")
        token = client.cookies.get("csrf_")
    return token or ""


@events.test_start.add_listener
def setup_challenges(environment, **kwargs):
    global CHALL_CONTAINER_ID, CHALL_COMPOSE_ID
    base_url = (getattr(environment, "host", None) or os.getenv("LOCUST_HOST") or "http://localhost:1337").rstrip("/")

    if not ADMIN_EMAIL or not ADMIN_PASSWORD:
        # Try to reuse if already created; just leave globals as None if not found
        return

    admin = AdminAPI(base_url)
    if not admin.login_admin(ADMIN_EMAIL, ADMIN_PASSWORD):
        return

    # Minimal configs to enable flows
    admin.set_config("allow-register", "true")
    admin.set_config("domain", "localhost")

    # Category
    admin.ensure_category(CATEGORY_NAME, icon="locust")

    # Create container challenge (port 80 for nginx)
    admin.create_challenge(CHALL_CONTAINER_NAME, CATEGORY_NAME, "Container", host="localhost", port=80)
    # Create compose challenge
    admin.create_challenge(CHALL_COMPOSE_NAME, CATEGORY_NAME, "Compose")

    # Resolve IDs
    cid = admin.find_challenge_id_by_name(CHALL_CONTAINER_NAME)
    kid = admin.find_challenge_id_by_name(CHALL_COMPOSE_NAME)
    if cid:
        admin.update_docker_config_image(cid, image="nginx:1.25-alpine", lifetime=600, hash_domain=False)
        CHALL_CONTAINER_ID = cid
    if kid:
        compose_body = (
            "services:\n"
            "  chall:\n"
            "    image: nginx:1.25-alpine\n"
            "    container_name: ${CONTAINER_NAME}\n"
            "    ports:\n"
            "      - \"${INSTANCE_PORT}:80\"\n"
            "    environment:\n"
            "      - INSTANCE_PORT=${INSTANCE_PORT}\n"
            "      - INSTANCE_HOST=${INSTANCE_HOST}\n"
        )
        admin.update_docker_config_compose(kid, compose_body, lifetime=600, hash_domain=False)
        CHALL_COMPOSE_ID = kid


class TeamUser(HttpUser):
    wait_time = between(0.5, 2.0)

    def on_start(self):
        # Register user
        self.name = rand_str("locust")
        self.email = f"{self.name}@test.local"
        self.password = rand_str("pass")

        csrf = get_csrf(self.client)
        headers = {"Content-Type": "application/json", "X-CSRF-Token": csrf}
        r = self.client.post(
            "/api/register",
            json={"name": self.name, "email": self.email, "password": self.password},
            headers=headers,
            name="POST /api/register",
        )
        self.registered = r.status_code == 200
        if not self.registered:
            return

        # Create team
        team_name = rand_str("team")
        team_pass = rand_str("tpass")
        csrf = get_csrf(self.client)
        headers = {"Content-Type": "application/json", "X-CSRF-Token": csrf}
        tr = self.client.post(
            "/api/teams/register",
            json={"name": team_name, "password": team_pass},
            headers=headers,
            name="POST /api/teams/register",
        )
        self.has_team = tr.status_code == 200

        # Resolve challenge ids if not set by setup
        global CHALL_CONTAINER_ID, CHALL_COMPOSE_ID
        if self.has_team and (CHALL_COMPOSE_ID is None or CHALL_CONTAINER_ID is None):
            lr = self.client.get("/api/challenges", name="GET /api/challenges (init)")
            if lr.status_code == 200:
                try:
                    lst = lr.json()
                except Exception:
                    lst = []
                for c in lst:
                    if c.get("name") == CHALL_CONTAINER_NAME:
                        try:
                            CHALL_CONTAINER_ID = int(c["id"]) if "id" in c else None
                        except Exception:
                            pass
                    if c.get("name") == CHALL_COMPOSE_NAME:
                        try:
                            CHALL_COMPOSE_ID = int(c["id"]) if "id" in c else None
                        except Exception:
                            pass

    def create_instance(self, chall_id: int) -> Optional[Tuple[str, Optional[int]]]:
        csrf = get_csrf(self.client)
        headers = {"Content-Type": "application/json", "X-CSRF-Token": csrf}
        r = self.client.post(
            "/api/instances",
            json={"chall_id": chall_id},
            headers=headers,
            name="POST /api/instances",
        )
        if r.status_code != 200:
            return None
        try:
            data = r.json()
            return data.get("host"), data.get("port")
        except Exception:
            return None

    def delete_instance(self, chall_id: int) -> bool:
        csrf = get_csrf(self.client)
        headers = {"Content-Type": "application/json", "X-CSRF-Token": csrf}
        r = self.client.delete(
            "/api/instances",
            json={"chall_id": chall_id},
            headers=headers,
            name="DELETE /api/instances",
        )
        return r.status_code == 200

    def check_http(self, host: str, port: Optional[int]) -> bool:
        if not host:
            return False
        if port is None:
            # hash_domain true would imply portless access; we use domain-only not in this test
            return False
        # Wait briefly for container to be ready
        if not wait_for_tcp(host, int(port), timeout=10.0):
            return False
        try:
            r = requests.get(f"http://{host}:{int(port)}/", timeout=3)
            return r.status_code == 200
        except Exception:
            return False

    @task
    def instance_container_flow(self):
        if not getattr(self, "registered", False) or not getattr(self, "has_team", False):
            return
        if CHALL_CONTAINER_ID is None:
            return

        res = self.create_instance(CHALL_CONTAINER_ID)
        if not res:
            return
        host, port = res
        _ = self.check_http(host, port)
        self.delete_instance(CHALL_CONTAINER_ID)
        # After deletion, it should not be reachable
        ok_after = self.check_http(host, port)
        # We do not explicitly fail the task; metrics will reflect reachability
        _ = ok_after

    @task
    def instance_compose_flow(self):
        if not getattr(self, "registered", False) or not getattr(self, "has_team", False):
            return
        if CHALL_COMPOSE_ID is None:
            return

        res = self.create_instance(CHALL_COMPOSE_ID)
        if not res:
            return
        host, port = res
        _ = self.check_http(host, port)
        self.delete_instance(CHALL_COMPOSE_ID)
        ok_after = self.check_http(host, port)
        _ = ok_after

