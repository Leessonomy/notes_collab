import { Provider } from '@angular/core';
import { WorkspaceStore } from './application/workspace.store';
import { WorkspaceFacade } from './application/workspace.facade';

export const provideWorkspace = (): Provider[] => [WorkspaceStore, WorkspaceFacade];

export { WorkspaceFacade };
export { WorkspaceInviteService } from './api/workspace-invite.service';
export type { Workspace } from './domain/workspace.model';
