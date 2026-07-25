import { Component, inject, signal } from '@angular/core';
import { Router, RouterLink } from '@angular/router';
import {
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
import { AuthService } from '../../core/auth/auth.service';
import { FormFieldComponent } from '../../shared/ui/form-field/form-field.component';

@Component({
  selector: 'sign-up-page',
  imports: [RouterLink, HlmCardImports, HlmButtonImports, FormFieldComponent],
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
            <form-field [field]="form.name" label="Full Name" id="name" placeholder="John Doe" />

            <form-field
              [field]="form.email"
              label="Email"
              type="email"
              id="email"
              placeholder="m@example.com"
            />

            <form-field [field]="form.password" label="Password" type="password" id="password" />

            <form-field
              [field]="form.confirmPassword"
              label="Confirm Password"
              type="password"
              id="confirmPassword"
            />
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
        this.auth.signup({ email, password, name }).subscribe((user) => {
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
