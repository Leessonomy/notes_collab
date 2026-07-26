import { Routes } from '@angular/router';
import { MainLayoutComponent } from './layout/main-layout.component';
import { AuthLayoutComponent } from './layout/auth-layout.component';
import { WorkspaceStore, WorkspaceFacade } from './features/workspace';
import { NotesStore, NotesFacade, NoteTabsService, PresenceService } from './features/notes';
import { authGuard, guestGuard } from './core/auth/auth.guard';

export const routes: Routes = [
  {
    path: 'auth',
    canActivate: [guestGuard],
    component: AuthLayoutComponent,
    children: [
      {
        path: '',
        loadChildren: () => import('./pages/auth/auth.routes').then((m) => m.routes),
      },
    ],
  },
  {
    path: 'app',
    canActivate: [authGuard],
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
        loadChildren: () => import('./pages/workspace/workspace.routes').then((m) => m.routes),
      },
      {
        path: 'note',
        loadChildren: () => import('./pages/note/note.routes').then((m) => m.routes),
      },
    ],
  },
  { path: '', redirectTo: 'app', pathMatch: 'full' },
  { path: '**', redirectTo: 'app' },
];
