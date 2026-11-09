import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { CreatedInvite, Invite } from '../domain/invite.model';
import { Role } from '../domain/role.model';

@Injectable({ providedIn: 'root' })
export class InvitesApi {
  private readonly http = inject(HttpClient);

  list(workspaceId: string) {
    return this.http.get<Invite[]>(`/api/workspaces/${workspaceId}/invites`);
  }

  create(workspaceId: string, role: Role) {
    return this.http.post<CreatedInvite>(`/api/workspaces/${workspaceId}/invites`, { role });
  }

  revoke(workspaceId: string, inviteId: string) {
    return this.http.delete(`/api/workspaces/${workspaceId}/invites/${inviteId}`);
  }

  accept(token: string) {
    return this.http.post<{ workspaceId: string }>(`/api/invites/${token}/accept`, {});
  }
}
