import { Routes } from '@angular/router';
import { MainLayoutComponent } from './layout/main-layout.component';
import { WorkspaceStore, WorkspaceFacade } from './features/workspace';
import { NotesStore, NotesFacade, NoteTabsService, PresenceService } from './features/notes';

export const routes: Routes = [
  {
    path: '',
    component: MainLayoutComponent,
    providers: [
      WorkspaceStore,
      WorkspaceFacade,
      NotesStore,
      NotesFacade,
      NoteTabsService,
      PresenceService,
    ],
    children: [
      {
        path: '',
        redirectTo: 'workspace',
        pathMatch: 'full',
      },
      {
        path: 'workspace',
        loadChildren: () => import('./pages/workspace/workspace-page.routes').then((m) => m.routes),
      },
      {
        path: 'note',
        loadChildren: () => import('./pages/note/note-page.routes').then((m) => m.routes),
      },
    ],
  },
];
