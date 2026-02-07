-- tr_integrity_solve

CREATE OR REPLACE FUNCTION fn_integrity_solve()
RETURNS TRIGGER AS $$
DECLARE
  team INTEGER;
  existing_correct_count INTEGER;
BEGIN
  IF (SELECT role FROM users WHERE id = NEW.user_id) != 'Player' THEN
    RETURN NEW;
  END IF;

  SELECT team_id INTO team
    FROM users
    WHERE id = NEW.user_id;

  IF team IS NULL THEN
    NEW.status = 'Invalid';
    RETURN NEW;
  END IF;

  SELECT COUNT(*) INTO existing_correct_count
    FROM submissions
    JOIN users ON users.id = submissions.user_id
    WHERE users.team_id = team
      AND users.role = 'Player'
      AND submissions.chall_id = NEW.chall_id
      AND submissions.status = 'Correct';

  IF existing_correct_count > 0 THEN
    NEW.status = 'Repeated';
  END IF;

  IF NEW.status = 'Correct' THEN
    IF NOT EXISTS (
      SELECT 1 FROM submissions
        WHERE chall_id = NEW.chall_id
          AND first_blood = TRUE
    ) THEN
      NEW.first_blood = TRUE;
    END IF;
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_integrity_solve
BEFORE INSERT ON submissions
FOR EACH ROW
WHEN (NEW.status = 'Correct')
EXECUTE FUNCTION fn_integrity_solve();


-- tr_integrity_delete_solve

CREATE OR REPLACE FUNCTION fn_integrity_delete_solve()
RETURNS TRIGGER AS $$
DECLARE
  next_id INTEGER;
BEGIN
  IF (SELECT role FROM users WHERE id = OLD.user_id) != 'Player' THEN
    RETURN OLD;
  END IF;

  SELECT s.id INTO next_id
    FROM submissions s
    JOIN users u ON u.id = s.user_id
    WHERE s.chall_id = OLD.chall_id
      AND s.status = 'Correct'
      AND u.role = 'Player'
    ORDER BY s.timestamp ASC
    LIMIT 1;

  IF next_id IS NOT NULL THEN
    UPDATE submissions
      SET first_blood = TRUE
      WHERE id = next_id;
  END IF;

  RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_integrity_delete_solve
AFTER DELETE ON submissions
FOR EACH ROW
WHEN (OLD.status = 'Correct' AND OLD.first_blood = TRUE)
EXECUTE FUNCTION fn_integrity_delete_solve();


-- tr_integrity_chall_default_points

CREATE OR REPLACE FUNCTION fn_integrity_chall_default_points()
RETURNS TRIGGER AS $$
BEGIN
  NEW.points = NEW.max_points;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_integrity_chall_default_points
BEFORE INSERT ON challenges
FOR EACH ROW
EXECUTE FUNCTION fn_integrity_chall_default_points();


-- tr_integrity_chall_docker_configs_add

CREATE OR REPLACE FUNCTION fn_integrity_chall_docker_configs_add()
RETURNS TRIGGER AS $$
BEGIN
  INSERT INTO docker_configs (chall_id) VALUES (NEW.id);
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_integrity_chall_docker_configs_add
AFTER INSERT ON challenges
FOR EACH ROW
WHEN (NEW.type != 'Normal')
EXECUTE FUNCTION fn_integrity_chall_docker_configs_add();


-- tr_integrity_chall_docker_configs_add_on_update

CREATE OR REPLACE FUNCTION fn_integrity_chall_docker_configs_add_on_update()
RETURNS TRIGGER AS $$
BEGIN
  IF NOT EXISTS(SELECT * FROM docker_configs WHERE chall_id = NEW.id) THEN
    INSERT INTO docker_configs (chall_id) VALUES (NEW.id);
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_integrity_chall_docker_configs_add_on_update
AFTER UPDATE ON challenges
FOR EACH ROW
WHEN ((OLD.type = 'Normal') AND (NEW.type != 'Normal'))
EXECUTE FUNCTION fn_integrity_chall_docker_configs_add_on_update();
