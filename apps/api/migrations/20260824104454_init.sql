-- +goose Up

ALTER TABLE notes
ADD CONSTRAINT notes_workspace_id_fkey
FOREIGN KEY (workspace_id)
REFERENCES workspaces(id)
ON DELETE CASCADE;


-- +goose Down

ALTER TABLE notes
DROP CONSTRAINT notes_workspace_id_fkey;