import { computed, inject, Injectable } from '@angular/core';
import { appUrls } from '../../../core/app-urls';
import { MembersStore } from './members.store';
import { MembersApi } from '../api/members.api';
import { InvitesApi } from '../api/invites.api';
import { Workspace } from '../domain/workspace.model';
import { Member } from '../domain/member.model';
import { Invite } from '../domain/invite.model';
import { Role } from '../domain/role.model';

@Injectable()
export class MembersFacade {
  private readonly store = inject(MembersStore);
  private readonly membersApi = inject(MembersApi);
  private readonly invitesApi = inject(InvitesApi);

  readonly members = this.store.members;
  readonly invites = this.store.invites;
  readonly error = this.store.error;

  readonly newInviteLink = computed(() => {
    const invite = this.store.newInvite();
    return invite ? window.location.origin + appUrls.invite(invite.token) : null;
  });

  load(workspace: Workspace) {
    this.store.reset(workspace.id);

    this.membersApi.list(workspace.id).subscribe({
      next: (members) => this.store.setMembers(members),
      error: () => this.store.setError("Couldn't load members."),
    });

    if (workspace.rights.manage) {
      this.invitesApi.list(workspace.id).subscribe({
        next: (invites) => this.store.setInvites(invites),
        error: () => this.store.setError("Couldn't load invites."),
      });
    }
  }

  changeRole(member: Member, role: Role) {
    this.store.setRole(member.userId, role);

    this.membersApi.changeRole(this.store.workspaceId(), member.userId, role).subscribe({
      error: () => {
        this.store.setRole(member.userId, member.role);
        this.store.setError(`Couldn't change ${member.name}'s role.`);
      },
    });
  }

  removeMember(member: Member) {
    this.membersApi.remove(this.store.workspaceId(), member.userId).subscribe({
      next: () => this.store.removeMember(member.userId),
      error: () => this.store.setError(`Couldn't remove ${member.name}.`),
    });
  }

  createInvite(role: Role) {
    this.invitesApi.create(this.store.workspaceId(), role).subscribe({
      next: (invite) => this.store.addInvite(invite),
      error: () => this.store.setError("Couldn't create an invite link."),
    });
  }

  revokeInvite(invite: Invite) {
    this.invitesApi.revoke(this.store.workspaceId(), invite.id).subscribe({
      next: () => this.store.removeInvite(invite.id),
      error: () => this.store.setError("Couldn't revoke the invite."),
    });
  }
}
