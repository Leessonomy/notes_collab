import { Component, inject } from '@angular/core';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideMoon, lucideSun, lucideMenu, lucideX, lucideLogOut } from '@ng-icons/lucide';
import { HlmButton } from '@spartan-ng/helm/button';
import { HlmIcon } from '@spartan-ng/helm/icon';
import { ThemeService } from '../core/theme/theme.service';
import { HlmSidebarService } from '@spartan-ng/helm/sidebar';
import { NoteTabsService, NoteTab } from '../features/notes';
import { AuthService } from '../core/auth/auth.service';

@Component({
  selector: 'app-header',
  imports: [HlmButton, NgIcon, HlmIcon],
  providers: [provideIcons({ lucideMoon, lucideSun, lucideLogOut, lucideMenu, lucideX })],
  template: `
    <header
      class="h-[52px] flex items-center border-b border-border bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60"
    >
      <div class="flex items-center gap-2 px-4 border-r border-border">
        <button hlmBtn variant="ghost" size="icon" class="md:hidden" (click)="toggleSidebar()">
          <ng-icon hlm name="lucideMenu" size="sm" />
        </button>
      </div>

      <div class="flex-1 flex items-center overflow-x-auto">
        <div class="flex items-center h-full">
          @for (tab of tabs(); track tab.noteId) {
            <button
              (click)="switchTab(tab)"
              class="group relative h-full px-4 flex items-center gap-2 border-r border-border hover:bg-accent/50 transition-colors min-w-[120px] max-w-[200px]"
              [class.bg-accent]="activeTab()?.noteId === tab.noteId"
              [class.border-b-2]="activeTab()?.noteId === tab.noteId"
              [class.border-b-primary]="activeTab()?.noteId === tab.noteId"
            >
              <span class="text-sm truncate flex-1">{{ tab.title }}</span>
              <button
                (click)="closeTab(tab.noteId, $event)"
                class="opacity-0 group-hover:opacity-100 hover:bg-muted rounded p-0.5 transition-opacity"
              >
                <ng-icon hlm name="lucideX" size="xs" />
              </button>
            </button>
          }
        </div>
      </div>

      <div class="flex items-center gap-2 px-4 border-l border-border">
        <button hlmBtn variant="ghost" size="icon" (click)="toggleTheme()" class="rounded-full">
          @if (isDarkModeOn) {
            <ng-icon hlm name="lucideSun" size="sm" />
          } @else {
            <ng-icon hlm name="lucideMoon" size="sm" />
          }
        </button>
        <button hlmBtn variant="ghost" size="icon" (click)="logout()" class="rounded-full">
          <ng-icon hlm name="lucideLogOut" size="sm" />
        </button>
      </div>
    </header>
  `,
})
export class AppHeaderComponent {
  private readonly themeService = inject(ThemeService);
  private readonly authService = inject(AuthService);
  private readonly sidebarService = inject(HlmSidebarService);
  private readonly noteTabs = inject(NoteTabsService);

  readonly tabs = this.noteTabs.tabs;
  readonly isDarkModeOn = this.themeService.isDarkModeOn;
  readonly activeTab = this.noteTabs.activeTab;

  logout() {
    this.authService.logout().subscribe();
  }

  toggleTheme() {
    this.themeService.toggle();
  }

  toggleSidebar() {
    this.sidebarService.toggleSidebar();
  }

  switchTab(tab: NoteTab) {
    this.noteTabs.openNote({ id: tab.noteId, title: tab.title });
  }

  closeTab(noteId: string, event: Event) {
    event.stopPropagation();
    this.noteTabs.closeTab(noteId);
  }
}
