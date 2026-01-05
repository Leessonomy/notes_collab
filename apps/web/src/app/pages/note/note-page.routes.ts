import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: ':id',
    loadComponent: () => import('./note-page.component').then((m) => m.NotePageComponent),
  },
];
