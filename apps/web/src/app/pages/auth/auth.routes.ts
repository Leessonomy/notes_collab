import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: '',
    loadComponent: () => import('./sign-in-page.component').then((m) => m.SignInPageComponent),
  },
  {
    path: 'signup',
    loadComponent: () => import('./sign-up-page.component').then((m) => m.SignUpPageComponent),
  },
];
