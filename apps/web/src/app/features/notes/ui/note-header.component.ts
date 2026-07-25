import { Component, input } from '@angular/core';
import { Note } from '../domain/note.model';
import { OnlineUser } from '../domain/presence.model';

@Component({
  selector: 'app-note-header',
  template: `
    <div
      class="flex items-center justify-between h-10 px-4 border-b border-border bg-background/80 backdrop-blur-sm"
    >
      <div class="flex items-center gap-3 min-w-0">
        <h2 class="text-sm font-medium text-foreground truncate">
          {{ note()?.title || 'Untitled' }}
        </h2>
        @if (note()) {
          <span class="text-xs text-muted-foreground hidden sm:inline">
            Edited {{ formatDate(note()!.updatedAt) }}
          </span>
        }
      </div>

      <div class="flex items-center gap-1">
        @for (user of onlineUsers(); track user.user.id) {
          <div
            class="w-6 h-6 rounded-full flex items-center justify-center text-xs font-medium text-white ring-2 ring-background -ml-1 first:ml-0"
            [title]="user.user.name"
          >
            {{ getInitials(user.user.name) }}
          </div>
        }
        @if (onlineUsers().length > 0) {
          <span class="text-xs text-muted-foreground ml-2">
            {{ onlineUsers().length }} online
          </span>
        }
      </div>
    </div>
  `,
})
export class NoteHeaderComponent {
  readonly note = input<Note | null>(null);
  readonly onlineUsers = input<OnlineUser[]>([]);

  getInitials(name: string): string {
    return name
      .split(' ')
      .map((n) => n[0])
      .join('')
      .toUpperCase()
      .slice(0, 2);
  }

  formatDate(date: Date | string): string {
    const d = new Date(date);
    const now = new Date();
    const diff = now.getTime() - d.getTime();
    const minutes = Math.floor(diff / 60000);
    const hours = Math.floor(diff / 3600000);
    const days = Math.floor(diff / 86400000);

    if (minutes < 1) return 'just now';
    if (minutes < 60) return `${minutes}m ago`;
    if (hours < 24) return `${hours}h ago`;
    if (days < 7) return `${days}d ago`;
    return d.toLocaleDateString();
  }
}
