CREATE OR REPLACE FUNCTION fn_categories_add_chall()
RETURNS TRIGGER AS $$
BEGIN
	UPDATE categories
		SET chall_count = chall_count + 1
		WHERE name = NEW.category;
	RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION fn_categories_del_chall()
RETURNS TRIGGER AS $$
BEGIN
	UPDATE categories
		SET chall_count = chall_count - 1
		WHERE name = OLD.category;
	RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION fn_categories_update_chall()
RETURNS TRIGGER AS $$
BEGIN
	UPDATE categories
		SET chall_count = chall_count - 1
		WHERE name = OLD.category;
	UPDATE categories
		SET chall_count = chall_count + 1
		WHERE name = NEW.category;
	RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_categories_add_chall
AFTER INSERT ON challenges
FOR EACH ROW
EXECUTE FUNCTION fn_categories_add_chall();

CREATE TRIGGER tr_categories_del_chall
AFTER DELETE ON challenges
FOR EACH ROW
EXECUTE FUNCTION fn_categories_del_chall();

CREATE TRIGGER tr_categories_update_chall
AFTER UPDATE ON challenges
FOR EACH ROW
EXECUTE FUNCTION fn_categories_update_chall();
