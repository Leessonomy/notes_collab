import { Component, inject, signal } from '@angular/core';
import { DatePipe, TitleCasePipe } from '@angular/common';
import { CdkCopyToClipboard } from '@angular/cdk/clipboard';
import { HlmButtonImports } from '@spartan-ng/helm/button';
import { HlmInputImports } from '@spartan-ng/helm/input';
import { MembersFacade } from '../application/members.facade';
import { assignableRoles, Role } from '../domain/role.model';

@Component({
  selector: 'app-invite-links',
  imports: [HlmButtonImports, HlmInputImports, CdkCopyToClipboard, DatePipe, TitleCasePipe],
  template: `
    <section class="flex flex-col gap-3">
      <div>
        <h4 class="text-sm font-medium">Invite people</h4>
        <p class="text-xs text-muted-foreground">Each link works once and expires in 7 days.</p>
      </div>

      <div class="flex gap-2">
        <select
          #roleSelect
          aria-label="Role"
          class="h-8 flex-1 rounded-md border border-input bg-transparent px-2 text-sm shadow-xs outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 dark:bg-input/30"
        >
          @for (role of assignableRoles; track role) {
            <option [value]="role">{{ role | titlecase }}</option>
          }
        </select>
        <button hlmBtn size="sm" (click)="createInvite(roleSelect.value)">Create link</button>
      </div>

      @if (membersFacade.newInviteLink(); as link) {
        <div class="flex gap-2">
          <input
            #linkInput
            hlmInput
            readonly
            aria-label="Invite link"
            class="h-8"
            [value]="link"
            (focus)="linkInput.select()"
          />
          <button
            hlmBtn
            variant="outline"
            size="sm"
            [cdkCopyToClipboard]="link"
            (cdkCopyToClipboardCopied)="copied.set($event)"
          >
            {{ copied() ? 'Copied' : 'Copy' }}
          </button>
        </div>
        <p class="text-xs text-muted-foreground">Copy it now. You won't see this link again.</p>
      }

      @if (membersFacade.invites().length) {
        <ul class="divide-y divide-border">
          @for (invite of membersFacade.invites(); track invite.id) {
            <li class="flex items-center justify-between gap-2 py-1.5 text-sm">
              <span>
                {{ invite.role | titlecase }}
                <span class="text-muted-foreground">
                  · expires {{ invite.expiresAt | date: 'MMM d' }}
                </span>
              </span>
              <button hlmBtn variant="ghost" size="sm" (click)="membersFacade.revokeInvite(invite)">
                Revoke
              </button>
            </li>
          }
        </ul>
      }
    </section>
  `,
})
export class InviteLinksComponent {
  protected readonly membersFacade = inject(MembersFacade);
  protected readonly assignableRoles = assignableRoles;
  protected readonly copied = signal(false);

  createInvite(role: string) {
    this.copied.set(false);
    this.membersFacade.createInvite(role as Role);
  }
}
