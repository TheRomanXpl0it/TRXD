-- get_random_available_port

CREATE OR REPLACE FUNCTION get_random_available_port_from_range(min_port INTEGER, max_port INTEGER)
RETURNS INTEGER AS $$
DECLARE
  candidate INTEGER;
BEGIN
  PERFORM pg_advisory_xact_lock(1337);

  SELECT port INTO candidate
    FROM generate_series(min_port, max_port) AS g(port)
    WHERE port NOT IN (SELECT i.port FROM instances i WHERE i.port IS NOT NULL)
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


-- generate_instance_remote

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE OR REPLACE FUNCTION generate_instance_remote(chall_id INTEGER, hash_domain BOOLEAN)
RETURNS TABLE(host TEXT, port INTEGER) AS $$
DECLARE
  len INTEGER;
  hash TEXT;
BEGIN
  host := (SELECT challenges.host FROM challenges WHERE id = chall_id);

  IF host IS NULL OR host = '' THEN
    host := (SELECT value FROM configs WHERE key = 'domain');
  END IF;

  IF hash_domain THEN
    len := CAST((SELECT value FROM configs WHERE key = 'hash-len') AS INT);
    hash := encode(gen_random_bytes(len/2),'hex');
    host := hash || '.' || host;
  ELSE
    port := get_random_available_port();
  END IF;

  RETURN NEXT;
END;
$$ LANGUAGE plpgsql;
