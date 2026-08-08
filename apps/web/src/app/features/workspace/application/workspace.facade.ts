import { computed, inject, Injectable } from '@angular/core';
import { tap } from 'rxjs';
import { WorkspaceStore } from './workspace.store';
import { WorkspaceApi } from '../api/workspace.api';
import { Workspace } from '../domain/workspace.model';

@Injectable()
export class WorkspaceFacade {
  private readonly store = inject(WorkspaceStore);
  private readonly api = inject(WorkspaceApi);

  readonly workspaces = this.store.workspaces;
  readonly error = this.store.error;

  readonly currentWorkspace = computed<Workspace | null>(
    () => this.store.activeWorkspace() ?? this.store.workspaces()[0] ?? null,
  );

  set(workspaces: Workspace[]) {
    this.store.setWorkspaces(workspaces);
  }

  create(name: string) {
    return this.api.create({ name }).pipe(
      tap((workspace) => {
        this.store.add(workspace);
        this.store.switchActive(workspace);
      }),
    );
  }

  switch(id: string) {
    const workspace = this.store.workspaces().find((w) => w.id === id);
    if (workspace) this.store.switchActive(workspace);
  }
}
