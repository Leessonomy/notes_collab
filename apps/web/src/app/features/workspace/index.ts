import { Provider } from '@angular/core';
import { WorkspaceStore } from './application/workspace.store';
import { WorkspaceFacade } from './application/workspace.facade';

export const provideWorkspace = (): Provider[] => [WorkspaceStore, WorkspaceFacade];

export { WorkspaceFacade };
export { InvitesApi } from './api/invites.api';
export { MembersSheetComponent } from './ui/members-sheet.component';
export type { Workspace } from './domain/workspace.model';
