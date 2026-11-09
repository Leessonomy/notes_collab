import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: ':token',
    loadComponent: () => import('./invite-page.component').then((m) => m.InvitePageComponent),
  },
];
