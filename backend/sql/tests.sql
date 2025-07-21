CREATE OR REPLACE FUNCTION assert(result BOOLEAN)
RETURNS VOID AS $$
BEGIN
  IF result != TRUE THEN
    RAISE EXCEPTION 'Assertion failed: Got: %', result;
  END IF;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION delete_all()
RETURNS VOID AS $$
BEGIN
  delete from submissions;
  delete from tags;
  delete from instances;
  delete from flags;
  delete from challenges;
  delete from team_category_solves;
  delete from categories;
  delete from badges;
  delete from users;
  delete from teams;
  delete from configs;
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
  insert into configs (key, type, value) values ('chall-min-points', 'int', '100');
  insert into configs (key, type, value) values ('chall-points-decay', 'int', '5');
  insert into categories (name, icon) values ('cat-1', 'cat-1');
  insert into categories (name, icon) values ('cat-2', 'cat-2');
  insert into challenges (name, category, description, type, max_points, score_type, hidden) values ('chall-1', 'cat-1', 'TEST', 'N', 500, 'D', false);
  insert into challenges (name, category, description, type, max_points, score_type, hidden) values ('chall-2', 'cat-2', 'TEST', 'N', 500, 'D', false);
  insert into challenges (name, category, description, type, max_points, score_type, hidden) values ('chall-3', 'cat-1', 'TEST', 'N', 500, 'D', false);
  insert into challenges (name, category, description, type, max_points, score_type, hidden) values ('chall-4', 'cat-1', 'TEST', 'N', 500, 'D', false);
  insert into teams (name, password_hash) values ('A', 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa');
  insert into teams (name, password_hash) values ('B', 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa');
  insert into teams (name, password_hash) values ('C', 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa');
  insert into users (name, email, password_hash, role, team_id) values ('a', 'a@a', 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa', 'P', (select id from teams where name='A'));
  insert into users (name, email, password_hash, role, team_id) values ('b', 'b@b', 'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb', 'P', (select id from teams where name='A'));
  insert into users (name, email, password_hash, role, team_id) values ('c', 'c@c', 'cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc', 'P', (select id from teams where name='B'));
  insert into users (name, email, password_hash, role) values ('d', 'd@d', 'dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd', 'P');
  insert into users (name, email, password_hash, role, team_id) values ('e', 'e@e', 'eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee', 'M', (select id from teams where name='A'));
  insert into users (name, email, password_hash, role, team_id) values ('f', 'f@f', 'ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff', 'A', (select id from teams where name='C'));
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION test()
RETURNS VOID AS $$
DECLARE
  tmp INTEGER;
BEGIN
  perform delete_all();
  perform insert_base_data();

  -- insert a wrong submission from 'a' to 'chal-1'
  insert into submissions (user_id, chall_id, status, flag) values ((select id from users where name='a'), (select id from challenges where name='chall-1'), 'W', 'flag');
  perform assert(count(b)=0) from badges b;
  perform assert(count(s)=1) from submissions s where s.status='W';
  perform assert(c.solves=0) from challenges c where c.name='chall-1';
  perform assert(c.points=c.max_points) from challenges c where c.name='chall-1';

  -- insert a repeated submission from 'a' to 'chal-1'
  insert into submissions (user_id, chall_id, status, flag) values ((select id from users where name='a'), (select id from challenges where name='chall-1'), 'R', 'flag');
  perform assert(count(b)=0) from badges b;
  perform assert(count(s)=1) from submissions s where s.status='R';
  perform assert(c.solves=0) from challenges c where c.name='chall-1';
  perform assert(c.points=c.max_points) from challenges c where c.name='chall-1';

  -- insert a valid submission from 'a' to 'chal-1'
  insert into submissions (user_id, chall_id, status, flag) values ((select id from users where name='a'), (select id from challenges where name='chall-1'), 'C', 'flag');
  perform assert(count(b)=0) from badges b;
  perform assert(count(s)=1) from submissions s where s.status='C';
  perform assert(c.solves=1) from challenges c where c.name='chall-1';
  perform assert(c.points=c.max_points) from challenges c where c.name='chall-1';
  
  -- insert a valid submission from 'a' to 'chal-3'
  insert into submissions (user_id, chall_id, status, flag) values ((select id from users where name='a'), (select id from challenges where name='chall-3'), 'C', 'flag');
  perform assert(count(b)=0) from badges b;
  perform assert(count(s)=2) from submissions s where s.status='C';
  perform assert(c.solves=1) from challenges c where c.name='chall-3';
  perform assert(c.points=c.max_points) from challenges c where c.name='chall-3';

  -- insert a valid submission from 'a' to 'chal-4' (should give also a badge)
  insert into submissions (user_id, chall_id, status, flag) values ((select id from users where name='a'), (select id from challenges where name='chall-4'), 'C', 'flag');
  perform assert(count(b)=1) from badges b;
  perform assert(count(s)=3) from submissions s where s.status='C';
  perform assert(c.solves=1) from challenges c where c.name='chall-4';
  perform assert(c.points=c.max_points) from challenges c where c.name='chall-4';
  
  -- insert a valid submission from 'a' to 'chal-1' but it's repeated
  insert into submissions (user_id, chall_id, status, flag) values ((select id from users where name='a'), (select id from challenges where name='chall-1'), 'C', 'flag');
  perform assert(count(b)=1) from badges b;
  perform assert(count(s)=2) from submissions s where s.status='R';
  perform assert(c.solves=1) from challenges c where c.name='chall-1';
  perform assert(c.points=c.max_points) from challenges c where c.name='chall-1';

  -- insert a valid submission from 'b' to 'chal-3' but it's already solved by 'a' (its teammate)
  insert into submissions (user_id, chall_id, status, flag) values ((select id from users where name='b'), (select id from challenges where name='chall-3'), 'C', 'flag');
  perform assert(count(b)=1) from badges b;
  perform assert(count(s)=3) from submissions s where s.status='R';
  perform assert(c.solves=1) from challenges c where c.name='chall-3';
  perform assert(c.points=c.max_points) from challenges c where c.name='chall-3';

  -- insert a valid submission from 'c' to 'chal-4'
  insert into submissions (user_id, chall_id, status, flag) values ((select id from users where name='c'), (select id from challenges where name='chall-4'), 'C', 'flag');
  perform assert(count(b)=1) from badges b;
  perform assert(count(s)=4) from submissions s where s.status='C';
  perform assert(c.solves=2) from challenges c where c.name='chall-4';
  perform assert(c.points<c.max_points) from challenges c where c.name='chall-4';

  -- insert a valid submission from 'c' to 'chal-2' (should give another badge)
  insert into submissions (user_id, chall_id, status, flag) values ((select id from users where name='c'), (select id from challenges where name='chall-2'), 'C', 'flag');
  perform assert(count(b)=2) from badges b;
  perform assert(count(s)=5) from submissions s where s.status='C';
  perform assert(c.solves=1) from challenges c where c.name='chall-2';
  perform assert(c.points=c.max_points) from challenges c where c.name='chall-2';
  
  -- insert a valid submission from 'f' to 'chal-3' but it's not a player so it doesn't add up
  insert into submissions (user_id, chall_id, status, flag) values ((select id from users where name='f'), (select id from challenges where name='chall-3'), 'C', 'flag');
  perform assert(count(b)=2) from badges b;
  perform assert(count(s)=6) from submissions s where s.status='C';
  perform assert(c.solves=1) from challenges c where c.name='chall-3';

  -- deletes all correct submissions from 'a' to 'chall-1' (should also remove the badge)
  delete from submissions where user_id=(select id from users where name='a') and chall_id=(select id from challenges where name='chall-1') and status='C';
  perform assert(count(b)=1) from badges b;
  perform assert(count(s)=5) from submissions s where s.status='C';

  -- insert a valid submission from 'a' to 'chal-1' (should give back the badge)
  insert into submissions (user_id, chall_id, status, flag) values ((select id from users where name='a'), (select id from challenges where name='chall-1'), 'C', 'flag');
  perform assert(count(b)=2) from badges b;
  perform assert(count(s)=6) from submissions s where s.status='C';
  perform assert(c.solves=1) from challenges c where c.name='chall-1';

  -- insert a valid submission from 'c' to 'chal-1' (should not give a badge with 2 solves of 3)
  insert into submissions (user_id, chall_id, status, flag) values ((select id from users where name='c'), (select id from challenges where name='chall-1'), 'C', 'flag');
  perform assert(count(b)=2) from badges b;
  perform assert(count(s)=7) from submissions s where s.status='C';
  perform assert(c.solves=2) from challenges c where c.name='chall-1';

  -- check that the scores are the same
  select score into tmp from teams where name='A';
  perform assert(score=tmp) from teams where name='B';

  -- after recomputing the scores, they should be lower
  update configs set value='3' where key='chall-points-decay';
  perform assert(score<tmp) from teams where name='B';

  -- after recomputing the scores, they should be greater
  update configs set value='7' where key='chall-points-decay';
  perform assert(score>tmp) from teams where name='B';

  -- after recomputing the scores, they should be all the same as 500
  update configs set value='500' where key='chall-min-points';
  perform assert(count(c)=4) from challenges c where points=500;

  update challenges set max_points=1000;
  perform assert(count(c)=0) from challenges c where points=500;

  -- deletes 'chall-3' (should now give the badge to team 'B' that hadn't solved only that)
  delete from challenges c where c.name='chall-3';
  perform assert(count(b)=3) from badges b;
  perform assert(count(s)=5) from submissions s where s.status='C';

  -- inserts back 'chall-3' as hidden (nothing should change)
  insert into challenges (name, category, description, type, max_points, score_type) values ('chall-3', 'cat-1', 'TEST', 'N', 1000, 'D');
  perform assert(count(b)=3) from badges b;

  -- makes 'chall-3' visible again, now shoud remove the badges (previous solves were deleted on the delete as cascade)
  update challenges set hidden=false where name='chall-3';
  perform assert(count(b)=1) from badges b;

  -- inserts again new valid submissions, should give aìbadges and lower the challenge points
  insert into submissions (user_id, chall_id, status, flag) values ((select id from users where name='a'), (select id from challenges where name='chall-3'), 'C', 'flag');
  insert into submissions (user_id, chall_id, status, flag) values ((select id from users where name='c'), (select id from challenges where name='chall-3'), 'C', 'flag');
  perform assert(count(b)=3) from badges b;
  perform assert(count(s)=7) from submissions s where s.status='C';
  perform assert(points!=max_points) from challenges where name='chall-3';

  -- changes 'chall-3' scoring type as static, so the points should reset
  update challenges set score_type='S' where name='chall-3';
  perform assert(points=max_points) from challenges where name='chall-3';

  -- inserts a correct submission from 'd' on 'chall-1', but it doesn't have a team so it's an invalid submission
  insert into submissions (user_id, chall_id, status, flag) values ((select id from users where name='d'), (select id from challenges where name='chall-1'), 'C', 'flag');
  perform assert(count(s)=7) from submissions s where s.status='C';
  perform assert(count(s)=1) from submissions s where s.status='I';

END;
$$ LANGUAGE plpgsql;

/*
select test();

select * from submissions;
select * from tags;
select * from instances;
select * from flags;
select * from challenges;
select * from team_category_solves;
select * from categories;
select * from badges;
select * from users;
select * from teams;
select * from configs;
*/
