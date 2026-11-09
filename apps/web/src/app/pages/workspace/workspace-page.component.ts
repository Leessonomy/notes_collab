import { Component, computed, effect, inject, input } from '@angular/core';
import { MembersSheetComponent, WorkspaceFacade } from '../../features/workspace';
import { NotesFacade, NoteGridComponent, Note } from '../../features/notes';
import { ContentLayoutComponent } from '../../shared/ui/content-layout/content-layout.component';

@Component({
  selector: 'app-workspace-page',
  imports: [ContentLayoutComponent, NoteGridComponent, MembersSheetComponent],
  template: `
    @if (currentWorkspace(); as workspace) {
      @let notes = workspaceNotes();
      <div class="h-full w-full overflow-auto bg-background">
        <content-layout>
          <div class="mb-8 flex items-start justify-between gap-4">
            <div>
              <h1 class="text-3xl font-bold text-foreground mb-2">
                {{ workspace.name }}
              </h1>
              <p class="text-muted-foreground">
                {{ notes.length }} {{ notes.length === 1 ? 'note' : 'notes' }}
              </p>
            </div>
            <app-members-sheet [workspace]="workspace" />
          </div>

          <app-note-grid
            [notes]="notes"
            [readonly]="!workspace.rights.write"
            (open)="openNote($event)"
            (create)="createNote()"
            (delete)="deleteNote($event)"
          />
        </content-layout>
      </div>
    }
  `,
})
export class WorkspacePageComponent {
  private readonly workspaceFacade = inject(WorkspaceFacade);
  private readonly notesFacade = inject(NotesFacade);

  readonly selectedWorkspaceId = input<string | undefined>(undefined, { alias: 'workspace' });

  readonly currentWorkspace = this.workspaceFacade.currentWorkspace;

  readonly workspaceNotes = computed(() => {
    const workspace = this.currentWorkspace();
    return workspace ? this.notesFacade.byWorkspace(workspace.id) : [];
  });

  constructor() {
    effect(() => {
      const workspaceId = this.selectedWorkspaceId();
      if (workspaceId) this.workspaceFacade.switch(workspaceId);
    });
  }

  createNote() {
    const workspace = this.currentWorkspace();
    if (workspace) this.notesFacade.createAndOpen(workspace.id);
  }

  openNote(note: Note) {
    this.notesFacade.openNote(note.id);
  }

  deleteNote(note: Note) {
    this.notesFacade.deleteNote(note.id);
  }
}
