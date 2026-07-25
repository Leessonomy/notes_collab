import { Component, inject, signal } from '@angular/core';
import { Field, email, form, required, submit, validate } from '@angular/forms/signals';
import { Router, RouterLink } from '@angular/router';
import { HlmButtonImports } from '@spartan-ng/helm/button';
import { HlmCardImports } from '@spartan-ng/helm/card';
import { HlmInputImports } from '@spartan-ng/helm/input';
import { HlmLabelImports } from '@spartan-ng/helm/label';
import { AuthService } from '../../core/auth/auth.service';

@Component({
  selector: 'sign-in-page',
  imports: [Field, RouterLink, HlmCardImports, HlmLabelImports, HlmInputImports, HlmButtonImports],
  host: { class: 'w-full max-w-md' },
  template: `
    <div hlmCard class="w-full max-w-sm">
      <div hlmCardHeader>
        <h3 hlmCardTitle>Login to your account</h3>
        <p hlmCardDescription>Enter your email below to login to your account</p>

        <div hlmCardAction>
          <a hlmBtn variant="link" routerLink="/auth/signup">Sign Up</a>
        </div>
      </div>

      <div hlmCardContent>
        <form id="login-form" (submit)="onSubmit($event)">
          <div class="flex flex-col gap-6">
            <div class="grid gap-2">
              <label hlmLabel for="email">Email</label>
              <input
                type="email"
                id="email"
                placeholder="m@example.com"
                hlmInput
                [field]="form.email"
              />
              @if (form.email().touched()) {
                @for (error of form.email().errors(); track error.kind) {
                  <p class="text-destructive text-sm">{{ error.message }}</p>
                }
              }
            </div>

            <div class="grid gap-2">
              <div class="flex items-center">
                <label hlmLabel for="password">Password</label>
                <a href="#" class="ml-auto inline-block text-sm underline-offset-4 hover:underline">
                  Forgot your password?
                </a>
              </div>
              <input type="password" id="password" hlmInput [field]="form.password" />
              @if (form.password().touched()) {
                @for (error of form.password().errors(); track error.kind) {
                  <p class="text-destructive text-sm">{{ error.message }}</p>
                }
              }
            </div>
          </div>

          <div hlmCardFooter class="flex-col gap-2 px-0 pt-6">
            <button hlmBtn type="submit" class="w-full" [disabled]="form().submitting()">
              Login
            </button>
            <button hlmBtn type="button" variant="outline" class="w-full">Login with Google</button>
          </div>
        </form>
      </div>
    </div>
  `,
})
export class SignInPageComponent {
  private readonly auth = inject(AuthService);
  private readonly router = inject(Router);

  private readonly model = signal({ email: '', password: '' });

  readonly form = form(this.model, (schema) => {
    required(schema.email, { message: 'Email is required.' });
    email(schema.email, { message: 'Enter a valid email address.' });
    required(schema.password, { message: 'Password is required.' });
  });

  onSubmit(event: Event) {
    event.preventDefault();
    submit(this.form, async (field) => {
      const { email, password } = this.model();
      try {
        this.auth.login({ email, password }).subscribe(() => {
          this.router.navigateByUrl('/app');
        });
      } catch (e) {}
    });
  }
}
