import { Component, inject, input } from '@angular/core';
import { BrnSheetContent } from '@spartan-ng/brain/sheet';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideUsers } from '@ng-icons/lucide';
import { HlmButtonImports } from '@spartan-ng/helm/button';
import { HlmIcon } from '@spartan-ng/helm/icon';
import { HlmSheetImports } from '@spartan-ng/helm/sheet';
import { MembersStore } from '../application/members.store';
import { MembersFacade } from '../application/members.facade';
import { Workspace } from '../domain/workspace.model';
import { InviteLinksComponent } from './invite-links.component';
import { MemberListComponent } from './member-list.component';

@Component({
  selector: 'app-members-sheet',
  imports: [
    BrnSheetContent,
    HlmSheetImports,
    HlmButtonImports,
    NgIcon,
    HlmIcon,
    InviteLinksComponent,
    MemberListComponent,
  ],
  providers: [MembersStore, MembersFacade, provideIcons({ lucideUsers })],
  template: `
    <hlm-sheet side="right">
      <button
        hlmSheetTrigger
        hlmBtn
        variant="outline"
        size="sm"
        (click)="membersFacade.load(workspace())"
      >
        <ng-icon hlm name="lucideUsers" size="sm" />
        Members
      </button>

      <hlm-sheet-content *brnSheetContent>
        <hlm-sheet-header>
          <h3 hlmSheetTitle>{{ workspace().name }}</h3>
          <p hlmSheetDescription>Who has access to this workspace</p>
        </hlm-sheet-header>

        <div class="flex flex-1 flex-col gap-8 overflow-y-auto px-4 pb-6">
          @if (workspace().rights.manage) {
            <app-invite-links />
          }
          <app-member-list [readonly]="!workspace().rights.manage" />

          @if (membersFacade.error(); as message) {
            <p class="text-sm text-destructive">{{ message }}</p>
          }
        </div>
      </hlm-sheet-content>
    </hlm-sheet>
  `,
})
export class MembersSheetComponent {
  protected readonly membersFacade = inject(MembersFacade);

  readonly workspace = input.required<Workspace>();
}
