import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Member } from '../domain/member.model';
import { Membership } from '../domain/workspace.model';
import { Role } from '../domain/role.model';

@Injectable({ providedIn: 'root' })
export class MembersApi {
  private readonly http = inject(HttpClient);

  list(workspaceId: string) {
    return this.http.get<Member[]>(`/api/workspaces/${workspaceId}/members`);
  }

  changeRole(workspaceId: string, memberId: string, role: Role) {
    return this.http.patch<Membership>(`/api/workspaces/${workspaceId}/members/${memberId}`, {
      role,
    });
  }

  remove(workspaceId: string, memberId: string) {
    return this.http.delete(`/api/workspaces/${workspaceId}/members/${memberId}`);
  }
}
