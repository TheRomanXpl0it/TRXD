-- tr_categories_add_chall

CREATE OR REPLACE FUNCTION fn_categories_add_chall()
RETURNS TRIGGER AS $$
BEGIN
  UPDATE categories
    SET visible_challs = visible_challs + 1
    WHERE name = NEW.category;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_categories_add_chall
AFTER INSERT ON challenges
FOR EACH ROW
WHEN (NEW.hidden = FALSE)
EXECUTE FUNCTION fn_categories_add_chall();


-- tr_categories_del_chall

CREATE OR REPLACE FUNCTION fn_categories_del_chall()
RETURNS TRIGGER AS $$
BEGIN
  UPDATE categories
    SET visible_challs = visible_challs - 1
    WHERE name = OLD.category;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_categories_del_chall
AFTER DELETE ON challenges
FOR EACH ROW
WHEN (OLD.hidden = FALSE)
EXECUTE FUNCTION fn_categories_del_chall();


-- tr_categories_update_chall

CREATE OR REPLACE FUNCTION fn_categories_update_chall()
RETURNS TRIGGER AS $$
BEGIN
  IF NEW.hidden = FALSE THEN
    UPDATE categories
      SET visible_challs = visible_challs + 1
      WHERE name = NEW.category;
  END IF;
  
  IF OLD.hidden = FALSE THEN
    UPDATE categories
      SET visible_challs = visible_challs - 1
      WHERE name = OLD.category;
  END IF;
  
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_categories_update_chall
AFTER UPDATE ON challenges
FOR EACH ROW
WHEN (NEW.hidden != OLD.hidden OR NEW.category != OLD.category)
EXECUTE FUNCTION fn_categories_update_chall();


-- tr_categories_add

CREATE OR REPLACE FUNCTION fn_categories_add()
RETURNS TRIGGER AS $$
BEGIN
  INSERT INTO team_category_solves (team_id, category)
    SELECT id, NEW.name
    FROM teams;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_categories_add
AFTER INSERT ON categories
FOR EACH ROW
EXECUTE FUNCTION fn_categories_add();


-- tr_categories_add_team

CREATE OR REPLACE FUNCTION fn_categories_add_team()
RETURNS TRIGGER AS $$
BEGIN
  INSERT INTO team_category_solves (team_id, category)
    SELECT NEW.id, name
    FROM categories;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_categories_add_team
AFTER INSERT ON teams
FOR EACH ROW
EXECUTE FUNCTION fn_categories_add_team();
