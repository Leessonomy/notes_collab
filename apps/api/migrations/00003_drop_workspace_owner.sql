-- +goose Up

ALTER TABLE workspaces DROP COLUMN owner_id;

-- +goose Down

ALTER TABLE workspaces ADD COLUMN owner_id text REFERENCES users(id) ON DELETE CASCADE;

UPDATE workspaces w
SET owner_id = m.user_id
FROM workspace_members m
WHERE m.workspace_id = w.id AND m.role = 'owner';

ALTER TABLE workspaces ALTER COLUMN owner_id SET NOT NULL;
CREATE INDEX workspaces_owner_id_idx ON workspaces(owner_id);
