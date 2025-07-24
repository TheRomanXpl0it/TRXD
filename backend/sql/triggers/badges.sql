-- utils

CREATE OR REPLACE FUNCTION fn_badges_handler(category VARCHAR, team INTEGER, category_solves INTEGER)
RETURNS VOID AS $$
DECLARE
  challs INTEGER;
BEGIN
  SELECT categories.visible_challs INTO challs
    FROM categories
    WHERE categories.name = category;
  
  IF category_solves >= challs THEN
    IF NOT EXISTS(SELECT 1 FROM badges WHERE name = category AND team_id = team) THEN
      INSERT INTO badges (name, description, team_id)
        VALUES (category, 'Completed all ' || category || ' challenges', team);
    END IF;
  ELSE
    DELETE FROM badges
      WHERE name = category
        AND team_id = team;
  END IF;
END;
$$ LANGUAGE plpgsql;


-- tr_badges_solve_insert

CREATE OR REPLACE FUNCTION fn_badges_solve_insert()
RETURNS TRIGGER AS $$
DECLARE
  team INTEGER;
  category_name VARCHAR;
BEGIN
  IF (SELECT role FROM users WHERE id = NEW.user_id) != 'Player' THEN
    RETURN NEW;
  END IF;
  
  SELECT users.team_id, challenges.category
    INTO team, category_name
    FROM users
    JOIN challenges ON challenges.id = NEW.chall_id
    WHERE users.id = NEW.user_id;
  
  UPDATE team_category_solves
    SET solves = solves + 1
    WHERE team_id = team
      AND category = category_name;
  
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_badges_solve_insert
AFTER INSERT ON submissions
FOR EACH ROW
WHEN (NEW.status = 'Correct')
EXECUTE FUNCTION fn_badges_solve_insert();


-- tr_badges_solve_del

CREATE OR REPLACE FUNCTION fn_badges_solve_del()
RETURNS TRIGGER AS $$
DECLARE
  team INTEGER;
  category_name VARCHAR;
BEGIN
  IF (SELECT role FROM users WHERE id = OLD.user_id) != 'Player' THEN
    RETURN OLD;
  END IF;

  SELECT users.team_id, challenges.category
    INTO team, category_name
    FROM users
    JOIN challenges ON challenges.id = OLD.chall_id
    WHERE users.id = OLD.user_id;

  UPDATE team_category_solves
    SET solves = solves - 1
    WHERE team_id = team
      AND category = category_name;

  RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_badges_solve_del
AFTER DELETE ON submissions
FOR EACH ROW
WHEN (OLD.status = 'Correct')
EXECUTE FUNCTION fn_badges_solve_del();


-- tr_badges_chall_del

CREATE OR REPLACE FUNCTION fn_badges_chall_del()
RETURNS TRIGGER AS $$
BEGIN
  UPDATE team_category_solves
    SET solves = team_category_solves.solves - 1
    FROM users
    JOIN submissions ON submissions.user_id = users.id
      AND submissions.chall_id = OLD.id
    WHERE team_category_solves.category = OLD.category
      AND team_category_solves.team_id = users.team_id
      AND users.role = 'Player'
      AND submissions.status = 'Correct';
  RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_badges_chall_del
BEFORE DELETE ON challenges
FOR EACH ROW
EXECUTE FUNCTION fn_badges_chall_del();


-- tr_badges_user_del

CREATE OR REPLACE FUNCTION fn_badges_user_del()
RETURNS TRIGGER AS $$
BEGIN
  IF OLD.role != 'Player' THEN
    RETURN OLD;
  END IF;

  UPDATE team_category_solves
    SET solves = team_category_solves.solves - 1
    FROM challenges
    JOIN submissions ON submissions.chall_id = challenges.id
      AND submissions.user_id = OLD.id
    WHERE team_category_solves.category = challenges.category
      AND team_category_solves.team_id = OLD.team_id
      AND submissions.status = 'Correct';

  RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_badges_user_del
BEFORE DELETE ON users
FOR EACH ROW
EXECUTE FUNCTION fn_badges_user_del();


-- tr_badges_add_and_del

CREATE OR REPLACE FUNCTION fn_badges_add_and_del()
RETURNS TRIGGER AS $$
BEGIN
  PERFORM fn_badges_handler(NEW.category, NEW.team_id, NEW.solves);
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_badges_add_and_del
AFTER UPDATE ON team_category_solves
FOR EACH ROW
WHEN (NEW.solves != OLD.solves)
EXECUTE FUNCTION fn_badges_add_and_del();


-- tr_badges_recompute

CREATE OR REPLACE FUNCTION fn_badges_recompute()
RETURNS TRIGGER AS $$
DECLARE
  team INTEGER;
  category_solves INTEGER;
BEGIN
  FOR team, category_solves IN (SELECT team_id, solves FROM team_category_solves WHERE category = NEW.name)
  LOOP
    PERFORM fn_badges_handler(NEW.name, team, category_solves);
  END LOOP;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_badges_recompute
AFTER UPDATE ON categories
FOR EACH ROW
WHEN (NEW.visible_challs != OLD.visible_challs)
EXECUTE FUNCTION fn_badges_recompute();
