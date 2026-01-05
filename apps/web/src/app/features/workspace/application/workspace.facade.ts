import { computed, inject, Injectable } from '@angular/core';
import { tap } from 'rxjs';
import { WorkspaceStore } from './workspace.store';
import { WorkspaceService } from '../api/workspace.service';
import { AuthService } from '../../../core/auth/auth.service';

@Injectable()
export class WorkspaceFacade {
  private readonly store = inject(WorkspaceStore);
  private readonly api = inject(WorkspaceService);
  private readonly auth = inject(AuthService);

  readonly workspaces = this.store.workspaces;
  readonly isLoading = this.store.isLoading;
  readonly error = this.store.error;

  readonly currentWorkspace = computed(
    () => this.store.activeWorkspace() ?? this.store.workspaces()[0] ?? null,
  );

  load() {
    this.store.setLoading(true);
    this.api.getAll().subscribe({
      next: (workspaces) => {
        this.store.setWorkspaces(workspaces);
        this.store.setError(null);
        this.store.setLoading(false);
      },
      error: () => {
        this.store.setError('Failed to load workspaces');
        this.store.setLoading(false);
      },
    });
  }

  create(name: string) {
    return this.api.create({ name, ownerId: this.auth.getCurrentUserId() }).pipe(
      tap((workspace) => {
        this.store.addNew(workspace);
        this.store.switchActive(workspace);
      }),
    );
  }

  switch(id: string) {
    const workspace = this.store.workspaces().find((w) => w.id === id);
    if (workspace) this.store.switchActive(workspace);
  }
}
