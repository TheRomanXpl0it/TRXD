-- get_random_available_port

CREATE OR REPLACE FUNCTION get_random_available_port_from_range(min_port INTEGER, max_port INTEGER)
RETURNS INTEGER AS $$
DECLARE
  candidate INTEGER;
BEGIN
  SELECT port INTO candidate
    FROM generate_series(min_port, max_port) AS g(port)
    WHERE port NOT IN (SELECT i.port FROM instances i)
    ORDER BY random()
    LIMIT 1;

  IF candidate IS NULL THEN
    RAISE EXCEPTION 'No available ports in range % - %', min_port, max_port;
  END IF;

  RETURN candidate;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_random_available_port()
RETURNS INTEGER AS $$
DECLARE
  min_port INTEGER;
  max_port INTEGER;
BEGIN
  min_port := CAST((SELECT value FROM configs WHERE key = 'min-port') AS INT);
  max_port := CAST((SELECT value FROM configs WHERE key = 'max-port') AS INT);
  RETURN get_random_available_port_from_range(min_port, max_port);
END;
$$ LANGUAGE plpgsql;


-- generate_team_hash

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE OR REPLACE FUNCTION generate_team_hash_trimmed(secret TEXT, team_id INTEGER, port INTEGER, len INTEGER)
RETURNS TEXT AS $$
DECLARE
  full_hash TEXT;
BEGIN
  full_hash := encode(
    digest(secret || team_id::TEXT || port::TEXT, 'sha256'),
    'hex'
  );

  RETURN substring(full_hash FROM 1 FOR len);
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION generate_team_hash(secret TEXT, team_id INTEGER, port INTEGER)
RETURNS TEXT AS $$
DECLARE
  len INTEGER;
BEGIN
  len := CAST((SELECT value FROM configs WHERE key = 'hash-len') AS INT);
  RETURN generate_team_hash_trimmed(secret, team_id, port, len);
END;
$$ LANGUAGE plpgsql;


-- generate_instance_remote

CREATE OR REPLACE FUNCTION generate_instance_remote(
  secret TEXT, team_id INTEGER, chall_id INTEGER, hash_domain BOOLEAN)
RETURNS TABLE(host TEXT, port INTEGER) AS $$
DECLARE
  hash TEXT;
BEGIN
  port := get_random_available_port();
  host := (SELECT challenges.host FROM challenges WHERE id = chall_id);

  IF host IS NULL OR host = '' THEN
    host := (SELECT value FROM configs WHERE key = 'domain');
  END IF;

  IF hash_domain THEN
    hash := generate_team_hash(secret, team_id, port);
    host := hash || '.' || host;
  END IF;

  RETURN NEXT;
END;
$$ LANGUAGE plpgsql;
