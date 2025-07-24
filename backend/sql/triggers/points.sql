-- utils

CREATE OR REPLACE FUNCTION fn_compute_chall_points(
  min_points INTEGER, decay REAL, chall_max_points INTEGER, chall_solves INTEGER)
RETURNS INTEGER AS $$
BEGIN
  IF chall_max_points <= min_points THEN
    RETURN chall_max_points;
  END IF;

  RETURN GREATEST(
    min_points,
    CAST((chall_max_points + (min_points - chall_max_points) / (decay ^ 2) *
      (CASE WHEN chall_solves > 0 THEN (chall_solves - 1) ^ 2 ELSE 0 END)) AS INT)
  );
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION fn_points_propagate_points(diff INTEGER, chall_id INTEGER)
RETURNS VOID AS $$
BEGIN
  UPDATE users
    SET score = score - diff
    FROM submissions
    WHERE submissions.chall_id = fn_points_propagate_points.chall_id
      AND submissions.status = 'Correct'
      AND users.id = submissions.user_id
      AND users.role = 'Player';
END;
$$ LANGUAGE plpgsql;


-- tr_points_add_solve

CREATE OR REPLACE FUNCTION fn_points_add_solve()
RETURNS TRIGGER AS $$
BEGIN
  IF (SELECT role FROM users WHERE id = NEW.user_id) != 'Player' THEN
    RETURN NEW;
  END IF;

  UPDATE users
    SET score = score + challenges.points
    FROM challenges
    WHERE challenges.id = NEW.chall_id
      AND users.id = NEW.user_id;

  UPDATE challenges
    SET solves = solves + 1
    WHERE id = NEW.chall_id;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_points_add_solve
AFTER INSERT ON submissions
FOR EACH ROW
WHEN (NEW.status = 'Correct')
EXECUTE FUNCTION fn_points_add_solve();


-- tr_points_del_solve

CREATE OR REPLACE FUNCTION fn_points_del_solve()
RETURNS TRIGGER AS $$
BEGIN
  IF (SELECT role FROM users WHERE id = OLD.user_id) != 'Player' THEN
    RETURN OLD;
  END IF;

  UPDATE challenges
    SET solves = solves - 1
    WHERE id = OLD.chall_id;

  UPDATE users
    SET score = score - challenges.points
    FROM challenges
    WHERE challenges.id = OLD.chall_id
      AND users.id = OLD.user_id;

  RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_points_del_solve
BEFORE DELETE ON submissions
FOR EACH ROW
WHEN (OLD.status = 'Correct')
EXECUTE FUNCTION fn_points_del_solve();


-- tr_points_chall_update

CREATE OR REPLACE FUNCTION fn_points_chall_update()
RETURNS TRIGGER AS $$
DECLARE
  min_points INTEGER;
  decay REAL;
BEGIN
  min_points = CAST((SELECT value FROM configs WHERE key = 'chall-min-points') AS INT);
  decay = CAST((SELECT value FROM configs WHERE key = 'chall-points-decay') AS REAL);
  NEW.points = fn_compute_chall_points(min_points, decay, NEW.max_points, NEW.solves);
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_points_chall_update
BEFORE UPDATE ON challenges
FOR EACH ROW
WHEN ((NEW.score_type = 'Dynamic') AND 
  ((NEW.solves != OLD.solves) OR (NEW.max_points != OLD.max_points)))
EXECUTE FUNCTION fn_points_chall_update();


-- tr_points_chall_update_static

CREATE OR REPLACE FUNCTION fn_points_chall_update_static()
RETURNS TRIGGER AS $$
BEGIN
  NEW.points = NEW.max_points;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_points_chall_update_static
BEFORE UPDATE ON challenges
FOR EACH ROW
WHEN (((OLD.score_type = 'Dynamic') AND (NEW.score_type = 'Static'))
  OR (NEW.score_type = 'Static' AND NEW.max_points != OLD.max_points))
EXECUTE FUNCTION fn_points_chall_update_static();


-- tr_points_propagate_config

CREATE OR REPLACE FUNCTION fn_points_propagate_config()
RETURNS TRIGGER AS $$
DECLARE
  min_points INTEGER;
  decay REAL;
BEGIN
  min_points = CAST((SELECT value FROM configs WHERE key = 'chall-min-points') AS INT);
  decay = CAST((SELECT value FROM configs WHERE key = 'chall-points-decay') AS REAL);
  UPDATE challenges
    SET points = fn_compute_chall_points(min_points, decay, max_points, solves)
    WHERE score_type = 'Dynamic';
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_points_propagate_config
AFTER UPDATE ON configs
FOR EACH ROW
WHEN (NEW.key = 'chall-min-points' OR NEW.key = 'chall-points-decay')
EXECUTE FUNCTION fn_points_propagate_config();


-- tr_points_chall_del

CREATE OR REPLACE FUNCTION fn_points_chall_del()
RETURNS TRIGGER AS $$
BEGIN
  PERFORM fn_points_propagate_points(OLD.points, OLD.id);
  RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_points_chall_del
BEFORE DELETE ON challenges
FOR EACH ROW
EXECUTE FUNCTION fn_points_chall_del();


-- tr_points_propagate_chall

CREATE OR REPLACE FUNCTION fn_points_propagate_chall()
RETURNS TRIGGER AS $$
BEGIN
  PERFORM fn_points_propagate_points(OLD.points - NEW.points, OLD.id);
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_points_propagate_chall
AFTER UPDATE ON challenges
FOR EACH ROW
WHEN (NEW.points != OLD.points)
EXECUTE FUNCTION fn_points_propagate_chall();


-- tr_points_propagate_user

CREATE OR REPLACE FUNCTION fn_points_propagate_user()
RETURNS TRIGGER AS $$
DECLARE
  diff INTEGER;
BEGIN
  diff = NEW.score - OLD.score;
  UPDATE teams
    SET score = score + diff
    WHERE id = NEW.team_id;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_points_propagate_user
AFTER UPDATE ON users
FOR EACH ROW
WHEN (NEW.score != OLD.score)
EXECUTE FUNCTION fn_points_propagate_user();


-- tr_points_user_del

CREATE OR REPLACE FUNCTION fn_points_user_del()
RETURNS TRIGGER AS $$
BEGIN
  UPDATE teams
    SET score = score - OLD.score
    WHERE id = OLD.team_id;
  RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_points_user_del
BEFORE DELETE ON users
FOR EACH ROW
EXECUTE FUNCTION fn_points_user_del();
