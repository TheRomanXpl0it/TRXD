envs:
- `POSTGRES_HOST`: the postgres host (default localhost)
- `POSTGRES_PORT`: the postgres port (default 5432)
- `POSTGRES_DB`: the postgres database name
- `POSTGRES_USER`: the postgres user
- `POSTGRES_PASSWORD`: the postgres password
- `REDIS_HOST`: the redis host (default localhost)
- `REDIS_PORT`: the redis port (default 6379)
- `REDIS_PASSWORD`: the redis password (optional or empty if not needed)
- `REDIS_DISABLE`: (optional) set to something different from empty string to disable redis
- `TESTING`: set to "1" will disable rate-limiter and anti-panic
- `PROJECT_NAME`: use to set the project name for the compose inside the backend (default is "trxd")

flags:
- `-help`: Show help
- `-h`: Show help
- `-t`: Toggle the allow-register config
- `-r`: Register a new admin user with 'username:email:password'
- `-test-data-WARNING-DO-NOT-USE-IN-PRODUCTION`: Inserts mocks data into the db

dev:
create `DEV` file in backend directory to enable dev mode (log level set to debug + debug middleware for api)

roles:
- `spectator`: can only see challenges and not submit (Read Only)
- `player`: the default role
- `author`: can create/update/delete challenges
- `admin`: can edit the platform configs

middlewares:
- `attachments`: only used for attachments, filters the files
- `noAuth`: allows everyone to access the endpoint
- `spectator`: allows only users with spectator role or above
- `player`: allows only users with player role or above
- `author`: allows only users with author role or above
- `admin`: allows only users with admin role or above
- `team`: if the user is a player, requires to be in a team

quick notes on endpoints:
every endpoint returns 200 if everything goes well, otherwise it will return an error code with a json: `{"error": "error message here"}`

endpoints:
- monitor: `/monitor`, admin
- static:
	- `/`: `./frontend`
	- `/static`: `./static`
	- `/attachments`: spectator, team, attachments `./attachments`
	- `/favicon.ico`: `./static/favicon.ico`
- /api:
	- Post(`/register`, noAuth, users_register)
	- Post(`/login`, noAuth, users_login)
	- Post(`/logout`, noAuth, users_logout)
	- Get(`/info`, noAuth, users_info)
	- Get(`/scoreboard`, noAuth, teams_scoreboard)

	- Patch(`/users`, player, users_update)
	- Patch(`/users/password`, admin, users_password)
	- Get(`/users`, noAuth, users_all_get)
	- Get(`/users/:id`, noAuth, users_get)

	- Post(`/teams/register`, player, teams_register)
	- Post(`/teams/join`, player, teams_join)
	- Patch(`/teams`, player, team, teams_update)
	- Patch(`/teams/password`, admin, teams_password)
	- Get(`/teams`, noAuth, teams_all_get)
	- Get(`/teams/:id`, noAuth, teams_get)

	- Post(`/categories`, author, categories_create)
	- Patch(`/categories`, author, categories_update)
	- Delete(`/categories`, author, categories_delete)

	- Post(`/challenges`, author, challenges_create)
	- Patch(`/challenges`, author, challenges_update)
	- Delete(`/challenges`, author, challenges_delete)
	- Get(`/challenges`, spectator, team, challenges_all_get)
	- Get(`/challenges/:id`, spectator, team, challenges_get)

	- Post(`/instances`, player, team, instances_create)
	- Patch(`/instances`, player, team, instances_update)
	- Delete(`/instances`, player, team, instances_delete)

	- Post(`/submissions`, spectator, team, submissions_create)

	- Post(`/tags`, author, tags_create)
	- Patch(`/tags`, author, tags_update)
	- Delete(`/tags`, author, tags_delete)

	- Post(`/flags`, author, flags_create)
	- Patch(`/flags`, author, flags_update)
	- Delete(`/flags`, author, flags_delete)

	- Get(`/configs`, admin, configs_get)
	- Patch(`/configs`, admin, configs_update)