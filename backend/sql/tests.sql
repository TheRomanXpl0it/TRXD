CREATE OR REPLACE FUNCTION assert(result BOOLEAN, msg TEXT DEFAULT NULL)
RETURNS VOID AS $$
BEGIN
  IF result != TRUE THEN
    RAISE EXCEPTION 'Assertion failed%Got: %',
      CASE WHEN msg IS NOT NULL THEN E':\n' || msg || E':\n' ELSE E'\n' END,
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


CREATE OR REPLACE FUNCTION insert_mock_configs()
RETURNS VOID AS $$
BEGIN
  INSERT INTO configs (key, type, value) VALUES ('chall-min-points', 'int', '100');
  INSERT INTO configs (key, type, value) VALUES ('chall-points-decay', 'int', '5');
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION insert_mock_data()
RETURNS VOID AS $$
BEGIN
  /*
  categories:
    cat-1: chall-1, chall-3, chall-4
    cat-2: chall-2, chall-5
  teams:
    A: a (player), b (player), e (master)
    B: c (player)
    C: f (author)
    no-team: d (player)
  */
  INSERT INTO categories (name, icon) VALUES ('cat-1', 'cat-1');
  INSERT INTO categories (name, icon) VALUES ('cat-2', 'cat-2');
  INSERT INTO challenges (name, category, description, difficulty, authors, type, max_points, score_type, host, port, hidden) VALUES ('chall-1', 'cat-1', 'TEST chall-1 DESC', 'Easy', E'author1\x01author2', 'Normal', 500, 'Dynamic', 'http://theromanxpl0.it', 1234, false);
  INSERT INTO challenges (name, category, description, difficulty, authors, type, max_points, score_type, hidden) VALUES ('chall-2', 'cat-2', 'TEST chall-2 DESC', 'Medium', E'author1\x01author2\x01author3', 'Normal', 500, 'Dynamic', false);
  INSERT INTO challenges (name, category, description, difficulty, authors, type, max_points, score_type, host, port, hidden) VALUES ('chall-3', 'cat-1', 'TEST chall-3 DESC', 'Hard', 'author1', 'Container', 500, 'Dynamic', 'chall-3.test.com', 1337, false);
  INSERT INTO challenges (name, category, description, difficulty, authors, type, max_points, score_type, hidden) VALUES ('chall-4', 'cat-1', 'TEST chall-4 DESC', 'Insane', 'author2', 'Compose', 500, 'Dynamic', false);
  INSERT INTO challenges (name, category, description, difficulty, authors, type, max_points, score_type) VALUES ('chall-5', 'cat-2', 'TEST chall-5 DESC', 'Easy', 'author3', 'Normal', 500, 'Static');
  UPDATE docker_configs SET image='echo-server:latest', hash_domain=TRUE WHERE chall_id=(SELECT id FROM challenges WHERE name='chall-3');
  UPDATE docker_configs SET compose='
services:
  app:
    image: echo-server:latest
    container_name: chall_333333333333
    ports:
      - "1337:1337"
    environment:
      - ECHO_MESSAGE=Hello from app
      - INSTANCE_PORT=${INSTANCE_PORT}
      - INSTANCE_HOST=${INSTANCE_HOST}
    ', hash_domain=TRUE WHERE chall_id=(SELECT id FROM challenges WHERE name='chall-4');
  INSERT INTO tags (name, chall_id) VALUES ('tag-1', (SELECT id FROM challenges WHERE name='chall-1'));
  INSERT INTO tags (name, chall_id) VALUES ('test-tag', (SELECT id FROM challenges WHERE name='chall-1'));
  INSERT INTO tags (name, chall_id) VALUES ('tag-2', (SELECT id FROM challenges WHERE name='chall-2'));
  INSERT INTO tags (name, chall_id) VALUES ('tag-3', (SELECT id FROM challenges WHERE name='chall-3'));
  INSERT INTO tags (name, chall_id) VALUES ('tag-4', (SELECT id FROM challenges WHERE name='chall-4'));
  INSERT INTO tags (name, chall_id) VALUES ('tag-5', (SELECT id FROM challenges WHERE name='chall-5'));
  INSERT INTO flags (flag, chall_id) VALUES ('flag{test-1}', (SELECT id FROM challenges WHERE name='chall-1'));
  INSERT INTO flags (flag, chall_id, regex) VALUES ('flag\{test-[a-z]{2}\}', (SELECT id FROM challenges WHERE name='chall-1'), true);
  INSERT INTO flags (flag, chall_id) VALUES ('flag{test-2}', (SELECT id FROM challenges WHERE name='chall-2'));
  INSERT INTO flags (flag, chall_id) VALUES ('flag{test-3}', (SELECT id FROM challenges WHERE name='chall-3'));
  INSERT INTO flags (flag, chall_id) VALUES ('flag{test-4}', (SELECT id FROM challenges WHERE name='chall-4'));
  INSERT INTO flags (flag, chall_id) VALUES ('flag{test-5}', (SELECT id FROM challenges WHERE name='chall-5'));
  -- password:'testpass' 
  INSERT INTO teams (name, password_hash, password_salt) VALUES ('A', '41d65efe433e60755bef957e56ed6466b24c44a86b8ec595df4c9cdfa9c3aca9', '1a5e93869fa3c2ee04139db8834f8808');
  INSERT INTO teams (name, password_hash, password_salt) VALUES ('B', '41d65efe433e60755bef957e56ed6466b24c44a86b8ec595df4c9cdfa9c3aca9', '1a5e93869fa3c2ee04139db8834f8808');
  INSERT INTO teams (name, password_hash, password_salt) VALUES ('C', '41d65efe433e60755bef957e56ed6466b24c44a86b8ec595df4c9cdfa9c3aca9', '1a5e93869fa3c2ee04139db8834f8808');
  INSERT INTO users (name, email, password_hash, password_salt, role, team_id) VALUES ('a', 'a@a.a', '41d65efe433e60755bef957e56ed6466b24c44a86b8ec595df4c9cdfa9c3aca9', '1a5e93869fa3c2ee04139db8834f8808', 'Player', (SELECT id FROM teams WHERE name='A'));
  INSERT INTO users (name, email, password_hash, password_salt, role, team_id) VALUES ('b', 'b@b.b', '41d65efe433e60755bef957e56ed6466b24c44a86b8ec595df4c9cdfa9c3aca9', '1a5e93869fa3c2ee04139db8834f8808', 'Player', (SELECT id FROM teams WHERE name='A'));
  INSERT INTO users (name, email, password_hash, password_salt, role, team_id) VALUES ('c', 'c@c.c', '41d65efe433e60755bef957e56ed6466b24c44a86b8ec595df4c9cdfa9c3aca9', '1a5e93869fa3c2ee04139db8834f8808', 'Player', (SELECT id FROM teams WHERE name='B'));
  INSERT INTO users (name, email, password_hash, password_salt, role) VALUES ('d', 'd@d.d', '41d65efe433e60755bef957e56ed6466b24c44a86b8ec595df4c9cdfa9c3aca9', '1a5e93869fa3c2ee04139db8834f8808', 'Player');
  INSERT INTO users (name, email, password_hash, password_salt, role, team_id) VALUES ('e', 'admin@email.com', '41d65efe433e60755bef957e56ed6466b24c44a86b8ec595df4c9cdfa9c3aca9', '1a5e93869fa3c2ee04139db8834f8808', 'Admin', (SELECT id FROM teams WHERE name='A'));
  INSERT INTO users (name, email, password_hash, password_salt, role, team_id) VALUES ('f', 'f@f.f', '41d65efe433e60755bef957e56ed6466b24c44a86b8ec595df4c9cdfa9c3aca9', '1a5e93869fa3c2ee04139db8834f8808', 'Author', (SELECT id FROM teams WHERE name='C'));
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION insert_mock_submissions()
RETURNS VOID AS $$
BEGIN
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='a'),
    (SELECT id FROM challenges WHERE name='chall-1'),
    'Wrong', 'flag');
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='a'),
    (SELECT id FROM challenges WHERE name='chall-1'),
    'Repeated', 'flag');
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='a'),
    (SELECT id FROM challenges WHERE name='chall-1'),
    'Correct', 'flag');
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='a'),
    (SELECT id FROM challenges WHERE name='chall-3'),
    'Correct', 'flag');
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='a'),
    (SELECT id FROM challenges WHERE name='chall-4'),
    'Correct', 'flag');
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='a'),
    (SELECT id FROM challenges WHERE name='chall-4'),
    'Correct', 'flag');
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='a'),
    (SELECT id FROM challenges WHERE name='chall-1'),
    'Correct', 'flag');
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='b'),
    (SELECT id FROM challenges WHERE name='chall-3'),
    'Correct', 'flag');
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='c'),
    (SELECT id FROM challenges WHERE name='chall-4'),
    'Correct', 'flag');
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='c'),
    (SELECT id FROM challenges WHERE name='chall-2'),
    'Correct', 'flag');
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='f'),
    (SELECT id FROM challenges WHERE name='chall-3'),
    'Correct', 'flag');
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='d'),
    (SELECT id FROM challenges WHERE name='chall-1'),
    'Correct', 'flag');
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION tests()
RETURNS VOID AS $$
DECLARE
  tmp INTEGER;
BEGIN
  PERFORM delete_all();
  PERFORM insert_mock_configs();
  PERFORM insert_mock_data();

  -- insert a wrong submission from 'a' to 'chal-1'
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='a'),
    (SELECT id FROM challenges WHERE name='chall-1'),
    'Wrong', 'flag');
  PERFORM assert(COUNT(b)=0, 'check 1') FROM badges b;
  PERFORM assert(COUNT(s)=1, 'check 2') FROM submissions s WHERE s.status='Wrong';
  PERFORM assert(c.solves=0, 'check 3') FROM challenges c WHERE c.name='chall-1';
  PERFORM assert(c.points=c.max_points, 'check 4') FROM challenges c WHERE c.name='chall-1';

  -- insert a repeated submission from 'a' to 'chal-1'
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='a'),
    (SELECT id FROM challenges WHERE name='chall-1'),
    'Repeated', 'flag');
  PERFORM assert(COUNT(b)=0, 'check 6') FROM badges b;
  PERFORM assert(COUNT(s)=1, 'check 7') FROM submissions s WHERE s.status='Repeated';
  PERFORM assert(c.solves=0, 'check 8') FROM challenges c WHERE c.name='chall-1';
  PERFORM assert(c.points=c.max_points, 'check 9') FROM challenges c WHERE c.name='chall-1';

  -- insert a valid submission from 'a' to 'chal-1'
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='a'),
    (SELECT id FROM challenges WHERE name='chall-1'),
    'Correct', 'flag');
  PERFORM assert(COUNT(b)=0, 'check 10') FROM badges b;
  PERFORM assert(COUNT(s)=1, 'check 11') FROM submissions s WHERE s.status='Correct';
  PERFORM assert(c.solves=1, 'check 12') FROM challenges c WHERE c.name='chall-1';
  PERFORM assert(c.points=c.max_points, 'check 13') FROM challenges c WHERE c.name='chall-1';

  -- insert a valid submission from 'a' to 'chal-3'
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='a'),
    (SELECT id FROM challenges WHERE name='chall-3'),
    'Correct', 'flag');
  PERFORM assert(COUNT(b)=0, 'check 14') FROM badges b;
  PERFORM assert(COUNT(s)=2, 'check 15') FROM submissions s WHERE s.status='Correct';
  PERFORM assert(c.solves=1, 'check 16') FROM challenges c WHERE c.name='chall-3';
  PERFORM assert(c.points=c.max_points, 'check 17') FROM challenges c WHERE c.name='chall-3';

  -- insert a valid submission from 'a' to 'chal-4' (should give also a badge)
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='a'),
    (SELECT id FROM challenges WHERE name='chall-4'),
    'Correct', 'flag');
  PERFORM assert(COUNT(b)=1, 'check 18') FROM badges b;
  PERFORM assert(COUNT(s)=3, 'check 19') FROM submissions s WHERE s.status='Correct';
  PERFORM assert(c.solves=1, 'check 20') FROM challenges c WHERE c.name='chall-4';
  PERFORM assert(c.points=c.max_points, 'check 21') FROM challenges c WHERE c.name='chall-4';

  -- insert a valid submission from 'a' to 'chal-1' but it's repeated
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='a'),
    (SELECT id FROM challenges WHERE name='chall-1'),
    'Correct', 'flag');
  PERFORM assert(COUNT(b)=1, 'check 22') FROM badges b;
  PERFORM assert(COUNT(s)=2, 'check 23') FROM submissions s WHERE s.status='Repeated';
  PERFORM assert(c.solves=1, 'check 24') FROM challenges c WHERE c.name='chall-1';
  PERFORM assert(c.points=c.max_points, 'check 25') FROM challenges c WHERE c.name='chall-1';

  -- insert a valid submission from 'b' to 'chal-3' but it's already solved by 'a' (its teammate)
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='b'),
    (SELECT id FROM challenges WHERE name='chall-3'),
    'Correct', 'flag');
  PERFORM assert(COUNT(b)=1, 'check 26') FROM badges b;
  PERFORM assert(COUNT(s)=3, 'check 27') FROM submissions s WHERE s.status='Repeated';
  PERFORM assert(c.solves=1, 'check 28') FROM challenges c WHERE c.name='chall-3';
  PERFORM assert(c.points=c.max_points, 'check 29') FROM challenges c WHERE c.name='chall-3';

  -- insert a valid submission from 'c' to 'chal-4'
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='c'),
    (SELECT id FROM challenges WHERE name='chall-4'),
    'Correct', 'flag');
  PERFORM assert(COUNT(b)=1, 'check 30') FROM badges b;
  PERFORM assert(COUNT(s)=4, 'check 31') FROM submissions s WHERE s.status='Correct';
  PERFORM assert(c.solves=2, 'check 32') FROM challenges c WHERE c.name='chall-4';
  PERFORM assert(c.points<c.max_points, 'check 33') FROM challenges c WHERE c.name='chall-4';

  -- insert a valid submission from 'c' to 'chal-2' (should give another badge)
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='c'),
    (SELECT id FROM challenges WHERE name='chall-2'),
    'Correct', 'flag');
  PERFORM assert(COUNT(b)=2, 'check 34') FROM badges b;
  PERFORM assert(COUNT(s)=5, 'check 35') FROM submissions s WHERE s.status='Correct';
  PERFORM assert(c.solves=1, 'check 36') FROM challenges c WHERE c.name='chall-2';
  PERFORM assert(c.points=c.max_points, 'check 37') FROM challenges c WHERE c.name='chall-2';

  -- insert a valid submission from 'f' to 'chal-3' but it's not a player so it doesn't add up
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='f'),
    (SELECT id FROM challenges WHERE name='chall-3'),
    'Correct', 'flag');
  PERFORM assert(COUNT(b)=2, 'check 38') FROM badges b;
  PERFORM assert(COUNT(s)=6, 'check 39') FROM submissions s WHERE s.status='Correct';
  PERFORM assert(c.solves=1, 'check 40') FROM challenges c WHERE c.name='chall-3';

  -- deletes all correct submissions from 'a' to 'chall-1' (should also remove the badge)
  DELETE FROM submissions WHERE user_id=(SELECT id FROM users WHERE name='a')
    AND chall_id=(SELECT id FROM challenges WHERE name='chall-1')
    AND status='Correct';
  PERFORM assert(COUNT(b)=1, 'check 41') FROM badges b;
  PERFORM assert(COUNT(s)=5, 'check 42') FROM submissions s WHERE s.status='Correct';

  -- insert a valid submission from 'a' to 'chal-1' (should give back the badge)
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='a'),
    (SELECT id FROM challenges WHERE name='chall-1'),
    'Correct', 'flag');
  PERFORM assert(COUNT(b)=2, 'check 43') FROM badges b;
  PERFORM assert(COUNT(s)=6, 'check 44') FROM submissions s WHERE s.status='Correct';
  PERFORM assert(c.solves=1, 'check 45') FROM challenges c WHERE c.name='chall-1';

  -- insert a valid submission from 'c' to 'chal-1' (should not give a badge with 2 solves of 3)
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='c'),
    (SELECT id FROM challenges WHERE name='chall-1'),
    'Correct', 'flag');
  PERFORM assert(COUNT(b)=2, 'check 46') FROM badges b;
  PERFORM assert(COUNT(s)=7, 'check 47') FROM submissions s WHERE s.status='Correct';
  PERFORM assert(c.solves=2, 'check 48') FROM challenges c WHERE c.name='chall-1';

  -- check that the scores are the same
  SELECT score INTO tmp FROM teams WHERE name='A';
  PERFORM assert(score=tmp, 'check 49') FROM teams WHERE name='B';

  -- after recomputing the scores, they should be lower
  UPDATE configs SET value='3' WHERE key='chall-points-decay';
  PERFORM assert(score<tmp, 'check 50') FROM teams WHERE name='B';

  -- after recomputing the scores, they should be greater
  UPDATE configs SET value='7' WHERE key='chall-points-decay';
  PERFORM assert(score>tmp, 'check 51') FROM teams WHERE name='B';

  -- after recomputing the scores, they should be all the same as 500
  UPDATE configs SET value='500' WHERE key='chall-min-points';
  PERFORM assert(COUNT(c)=5, 'check 52') FROM challenges c WHERE points=500;

  -- after changing the score to 1000, there should be no challenge left at 500 after recomputing the scores
  UPDATE challenges SET max_points=1000;
  PERFORM assert(COUNT(c)=0, 'check 53') FROM challenges c WHERE points=500;

  -- deletes 'chall-3' (should now give the badge to team 'B' that hadn't solved only that)
  DELETE FROM challenges c WHERE c.name='chall-3';
  PERFORM assert(COUNT(b)=3, 'check 54') FROM badges b;
  PERFORM assert(COUNT(s)=5, 'check 55') FROM submissions s WHERE s.status='Correct';

  -- inserts back 'chall-3' as hidden (nothing should change)
  INSERT INTO challenges (name, category, description, type, max_points, score_type)
    VALUES ('chall-3', 'cat-1', 'TEST', 'Container', 1000, 'Dynamic');
  PERFORM assert(COUNT(b)=3, 'check 56') FROM badges b;

  -- makes 'chall-3' visible again, now should remove the badges (previous solves were deleted on the DELETE as cascade)
  UPDATE challenges SET hidden=false WHERE name='chall-3';
  PERFORM assert(COUNT(b)=1, 'check 57') FROM badges b;

  -- inserts again new valid submissions, should give badges and lower the challenge points
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='a'),
    (SELECT id FROM challenges WHERE name='chall-3'),
    'Correct', 'flag');
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='c'),
    (SELECT id FROM challenges WHERE name='chall-3'),
    'Correct', 'flag');
  PERFORM assert(COUNT(b)=3, 'check 58') FROM badges b;
  PERFORM assert(COUNT(s)=7, 'check 59') FROM submissions s WHERE s.status='Correct';
  PERFORM assert(points!=max_points, 'check 60') FROM challenges WHERE name='chall-3';

  -- changes 'chall-3' scoring type as static, so the points should reset
  UPDATE challenges SET score_type='Static' WHERE name='chall-3';
  PERFORM assert(points=max_points, 'check 61') FROM challenges WHERE name='chall-3';

  -- changes 'chall-3' max_points so the points should change as them
  UPDATE challenges SET max_points=500 WHERE name='chall-3';
  PERFORM assert(points=max_points, 'check 62') FROM challenges WHERE name='chall-3';

  -- inserts a correct submission from 'd' on 'chall-1', but it doesn't have a team so it's an invalid submission
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='d'),
    (SELECT id FROM challenges WHERE name='chall-1'),
    'Correct', 'flag');
  PERFORM assert(COUNT(s)=7, 'check 63') FROM submissions s WHERE s.status='Correct';
  PERFORM assert(COUNT(s)=1, 'check 64') FROM submissions s WHERE s.status='Invalid';

  -- removes the user 'Correct', so all his submissions should be deleted and also badges and team points removed
  DELETE FROM users WHERE name='c';
  PERFORM assert(COUNT(b)=1, 'check 65') FROM badges b;
  PERFORM assert(COUNT(s)=3, 'check 66') FROM submissions s WHERE s.status='Correct';
  PERFORM assert(score=0, 'check 67') FROM teams WHERE name='B';

  -- checks that 'chall-3' and 'chall-4' have their configs created
  PERFORM assert(count(d)=1, 'check 68') FROM docker_configs d WHERE d.chall_id=(SELECT id FROM challenges WHERE name='chall-3');
  PERFORM assert(count(d)=1, 'check 69') FROM docker_configs d WHERE d.chall_id=(SELECT id FROM challenges WHERE name='chall-4');

  -- checks that docker configs are created on type update and not duplicated
  INSERT INTO challenges (name, category, description, type, max_points, score_type) VALUES ('chall-test', 'cat-1', 'TEST', 'Normal', 500, 'Dynamic');
  PERFORM assert(count(d)=0, 'check 70') FROM docker_configs d WHERE d.chall_id=(SELECT id FROM challenges WHERE name='chall-test');
  UPDATE challenges SET type='Container' WHERE id=(SELECT id FROM challenges WHERE name='chall-test');
  PERFORM assert(count(d)=1, 'check 71') FROM docker_configs d WHERE d.chall_id=(SELECT id FROM challenges WHERE name='chall-test');
  UPDATE challenges SET type='Normal' WHERE id=(SELECT id FROM challenges WHERE name='chall-test');
  PERFORM assert(count(d)=1, 'check 72') FROM docker_configs d WHERE d.chall_id=(SELECT id FROM challenges WHERE name='chall-test');
  UPDATE challenges SET type='Compose' WHERE id=(SELECT id FROM challenges WHERE name='chall-test');
  PERFORM assert(count(d)=1, 'check 73') FROM docker_configs d WHERE d.chall_id=(SELECT id FROM challenges WHERE name='chall-test');

END;
$$ LANGUAGE plpgsql;

/*
SELECT tests();
SELECT insert_mock_data();

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
