import { computed, Injectable, signal } from '@angular/core';
import { Note } from '../domain/note.model';

interface NotesState {
  notes: Note[];
  error: string | null;
}

@Injectable()
export class NotesStore {
  private readonly state = signal<NotesState>({
    notes: [],
    error: null,
  });

  readonly notes = computed(() => this.state().notes);
  readonly error = computed(() => this.state().error);

  byWorkspace(workspaceId: string): Note[] {
    return this.state().notes.filter((n) => n.workspaceId === workspaceId);
  }

  noteById(id: string | null): Note | null {
    return id ? (this.state().notes.find((n) => n.id === id) ?? null) : null;
  }

  setNotes(notes: Note[]) {
    this.state.update((state) => ({ ...state, notes }));
  }

  add(note: Note) {
    this.state.update((state) => ({ ...state, notes: [...state.notes, note] }));
  }

  delete(noteId: string) {
    this.state.update((state) => ({...state, notes: state.notes.filter(({id}) => id !== noteId)}))
  }

  upsert(note: Note) {
    this.state.update((state) => ({
      ...state,
      notes: state.notes.some((n) => n.id === note.id)
        ? state.notes.map((n) => (n.id === note.id ? note : n))
        : [...state.notes, note],
    }));
  }

  setError(error: string | null) {
    this.state.update((state) => ({ ...state, error }));
  }
}
