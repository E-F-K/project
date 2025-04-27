CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    token TEXT NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(email),
    UNIQUE(token)
);

CREATE TABLE IF NOT EXISTS lists (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    name TEXT NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TYPE priority AS ENUM ('low', 'normal', 'high');

CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY,
    list_id UUID NOT NULL,
    priority priority NOT NULL,
    deadline TIMESTAMP WITH TIME ZONE NULL,
    done BOOL NOT NULL,
    name TEXT NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY(list_id) REFERENCES lists(id) ON DELETE CASCADE
);
