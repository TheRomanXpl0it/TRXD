-- tr_integrity_repeat_solve

CREATE OR REPLACE FUNCTION fn_integrity_repeat_solve()
RETURNS TRIGGER AS $$
DECLARE
  team INTEGER;
  existing_correct_count INTEGER;
BEGIN
  SELECT team_id INTO team
    FROM users
    WHERE id = NEW.user_id;
  IF team IS NULL THEN
    NEW.status = 'I';
    RETURN NEW;
  END IF;
  SELECT COUNT(*) INTO existing_correct_count
    FROM submissions
    JOIN users ON users.id = submissions.user_id
    WHERE users.team_id = team
      AND users.role = 'P'
      AND submissions.chall_id = NEW.chall_id
      AND submissions.status = 'C';
  IF existing_correct_count > 0 THEN
    NEW.status = 'R';
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_integrity_repeat_solve
BEFORE INSERT ON submissions
FOR EACH ROW
WHEN (NEW.status = 'C')
EXECUTE FUNCTION fn_integrity_repeat_solve();
