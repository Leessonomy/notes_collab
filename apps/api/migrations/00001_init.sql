-- +goose Up

CREATE TABLE users (
    id text PRIMARY KEY,
    name text NOT NULL,
    email text NOT NULL UNIQUE,
    password text NOT NULL,
    created_at timestamptz NOT NULL
);

CREATE TABLE refresh_tokens (
    token text PRIMARY KEY,
    user_id text NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL
);
CREATE INDEX refresh_tokens_user_id_idx ON refresh_tokens(user_id);

CREATE TABLE workspaces (
    id text PRIMARY KEY,
    name text NOT NULL,
    owner_id text NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);
CREATE INDEX workspaces_owner_id_idx ON workspaces(owner_id);

CREATE TABLE notes (
    id text PRIMARY KEY,
    workspace_id text NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    title text NOT NULL DEFAULT '',
    content text NOT NULL DEFAULT '',
    owner_id text NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);
CREATE INDEX notes_workspace_id_idx ON notes(workspace_id);
CREATE INDEX notes_owner_id_idx ON notes(owner_id);

-- +goose Down

DROP TABLE notes;
DROP TABLE workspaces;
DROP TABLE refresh_tokens;
DROP TABLE users;
