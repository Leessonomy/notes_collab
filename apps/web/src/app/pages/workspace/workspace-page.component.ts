import { Component, computed, inject } from '@angular/core';
import { WorkspaceFacade } from '../../features/workspace';
import { NotesFacade, NoteGridComponent, Note } from '../../features/notes';
import { ContentLayoutComponent } from '../../shared/ui/content-layout/content-layout.component';

@Component({
  selector: 'app-workspace-page',
  imports: [ContentLayoutComponent, NoteGridComponent],
  template: `
    @if (currentWorkspace(); as workspace) {
      @let notes = workspaceNotes();
      <div class="h-full w-full overflow-auto bg-background">
        <content-layout>
          <div class="mb-8">
            <h1 class="text-3xl font-bold text-foreground mb-2">
              {{ workspace.name }}
            </h1>
            <p class="text-muted-foreground">
              {{ notes.length }} {{ notes.length === 1 ? 'note' : 'notes' }}
            </p>
          </div>

          <app-note-grid
            [notes]="notes"
            [canCreate]="isAdmin"
            (open)="openNote($event)"
            (create)="createNote()"
          />
        </content-layout>
      </div>
    }
  `,
})
export class WorkspacePageComponent {
  private readonly workspaceFacade = inject(WorkspaceFacade);
  private readonly notesFacade = inject(NotesFacade);

  readonly currentWorkspace = this.workspaceFacade.currentWorkspace;
  readonly isAdmin = true;

  readonly workspaceNotes = computed(() => {
    const workspace = this.currentWorkspace();
    return workspace ? this.notesFacade.byWorkspace(workspace.id) : [];
  });

  createNote() {
    const workspace = this.currentWorkspace();
    if (workspace) this.notesFacade.createAndOpen(workspace.id);
  }

  openNote(note: Note) {
    this.notesFacade.openNote(note);
  }
}
