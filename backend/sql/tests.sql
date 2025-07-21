CREATE OR REPLACE FUNCTION assert(result BOOLEAN)
RETURNS VOID AS $$
BEGIN
  IF result != TRUE THEN
    RAISE EXCEPTION 'Assertion failed: Got: %', result;
  ELSE
    RAISE NOTICE 'PASS';
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

  insert into submissions (user_id, chall_id, status, flag) values ((select id from users where name='a'), (select id from challenges where name='chall-1'), 'C', 'flag');
  perform assert(count(b)=0) from badges b;
  perform assert(count(s)=1) from submissions s where s.status='C';
  perform assert(c.solves=1) from challenges c where c.name='chall-1';
  perform assert(c.points=c.max_points) from challenges c where c.name='chall-1';
  
  insert into submissions (user_id, chall_id, status, flag) values ((select id from users where name='a'), (select id from challenges where name='chall-3'), 'C', 'flag');
  perform assert(count(b)=0) from badges b;
  perform assert(count(s)=2) from submissions s where s.status='C';
  perform assert(c.solves=1) from challenges c where c.name='chall-3';
  perform assert(c.points=c.max_points) from challenges c where c.name='chall-3';

  insert into submissions (user_id, chall_id, status, flag) values ((select id from users where name='a'), (select id from challenges where name='chall-4'), 'C', 'flag');
  perform assert(count(b)=1) from badges b;
  perform assert(count(s)=3) from submissions s where s.status='C';
  perform assert(c.solves=1) from challenges c where c.name='chall-4';
  perform assert(c.points=c.max_points) from challenges c where c.name='chall-4';
  
  insert into submissions (user_id, chall_id, status, flag) values ((select id from users where name='a'), (select id from challenges where name='chall-1'), 'C', 'flag');
  perform assert(count(b)=1) from badges b;
  perform assert(count(s)=1) from submissions s where s.status='R';
  perform assert(c.solves=1) from challenges c where c.name='chall-1';
  perform assert(c.points=c.max_points) from challenges c where c.name='chall-1';

  insert into submissions (user_id, chall_id, status, flag) values ((select id from users where name='b'), (select id from challenges where name='chall-3'), 'C', 'flag');
  perform assert(count(b)=1) from badges b;
  perform assert(count(s)=2) from submissions s where s.status='R';
  perform assert(c.solves=1) from challenges c where c.name='chall-3';
  perform assert(c.points=c.max_points) from challenges c where c.name='chall-3';

  insert into submissions (user_id, chall_id, status, flag) values ((select id from users where name='c'), (select id from challenges where name='chall-4'), 'C', 'flag');
  perform assert(count(b)=1) from badges b;
  perform assert(count(s)=4) from submissions s where s.status='C';
  perform assert(c.solves=2) from challenges c where c.name='chall-4';
  perform assert(c.points<c.max_points) from challenges c where c.name='chall-4';

  insert into submissions (user_id, chall_id, status, flag) values ((select id from users where name='c'), (select id from challenges where name='chall-2'), 'C', 'flag');
  perform assert(count(b)=2) from badges b;
  perform assert(count(s)=5) from submissions s where s.status='C';
  perform assert(c.solves=1) from challenges c where c.name='chall-2';
  perform assert(c.points=c.max_points) from challenges c where c.name='chall-2';
  
  insert into submissions (user_id, chall_id, status, flag) values ((select id from users where name='f'), (select id from challenges where name='chall-3'), 'C', 'flag');
  perform assert(count(b)=2) from badges b;
  perform assert(count(s)=6) from submissions s where s.status='C';
  perform assert(c.solves=1) from challenges c where c.name='chall-3';

  delete from submissions where user_id=(select id from users where name='a') and chall_id=(select id from challenges where name='chall-1') and status='C';
  perform assert(count(b)=1) from badges b;
  perform assert(count(s)=5) from submissions s where s.status='C';

  insert into submissions (user_id, chall_id, status, flag) values ((select id from users where name='a'), (select id from challenges where name='chall-1'), 'C', 'flag');
  perform assert(count(b)=2) from badges b;
  perform assert(count(s)=6) from submissions s where s.status='C';
  perform assert(c.solves=1) from challenges c where c.name='chall-1';

  insert into submissions (user_id, chall_id, status, flag) values ((select id from users where name='c'), (select id from challenges where name='chall-1'), 'C', 'flag');
  perform assert(count(b)=2) from badges b;
  perform assert(count(s)=7) from submissions s where s.status='C';
  perform assert(c.solves=2) from challenges c where c.name='chall-1';

  select score into tmp from teams where name='A';
  perform assert(score=tmp) from teams where name='B';

  update configs set value='3' where key='chall-points-decay';
  perform assert(score<tmp) from teams where name='B';

  delete from challenges c where c.name='chall-3';
  perform assert(count(b)=3) from badges b;
  perform assert(count(s)=5) from submissions s where s.status='C';

  insert into challenges (name, category, description, type, max_points, score_type) values ('chall-3', 'cat-1', 'TEST', 'N', 500, 'D');
  perform assert(count(b)=3) from badges b;

  update challenges set hidden=false where name='chall-3';
  perform assert(count(b)=1) from badges b;

  insert into submissions (user_id, chall_id, status, flag) values ((select id from users where name='d'), (select id from challenges where name='chall-1'), 'C', 'flag');
  perform assert(count(s)=5) from submissions s where s.status='C';
  perform assert(count(s)=1) from submissions s where s.status='I';

END;
$$ LANGUAGE plpgsql;

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
