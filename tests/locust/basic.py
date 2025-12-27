import os
import random
import string
from locust import HttpUser, task, between, events


def rand_str(prefix: str, n: int = 8) -> str:
    return f"{prefix}-" + "".join(random.choices(string.ascii_lowercase + string.digits, k=n))


def get_csrf(client) -> str:
    token = client.cookies.get("csrf_")
    if not token:
        client.get("/api/info", name="GET /api/info")
        token = client.cookies.get("csrf_")
    return token or ""


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
        # Try to register a new user
        name = rand_str("locust")
        email = f"{name}@test.local"
        password = rand_str("pass")

        csrf = get_csrf(self.client)
        headers = {"Content-Type": "application/json", "X-CSRF-Token": csrf}
        resp = self.client.post(
            "/api/register",
            json={"name": name, "email": email, "password": password},
            headers=headers,
            name="POST /api/register",
        )

        # If registration is disabled (403), skip further auth-required flows in tasks
        self._registered = resp.status_code == 200
        self._name = name
        self._email = email
        self._password = password
        self._has_team = False

        if not self._registered:
            return

        # Create a team so we can access /api/challenges (requires team for Player)
        team_name = rand_str("team")
        team_pass = rand_str("tpass")
        csrf = get_csrf(self.client)
        headers = {"Content-Type": "application/json", "X-CSRF-Token": csrf}
        resp = self.client.post(
            "/api/teams/register",
            json={"name": team_name, "password": team_pass},
            headers=headers,
            name="POST /api/teams/register",
        )
        self._has_team = resp.status_code == 200

    @task
    def login(self):
        if not getattr(self, "_registered", False):
            return
        csrf = get_csrf(self.client)
        headers = {"Content-Type": "application/json", "X-CSRF-Token": csrf}
        self.client.post(
            "/api/login",
            json={"email": self._email, "password": self._password},
            headers=headers,
            name="POST /api/login",
        )

    @task
    def get_challenges(self):
        if not getattr(self, "_registered", False) or not getattr(self, "_has_team", False):
            return
        self.client.get("/api/challenges", name="GET /api/challenges")

    @task
    def get_users(self):
        self.client.get("/api/users", name="GET /api/users")
