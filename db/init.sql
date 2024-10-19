create table if not exists users (
    id UUID PRIMARY KEY,
    name TEXT not null,
    email TEXT not null,
    password_hash TEXT not null,
    token TEXT not null,
    updated_at timestamp with time zone,
    unique(email),
    unique(token)
);

create table if not exists lists (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    name TEXT not null,
    email TEXT not null,
    updated_at timestamp with time zone
);

create type priority AS ENUM ('low', 'normal', 'high');

create table if not exists tasks (
    id UUID PRIMARY KEY,
    list_id UUID REFERENCES lists(id),
    priority priority,
    deadline timestamp with time zone NULL,
    done BOOL NOT NULL,
    name TEXT not null,
    updated_at timestamp with time zone
);
