import { computed, Injectable, signal } from '@angular/core';
import { Workspace } from '../domain/workspace.model';

interface WorkspaceState {
  workspaces: Workspace[];
  activeWorkspaceId: string | null;
  error: string | null;
}

@Injectable()
export class WorkspaceStore {
  private readonly state = signal<WorkspaceState>({
    workspaces: [],
    activeWorkspaceId: null,
    error: null,
  });

  readonly workspaces = computed(() => this.state().workspaces);
  readonly error = computed(() => this.state().error);

  readonly activeWorkspace = computed(() => {
    const { workspaces, activeWorkspaceId } = this.state();
    return workspaces.find((w) => w.id === activeWorkspaceId) ?? null;
  });

  setWorkspaces(workspaces: Workspace[]) {
    this.state.update((state) => ({ ...state, workspaces }));
  }

  add(workspace: Workspace) {
    this.state.update((state) => ({
      ...state,
      workspaces: [...state.workspaces, workspace],
    }));
  }

  delete(workspaceId: string) {
    this.state.update((state) => ({
      ...state,
      workspaces: state.workspaces.filter((w) => w.id !== workspaceId),
    }));
  }

  switchActive(workspaceId: string) {
    this.state.update((state) => ({ ...state, activeWorkspaceId: workspaceId }));
  }

  setError(error: string | null) {
    this.state.update((state) => ({ ...state, error }));
  }
}
