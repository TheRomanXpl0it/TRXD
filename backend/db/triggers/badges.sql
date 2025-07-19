-- tr_badges_solve_insert

CREATE OR REPLACE FUNCTION fn_badges_solve_insert()
RETURNS TRIGGER AS $$
DECLARE
	team INTEGER;
	category_name TEXT;
BEGIN
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
WHEN (NEW.status = 'C')
EXECUTE FUNCTION fn_badges_solve_insert();


-- tr_badges_solve_del

CREATE OR REPLACE FUNCTION fn_badges_solve_del()
RETURNS TRIGGER AS $$
DECLARE
	team INTEGER;
	category_name TEXT;
BEGIN
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
WHEN (OLD.status = 'C')
EXECUTE FUNCTION fn_badges_solve_del();


-- tr_badges_add_and_del

CREATE OR REPLACE FUNCTION fn_badges_add_and_del()
RETURNS TRIGGER AS $$
DECLARE
	challs INTEGER;
BEGIN
	SELECT categories.chall_count INTO challs
		FROM categories
		WHERE categories.name = NEW.category;
	IF NEW.solves >= challs THEN
		INSERT INTO badges (name, description, team_id)
			VALUES (NEW.category, 'Completed all challenges', NEW.team_id);
	ELSE
		DELETE FROM badges
			WHERE name = NEW.category
				AND team_id = NEW.team_id;
	END IF;
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
	challs INTEGER;
	category_solves INTEGER;
BEGIN
	FOR team, category_solves IN (SELECT team_id, solves FROM team_category_solves WHERE category = NEW.name)
	LOOP
		SELECT chall_count INTO challs
			FROM categories
			WHERE categories.name = NEW.name;
		IF category_solves >= challs THEN
			INSERT INTO badges (name, description, team_id)
				VALUES (NEW.name, 'Completed all challenges', team);
		ELSE
			DELETE FROM badges
				WHERE name = NEW.name
					AND team_id = team;
		END IF;
	END LOOP;
	RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_badges_recompute
AFTER UPDATE ON categories
FOR EACH ROW
WHEN (NEW.chall_count != OLD.chall_count)
EXECUTE FUNCTION fn_badges_recompute();
