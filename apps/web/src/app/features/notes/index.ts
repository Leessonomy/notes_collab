import { Provider } from '@angular/core';
import { NotesStore } from './application/notes.store';
import { NotesFacade } from './application/notes.facade';
import { NoteTabsService } from './application/note-tabs.service';
import { PresenceService } from './application/presence.service';

export const provideNotes = (): Provider[] => [
  NotesStore,
  NotesFacade,
  NoteTabsService,
  PresenceService,
];

export { NotesFacade };
export { NoteTabsService, type NoteTab } from './application/note-tabs.service';
export { PresenceService };
export { NoteHeaderComponent } from './ui/note-header.component';
export { NoteGridComponent } from './ui/note-grid.component';
export type { Note } from './domain/note.model';
export type { OnlineUser } from './domain/presence.model';
