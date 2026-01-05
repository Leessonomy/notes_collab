import { Injectable } from '@angular/core';

@Injectable({ providedIn: 'root' })
export class WorkspaceInviteService {
  generateInviteLink(workspaceId: string): string {
    const token = btoa(`workspace:${workspaceId}:${Date.now()}`);
    return `${window.location.origin}/invite/${token}`;
  }
}
