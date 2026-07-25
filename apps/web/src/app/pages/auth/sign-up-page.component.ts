import { Component, inject, signal } from '@angular/core';
import { Router, RouterLink } from '@angular/router';
import {
  Field,
  customError,
  email,
  form,
  minLength,
  required,
  submit,
  validate,
} from '@angular/forms/signals';
import { HlmButtonImports } from '@spartan-ng/helm/button';
import { HlmCardImports } from '@spartan-ng/helm/card';
import { HlmInputImports } from '@spartan-ng/helm/input';
import { HlmLabelImports } from '@spartan-ng/helm/label';
import { AuthService } from '../../core/auth/auth.service';

@Component({
  selector: 'sign-up-page',
  imports: [Field, RouterLink, HlmCardImports, HlmLabelImports, HlmInputImports, HlmButtonImports],
  host: { class: 'w-full max-w-md' },
  template: `
    <div hlmCard class="w-full max-w-sm">
      <div hlmCardHeader>
        <h3 hlmCardTitle>Create an account</h3>
        <p hlmCardDescription>Enter your information below to create your account</p>
      </div>

      <div hlmCardContent>
        <form id="signup-form" (submit)="onSubmit($event)">
          <div class="flex flex-col gap-6">
            <div class="grid gap-2">
              <label hlmLabel for="name">Full Name</label>
              <input type="text" id="name" placeholder="John Doe" hlmInput [field]="form.name" />
              @if (form.name().touched()) {
                @for (error of form.name().errors(); track error.kind) {
                  <p class="text-destructive text-sm">{{ error.message }}</p>
                }
              }
            </div>

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
              <label hlmLabel for="password">Password</label>
              <input type="password" id="password" hlmInput [field]="form.password" />
              @if (form.password().touched()) {
                @for (error of form.password().errors(); track error.kind) {
                  <p class="text-destructive text-sm">{{ error.message }}</p>
                }
              }
            </div>

            <div class="grid gap-2">
              <label hlmLabel for="confirmPassword">Confirm Password</label>
              <input type="password" id="confirmPassword" hlmInput [field]="form.confirmPassword" />
              @if (form.confirmPassword().touched()) {
                @for (error of form.confirmPassword().errors(); track error.kind) {
                  <p class="text-destructive text-sm">{{ error.message }}</p>
                }
              }
            </div>
          </div>

          <div hlmCardFooter class="flex-col gap-2 px-0 pt-6">
            <button hlmBtn type="submit" class="w-full" [disabled]="form().submitting()">
              Create Account
            </button>
            <p class="text-center text-sm">
              Already have an account?
              <a hlmBtn variant="link" routerLink="/auth">Sign in</a>
            </p>
          </div>
        </form>
      </div>
    </div>
  `,
})
export class SignUpPageComponent {
  private readonly auth = inject(AuthService);
  private readonly router = inject(Router);

  private readonly model = signal({ name: '', email: '', password: '', confirmPassword: '' });

  readonly form = form(this.model, (schema) => {
    required(schema.name, { message: 'Name is required.' });
    required(schema.email, { message: 'Email is required.' });
    email(schema.email, { message: 'Enter a valid email address.' });
    required(schema.password, { message: 'Password is required.' });
    minLength(schema.password, 8, { message: 'Password must be at least 8 characters long.' });
    required(schema.confirmPassword, { message: 'Confirming your password is required.' });

    validate(schema.confirmPassword, ({ value, valueOf, stateOf }) => {
      if (!stateOf(schema.password).touched()) {
        return null;
      }
      if (value() !== valueOf(schema.password)) {
        return { kind: 'passwordMismatch', message: 'Passwords must match.' };
      }
      return null;
    });
  });

  onSubmit(event: Event) {
    event.preventDefault();
    submit(this.form, async (field) => {
      const { email, password, name } = this.model();
      try {
        this.auth.signUp({ email, password, name }).subscribe((user) => {
          this.router.navigateByUrl('/app');
        });
        return;
      } catch (e) {
        return customError({
          kind: 'server',
          message: 'This email is already registered.',
          fieldTree: field.email,
        });
      }
    });
  }
}
