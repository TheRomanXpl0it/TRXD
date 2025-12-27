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
  DELETE FROM instances;
  DELETE FROM flags;
  DELETE FROM attachments;
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
  INSERT INTO categories (name) VALUES ('cat-1');
  INSERT INTO categories (name) VALUES ('cat-2');
  INSERT INTO challenges (name, category, description, difficulty, authors, tags, type, max_points, score_type, host, port, conn_type, hidden) VALUES ('chall-1', 'cat-1', 'TEST chall-1 DESC', 'Easy', ARRAY['author1', 'author2'], ARRAY['tag-1', 'test-tag'], 'Normal', 500, 'Dynamic', 'ctf.theromanxpl0.it', 1234, 'TCP', false);
  INSERT INTO challenges (name, category, description, difficulty, authors, tags, type, max_points, score_type, hidden) VALUES ('chall-2', 'cat-2', 'TEST chall-2 DESC', 'Medium', ARRAY['author1', 'author2', 'author3'], ARRAY['tag-2'], 'Normal', 500, 'Dynamic', false);
  INSERT INTO challenges (name, category, description, difficulty, authors, tags, type, max_points, score_type, host, port, conn_type, hidden) VALUES ('chall-3', 'cat-1', 'TEST chall-3 DESC', 'Hard', ARRAY['author1'], ARRAY['tag-3'], 'Container', 500, 'Dynamic', 'chall-3.test.com', 1337, 'HTTP', false);
  INSERT INTO challenges (name, category, description, difficulty, authors, tags, type, max_points, score_type, conn_type, hidden) VALUES ('chall-4', 'cat-1', 'TEST chall-4 DESC', 'Insane', ARRAY['author2'], ARRAY['tag-4'], 'Compose', 500, 'Dynamic', 'HTTP', false);
  INSERT INTO challenges (name, category, description, difficulty, authors, tags, type, max_points, score_type) VALUES ('chall-5', 'cat-2', 'TEST chall-5 DESC', 'Easy', ARRAY['author3'], ARRAY['tag-5'], 'Normal', 500, 'Static');
  UPDATE docker_configs SET image='echo-server:latest', hash_domain=TRUE WHERE chall_id=(SELECT id FROM challenges WHERE name='chall-3');
  UPDATE docker_configs SET compose='
services:
  chall:
    image: echo-server:latest
    container_name: ${CONTAINER_NAME}
    ports:
      - "${INSTANCE_PORT}:1337"
    environment:
      - ECHO_MESSAGE=Hello from app
      - INSTANCE_PORT=${INSTANCE_PORT}
      - INSTANCE_HOST=${INSTANCE_HOST}
    ', hash_domain=TRUE WHERE chall_id=(SELECT id FROM challenges WHERE name='chall-4');
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
  PERFORM assert(COUNT(b)=0, 'wrong_submission_no_badges') FROM badges b;
  PERFORM assert(COUNT(s)=1, 'wrong_submission_recorded') FROM submissions s WHERE s.status='Wrong';
  PERFORM assert(c.solves=0, 'wrong_submission_no_solves') FROM challenges c WHERE c.name='chall-1';
  PERFORM assert(c.points=c.max_points, 'wrong_submission_points_unchanged') FROM challenges c WHERE c.name='chall-1';

  -- insert a repeated submission from 'a' to 'chal-1'
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='a'),
    (SELECT id FROM challenges WHERE name='chall-1'),
    'Repeated', 'flag');
  PERFORM assert(COUNT(b)=0, 'repeated_submission_no_badges') FROM badges b;
  PERFORM assert(COUNT(s)=1, 'repeated_submission_recorded') FROM submissions s WHERE s.status='Repeated';
  PERFORM assert(c.solves=0, 'repeated_submission_no_solves') FROM challenges c WHERE c.name='chall-1';
  PERFORM assert(c.points=c.max_points, 'repeated_submission_points_unchanged') FROM challenges c WHERE c.name='chall-1';

  -- insert a valid submission from 'a' to 'chal-1'
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='a'),
    (SELECT id FROM challenges WHERE name='chall-1'),
    'Correct', 'flag');
  PERFORM assert(COUNT(b)=0, 'correct_submission_no_badges_yet') FROM badges b;
  PERFORM assert(COUNT(s)=1, 'correct_submission_recorded') FROM submissions s WHERE s.status='Correct';
  PERFORM assert(c.solves=1, 'correct_submission_solves_increment') FROM challenges c WHERE c.name='chall-1';
  PERFORM assert(c.points=c.max_points, 'correct_submission_points_full') FROM challenges c WHERE c.name='chall-1';

  -- insert a valid submission from 'a' to 'chal-3'
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='a'),
    (SELECT id FROM challenges WHERE name='chall-3'),
    'Correct', 'flag');
  PERFORM assert(COUNT(b)=0, 'second_correct_no_badges') FROM badges b;
  PERFORM assert(COUNT(s)=2, 'total_correct_submissions_2') FROM submissions s WHERE s.status='Correct';
  PERFORM assert(c.solves=1, 'chall3_first_solve') FROM challenges c WHERE c.name='chall-3';
  PERFORM assert(c.points=c.max_points, 'chall3_points_full') FROM challenges c WHERE c.name='chall-3';

  -- insert a valid submission from 'a' to 'chal-4' (should give also a badge)
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='a'),
    (SELECT id FROM challenges WHERE name='chall-4'),
    'Correct', 'flag');
  PERFORM assert(COUNT(b)=1, 'third_correct_awards_badge') FROM badges b;
  PERFORM assert(COUNT(s)=3, 'total_correct_submissions_3') FROM submissions s WHERE s.status='Correct';
  PERFORM assert(c.solves=1, 'chall4_first_solve') FROM challenges c WHERE c.name='chall-4';
  PERFORM assert(c.points=c.max_points, 'chall4_points_full') FROM challenges c WHERE c.name='chall-4';

  -- insert a valid submission from 'a' to 'chal-1' but it's repeated
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='a'),
    (SELECT id FROM challenges WHERE name='chall-1'),
    'Correct', 'flag');
  PERFORM assert(COUNT(b)=1, 'repeated_correct_no_new_badges') FROM badges b;
  PERFORM assert(COUNT(s)=2, 'repeated_correct_count') FROM submissions s WHERE s.status='Repeated';
  PERFORM assert(c.solves=1, 'chall1_solves_remain') FROM challenges c WHERE c.name='chall-1';
  PERFORM assert(c.points=c.max_points, 'chall1_points_unchanged') FROM challenges c WHERE c.name='chall-1';

  -- insert a valid submission from 'b' to 'chal-3' but it's already solved by 'a' (its teammate)
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='b'),
    (SELECT id FROM challenges WHERE name='chall-3'),
    'Correct', 'flag');
  PERFORM assert(COUNT(b)=1, 'teammate_repeated_no_badges') FROM badges b;
  PERFORM assert(COUNT(s)=3, 'repeated_total_3') FROM submissions s WHERE s.status='Repeated';
  PERFORM assert(c.solves=1, 'chall3_solves_still_1') FROM challenges c WHERE c.name='chall-3';
  PERFORM assert(c.points=c.max_points, 'chall3_points_unchanged') FROM challenges c WHERE c.name='chall-3';

  -- insert a valid submission from 'c' to 'chal-4'
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='c'),
    (SELECT id FROM challenges WHERE name='chall-4'),
    'Correct', 'flag');
  PERFORM assert(COUNT(b)=1, 'chall4_second_solve_no_badge') FROM badges b;
  PERFORM assert(COUNT(s)=4, 'total_correct_submissions_4') FROM submissions s WHERE s.status='Correct';
  PERFORM assert(c.solves=2, 'chall4_solves_increment') FROM challenges c WHERE c.name='chall-4';
  PERFORM assert(c.points<c.max_points, 'chall4_points_decrease') FROM challenges c WHERE c.name='chall-4';

  -- insert a valid submission from 'c' to 'chal-2' (should give another badge)
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='c'),
    (SELECT id FROM challenges WHERE name='chall-2'),
    'Correct', 'flag');
  PERFORM assert(COUNT(b)=2, 'chall2_correct_badge_awarded') FROM badges b;
  PERFORM assert(COUNT(s)=5, 'total_correct_submissions_5') FROM submissions s WHERE s.status='Correct';
  PERFORM assert(c.solves=1, 'chall2_first_solve') FROM challenges c WHERE c.name='chall-2';
  PERFORM assert(c.points=c.max_points, 'chall2_points_full') FROM challenges c WHERE c.name='chall-2';

  -- insert a valid submission from 'f' to 'chal-3' but it's not a player so it doesn't add up
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='f'),
    (SELECT id FROM challenges WHERE name='chall-3'),
    'Correct', 'flag');
  PERFORM assert(COUNT(b)=2, 'non_player_submission_no_badges') FROM badges b;
  PERFORM assert(COUNT(s)=6, 'total_correct_submissions_6') FROM submissions s WHERE s.status='Correct';
  PERFORM assert(c.solves=1, 'chall3_solves_still_1_nonplayer') FROM challenges c WHERE c.name='chall-3';

  -- deletes all correct submissions from 'a' to 'chall-1' (should also remove the badge)
  DELETE FROM submissions WHERE user_id=(SELECT id FROM users WHERE name='a')
    AND chall_id=(SELECT id FROM challenges WHERE name='chall-1')
    AND status='Correct';
  PERFORM assert(COUNT(b)=1, 'delete_correct_removes_badge') FROM badges b;
  PERFORM assert(COUNT(s)=5, 'total_correct_submissions_5_after_delete') FROM submissions s WHERE s.status='Correct';

  -- insert a valid submission from 'a' to 'chal-1' (should give back the badge)
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='a'),
    (SELECT id FROM challenges WHERE name='chall-1'),
    'Correct', 'flag');
  PERFORM assert(COUNT(b)=2, 'correct_submission_restores_badge') FROM badges b;
  PERFORM assert(COUNT(s)=6, 'total_correct_submissions_6') FROM submissions s WHERE s.status='Correct';
  PERFORM assert(c.solves=1, 'chall1_solves_after_restore') FROM challenges c WHERE c.name='chall-1';

  -- insert a valid submission from 'c' to 'chal-1' (should not give a badge with 2 solves of 3)
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='c'),
    (SELECT id FROM challenges WHERE name='chall-1'),
    'Correct', 'flag');
  PERFORM assert(COUNT(b)=2, 'chall1_second_solve_no_badge') FROM badges b;
  PERFORM assert(COUNT(s)=7, 'total_correct_submissions_7') FROM submissions s WHERE s.status='Correct';
  PERFORM assert(c.solves=2, 'chall1_solves_increment_to_2') FROM challenges c WHERE c.name='chall-1';

  -- check that the scores are the same
  SELECT score INTO tmp FROM teams WHERE name='A';
  PERFORM assert(score=tmp, 'teams_scores_equal') FROM teams WHERE name='B';

  -- after recomputing the scores, they should be lower
  UPDATE configs SET value='3' WHERE key='chall-points-decay';
  PERFORM assert(score<tmp, 'decay_lower_scores') FROM teams WHERE name='B';

  -- after recomputing the scores, they should be greater
  UPDATE configs SET value='7' WHERE key='chall-points-decay';
  PERFORM assert(score>tmp, 'decay_higher_scores') FROM teams WHERE name='B';

  -- after recomputing the scores, they should be all the same as 500
  UPDATE configs SET value='500' WHERE key='chall-min-points';
  PERFORM assert(COUNT(c)=5, 'all_challs_min_points_500') FROM challenges c WHERE points=500;

  -- after changing the score to 1000, there should be no challenge left at 500 after recomputing the scores
  UPDATE challenges SET max_points=1000;
  PERFORM assert(COUNT(c)=0, 'no_challs_left_at_500') FROM challenges c WHERE points=500;

  -- deletes 'chall-3' (should now give the badge to team 'B' that hadn't solved only that)
  DELETE FROM challenges c WHERE c.name='chall-3';
  PERFORM assert(COUNT(b)=3, 'delete_chall3_awards_badge_B') FROM badges b;
  PERFORM assert(COUNT(s)=5, 'correct_submissions_drop_after_delete') FROM submissions s WHERE s.status='Correct';

  -- inserts back 'chall-3' as hidden (nothing should change)
  INSERT INTO challenges (name, category, description, type, max_points, score_type)
    VALUES ('chall-3', 'cat-1', 'TEST', 'Container', 1000, 'Dynamic');
  PERFORM assert(COUNT(b)=3, 'reinstated_hidden_chall_no_changes') FROM badges b;

  -- makes 'chall-3' visible again, now should remove the badges (previous solves were deleted on the DELETE as cascade)
  UPDATE challenges SET hidden=false WHERE name='chall-3';
  PERFORM assert(COUNT(b)=1, 'chall3_visible_badges_removed') FROM badges b;

  -- inserts again new valid submissions, should give badges and lower the challenge points
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='a'),
    (SELECT id FROM challenges WHERE name='chall-3'),
    'Correct', 'flag');
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='c'),
    (SELECT id FROM challenges WHERE name='chall-3'),
    'Correct', 'flag');
  PERFORM assert(COUNT(b)=3, 'new_solves_restore_badges') FROM badges b;
  PERFORM assert(COUNT(s)=7, 'total_correct_submissions_7_again') FROM submissions s WHERE s.status='Correct';
  PERFORM assert(points!=max_points, 'chall3_dynamic_points_decrease') FROM challenges WHERE name='chall-3';

  -- changes 'chall-3' scoring type as static, so the points should reset
  UPDATE challenges SET score_type='Static' WHERE name='chall-3';
  PERFORM assert(points=max_points, 'chall3_static_points_full_reset') FROM challenges WHERE name='chall-3';

  -- changes 'chall-3' max_points so the points should change as them
  UPDATE challenges SET max_points=500 WHERE name='chall-3';
  PERFORM assert(points=max_points, 'chall3_maxpoints_update_reflects') FROM challenges WHERE name='chall-3';

  -- changes 'chall-3' score type back to dynamic and recomputes the score
  UPDATE configs SET value='50' WHERE key='chall-min-points';
  UPDATE challenges SET score_type='Dynamic' WHERE name='chall-3';
  PERFORM assert(points!=max_points, 'chall3_dynamic_recompute_changes_points') FROM challenges WHERE name='chall-3';

  -- inserts a correct submission from 'd' on 'chall-1', but it doesn't have a team so it's an invalid submission
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='d'),
    (SELECT id FROM challenges WHERE name='chall-1'),
    'Correct', 'flag');
  PERFORM assert(COUNT(s)=7, 'invalid_submission_does_not_increase_correct_count') FROM submissions s WHERE s.status='Correct';
  PERFORM assert(COUNT(s)=1, 'invalid_submission_recorded') FROM submissions s WHERE s.status='Invalid';

  -- removes the user 'Correct', so all his submissions should be deleted and also badges and team points removed
  DELETE FROM users WHERE name='c';
  PERFORM assert(COUNT(b)=1, 'delete_user_removes_badges') FROM badges b;
  PERFORM assert(COUNT(s)=3, 'delete_user_removes_submissions') FROM submissions s WHERE s.status='Correct';
  PERFORM assert(score=0, 'teamB_score_zero_after_user_delete') FROM teams WHERE name='B';

  -- checks that 'chall-3' and 'chall-4' have their configs created
  PERFORM assert(count(d)=1, 'chall3_config_exists') FROM docker_configs d WHERE d.chall_id=(SELECT id FROM challenges WHERE name='chall-3');
  PERFORM assert(count(d)=1, 'chall4_config_exists') FROM docker_configs d WHERE d.chall_id=(SELECT id FROM challenges WHERE name='chall-4');

  -- checks that docker configs are created on type update and not duplicated
  INSERT INTO challenges (name, category, description, type, max_points, score_type) VALUES ('chall-test', 'cat-1', 'TEST', 'Normal', 500, 'Dynamic');
  PERFORM assert(count(d)=0, 'chall_test_no_config_initial') FROM docker_configs d WHERE d.chall_id=(SELECT id FROM challenges WHERE name='chall-test');
  UPDATE challenges SET type='Container' WHERE id=(SELECT id FROM challenges WHERE name='chall-test');
  PERFORM assert(count(d)=1, 'chall_test_config_created_on_type_container') FROM docker_configs d WHERE d.chall_id=(SELECT id FROM challenges WHERE name='chall-test');
  UPDATE challenges SET type='Normal' WHERE id=(SELECT id FROM challenges WHERE name='chall-test');
  PERFORM assert(count(d)=1, 'chall_test_config_not_duplicated_normal') FROM docker_configs d WHERE d.chall_id=(SELECT id FROM challenges WHERE name='chall-test');
  UPDATE challenges SET type='Compose' WHERE id=(SELECT id FROM challenges WHERE name='chall-test');
  PERFORM assert(count(d)=1, 'chall_test_config_not_duplicated_compose') FROM docker_configs d WHERE d.chall_id=(SELECT id FROM challenges WHERE name='chall-test');

  -- checks for the first bloods does not change on already blooded challs
  PERFORM assert(count(s)=3, 'first_blood_count_initial') FROM submissions s WHERE s.first_blood=TRUE;
  INSERT INTO users (name, email, password_hash, password_salt, role, team_id) VALUES ('c', 'c@c.c',
    '41d65efe433e60755bef957e56ed6466b24c44a86b8ec595df4c9cdfa9c3aca9', '1a5e93869fa3c2ee04139db8834f8808',
    'Player', (SELECT id FROM teams WHERE name='B'));
  INSERT INTO submissions (user_id, chall_id, status, flag) VALUES (
    (SELECT id FROM users WHERE name='c'),
    (SELECT id FROM challenges WHERE name='chall-1'),
    'Correct', 'flag');
  PERFORM assert(count(s)=3, 'first_blood_not_changed_by_new_user') FROM submissions s WHERE s.first_blood=TRUE;
  
  -- checks that the first blood is transferred on delete of the first blood submission
  PERFORM assert(count(s)=0, 'new_user_has_no_firstblood') FROM submissions s WHERE s.first_blood=TRUE AND s.user_id=(SELECT id FROM users WHERE name='c');
  DELETE FROM submissions s WHERE s.first_blood=TRUE AND s.chall_id=(SELECT id FROM challenges WHERE name='chall-1');
  PERFORM assert(count(s)=3, 'firstblood_transferred_after_delete') FROM submissions s WHERE s.first_blood=TRUE;
  PERFORM assert(count(s)=1, 'new_user_firstblood_count') FROM submissions s WHERE s.first_blood=TRUE AND s.user_id=(SELECT id FROM users WHERE name='c');

  -- changes user 'a' role to non-player so the solves should be subtracted and points removed
  PERFORM assert(score>0) FROM teams WHERE name='A';
  PERFORM assert(solves=1) FROM challenges WHERE name='chall-3';
  PERFORM assert(solves=1) FROM challenges WHERE name='chall-4';
  PERFORM assert(solves=1) FROM challenges WHERE name='chall-1';
  UPDATE users SET role='Author' WHERE id=(SELECT id FROM users WHERE name='a');
  PERFORM assert(score=0) FROM teams WHERE name='A';
  PERFORM assert(solves=0) FROM challenges WHERE name='chall-3';
  PERFORM assert(solves=0) FROM challenges WHERE name='chall-4';
  PERFORM assert(solves=1) FROM challenges WHERE name='chall-1';
  
  -- if user 'a' goes from a non-player role to another non-player role, nothing should change
  UPDATE users SET role='Spectator' WHERE id=(SELECT id FROM users WHERE name='a');
  PERFORM assert(score=0) FROM teams WHERE name='A';
  PERFORM assert(solves=0) FROM challenges WHERE name='chall-3';
  PERFORM assert(solves=0) FROM challenges WHERE name='chall-4';
  PERFORM assert(solves=1) FROM challenges WHERE name='chall-1';

  -- if user 'a' goes from a non-player role to another non-player role, nothing should change
  UPDATE users SET role='Admin' WHERE id=(SELECT id FROM users WHERE name='a');
  PERFORM assert(score=0) FROM teams WHERE name='A';
  PERFORM assert(solves=0) FROM challenges WHERE name='chall-3';
  PERFORM assert(solves=0) FROM challenges WHERE name='chall-4';
  PERFORM assert(solves=1) FROM challenges WHERE name='chall-1';

  -- changes user 'a' role back to player, so the solves should be added back and points restored
  UPDATE users SET role='Player' WHERE id=(SELECT id FROM users WHERE name='a');
  PERFORM assert(score>0) FROM teams WHERE name='A';
  PERFORM assert(solves=1) FROM challenges WHERE name='chall-3';
  PERFORM assert(solves=1) FROM challenges WHERE name='chall-4';
  PERFORM assert(solves=1) FROM challenges WHERE name='chall-1';
END;
$$ LANGUAGE plpgsql;

/*
SELECT delete_all();
SELECT insert_mock_configs();
SELECT insert_mock_data();
SELECT insert_mock_submissions();
SELECT tests();

SELECT * FROM submissions;
SELECT * FROM instances;
SELECT * FROM flags;
SELECT * FROM attachments;
SELECT * FROM docker_configs;
SELECT * FROM challenges;
SELECT * FROM team_category_solves;
SELECT * FROM categories;
SELECT * FROM badges;
SELECT * FROM users;
SELECT * FROM teams;
SELECT * FROM configs;
--*/
