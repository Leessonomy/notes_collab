import { computed, Injectable, signal } from '@angular/core';
import { Member } from '../domain/member.model';
import { CreatedInvite, Invite } from '../domain/invite.model';
import { Role } from '../domain/role.model';

interface MembersState {
  workspaceId: string;
  members: Member[];
  invites: Invite[];
  newInvite: CreatedInvite | null;
  error: string | null;
}

@Injectable()
export class MembersStore {
  private readonly state = signal<MembersState>({
    workspaceId: '',
    members: [],
    invites: [],
    newInvite: null,
    error: null,
  });

  readonly workspaceId = computed(() => this.state().workspaceId);
  readonly members = computed(() => this.state().members);
  readonly invites = computed(() => this.state().invites);
  readonly newInvite = computed(() => this.state().newInvite);
  readonly error = computed(() => this.state().error);

  reset(workspaceId: string) {
    this.state.set({ workspaceId, members: [], invites: [], newInvite: null, error: null });
  }

  setMembers(members: Member[]) {
    this.state.update((state) => ({ ...state, members }));
  }

  setRole(userId: string, role: Role) {
    this.state.update((state) => ({
      ...state,
      members: state.members.map((m) => (m.userId === userId ? { ...m, role } : m)),
    }));
  }

  removeMember(userId: string) {
    this.state.update((state) => ({
      ...state,
      members: state.members.filter((m) => m.userId !== userId),
    }));
  }

  setInvites(invites: Invite[]) {
    this.state.update((state) => ({ ...state, invites }));
  }

  addInvite(invite: CreatedInvite) {
    this.state.update((state) => ({
      ...state,
      invites: [...state.invites, invite],
      newInvite: invite,
    }));
  }

  removeInvite(inviteId: string) {
    this.state.update((state) => ({
      ...state,
      invites: state.invites.filter((i) => i.id !== inviteId),
    }));
  }

  setError(error: string | null) {
    this.state.update((state) => ({ ...state, error }));
  }
}
