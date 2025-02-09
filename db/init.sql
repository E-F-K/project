create table if not exists users (
    id UUID PRIMARY KEY,
    name TEXT not null,
    email TEXT not null,
    password_hash TEXT not null,
    token TEXT not null,
    updated_at timestamp with time zone not null,
    unique(email),
    unique(token)
);

create table if not exists lists (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) not null,
    name TEXT not null,
    updated_at timestamp with time zone not null
);

create type priority AS ENUM ('low', 'normal', 'high');

create table if not exists tasks (
    id UUID PRIMARY KEY,
    list_id UUID REFERENCES lists(id) not null,
    priority priority not null,
    deadline timestamp with time zone NULL,
    done BOOL NOT NULL,
    name TEXT not null,
    updated_at timestamp with time zone not null
);
