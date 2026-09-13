-- +goose Up

CREATE TABLE workspace_members (
    workspace_id text NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    user_id text NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role text NOT NULL CHECK (role IN ('owner', 'collaborator', 'guest')),
    joined_at timestamptz NOT NULL,
    PRIMARY KEY (workspace_id, user_id)
);
CREATE INDEX workspace_members_user_id_idx ON workspace_members(user_id);
CREATE UNIQUE INDEX workspace_members_single_owner_idx ON workspace_members(workspace_id) WHERE role = 'owner';

INSERT INTO workspace_members (workspace_id, user_id, role, joined_at)
SELECT id, owner_id, 'owner', created_at FROM workspaces;

-- +goose Down

DROP TABLE workspace_members;
