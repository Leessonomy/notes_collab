import { Component, input, output } from '@angular/core';
import { DatePipe } from '@angular/common';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucidePlus, lucideFileText } from '@ng-icons/lucide';
import { HlmIcon } from '@spartan-ng/helm/icon';
import { Note } from '../domain/note.model';

@Component({
  selector: 'app-note-grid',
  imports: [NgIcon, HlmIcon, DatePipe],
  providers: [provideIcons({ lucidePlus, lucideFileText })],
  template: `
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      @if (canCreate()) {
        <button
          (click)="create.emit()"
          class="group relative h-48 rounded-lg border-2 border-dashed border-border hover:border-primary hover:bg-accent/50 transition-all duration-200 flex flex-col items-center justify-center gap-3 cursor-pointer"
        >
          <div
            class="w-12 h-12 rounded-full bg-primary/10 group-hover:bg-primary/20 flex items-center justify-center transition-colors"
          >
            <ng-icon hlm name="lucidePlus" size="lg" class="text-primary" />
          </div>
          <span
            class="text-sm font-medium text-muted-foreground group-hover:text-foreground transition-colors"
          >
            New Note
          </span>
        </button>
      }
      @for (note of notes(); track note.id) {
        <button
          (click)="open.emit(note)"
          class="group relative h-48 rounded-lg border border-border hover:border-primary hover:shadow-md bg-card transition-all duration-200 p-4 flex flex-col cursor-pointer text-left overflow-hidden"
        >
          <div class="flex items-start gap-3 mb-3">
            <div
              class="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center flex-shrink-0"
            >
              <ng-icon hlm name="lucideFileText" size="sm" class="text-primary" />
            </div>
            <div class="flex-1 min-w-0">
              <h3
                class="font-semibold text-foreground truncate group-hover:text-primary transition-colors"
              >
                {{ note.title || 'Untitled' }}
              </h3>
              <p class="text-xs text-muted-foreground mt-1">
                {{ note.updatedAt | date }}
              </p>
            </div>
          </div>
          <div class="flex-1 overflow-hidden">
            <p class="text-sm text-muted-foreground line-clamp-4">
              {{ getPreview(note.content) }}
            </p>
          </div>
        </button>
      }
    </div>

    @if (!notes().length) {
      <div class="text-center py-16">
        <div class="w-16 h-16 rounded-full bg-muted mx-auto mb-4 flex items-center justify-center">
          <ng-icon hlm name="lucideFileText" size="xl" class="text-muted-foreground" />
        </div>
        <h3 class="text-lg font-medium text-foreground mb-2">No notes yet</h3>
        <p class="text-sm text-muted-foreground">This workspace doesn't have any notes yet</p>
      </div>
    }
  `,
})
export class NoteGridComponent {
  readonly notes = input.required<Note[]>();
  readonly canCreate = input(false);

  readonly open = output<Note>();
  readonly create = output<void>();

  getPreview(content: string): string {
    const text = content
      .replace(/<[^>]*>/g, '')
      .replace(/\n/g, ' ')
      .trim();
    return text || 'No content';
  }
}
