import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Workspace } from '../domain/workspace.model';

@Injectable({ providedIn: 'root' })
export class WorkspaceApi {
  private readonly http = inject(HttpClient);

  create({ name }: { name: string }) {
    return this.http.post<Workspace>('/api/workspaces', { name });
  }
}
