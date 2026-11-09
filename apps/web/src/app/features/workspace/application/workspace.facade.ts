import { computed, inject, Injectable } from '@angular/core';
import { tap } from 'rxjs';
import { AuthService } from '../../../core/auth/auth.service';
import { WorkspaceStore } from './workspace.store';
import { WorkspaceApi } from '../api/workspace.api';
import { MembersApi } from '../api/members.api';
import { Workspace } from '../domain/workspace.model';

@Injectable()
export class WorkspaceFacade {
  private readonly store = inject(WorkspaceStore);
  private readonly api = inject(WorkspaceApi);
  private readonly membersApi = inject(MembersApi);
  private readonly auth = inject(AuthService);

  readonly workspaces = this.store.workspaces;
  readonly error = this.store.error;

  readonly currentWorkspace = computed<Workspace | null>(
    () => this.store.activeWorkspace() ?? this.store.workspaces()[0] ?? null,
  );

  find(id: string): Workspace | null {
    return this.store.workspaces().find((w) => w.id === id) ?? null;
  }

  set(workspaces: Workspace[]) {
    this.store.setWorkspaces(workspaces);
  }

  create(name: string) {
    return this.api.create({ name }).pipe(
      tap((workspace) => {
        this.store.add(workspace);
        this.store.switchActive(workspace.id);
      }),
    );
  }

  delete(id: string) {
    return this.api.delete(id).pipe(tap(() => this.store.delete(id)));
  }

  leave(id: string) {
    const userId = this.auth.me()!.id;
    return this.membersApi.remove(id, userId).pipe(tap(() => this.store.delete(id)));
  }

  switch(id: string) {
    this.store.switchActive(id);
  }
}
