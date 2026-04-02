import os
import random
import string
from typing import Optional

from locust import HttpUser, between, events, task
import requests


ADMIN_EMAIL = os.getenv("ADMIN_EMAIL")
ADMIN_PASSWORD = os.getenv("ADMIN_PASSWORD")
CATEGORY_NAME = os.getenv("LOCUST_CATEGORY", "locust")
BASIC_CHALL_NAME = os.getenv("LOCUST_BASIC_CHALL", "locust-basic")
BASIC_CHALL_FLAG = os.getenv("LOCUST_BASIC_FLAG", "TRXD{locust_basic}")

BASIC_CHALL_ID: Optional[int] = None


def rand_str(prefix: str, n: int = 8) -> str:
    return f"{prefix}-" + "".join(random.choices(string.ascii_lowercase + string.digits, k=n))


def get_csrf(client) -> str:
    token = client.cookies.get("csrf_")
    if not token:
        client.get("/api/info", name="GET /api/info")
        token = client.cookies.get("csrf_")
    return token or ""


def resolve_challenge_id(challs: list[dict], name: str) -> Optional[int]:
    for chall in challs:
        if chall.get("name") != name:
            continue
        chall_id = chall.get("id")
        if isinstance(chall_id, int):
            return chall_id
    return None


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

    def get(self, path: str) -> requests.Response:
        return self.sess.get(self.base_url + path)

    def post_json(self, path: str, data: dict) -> requests.Response:
        headers = {"Content-Type": "application/json", "X-CSRF-Token": self._csrf()}
        return self.sess.post(self.base_url + path, headers=headers, json=data)

    def patch_json(self, path: str, data: dict) -> requests.Response:
        headers = {"Content-Type": "application/json", "X-CSRF-Token": self._csrf()}
        return self.sess.patch(self.base_url + path, headers=headers, json=data)

    def login_admin(self) -> bool:
        if not ADMIN_EMAIL or not ADMIN_PASSWORD:
            return False
        resp = self.post_json("/api/login", {"email": ADMIN_EMAIL, "password": ADMIN_PASSWORD})
        return resp.status_code == 200

    def set_config(self, key: str, value: str) -> bool:
        resp = self.patch_json("/api/configs", {"key": key, "value": value})
        return resp.status_code == 200

    def ensure_category(self, name: str) -> bool:
        resp = self.post_json("/api/categories", {"name": name})
        return resp.status_code in (200, 409)

    def list_challenges(self) -> list[dict]:
        resp = self.get("/api/challenges")
        if resp.status_code != 200:
            return []
        try:
            body = resp.json()
        except ValueError:
            return []
        return body if isinstance(body, list) else []

    def find_challenge_id_by_name(self, name: str) -> Optional[int]:
        return resolve_challenge_id(self.list_challenges(), name)

    def ensure_challenge(self, name: str, category: str, ctype: str, **patch_fields) -> Optional[int]:
        resp = self.post_json(
            "/api/challenges",
            {
                "name": name,
                "category": category,
                "description": f"{name} description",
                "type": ctype,
                "max_points": 500,
                "score_type": "Dynamic",
            },
        )
        if resp.status_code not in (200, 409):
            return None

        chall_id = self.find_challenge_id_by_name(name)
        if chall_id is None:
            return None

        if patch_fields:
            update = {"chall_id": chall_id}
            update.update(patch_fields)
            if self.patch_json("/api/challenges", update).status_code != 200:
                return None

        return chall_id

    def ensure_flag(self, chall_id: int, flag: str) -> bool:
        resp = self.post_json("/api/flags", {"chall_id": chall_id, "flag": flag, "regex": False})
        return resp.status_code in (200, 409)


@events.test_start.add_listener
def setup_basic_challenge(environment, **kwargs):
    global BASIC_CHALL_ID

    if not ADMIN_EMAIL or not ADMIN_PASSWORD:
        return

    base_url = (getattr(environment, "host", None) or os.getenv("LOCUST_HOST") or "http://localhost:1337").rstrip("/")
    admin = AdminAPI(base_url)
    if not admin.login_admin():
        raise RuntimeError("Locust basic setup failed: admin login failed")
    if not admin.set_config("allow-register", "true"):
        raise RuntimeError("Locust basic setup failed: allow-register update failed")
    if not admin.ensure_category(CATEGORY_NAME):
        raise RuntimeError("Locust basic setup failed: category setup failed")

    chall_id = admin.ensure_challenge(BASIC_CHALL_NAME, CATEGORY_NAME, "Normal", hidden=False)
    if chall_id is None:
        raise RuntimeError("Locust basic setup failed: challenge setup failed")
    if not admin.ensure_flag(chall_id, BASIC_CHALL_FLAG):
        raise RuntimeError("Locust basic setup failed: flag setup failed")

    BASIC_CHALL_ID = chall_id


class AnonymousUser(HttpUser):
    wait_time = between(0.5, 2.0)

    @task
    def get_scoreboard(self):
        self.client.get("/api/scoreboard", name="GET /api/scoreboard")

    @task
    def get_scoreboard_graph(self):
        self.client.get("/api/scoreboard/graph", name="GET /api/scoreboard/graph")

    @task
    def list_teams(self):
        self.client.get("/api/teams", name="GET /api/teams")

    @task
    def info(self):
        self.client.get("/api/info", name="GET /api/info")


class PlayerUser(HttpUser):
    wait_time = between(0.5, 2.0)

    def on_start(self):
        name = rand_str("locust")
        self._email = f"{name}@test.local"
        self._password = rand_str("pass")
        self._basic_chall_id = BASIC_CHALL_ID
        self._submission_stage = 0

        csrf = get_csrf(self.client)
        headers = {"Content-Type": "application/json", "X-CSRF-Token": csrf}
        resp = self.client.post(
            "/api/register",
            json={"name": name, "email": self._email, "password": self._password},
            headers=headers,
            name="POST /api/register",
        )
        self._registered = resp.status_code == 200
        self._has_team = False

        if not self._registered:
            return

        csrf = get_csrf(self.client)
        headers = {"Content-Type": "application/json", "X-CSRF-Token": csrf}
        resp = self.client.post(
            "/api/teams/register",
            json={"name": rand_str("team"), "password": rand_str("tpass")},
            headers=headers,
            name="POST /api/teams/register",
        )
        self._has_team = resp.status_code == 200

    def _resolve_basic_challenge_id(self) -> Optional[int]:
        resp = self.client.get("/api/challenges", name="GET /api/challenges")
        if resp.status_code != 200:
            return None
        try:
            body = resp.json()
        except ValueError:
            return None
        if not isinstance(body, list):
            return None

        chall_id = resolve_challenge_id(body, BASIC_CHALL_NAME)
        if chall_id is not None:
            self._basic_chall_id = chall_id
        return chall_id

    @task
    def get_challenges(self):
        if not self._registered or not self._has_team:
            return
        self._resolve_basic_challenge_id()

    @task
    def submit_basic_flag(self):
        if not self._registered or not self._has_team:
            return

        chall_id = self._basic_chall_id or BASIC_CHALL_ID or self._resolve_basic_challenge_id()
        if chall_id is None:
            return

        stage = min(self._submission_stage, 2)
        flag = BASIC_CHALL_FLAG if stage else rand_str("wrong-flag")
        expected = "Wrong" if stage == 0 else "Correct" if stage == 1 else "Repeated"

        csrf = get_csrf(self.client)
        headers = {"Content-Type": "application/json", "X-CSRF-Token": csrf}
        with self.client.post(
            "/api/submissions",
            json={"chall_id": chall_id, "flag": flag},
            headers=headers,
            name=f"POST /api/submissions [{expected}]",
            catch_response=True,
        ) as resp:
            if resp.status_code != 200:
                resp.failure(f"unexpected status code {resp.status_code}")
                return

            try:
                body = resp.json()
            except ValueError:
                resp.failure("invalid JSON response")
                return

            if body.get("status") != expected:
                resp.failure(f"expected {expected}, got {body.get('status')}")
                return

            resp.success()
            self._submission_stage += 1

    @task
    def get_users(self):
        self.client.get("/api/users", name="GET /api/users")
