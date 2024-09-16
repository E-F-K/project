------------------------------------------------------------------------------------------
-- create user
insert into users
    (id, name, email, token)
values
    ('00000000-0000-0000-0000-000000000001', 'vera', 'vera@localhost', 'secret-token');

-- create list
insert into lists
    (id, user_id, name, email)
values
    ('10000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'homework', 'vera@localhost');

--create task
insert into tasks
    (id, list_id, priority, deadline, done, name)
values
    ('20000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000001', 'high', NUll, false, 'probability theory');

insert into tasks
    (id, list_id, priority, deadline, done, name)
values
    ('20000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000001', 'high', NUll, false, 'matan');

insert into tasks
    (id, list_id, priority, deadline, done, name)
values
    ('20000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000001', 'low', NUll, false, 'graphs curs');


------------------------------------------------------------------------------------------
-- read user
SELECT * FROM users;

-- read list
SELECT * FROM lists;

-- read task
SELECT * FROM tasks;

-- update user
UPDATE users
SET name = 'new'
WHERE id = '00000000-0000-0000-0000-000000000001';


-- update task
UPDATE users
SET done = true
WHERE id = '20000000-0000-0000-0000-000000000001';


------------------------------------------------------------------------------------------
-- delete user
DELETE FROM users

DELETE FROM users WHERE id = '00000000-0000-0000-0000-000000000001';

-- delete list
DELETE FROM lists

DELETE FROM lists WHERE id = '10000000-0000-0000-0000-000000000001';

-- delete task
DELETE FROM tasks

DELETE FROM tasks WHERE id = '20000000-0000-0000-0000-000000000001';



