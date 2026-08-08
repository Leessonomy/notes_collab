import { computed, inject, Injectable, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { WorkspaceFacade, Workspace } from './features/workspace';
import { NotesFacade, Note } from './features/notes';

interface WorkspaceWithNotes extends Workspace {
  notes: Note[];
}

interface AppState {
  isLoading: boolean;
  error: string | null;
}

@Injectable()
export class AppFacade {
  private readonly http = inject(HttpClient);
  private readonly workspaces = inject(WorkspaceFacade);
  private readonly notes = inject(NotesFacade);

  private readonly state = signal<AppState>({
    isLoading: false,
    error: null,
  });

  readonly isLoading = computed(() => this.state().isLoading);
  readonly error = computed(() => this.state().error);

  load() {
    this.state.set({ isLoading: true, error: null });
    this.http.get<WorkspaceWithNotes[]>('/api/workspaces').subscribe({
      next: (data) => {
        this.workspaces.set(data.map(({ notes, ...workspace }) => workspace));
        this.notes.set(data.flatMap((w) => w.notes));
        this.state.set({ isLoading: false, error: null });
      },
      error: () => this.state.set({ isLoading: false, error: 'Failed to load workspaces' }),
    });
  }
}
