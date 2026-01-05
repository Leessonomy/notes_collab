import { computed, Injectable, signal } from '@angular/core';
import { Note } from '../domain/note.model';

interface NotesState {
  workspaceId: string | null;
  notes: Note[];
  isLoading: boolean;
  error: string | null;
}

@Injectable()
export class NotesStore {
  private readonly state = signal<NotesState>({
    workspaceId: null,
    notes: [],
    isLoading: false,
    error: null,
  });

  readonly notes = computed(() => this.state().notes);
  readonly isLoading = computed(() => this.state().isLoading);
  readonly error = computed(() => this.state().error);

  noteById(id: string | null): Note | null {
    return id ? (this.state().notes.find((n) => n.id === id) ?? null) : null;
  }

  setNotes(workspaceId: string, notes: Note[]) {
    this.state.update((state) => ({ ...state, workspaceId, notes }));
  }

  add(note: Note) {
    this.state.update((state) => ({ ...state, notes: [...state.notes, note] }));
  }

  upsert(note: Note) {
    this.state.update((state) => {
      if (note.workspaceId !== state.workspaceId) return state;
      return {
        ...state,
        notes: state.notes.some((n) => n.id === note.id)
          ? state.notes.map((n) => (n.id === note.id ? note : n))
          : [...state.notes, note],
      };
    });
  }

  setLoading(isLoading: boolean) {
    this.state.update((state) => ({ ...state, isLoading }));
  }

  setError(error: string | null) {
    this.state.update((state) => ({ ...state, error }));
  }
}
