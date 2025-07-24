CREATE OR REPLACE FUNCTION assert(result BOOLEAN, msg TEXT DEFAULT NULL)
RETURNS VOID AS $$
BEGIN
  IF result != TRUE THEN
    RAISE EXCEPTION 'Assertion failed%Got: %',
      CASE WHEN msg IS NOT NULL THEN E' at ' || msg || E':\n' ELSE E'\n' END,
      result;
  END IF;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION delete_all()
RETURNS VOID AS $$
BEGIN
  DELETE FROM submissions;
  DELETE FROM tags;
  DELETE FROM instances;
  DELETE FROM flags;
  DELETE FROM docker_configs;
  DELETE FROM challenges;
  DELETE FROM team_category_solves;
  DELETE FROM categories;
  DELETE FROM badges;
  DELETE FROM users;
  DELETE FROM teams;
  DELETE FROM configs;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION insert_base_data()
RETURNS VOID AS $$
BEGIN
  /*
  categories:
    cat-1: chall-1, chall-3, chall-4
    cat-2: chall-2
  teams:
    A: a (player), b (player), e (master)
    B: c (player)
    C: f (author)
    no-team: d (player)
  */
  INSERT INTO configs (key, type, value) VALUES ('chall-min-points', 'int', '100');
  INSERT INTO configs (key, type, value) VALUES ('chall-points-decay', 'int', '5');
  INSERT INTO categories (name, icon) VALUES ('cat-1', 'cat-1');
  INSERT INTO categories (name, icon) VALUES ('cat-2', 'cat-2');
  INSERT INTO challenges (name, category, description, type, max_points, score_type, hidden) VALUES ('chall-1', 'cat-1', 'TEST', 'Normal', 500, 'Dynamic', false);
  INSERT INTO challenges (name, category, description, type, max_points, score_type, hidden) VALUES ('chall-2', 'cat-2', 'TEST', 'Normal', 500, 'Dynamic', false);
  INSERT INTO challenges (name, category, description, type, max_points, score_type, hidden) VALUES ('chall-3', 'cat-1', 'TEST', 'Container', 500, 'Dynamic', false);
  INSERT INTO challenges (name, category, description, type, max_points, score_type, hidden) VALUES ('chall-4', 'cat-1', 'TEST', 'Compose', 500, 'Dynamic', false);
  INSERT INTO teams (name, password_hash) VALUES ('Author', 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa');
  INSERT INTO teams (name, password_hash) VALUES ('B', 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa');
  INSERT INTO teams (name, password_hash) VALUES ('Correct', 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa');
  INSERT INTO users (name, email, password_hash, role, team_id) VALUES ('Author', 'a@a', 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa', 'Player', (SELECT id FROM teams WHERE name='Author'));
  INSERT INTO users (name, email, password_hash, role, team_id) VALUES ('b', 'b@b', 'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb', 'Player', (SELECT id FROM teams WHERE name='Author'));
  INSERT INTO users (name, email, password_hash, role, team_id) VALUES ('Correct', 'c@c', 'cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc', 'Player', (SELECT id FROM teams WHERE name='B'));
  INSERT INTO users (name, email, password_hash, role) VALUES ('Dynamic', 'd@d', 'dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd', 'Player');
  INSERT INTO users (name, email, password_hash, role, team_id) VALUES ('e', 'e@e', 'eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee', 'Admin', (SELECT id FROM teams WHERE name='Author'));
  INSERT INTO users (name, email, password_hash, role, team_id) VALUES ('f', 'f@f', 'ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff', 'Author', (SELECT id FROM teams WHERE name='Correct'));
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION tests()
RETURNS VOID AS $$
DECLARE
  tmp INTEGER;
BEGIN
  PERFORM delete_all();
  PERFORM insert_base_data();

  -- insert a wrong submission from 'Author' to 'chal-1'
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='Author'),
    (SELECT id FROM challenges WHERE name='chall-1'),
    'Wrong', 'flag');
  PERFORM assert(COUNT(b)=0) FROM badges b;
  PERFORM assert(COUNT(s)=1) FROM submissions s WHERE s.status='Wrong';
  PERFORM assert(c.solves=0) FROM challenges c WHERE c.name='chall-1';
  PERFORM assert(c.points=c.max_points) FROM challenges c WHERE c.name='chall-1';

  -- insert a repeated submission from 'Author' to 'chal-1'
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='Author'),
    (SELECT id FROM challenges WHERE name='chall-1'),
    'Repeated', 'flag');
  PERFORM assert(COUNT(b)=0) FROM badges b;
  PERFORM assert(COUNT(s)=1) FROM submissions s WHERE s.status='Repeated';
  PERFORM assert(c.solves=0) FROM challenges c WHERE c.name='chall-1';
  PERFORM assert(c.points=c.max_points) FROM challenges c WHERE c.name='chall-1';

  -- insert a valid submission from 'Author' to 'chal-1'
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='Author'),
    (SELECT id FROM challenges WHERE name='chall-1'),
    'Correct', 'flag');
  PERFORM assert(COUNT(b)=0) FROM badges b;
  PERFORM assert(COUNT(s)=1) FROM submissions s WHERE s.status='Correct';
  PERFORM assert(c.solves=1) FROM challenges c WHERE c.name='chall-1';
  PERFORM assert(c.points=c.max_points) FROM challenges c WHERE c.name='chall-1';
  
  -- insert a valid submission from 'Author' to 'chal-3'
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='Author'),
    (SELECT id FROM challenges WHERE name='chall-3'),
    'Correct', 'flag');
  PERFORM assert(COUNT(b)=0) FROM badges b;
  PERFORM assert(COUNT(s)=2) FROM submissions s WHERE s.status='Correct';
  PERFORM assert(c.solves=1) FROM challenges c WHERE c.name='chall-3';
  PERFORM assert(c.points=c.max_points) FROM challenges c WHERE c.name='chall-3';

  -- insert a valid submission from 'Author' to 'chal-4' (should give also a badge)
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='Author'),
    (SELECT id FROM challenges WHERE name='chall-4'),
    'Correct', 'flag');
  PERFORM assert(COUNT(b)=1) FROM badges b;
  PERFORM assert(COUNT(s)=3) FROM submissions s WHERE s.status='Correct';
  PERFORM assert(c.solves=1) FROM challenges c WHERE c.name='chall-4';
  PERFORM assert(c.points=c.max_points) FROM challenges c WHERE c.name='chall-4';
  
  -- insert a valid submission from 'Author' to 'chal-1' but it's repeated
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='Author'),
    (SELECT id FROM challenges WHERE name='chall-1'),
    'Correct', 'flag');
  PERFORM assert(COUNT(b)=1) FROM badges b;
  PERFORM assert(COUNT(s)=2) FROM submissions s WHERE s.status='Repeated';
  PERFORM assert(c.solves=1) FROM challenges c WHERE c.name='chall-1';
  PERFORM assert(c.points=c.max_points) FROM challenges c WHERE c.name='chall-1';

  -- insert a valid submission from 'b' to 'chal-3' but it's already solved by 'Author' (its teammate)
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='b'),
    (SELECT id FROM challenges WHERE name='chall-3'),
    'Correct', 'flag');
  PERFORM assert(COUNT(b)=1) FROM badges b;
  PERFORM assert(COUNT(s)=3) FROM submissions s WHERE s.status='Repeated';
  PERFORM assert(c.solves=1) FROM challenges c WHERE c.name='chall-3';
  PERFORM assert(c.points=c.max_points) FROM challenges c WHERE c.name='chall-3';

  -- insert a valid submission from 'Correct' to 'chal-4'
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='Correct'),
    (SELECT id FROM challenges WHERE name='chall-4'),
    'Correct', 'flag');
  PERFORM assert(COUNT(b)=1) FROM badges b;
  PERFORM assert(COUNT(s)=4) FROM submissions s WHERE s.status='Correct';
  PERFORM assert(c.solves=2) FROM challenges c WHERE c.name='chall-4';
  PERFORM assert(c.points<c.max_points) FROM challenges c WHERE c.name='chall-4';

  -- insert a valid submission from 'Correct' to 'chal-2' (should give another badge)
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='Correct'),
    (SELECT id FROM challenges WHERE name='chall-2'),
    'Correct', 'flag');
  PERFORM assert(COUNT(b)=2) FROM badges b;
  PERFORM assert(COUNT(s)=5) FROM submissions s WHERE s.status='Correct';
  PERFORM assert(c.solves=1) FROM challenges c WHERE c.name='chall-2';
  PERFORM assert(c.points=c.max_points) FROM challenges c WHERE c.name='chall-2';

  -- insert a valid submission from 'f' to 'chal-3' but it's not a player so it doesn't add up
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='f'),
    (SELECT id FROM challenges WHERE name='chall-3'),
    'Correct', 'flag');
  PERFORM assert(COUNT(b)=2) FROM badges b;
  PERFORM assert(COUNT(s)=6) FROM submissions s WHERE s.status='Correct';
  PERFORM assert(c.solves=1) FROM challenges c WHERE c.name='chall-3';

  -- deletes all correct submissions from 'Author' to 'chall-1' (should also remove the badge)
  DELETE FROM submissions WHERE user_id=(SELECT id FROM users WHERE name='Author')
    AND chall_id=(SELECT id FROM challenges WHERE name='chall-1')
    AND status='Correct';
  PERFORM assert(COUNT(b)=1) FROM badges b;
  PERFORM assert(COUNT(s)=5) FROM submissions s WHERE s.status='Correct';

  -- insert a valid submission from 'Author' to 'chal-1' (should give back the badge)
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='Author'),
    (SELECT id FROM challenges WHERE name='chall-1'),
    'Correct', 'flag');
  PERFORM assert(COUNT(b)=2) FROM badges b;
  PERFORM assert(COUNT(s)=6) FROM submissions s WHERE s.status='Correct';
  PERFORM assert(c.solves=1) FROM challenges c WHERE c.name='chall-1';

  -- insert a valid submission from 'Correct' to 'chal-1' (should not give a badge with 2 solves of 3)
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='Correct'),
    (SELECT id FROM challenges WHERE name='chall-1'),
    'Correct', 'flag');
  PERFORM assert(COUNT(b)=2) FROM badges b;
  PERFORM assert(COUNT(s)=7) FROM submissions s WHERE s.status='Correct';
  PERFORM assert(c.solves=2) FROM challenges c WHERE c.name='chall-1';

  -- check that the scores are the same
  SELECT score INTO tmp FROM teams WHERE name='Author';
  PERFORM assert(score=tmp) FROM teams WHERE name='B';

  -- after recomputing the scores, they should be lower
  UPDATE configs SET value='3' WHERE key='chall-points-decay';
  PERFORM assert(score<tmp) FROM teams WHERE name='B';

  -- after recomputing the scores, they should be greater
  UPDATE configs SET value='7' WHERE key='chall-points-decay';
  PERFORM assert(score>tmp) FROM teams WHERE name='B';

  -- after recomputing the scores, they should be all the same as 500
  UPDATE configs SET value='500' WHERE key='chall-min-points';
  PERFORM assert(COUNT(c)=4) FROM challenges c WHERE points=500;

  -- after changing the score to 1000, there should be no shlallenge left at 500 after recomputing the scores
  UPDATE challenges SET max_points=1000;
  PERFORM assert(COUNT(c)=0) FROM challenges c WHERE points=500;

  -- deletes 'chall-3' (should now give the badge to team 'B' that hadn't solved only that)
  DELETE FROM challenges c WHERE c.name='chall-3';
  PERFORM assert(COUNT(b)=3) FROM badges b;
  PERFORM assert(COUNT(s)=5) FROM submissions s WHERE s.status='Correct';

  -- inserts back 'chall-3' as hidden (nothing should change)
  INSERT INTO challenges (name, category, description, type, max_points, score_type)
    VALUES ('chall-3', 'cat-1', 'TEST', 'Container', 1000, 'Dynamic');
  PERFORM assert(COUNT(b)=3) FROM badges b;

  -- makes 'chall-3' visible again, now should remove the badges (previous solves were deleted on the DELETE as cascade)
  UPDATE challenges SET hidden=false WHERE name='chall-3';
  PERFORM assert(COUNT(b)=1) FROM badges b;

  -- inserts again new valid submissions, should give badges and lower the challenge points
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='Author'),
    (SELECT id FROM challenges WHERE name='chall-3'),
    'Correct', 'flag');
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='Correct'),
    (SELECT id FROM challenges WHERE name='chall-3'),
    'Correct', 'flag');
  PERFORM assert(COUNT(b)=3) FROM badges b;
  PERFORM assert(COUNT(s)=7) FROM submissions s WHERE s.status='Correct';
  PERFORM assert(points!=max_points) FROM challenges WHERE name='chall-3';

  -- changes 'chall-3' scoring type as static, so the points should reset
  UPDATE challenges SET score_type='Static' WHERE name='chall-3';
  PERFORM assert(points=max_points) FROM challenges WHERE name='chall-3';

  -- changes 'chall-3' max_points so the points should change as them
  UPDATE challenges SET max_points=500 WHERE name='chall-3';
  PERFORM assert(points=max_points) FROM challenges WHERE name='chall-3';

  -- inserts a correct submission from 'Dynamic' on 'chall-1', but it doesn't have a team so it's an invalid submission
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='Dynamic'),
    (SELECT id FROM challenges WHERE name='chall-1'),
    'Correct', 'flag');
  PERFORM assert(COUNT(s)=7) FROM submissions s WHERE s.status='Correct';
  PERFORM assert(COUNT(s)=1) FROM submissions s WHERE s.status='Invalid';

  -- removes the user 'Correct', so all his submissions should be deleted and also badges and team points removed
  DELETE FROM users WHERE name='Correct';
  PERFORM assert(COUNT(b)=1) FROM badges b;
  PERFORM assert(COUNT(s)=3) FROM submissions s WHERE s.status='Correct';
  PERFORM assert(score=0) FROM teams WHERE name='B';

  -- checks that 'chall-3' and 'chall-4' have their configs created
  PERFORM assert(count(d)=1) FROM docker_configs d WHERE  d.chall_id=(SELECT id FROM challenges WHERE name='chall-3');
  PERFORM assert(count(d)=1) FROM docker_configs d WHERE  d.chall_id=(SELECT id FROM challenges WHERE name='chall-4');

END;
$$ LANGUAGE plpgsql;

/*
SELECT tests();

SELECT * FROM submissions;
SELECT * FROM tags;
SELECT * FROM instances;
SELECT * FROM flags;
SELECT * FROM docker_configs;
SELECT * FROM challenges;
SELECT * FROM team_category_solves;
SELECT * FROM categories;
SELECT * FROM badges;
SELECT * FROM users;
SELECT * FROM teams;
SELECT * FROM configs;
--*/
