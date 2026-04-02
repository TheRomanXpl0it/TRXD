import os
import random
import string
from threading import Lock
from typing import Optional, Tuple

import requests
from locust import HttpUser, between, events, task
from locust.exception import StopUser


ADMIN_EMAIL = os.getenv("ADMIN_EMAIL")
ADMIN_PASSWORD = os.getenv("ADMIN_PASSWORD")

CHALL_CONTAINER_NAME = os.getenv("LOCUST_CONTAINER_CHALL", "locust-container")
CHALL_COMPOSE_NAME = os.getenv("LOCUST_COMPOSE_CHALL", "locust-compose")
CATEGORY_NAME = os.getenv("LOCUST_CATEGORY", "locust")
CHALL_CONTAINER_FLAG = os.getenv("LOCUST_CONTAINER_FLAG", "TRXD{locust_container}")
CHALL_COMPOSE_FLAG = os.getenv("LOCUST_COMPOSE_FLAG", "TRXD{locust_compose}")
INSTANCE_HOST = os.getenv("LOCUST_INSTANCE_HOST", "127.0.0.1")

CHALL_CONTAINER_ID: Optional[int] = None
CHALL_COMPOSE_ID: Optional[int] = None
SEED_CHALLENGES: list[dict] = []
_USER_INDEX = 0
_USER_INDEX_LOCK = Lock()


def rand_str(prefix: str, n: int = 8) -> str:
    return f"{prefix}-" + "".join(random.choices(string.ascii_lowercase + string.digits, k=n))


def env_bool(name: str, default: bool = False) -> bool:
    raw = os.getenv(name)
    if raw is None:
        return default
    return raw.strip().lower() in {"1", "true", "yes", "on"}


def env_int(name: str, default: int) -> int:
    raw = os.getenv(name)
    if raw is None or raw.strip() == "":
        return default
    return int(raw)


SEED_MODE = env_bool("LOCUST_SEED_MODE", False)
SEED_PREFIX = os.getenv("LOCUST_SEED_PREFIX", rand_str("locust-seed", 6))
SEED_CATEGORY = os.getenv(
    "LOCUST_SEED_CATEGORY",
    CATEGORY_NAME if not SEED_MODE else f"locust-{SEED_PREFIX}",
)
SEED_CHALLENGE_COUNT = env_int("LOCUST_SEED_CHALLENGE_COUNT", 6)
SEED_KEEP_INSTANCES = env_bool("LOCUST_KEEP_INSTANCES", SEED_MODE)
SEED_HASH_DOMAIN = env_bool("LOCUST_SEED_HASH_DOMAIN", False)
SEED_VERIFY_INSTANCE_STATE = env_bool("LOCUST_VERIFY_INSTANCE_STATE", True)
SEED_INSTANCE_LIFETIME = env_int("LOCUST_SEED_CHALLENGE_LIFETIME", 3600)


def next_user_index() -> int:
    global _USER_INDEX
    with _USER_INDEX_LOCK:
        idx = _USER_INDEX
        _USER_INDEX += 1
        return idx


def build_compose_body() -> str:
    return (
        "services:\n"
        "  chall:\n"
        "    image: nginx:1.29-alpine\n"
        "    container_name: ${CONTAINER_NAME}\n"
        "    ports:\n"
        "      - \"${INSTANCE_PORT}:80\"\n"
    )


def build_seed_challenge_specs() -> list[dict]:
    base_prefix = SEED_PREFIX.replace("_", "-")
    flag_prefix = SEED_PREFIX.replace("-", "_")
    specs: list[dict] = []
    for idx in range(SEED_CHALLENGE_COUNT):
        number = idx + 1
        kind = "container" if idx % 2 == 0 else "compose"
        specs.append(
            {
                "index": idx,
                "kind": kind,
                "name": f"{base_prefix}-{kind}-{number:02d}",
                "flag": f"TRXD{{{flag_prefix}_{kind}_{number:02d}}}",
                "ctype": "Container" if kind == "container" else "Compose",
            }
        )
    return specs


def seed_challenge_fields(kind: str) -> dict:
    fields = {
        "hidden": False,
        "host": INSTANCE_HOST,
        "port": 80,
        "conn_type": "HTTP",
        "hash_domain": SEED_HASH_DOMAIN,
        "lifetime": SEED_INSTANCE_LIFETIME,
        "max_memory": 64,
        "max_cpu": "1",
    }
    if kind == "container":
        fields["image"] = "nginx:1.29-alpine"
    else:
        fields["compose"] = build_compose_body()
    return fields


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

    def get(self, path: str) -> requests.Response:
        return self.sess.get(self.base_url + path)

    def login_admin(self, email: str, password: str) -> bool:
        r = self.post_json("/api/login", {"email": email, "password": password})
        return r.status_code == 200

    def set_config(self, key: str, value: str) -> bool:
        r = self.patch_json("/api/configs", {"key": key, "value": value})
        return r.status_code == 200

    def ensure_category(self, name: str) -> bool:
        r = self.post_json("/api/categories", {"name": name})
        return r.status_code in (200, 409)

    def create_challenge(
        self,
        name: str,
        category: str,
        ctype: str,
        description: str = "",
        max_points: int = 500,
        score_type: str = "Dynamic",
    ) -> bool:
        data = {
            "name": name,
            "category": category,
            "description": description or f"{name} description",
            "type": ctype,
            "max_points": max_points,
            "score_type": score_type,
        }
        r = self.post_json("/api/challenges", data)
        return r.status_code in (200, 409)

    def find_challenge_id_by_name(self, name: str) -> Optional[int]:
        for chall in self.list_challenges():
            if chall.get("name") == name:
                return int(chall["id"]) if "id" in chall else None
        return None

    def list_challenges(self):
        r = self.get("/api/challenges")
        if r.status_code == 200:
            try:
                return r.json()
            except Exception:
                return []
        return []

    def update_challenge(self, chall_id: int, **fields) -> bool:
        data = {"chall_id": chall_id}
        data.update(fields)
        r = self.patch_json("/api/challenges", data)
        return r.status_code == 200

    def ensure_flag(self, chall_id: int, flag: str) -> bool:
        r = self.post_json("/api/flags", {"chall_id": chall_id, "flag": flag, "regex": False})
        return r.status_code in (200, 409)


def ensure_seed_challenges(admin: AdminAPI):
    global SEED_CHALLENGES
    challenges: list[dict] = []
    for spec in build_seed_challenge_specs():
        if not admin.create_challenge(spec["name"], SEED_CATEGORY, spec["ctype"]):
            raise RuntimeError(f"Locust seed setup failed: create {spec['name']}")

        chall_id = admin.find_challenge_id_by_name(spec["name"])
        if chall_id is None:
            raise RuntimeError(f"Locust seed setup failed: resolve id for {spec['name']}")

        if not admin.update_challenge(chall_id, **seed_challenge_fields(spec["kind"])):
            raise RuntimeError(f"Locust seed setup failed: update {spec['name']}")

        if not admin.ensure_flag(chall_id, spec["flag"]):
            raise RuntimeError(f"Locust seed setup failed: flag {spec['name']}")

        entry = dict(spec)
        entry["id"] = chall_id
        challenges.append(entry)

    SEED_CHALLENGES = challenges


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
        if SEED_MODE:
            raise RuntimeError("Locust seed setup failed: missing admin credentials")
        return

    admin = AdminAPI(base_url)
    if not admin.login_admin(ADMIN_EMAIL, ADMIN_PASSWORD):
        raise RuntimeError("Locust instance setup failed: admin login failed")

    if not admin.set_config("allow-register", "true"):
        raise RuntimeError("Locust instance setup failed: allow-register update failed")
    if not admin.set_config("domain", INSTANCE_HOST):
        raise RuntimeError("Locust instance setup failed: domain update failed")

    category_name = SEED_CATEGORY if SEED_MODE else CATEGORY_NAME
    if not admin.ensure_category(category_name):
        raise RuntimeError("Locust instance setup failed: category setup failed")

    if SEED_MODE:
        ensure_seed_challenges(admin)
        return

    if not admin.create_challenge(CHALL_CONTAINER_NAME, CATEGORY_NAME, "Container"):
        raise RuntimeError("Locust instance setup failed: container challenge setup failed")
    if not admin.create_challenge(CHALL_COMPOSE_NAME, CATEGORY_NAME, "Compose"):
        raise RuntimeError("Locust instance setup failed: compose challenge setup failed")

    cid = admin.find_challenge_id_by_name(CHALL_CONTAINER_NAME)
    kid = admin.find_challenge_id_by_name(CHALL_COMPOSE_NAME)

    if cid:
        ok = admin.update_challenge(
            cid,
            hidden=False,
            host=INSTANCE_HOST,
            port=80,
            conn_type="HTTP",
            image="nginx:1.29-alpine",
            hash_domain=False,
            lifetime=600,
            max_memory=128,
            max_cpu="0.5",
        )
        if not ok:
            raise RuntimeError("Locust instance setup failed: container challenge update failed")
        if not admin.ensure_flag(cid, CHALL_CONTAINER_FLAG):
            raise RuntimeError("Locust instance setup failed: container flag setup failed")
        CHALL_CONTAINER_ID = cid

    if kid:
        ok = admin.update_challenge(
            kid,
            hidden=False,
            host=INSTANCE_HOST,
            port=80,
            conn_type="HTTP",
            compose=build_compose_body(),
            hash_domain=False,
            lifetime=600,
            max_memory=128,
            max_cpu="0.5",
        )
        if not ok:
            raise RuntimeError("Locust instance setup failed: compose challenge update failed")
        if not admin.ensure_flag(kid, CHALL_COMPOSE_FLAG):
            raise RuntimeError("Locust instance setup failed: compose flag setup failed")
        CHALL_COMPOSE_ID = kid


class TeamUser(HttpUser):
    wait_time = between(0.5, 2.0)

    def on_start(self):
        self.user_index = next_user_index()
        self.name = rand_str("locust")
        self.email = f"{self.name}@test.local"
        self.password = rand_str("pass")
        self.seed_complete = False
        self.seed_done: set[int] = set()
        self.seed_challenges = list(SEED_CHALLENGES)

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
        self.submission_stage = {}

        global CHALL_CONTAINER_ID, CHALL_COMPOSE_ID
        if not self.has_team:
            return

        if not SEED_MODE and CHALL_COMPOSE_ID is not None and CHALL_CONTAINER_ID is not None:
            return

        lr = self.client.get("/api/challenges", name="GET /api/challenges (init)")
        if lr.status_code != 200:
            return
        try:
            lst = lr.json()
        except Exception:
            lst = []

        if not isinstance(lst, list):
            lst = []

        for chall in lst:
            if chall.get("name") == CHALL_CONTAINER_NAME:
                try:
                    CHALL_CONTAINER_ID = int(chall["id"]) if "id" in chall else None
                except Exception:
                    pass
            if chall.get("name") == CHALL_COMPOSE_NAME:
                try:
                    CHALL_COMPOSE_ID = int(chall["id"]) if "id" in chall else None
                except Exception:
                    pass

        if SEED_MODE and not self.seed_challenges:
            self.seed_challenges = self.resolve_seed_challenges(lst)

    def resolve_seed_challenges(self, known: Optional[list[dict]] = None) -> list[dict]:
        if known is None:
            lr = self.client.get("/api/challenges", name="GET /api/challenges [seed-init]")
            if lr.status_code != 200:
                return []
            try:
                known = lr.json()
            except Exception:
                return []

        known = known if isinstance(known, list) else []
        by_name = {
            chall.get("name"): chall
            for chall in known
            if isinstance(chall, dict) and isinstance(chall.get("name"), str)
        }

        resolved: list[dict] = []
        for spec in build_seed_challenge_specs():
            chall = by_name.get(spec["name"])
            if chall is None or "id" not in chall:
                continue
            try:
                chall_id = int(chall["id"])
            except Exception:
                continue
            entry = dict(spec)
            entry["id"] = chall_id
            resolved.append(entry)

        return resolved

    def create_instance(self, chall_id: int) -> Optional[Tuple[str, Optional[int]]]:
        csrf = get_csrf(self.client)
        headers = {"Content-Type": "application/json", "X-CSRF-Token": csrf}
        with self.client.post(
            "/api/instances",
            json={"chall_id": chall_id},
            headers=headers,
            name="POST /api/instances",
            catch_response=True,
        ) as r:
            if r.status_code != 200:
                r.failure(f"unexpected status code {r.status_code}")
                return None
            try:
                data = r.json()
            except ValueError:
                r.failure("invalid JSON response")
                return None

            host = data.get("host")
            port = data.get("port")
            timeout = data.get("timeout")
            if not host or timeout is None:
                r.failure("missing host or timeout")
                return None
            if port is not None and not isinstance(port, int):
                r.failure("invalid port in response")
                return None

            r.success()
            return host, port

    def delete_instance(self, chall_id: int) -> bool:
        csrf = get_csrf(self.client)
        headers = {"Content-Type": "application/json", "X-CSRF-Token": csrf}
        with self.client.delete(
            "/api/instances",
            json={"chall_id": chall_id},
            headers=headers,
            name="DELETE /api/instances",
            catch_response=True,
        ) as r:
            if r.status_code != 200:
                r.failure(f"unexpected status code {r.status_code}")
                return False
            r.success()
            return True

    def submit_flag(self, chall_id: int, flag: str, label: str) -> bool:
        stage = min(self.submission_stage.get(chall_id, 0), 2)
        sent_flag = flag if stage else rand_str("wrong-flag")
        expected = "Wrong" if stage == 0 else "Correct" if stage == 1 else "Repeated"

        csrf = get_csrf(self.client)
        headers = {"Content-Type": "application/json", "X-CSRF-Token": csrf}
        with self.client.post(
            "/api/submissions",
            json={"chall_id": chall_id, "flag": sent_flag},
            headers=headers,
            name=f"POST /api/submissions [{label}:{expected}]",
            catch_response=True,
        ) as r:
            if r.status_code != 200:
                r.failure(f"unexpected status code {r.status_code}")
                return False
            try:
                data = r.json()
            except ValueError:
                r.failure("invalid JSON response")
                return False
            if data.get("status") != expected:
                r.failure(f"expected {expected}, got {data.get('status')}")
                return False

            r.success()
            self.submission_stage[chall_id] = stage + 1
            return True

    def fetch_challenge_state(self, chall_id: int, name: str) -> Optional[dict]:
        with self.client.get(
            "/api/challenges",
            name=f"GET /api/challenges [{name}]",
            catch_response=True,
        ) as r:
            if r.status_code != 200:
                r.failure(f"unexpected status code {r.status_code}")
                return None
            try:
                data = r.json()
            except ValueError:
                r.failure("invalid JSON response")
                return None
            if not isinstance(data, list):
                r.failure("expected a challenge list")
                return None

            for chall in data:
                if chall.get("id") == chall_id:
                    r.success()
                    return chall

            r.failure(f"challenge {chall_id} not found")
            return None

    def assert_instance_state(self, chall_id: int, name: str, host: str, port: Optional[int]):
        chall = self.fetch_challenge_state(chall_id, name)
        if chall is None:
            return
        if chall.get("instance_host") != host:
            raise RuntimeError(f"{name} instance host mismatch: {chall.get('instance_host')} != {host}")
        if chall.get("instance_port") != port:
            raise RuntimeError(f"{name} instance port mismatch: {chall.get('instance_port')} != {port}")
        if int(chall.get("timeout", 0)) <= 0:
            raise RuntimeError(f"{name} instance timeout was not exposed")

    def should_solve_seed(self, chall_index: int) -> bool:
        return (self.user_index + chall_index) % 2 == 0

    def should_repeat_seed(self, chall_index: int) -> bool:
        return (self.user_index + chall_index) % 4 == 0

    @task(100)
    def seed_dataset(self):
        if not SEED_MODE:
            return
        if not getattr(self, "registered", False) or not getattr(self, "has_team", False):
            raise StopUser()
        if self.seed_complete:
            raise StopUser()

        if not self.seed_challenges:
            self.seed_challenges = self.resolve_seed_challenges()
        if not self.seed_challenges:
            raise RuntimeError("Locust seed mode could not resolve seeded challenge ids")

        for chall in self.seed_challenges:
            chall_id = chall["id"]
            if chall_id in self.seed_done:
                continue

            res = self.create_instance(chall_id)
            if not res:
                continue
            host, port = res

            if SEED_VERIFY_INSTANCE_STATE:
                self.assert_instance_state(chall_id, chall["name"], host, port)

            self.submit_flag(chall_id, chall["flag"], chall["kind"])
            if self.should_solve_seed(chall["index"]):
                self.submit_flag(chall_id, chall["flag"], chall["kind"])
                if self.should_repeat_seed(chall["index"]):
                    self.submit_flag(chall_id, chall["flag"], chall["kind"])

            if not SEED_KEEP_INSTANCES:
                if not self.delete_instance(chall_id):
                    raise RuntimeError(f"{chall['name']} instance deletion failed")
                if SEED_VERIFY_INSTANCE_STATE:
                    cleared = self.fetch_challenge_state(chall_id, chall["name"])
                    if cleared and (cleared.get("instance_host") or cleared.get("instance_port")):
                        raise RuntimeError(f"{chall['name']} instance state still present after deletion")

            self.seed_done.add(chall_id)

        self.seed_complete = True
        raise StopUser()

    @task
    def list_challenges(self):
        if SEED_MODE:
            return
        if not getattr(self, "registered", False) or not getattr(self, "has_team", False):
            return
        self.client.get("/api/challenges", name="GET /api/challenges")

    @task
    def submit_container_flag(self):
        if SEED_MODE:
            return
        if not getattr(self, "registered", False) or not getattr(self, "has_team", False):
            return
        if CHALL_CONTAINER_ID is None:
            return
        self.submit_flag(CHALL_CONTAINER_ID, CHALL_CONTAINER_FLAG, "container")

    @task
    def submit_compose_flag(self):
        if SEED_MODE:
            return
        if not getattr(self, "registered", False) or not getattr(self, "has_team", False):
            return
        if CHALL_COMPOSE_ID is None:
            return
        self.submit_flag(CHALL_COMPOSE_ID, CHALL_COMPOSE_FLAG, "compose")

    @task
    def instance_container_flow(self):
        if SEED_MODE:
            return
        if not getattr(self, "registered", False) or not getattr(self, "has_team", False):
            return
        if CHALL_CONTAINER_ID is None:
            return

        res = self.create_instance(CHALL_CONTAINER_ID)
        if not res:
            return
        host, port = res
        try:
            self.assert_instance_state(CHALL_CONTAINER_ID, CHALL_CONTAINER_NAME, host, port)
        finally:
            if not self.delete_instance(CHALL_CONTAINER_ID):
                raise RuntimeError("container instance deletion failed")

        chall = self.fetch_challenge_state(CHALL_CONTAINER_ID, CHALL_CONTAINER_NAME)
        if chall is None:
            return
        if chall.get("instance_host") or chall.get("instance_port"):
            raise RuntimeError("container instance state still present after deletion")

    @task
    def instance_compose_flow(self):
        if SEED_MODE:
            return
        if not getattr(self, "registered", False) or not getattr(self, "has_team", False):
            return
        if CHALL_COMPOSE_ID is None:
            return

        res = self.create_instance(CHALL_COMPOSE_ID)
        if not res:
            return
        host, port = res
        try:
            self.assert_instance_state(CHALL_COMPOSE_ID, CHALL_COMPOSE_NAME, host, port)
        finally:
            if not self.delete_instance(CHALL_COMPOSE_ID):
                raise RuntimeError("compose instance deletion failed")

        chall = self.fetch_challenge_state(CHALL_COMPOSE_ID, CHALL_COMPOSE_NAME)
        if chall is None:
            return
        if chall.get("instance_host") or chall.get("instance_port"):
            raise RuntimeError("compose instance state still present after deletion")
