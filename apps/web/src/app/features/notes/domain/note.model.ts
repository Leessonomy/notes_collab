export interface Note {
  id: string;
  workspaceId: string;
  title: string;
  content: string;
  ownerId: string;
  createdAt: Date;
  updatedAt: Date;
}
