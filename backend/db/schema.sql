CREATE TYPE user_role AS ENUM (
	'S', -- Spectator
	'P', -- Player
	'A', -- Author
	'M' -- Master (Admin)
);

CREATE TYPE deploy_type AS ENUM (
	'N', -- Normal
	'D', -- Docker Instance
	'C' -- Compose Instance
);

CREATE TYPE score_type AS ENUM (
	'S', -- Static
	'D' -- Dynamic (Logarithmic)
);

CREATE TYPE submission_status AS ENUM (
	'W', -- Wrong
	'C', -- Correct
	'R' -- Repeated
);

CREATE DOMAIN name AS VARCHAR(64);
CREATE DOMAIN short_name AS VARCHAR(32);
CREATE DOMAIN long_name AS VARCHAR(128);
CREATE DOMAIN bcrypt_hash AS CHAR(60);
CREATE DOMAIN port AS INTEGER CHECK (VALUE >= 1 AND VALUE <= 65535);


CREATE TABLE IF NOT EXISTS configs (
	key TEXT NOT NULL,
	type TEXT NOT NULL DEFAULT 'text',
	value TEXT NOT NULL DEFAULT '',
	description TEXT,
	PRIMARY KEY(key)
);

CREATE TABLE IF NOT EXISTS teams (
	id SERIAL NOT NULL,
	name name NOT NULL,
	password_hash bcrypt_hash NOT NULL,
	score INTEGER NOT NULL DEFAULT 0,
	nationality VARCHAR(3),
	image TEXT,
	bio TEXT,
	PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS users (
	id SERIAL NOT NULL,
	name name NOT NULL,
	email TEXT UNIQUE NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

	password_hash bcrypt_hash NOT NULL,
	apikey UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),

	score INTEGER NOT NULL DEFAULT 0,
	role user_role NOT NULL,

	team_id INTEGER,
	nationality VARCHAR(3),
	image TEXT,

	FOREIGN KEY(team_id) REFERENCES teams(id),
	PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS badges (
	name name NOT NULL,
	description VARCHAR(1024) NOT NULL,
	team_id INTEGER NOT NULL,
	FOREIGN KEY(team_id) REFERENCES teams(id) ON DELETE CASCADE,
	PRIMARY KEY(name, team_id)
);

CREATE TABLE IF NOT EXISTS categories (
	name short_name NOT NULL,
	chall_count INTEGER NOT NULL DEFAULT 0,
	icon short_name NOT NULL,
	PRIMARY KEY(name)
);

CREATE TABLE IF NOT EXISTS challenges (
	id SERIAL NOT NULL,
	name long_name UNIQUE NOT NULL,
	category short_name NOT NULL,
	description VARCHAR(1024) NOT NULL,
	difficulty VARCHAR(16),
	authors TEXT,
	type deploy_type NOT NULL,
	hidden BOOLEAN NOT NULL DEFAULT TRUE,

	max_points INTEGER NOT NULL,
	score_type score_type NOT NULL,
	points INTEGER NOT NULL,
	solves INTEGER NOT NULL DEFAULT 0,

	host TEXT,
	port port,
	attachments TEXT, -- List of attachments separated by nullbytes

	image TEXT, -- Docker image
	compose TEXT, -- Docker Compose file
	hash_domain BOOLEAN DEFAULT FALSE, -- Use a hash to generate the domain (e.g., 00112233AABB.example.com)
	lifetime INTEGER, -- Lifetime in seconds
	envs TEXT, -- Environment variables in JSON format
	max_memory INTEGER, -- Memory in MB (e.g., '512' for 512 MB)
	max_cpu VARCHAR(16), -- CPUs as float (e.g., '1.5' for 1.5 CPUs)

	FOREIGN KEY(category) REFERENCES categories(name),
	PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS flags (
	flag VARCHAR(128) UNIQUE NOT NULL,
	chall_id INTEGER NOT NULL,
	regex BOOLEAN NOT NULL DEFAULT FALSE,
	FOREIGN KEY(chall_id) REFERENCES challenges(id) ON DELETE CASCADE,
	PRIMARY KEY(flag, chall_id)
);

CREATE TABLE IF NOT EXISTS instances (
	id SERIAL NOT NULL,
	team_id INTEGER NOT NULL,
	chall_id INTEGER NOT NULL,
	expires_at TIMESTAMP NOT NULL,
	host TEXT NOT NULL,
	port port NOT NULL,
	FOREIGN KEY(team_id) REFERENCES teams(id) ON DELETE CASCADE,
	FOREIGN KEY(chall_id) REFERENCES challenges(id) ON DELETE CASCADE,
	PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS tags (
	name short_name NOT NULL,
	chall_id INTEGER NOT NULL,
	FOREIGN KEY(chall_id) REFERENCES challenges(id) ON DELETE CASCADE,
	PRIMARY KEY(name, chall_id)
);

CREATE TABLE IF NOT EXISTS submissions (
	id SERIAL NOT NULL,
	user_id INTEGER NOT NULL,
	chall_id INTEGER NOT NULL,
	status submission_status NOT NULL,
	flag TEXT NOT NULL,
	timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(user_id) REFERENCES users(id),
	FOREIGN KEY(chall_id) REFERENCES challenges(id) ON DELETE CASCADE,
	PRIMARY KEY(id)
);

CREATE INDEX IF NOT EXISTS idx_users_apikey ON users(apikey);
CREATE INDEX IF NOT EXISTS idx_users_team_id ON users(team_id);
CREATE INDEX IF NOT EXISTS idx_flags_flag ON flags(flag);
CREATE INDEX IF NOT EXISTS idx_flags_chall_id ON flags(chall_id);
CREATE INDEX IF NOT EXISTS idx_instances_team_id ON instances(team_id);
CREATE INDEX IF NOT EXISTS idx_instances_chall_id ON instances(chall_id);
CREATE INDEX IF NOT EXISTS idx_tags_chall_id ON tags(chall_id);
CREATE INDEX IF NOT EXISTS idx_submissions_user_id ON submissions(user_id);
CREATE INDEX IF NOT EXISTS idx_submissions_chall_id ON submissions(chall_id);
