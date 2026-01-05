import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Workspace } from '../domain/workspace.model';

@Injectable({ providedIn: 'root' })
export class WorkspaceService {
  private readonly http = inject(HttpClient);

  getAll() {
    return this.http.get<Workspace[]>('/api/workspaces');
  }

  create({ name, ownerId }: { name: string; ownerId: string }) {
    return this.http.post<Workspace>('/api/workspaces', {
      name,
      ownerId,
    });
  }
}
