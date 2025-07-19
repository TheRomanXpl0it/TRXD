
insert into configs (key, type, value) values ('chall-min-points', 'int', '100');
insert into configs (key, type, value) values ('chall-points-decay', 'int', '10');
insert into categories (name, icon) values ('test', 'test');
insert into challenges (name, category, description, type, max_points, score_type, points) values ('chall', 'test', 'TEST', 'N', 500, 'D', 500);
insert into challenges (name, category, description, type, max_points, score_type, points) values ('chall2', 'test', 'TEST', 'N', 500, 'D', 500);
insert into teams (name, password_hash) values ('A', 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa');
insert into teams (name, password_hash) values ('B', 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa');
insert into users (name, email, password_hash, role, team_id) values ('a', 'a@a', 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa', 'P', 1);
insert into users (name, email, password_hash, role, team_id) values ('b', 'b@b', 'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb', 'P', 1);
insert into users (name, email, password_hash, role, team_id) values ('c', 'c@c', 'cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc', 'P', 2);
insert into users (name, email, password_hash, role) values ('d', 'd@d', 'dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd', 'M');
insert into users (name, email, password_hash, role) values ('e', 'e@e', 'cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc', 'P');

insert into submissions (user_id, chall_id, status, flag) values (1, (select id from challenges where name='chall'), 'C', 'flag');
insert into submissions (user_id, chall_id, status, flag) values (3, (select id from challenges where name='chall2'), 'C', 'flag');

delete from submissions;

delete from challenges where name = 'chall';

SELECT name, score, team_id FROM users;
SELECT name, score FROM teams;
SELECT id, name, max_points, points, solves FROM challenges;
SELECT * FROM submissions;
select * from categories;
SELECT * FROM configs;

update challenges set max_points=500;
update challenges set solves=0;
update users set score=0;
update configs set value='100' where key='chall-min-points';
update configs set value='5' where key='chall-points-decay';

drop trigger tr_recompute_chall_points on challenges;
