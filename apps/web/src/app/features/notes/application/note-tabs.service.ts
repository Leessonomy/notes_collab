import { Injectable, computed, inject, signal } from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { NavigationEnd, Router } from '@angular/router';
import { filter, map, startWith } from 'rxjs';
import { NotesStore } from './notes.store';
import { appUrls } from '../../../core/app-urls';

export interface NoteTab {
  noteId: string;
  title: string;
}

@Injectable()
export class NoteTabsService {
  private readonly router = inject(Router);
  private readonly store = inject(NotesStore);

  private readonly _tabs = signal<NoteTab[]>([]);

  private readonly activeNoteId = toSignal(
    this.router.events.pipe(
      filter((e) => e instanceof NavigationEnd),
      map(() => this.noteIdFromUrl()),
      startWith(this.noteIdFromUrl()),
    ),
    { initialValue: this.noteIdFromUrl() },
  );

  readonly tabs = this._tabs.asReadonly();

  readonly activeTab = computed(
    () => this._tabs().find((tab) => tab.noteId === this.activeNoteId()) ?? null,
  );

  readonly currentNote = computed(() => this.store.noteById(this.activeNoteId()));

  openNote(note: { id: string; title: string }) {
    if (!this._tabs().some((tab) => tab.noteId === note.id)) {
      this._tabs.update((tabs) => [...tabs, { noteId: note.id, title: note.title || 'Untitled' }]);
    }
    this.router.navigateByUrl(appUrls.note(note.id));
  }

  closeTab(noteId: string) {
    const tabs = this._tabs();
    const index = tabs.findIndex((tab) => tab.noteId === noteId);
    const remaining = tabs.filter((tab) => tab.noteId !== noteId);
    this._tabs.set(remaining);

    if (this.activeNoteId() !== noteId) return;

    const next = remaining[Math.min(index, remaining.length - 1)];
    this.router.navigateByUrl(next ? appUrls.note(next.noteId) : appUrls.workspace);
  }

  private noteIdFromUrl(): string | null {
    const match = this.router.url.match(/^\/app\/note\/([^/?#]+)/);
    return match ? decodeURIComponent(match[1]) : null;
  }
}
