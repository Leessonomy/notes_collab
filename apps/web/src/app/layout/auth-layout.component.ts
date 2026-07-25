import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';

@Component({
  selector: 'auth-layout',
  imports: [RouterOutlet],
  host: { class: 'flex min-h-screen items-center justify-center p-4' },
  template: `<router-outlet />`,
})
export class AuthLayoutComponent {}
