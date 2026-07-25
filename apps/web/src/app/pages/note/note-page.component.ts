import { Component, inject, effect } from '@angular/core';
import { TiptapEditorComponent } from '../../shared/ui/editor/tiptap-editor.component';
// eslint-disable-next-line boundaries/element-types
import { ContentLayoutComponent } from '../../layout/content-layout.component';
import {
  NotesFacade,
  NoteTabsService,
  PresenceService,
  NoteHeaderComponent,
} from '../../features/notes';

@Component({
  selector: 'app-note-page',
  imports: [TiptapEditorComponent, NoteHeaderComponent, ContentLayoutComponent],
  template: `
    <div class="flex flex-col h-full w-full bg-background">
      @if (currentNote(); as note) {
        <app-note-header [note]="note" [onlineUsers]="onlineUsers()" />
        <div class="flex-1 overflow-auto w-full">
          <content-layout>
            <app-tiptap-editor [content]="note.content" (contentChange)="onContentChange($event)" />
          </content-layout>
        </div>
      } @else {
        <div class="flex-1 flex items-center justify-center w-full">
          <content-layout>
            <div class="text-center">
              <h3 class="text-lg font-medium text-foreground mb-2">No note selected</h3>
              <p class="text-sm text-muted-foreground">
                Select a note from the sidebar or create a new one
              </p>
            </div>
          </content-layout>
        </div>
      }
    </div>
  `,
})
export class NotePageComponent {
  private readonly notesFacade = inject(NotesFacade);
  private readonly noteTabs = inject(NoteTabsService);
  private readonly presenceService = inject(PresenceService);

  readonly currentNote = this.noteTabs.currentNote;
  readonly onlineUsers = this.presenceService.getOnlineUsersForNote();

  constructor() {
    effect(() => {
      const note = this.currentNote();
      if (note) {
        this.presenceService.joinNote(note.id);
      }
    });
  }

  onContentChange(content: string) {
    const note = this.currentNote();
    if (note) this.notesFacade.saveContent(note.id, content);
  }
}
