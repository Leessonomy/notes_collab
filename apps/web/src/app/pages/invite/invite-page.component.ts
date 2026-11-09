import { Component, inject, input, OnInit, signal } from '@angular/core';
import { HttpErrorResponse } from '@angular/common/http';
import { Router, RouterLink } from '@angular/router';
import { HlmButtonImports } from '@spartan-ng/helm/button';
import { HlmCardImports } from '@spartan-ng/helm/card';
import { InvitesApi } from '../../features/workspace';
import { appUrls } from '../../core/app-urls';

@Component({
  selector: 'app-invite-page',
  imports: [RouterLink, HlmCardImports, HlmButtonImports],
  host: { class: 'w-full max-w-sm' },
  template: `
    <div hlmCard>
      @switch (status()) {
        @case ('invalid') {
          <div hlmCardHeader>
            <h3 hlmCardTitle>This invite doesn't work</h3>
            <p hlmCardDescription>
              The link has expired, was already used, or was revoked. Ask for a new one.
            </p>
          </div>
          <div hlmCardFooter>
            <a hlmBtn class="w-full" [routerLink]="workspaceUrl">Go to my workspaces</a>
          </div>
        }
        @case ('failed') {
          <div hlmCardHeader>
            <h3 hlmCardTitle>Couldn't join the workspace</h3>
            <p hlmCardDescription>Something went wrong on our side. Try again in a moment.</p>
          </div>
          <div hlmCardFooter>
            <button hlmBtn class="w-full" (click)="join()">Try again</button>
          </div>
        }
        @default {
          <div hlmCardHeader>
            <h3 hlmCardTitle>Joining workspace…</h3>
          </div>
        }
      }
    </div>
  `,
})
export class InvitePageComponent implements OnInit {
  private readonly invites = inject(InvitesApi);
  private readonly router = inject(Router);

  readonly token = input.required<string>();

  protected readonly status = signal<'joining' | 'invalid' | 'failed'>('joining');
  protected readonly workspaceUrl = appUrls.workspace;

  ngOnInit() {
    this.join();
  }

  join() {
    this.status.set('joining');

    this.invites.accept(this.token()).subscribe({
      next: ({ workspaceId }) =>
        this.router.navigate([appUrls.workspace], {
          queryParams: { workspace: workspaceId },
          replaceUrl: true,
        }),
      error: (error: HttpErrorResponse) =>
        this.status.set(error.status === 404 ? 'invalid' : 'failed'),
    });
  }
}
