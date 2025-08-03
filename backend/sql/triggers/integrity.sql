-- tr_integrity_repeat_solve

CREATE OR REPLACE FUNCTION fn_integrity_repeat_solve()
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

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_integrity_repeat_solve
BEFORE INSERT ON submissions
FOR EACH ROW
WHEN (NEW.status = 'Correct')
EXECUTE FUNCTION fn_integrity_repeat_solve();


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
