-- +goose Up

CREATE TABLE invites (
    id text PRIMARY KEY,
    workspace_id text NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    token_hash bytea NOT NULL UNIQUE,
    role text NOT NULL CHECK (role IN ('collaborator', 'guest')),
    created_by text NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at timestamptz NOT NULL,
    accepted_by text REFERENCES users(id) ON DELETE SET NULL,
    accepted_at timestamptz,
    created_at timestamptz NOT NULL
);
CREATE INDEX invites_workspace_id_idx ON invites(workspace_id);

-- +goose Down

DROP TABLE invites;
