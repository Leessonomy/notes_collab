import { inject, Injectable } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { catchError, debounceTime, EMPTY, Subject, switchMap } from 'rxjs';
import { NotesStore } from './notes.store';
import { NotesApi } from '../api/notes.api';
import { NoteTabsService } from './note-tabs.service';
import { Note } from '../domain/note.model';

@Injectable()
export class NotesFacade {
  private readonly store = inject(NotesStore);
  private readonly api = inject(NotesApi);
  private readonly tabs = inject(NoteTabsService);

  readonly notes = this.store.notes;
  readonly error = this.store.error;

  private readonly saveRequests = new Subject<{ noteId: string; content: string }>();

  constructor() {
    this.saveRequests
      .pipe(
        debounceTime(1000),
        switchMap(({ noteId, content }) =>
          this.api.update(noteId, { content }).pipe(
            catchError(() => {
              this.store.setError('Failed to save note');
              return EMPTY;
            }),
          ),
        ),
        takeUntilDestroyed(),
      )
      .subscribe((updated) => this.store.upsert(updated));
  }

  set(notes: Note[]) {
    this.store.setNotes(notes);
  }

  byWorkspace(workspaceId: string) {
    return this.store.byWorkspace(workspaceId);
  }

  createAndOpen(workspaceId: string) {
    this.api
      .create({
        workspaceId,
        title: 'Untitled',
        content: '',
      })
      .subscribe({
        next: (note) => {
          this.store.add(note);
          this.tabs.openNote(note);
        },
        error: () => this.store.setError('Failed to create note'),
      });
  }

  openNote(note: { id: string; title: string }) {
    this.tabs.openNote(note);
  }

  saveContent(noteId: string, content: string) {
    this.saveRequests.next({ noteId, content });
  }
}
