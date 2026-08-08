import { computed, Injectable, signal } from '@angular/core';
import { Workspace } from '../domain/workspace.model';

interface WorkspaceState {
  workspaces: Workspace[];
  activeWorkspace: Workspace | null;
  error: string | null;
}

@Injectable()
export class WorkspaceStore {
  private readonly state = signal<WorkspaceState>({
    workspaces: [],
    activeWorkspace: null,
    error: null,
  });

  readonly workspaces = computed(() => this.state().workspaces);
  readonly activeWorkspace = computed(() => this.state().activeWorkspace);
  readonly error = computed(() => this.state().error);

  setWorkspaces(workspaces: Workspace[]) {
    this.state.update((state) => ({ ...state, workspaces }));
  }

  add(workspace: Workspace) {
    this.state.update((state) => ({
      ...state,
      workspaces: [...state.workspaces, workspace],
    }));
  }

  switchActive(workspace: Workspace) {
    this.state.update((state) => ({ ...state, activeWorkspace: workspace }));
  }

  setError(error: string | null) {
    this.state.update((state) => ({ ...state, error }));
  }
}
