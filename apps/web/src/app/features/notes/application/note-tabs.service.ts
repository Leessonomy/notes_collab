import { Injectable, computed, inject, signal } from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { NavigationEnd, Router } from '@angular/router';
import { filter, map, startWith } from 'rxjs';
import { NotesStore } from './notes.store';
import { Note } from '../domain/note.model';
import { appUrls } from '../../../core/app-urls';

@Injectable()
export class NoteTabsService {
  private readonly router = inject(Router);
  private readonly store = inject(NotesStore);

  private readonly openNoteIds = signal<string[]>([]);

  private readonly activeNoteId = toSignal(
    this.router.events.pipe(
      filter((e) => e instanceof NavigationEnd),
      map(() => this.noteIdFromUrl()),
      startWith(this.noteIdFromUrl()),
    ),
    { initialValue: this.noteIdFromUrl() },
  );

  readonly tabs = computed<Note[]>(() =>
    this.openNoteIds().flatMap((noteId) => this.store.noteById(noteId) ?? []),
  );

  readonly activeTab = computed(
    () => this.tabs().find((note) => note.id === this.activeNoteId()) ?? null,
  );

  readonly currentNote = computed(() => this.store.noteById(this.activeNoteId()));

  openNote(noteId: string) {
    this.openNoteIds.update((ids) => (ids.includes(noteId) ? ids : [...ids, noteId]));
    this.router.navigateByUrl(appUrls.note(noteId));
  }

  closeTab(noteId: string) {
    const ids = this.openNoteIds();
    const closedIndex = ids.indexOf(noteId);

    this.openNoteIds.set(ids.filter((id) => id !== noteId));

    if (this.activeNoteId() !== noteId) return;

    const neighbour = ids[closedIndex + 1] ?? ids[closedIndex - 1];
    this.router.navigateByUrl(neighbour ? appUrls.note(neighbour) : appUrls.workspace);
  }

  private noteIdFromUrl(): string | null {
    const match = this.router.url.match(/^\/app\/note\/([^/?#]+)/);
    return match ? decodeURIComponent(match[1]) : null;
  }
}
