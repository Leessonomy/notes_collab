import { Component, inject, input } from '@angular/core';
import { TitleCasePipe } from '@angular/common';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideX } from '@ng-icons/lucide';
import { HlmButtonImports } from '@spartan-ng/helm/button';
import { HlmIcon } from '@spartan-ng/helm/icon';
import { AuthService } from '../../../core/auth/auth.service';
import { MembersFacade } from '../application/members.facade';
import { Member } from '../domain/member.model';
import { assignableRoles, Role } from '../domain/role.model';

@Component({
  selector: 'app-member-list',
  imports: [HlmButtonImports, NgIcon, HlmIcon, TitleCasePipe],
  providers: [provideIcons({ lucideX })],
  template: `
    <section class="flex flex-col gap-3">
      <h4 class="text-sm font-medium">Members</h4>

      <ul class="divide-y divide-border">
        @for (member of membersFacade.members(); track member.userId) {
          <li class="flex items-center gap-2 py-2">
            <div class="min-w-0 flex-1">
              <p class="truncate text-sm font-medium">
                {{ member.name }}
                @if (member.userId === me()?.id) {
                  <span class="font-normal text-muted-foreground">(you)</span>
                }
              </p>
              <p class="truncate text-xs text-muted-foreground">{{ member.email }}</p>
            </div>

            @if (!readonly() && member.role !== 'owner') {
              <select
                #roleSelect
                [attr.aria-label]="member.name + ' role'"
                (change)="changeRole(member, roleSelect.value)"
                class="h-8 rounded-md border border-input bg-transparent px-2 text-sm shadow-xs outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 dark:bg-input/30"
              >
                @for (role of assignableRoles; track role) {
                  <option [value]="role" [selected]="role === member.role">
                    {{ role | titlecase }}
                  </option>
                }
              </select>
              <button
                hlmBtn
                variant="ghost"
                size="icon-sm"
                [attr.aria-label]="'Remove ' + member.name"
                (click)="removeMember(member)"
              >
                <ng-icon hlm name="lucideX" size="sm" />
              </button>
            } @else {
              <span class="text-xs text-muted-foreground">{{ member.role | titlecase }}</span>
            }
          </li>
        }
      </ul>
    </section>
  `,
})
export class MemberListComponent {
  protected readonly membersFacade = inject(MembersFacade);
  protected readonly me = inject(AuthService).me;
  protected readonly assignableRoles = assignableRoles;

  readonly readonly = input(false);

  changeRole(member: Member, role: string) {
    this.membersFacade.changeRole(member, role as Role);
  }

  removeMember(member: Member) {
    if (confirm(`Remove ${member.name} from this workspace?`)) {
      this.membersFacade.removeMember(member);
    }
  }
}
