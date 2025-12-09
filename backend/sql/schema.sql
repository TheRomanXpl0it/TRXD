CREATE TYPE user_role AS ENUM (
  'Spectator',
  'Player',
  'Author',
  'Admin'
);

CREATE TYPE deploy_type AS ENUM (
  'Normal',
  'Container',
  'Compose'
);

CREATE TYPE score_type AS ENUM (
  'Static',
  'Dynamic'
);

CREATE TYPE submission_status AS ENUM (
  'Wrong',
  'Correct',
  'Repeated',
  'Invalid'
);

-- CREATE TYPE conn_type AS ENUM (
--   'TCP',
--   'HTTP',
-- );

CREATE TABLE IF NOT EXISTS configs (
  key TEXT NOT NULL,
  type TEXT NOT NULL DEFAULT 'text',
  value TEXT NOT NULL DEFAULT '',
  description TEXT,
  PRIMARY KEY(key)
);

CREATE TABLE IF NOT EXISTS teams (
  id SERIAL NOT NULL,
  name VARCHAR(64) UNIQUE NOT NULL,
  password_hash CHAR(64) NOT NULL,
  password_salt CHAR(32) NOT NULL,
  score INTEGER NOT NULL DEFAULT 0,
  country VARCHAR(3),
  image VARCHAR(1024),
  PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS users (
  id SERIAL NOT NULL,
  name VARCHAR(64) UNIQUE NOT NULL,
  email VARCHAR(256) UNIQUE NOT NULL,
  password_hash CHAR(64) NOT NULL,
  password_salt CHAR(32) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

  score INTEGER NOT NULL DEFAULT 0,
  role user_role NOT NULL,

  team_id INTEGER,
  country VARCHAR(3),
  image VARCHAR(1024),

  FOREIGN KEY(team_id) REFERENCES teams(id),
  PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS badges (
  name VARCHAR(64) NOT NULL,
  description VARCHAR(1024) NOT NULL,
  team_id INTEGER NOT NULL,
  FOREIGN KEY(team_id) REFERENCES teams(id) ON DELETE CASCADE,
  PRIMARY KEY(name, team_id)
);

CREATE TABLE IF NOT EXISTS categories (
  name VARCHAR(32) NOT NULL,
  visible_challs INTEGER NOT NULL DEFAULT 0,
  PRIMARY KEY(name)
);

CREATE TABLE IF NOT EXISTS team_category_solves (
  team_id INTEGER NOT NULL,
  category VARCHAR(32) NOT NULL,
  solves INTEGER NOT NULL DEFAULT 0,
  PRIMARY KEY (team_id, category),
  FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
  FOREIGN KEY (category) REFERENCES categories(name) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS challenges (
  id SERIAL NOT NULL,
  name VARCHAR(128) UNIQUE NOT NULL,
  category VARCHAR(32) NOT NULL,
  description VARCHAR(1024) NOT NULL,
  difficulty VARCHAR(16) NOT NULL DEFAULT '',
  authors TEXT NOT NULL DEFAULT '', -- TODO: use a list
  type deploy_type NOT NULL,
  hidden BOOLEAN NOT NULL DEFAULT TRUE,

  max_points INTEGER NOT NULL,
  score_type score_type NOT NULL,
  points INTEGER NOT NULL,
  solves INTEGER NOT NULL DEFAULT 0,

  -- conn_type conn_type NOT NULL DEFAULT 'TCP'
  host TEXT NOT NULL DEFAULT '',
  port INTEGER NOT NULL CHECK (port >= 0 AND port <= 65535) DEFAULT 0,
  attachments TEXT NOT NULL DEFAULT '', -- TODO: use a list (maybe a table)

  FOREIGN KEY(category) REFERENCES categories(name),
  PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS docker_configs (
  chall_id INTEGER NOT NULL,
  image TEXT NOT NULL DEFAULT '', -- Docker image
  compose TEXT NOT NULL DEFAULT '', -- Docker Compose file
  hash_domain BOOLEAN NOT NULL DEFAULT FALSE, -- Use a hash to generate the domain (e.g., 00112233AABB.example.com)
  lifetime INTEGER NOT NULL DEFAULT 0, -- Lifetime in seconds
  envs TEXT NOT NULL DEFAULT '', -- Environment variables in JSON format
  max_memory INTEGER NOT NULL DEFAULT 0, -- Memory in MB (e.g., '512' for 512 MB)
  max_cpu VARCHAR(16) NOT NULL DEFAULT '', -- CPUs as float (e.g., '1.5' for 1.5 CPUs)
  FOREIGN KEY(chall_id) REFERENCES challenges(id) ON DELETE CASCADE,
  PRIMARY KEY(chall_id)
);

CREATE TABLE IF NOT EXISTS flags (
  flag VARCHAR(128) UNIQUE NOT NULL,
  chall_id INTEGER NOT NULL,
  regex BOOLEAN NOT NULL DEFAULT FALSE,
  FOREIGN KEY(chall_id) REFERENCES challenges(id) ON DELETE CASCADE,
  PRIMARY KEY(flag, chall_id)
);

CREATE TABLE IF NOT EXISTS instances (
  team_id INTEGER NOT NULL,
  chall_id INTEGER NOT NULL,
  expires_at TIMESTAMP NOT NULL,
  host TEXT NOT NULL,
  port INTEGER UNIQUE CHECK (port >= 0 AND port <= 65535),
  docker_id VARCHAR(64), -- Docker instance ID (container ID or compose project name)
  FOREIGN KEY(team_id) REFERENCES teams(id) ON DELETE CASCADE,
  FOREIGN KEY(chall_id) REFERENCES challenges(id) ON DELETE CASCADE,
  PRIMARY KEY(team_id, chall_id)
);

CREATE TABLE IF NOT EXISTS tags (
  chall_id INTEGER NOT NULL,
  name VARCHAR(32) NOT NULL,
  FOREIGN KEY(chall_id) REFERENCES challenges(id) ON DELETE CASCADE,
  PRIMARY KEY(chall_id, name)
);

CREATE TABLE IF NOT EXISTS submissions (
  id SERIAL NOT NULL,
  user_id INTEGER NOT NULL,
  chall_id INTEGER NOT NULL,
  status submission_status NOT NULL,
  first_blood BOOLEAN NOT NULL DEFAULT FALSE,
  flag TEXT NOT NULL,
  timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY(chall_id) REFERENCES challenges(id) ON DELETE CASCADE,
  PRIMARY KEY(id)
);


CREATE INDEX IF NOT EXISTS idx_teams_name ON teams(name);
CREATE INDEX IF NOT EXISTS idx_users_team_id ON users(team_id);
CREATE INDEX IF NOT EXISTS idx_challenges_category ON challenges(category);
CREATE INDEX IF NOT EXISTS idx_tags_chall_id ON tags(chall_id);
CREATE INDEX IF NOT EXISTS idx_submissions_user_id ON submissions(user_id);
CREATE INDEX IF NOT EXISTS idx_submissions_chall_id ON submissions(chall_id);
