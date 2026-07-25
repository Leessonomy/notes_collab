import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';
import { HlmButtonImports } from '@spartan-ng/helm/button';
import { HlmCardImports } from '@spartan-ng/helm/card';
import { HlmInputImports } from '@spartan-ng/helm/input';
import { HlmLabelImports } from '@spartan-ng/helm/label';

@Component({
  selector: 'sign-in-page',
  imports: [RouterLink, HlmCardImports, HlmLabelImports, HlmInputImports, HlmButtonImports],
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
        <form id="login-form">
          <div class="flex flex-col gap-6">
            <div class="grid gap-2">
              <label hlmLabel for="email">Email</label>
              <input type="email" id="email" placeholder="m@example.com" required hlmInput />
            </div>

            <div class="grid gap-2">
              <div class="flex items-center">
                <label hlmLabel for="password">Password</label>
                <a href="#" class="ml-auto inline-block text-sm underline-offset-4 hover:underline">
                  Forgot your password?
                </a>
              </div>
              <input type="password" id="password" hlmInput />
            </div>
          </div>
        </form>
      </div>

      <div hlmCardFooter class="flex-col gap-2">
        <button hlmBtn type="submit" class="w-full" form="login-form">Login</button>
        <button hlmBtn variant="outline" class="w-full">Login with Google</button>
      </div>
    </div>
  `,
})
export class SignInPageComponent {
  ngOnInit() {}
}
