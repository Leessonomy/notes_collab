import { Component, OnInit, effect, inject } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { HlmSidebarImports, provideHlmSidebarConfig } from '@spartan-ng/helm/sidebar';
import { AppHeaderComponent } from './app-header.component';
import { AppSidebarComponent } from './app-sidebar.component';
import { WorkspaceFacade } from '../features/workspace';
import { NotesFacade } from '../features/notes';

@Component({
  selector: 'app-main-layout',
  imports: [...HlmSidebarImports, AppHeaderComponent, AppSidebarComponent, RouterOutlet],
  providers: [
    provideHlmSidebarConfig({
      mobileBreakpoint: '768px',
      sidebarWidth: '320px',
      sidebarWidthIcon: '80px',
    }),
  ],
  template: `
    <div hlmSidebarWrapper class="min-h-screen">
      <app-sidebar />
      <div hlmSidebarInset class="flex flex-col w-full">
        <app-header />
        <main class="flex-1 overflow-hidden w-full">
          <router-outlet />
        </main>
      </div>
    </div>
  `,
})
export class MainLayoutComponent implements OnInit {
  private readonly workspaceFacade = inject(WorkspaceFacade);
  private readonly notesFacade = inject(NotesFacade);

  private lastLoadedWorkspaceId: string | null = null;

  constructor() {
    effect(() => {
      const id = this.workspaceFacade.currentWorkspace()?.id ?? null;
      if (id && id !== this.lastLoadedWorkspaceId) {
        this.lastLoadedWorkspaceId = id;
        this.notesFacade.loadForWorkspace(id);
      }
    });
  }

  ngOnInit() {
    this.workspaceFacade.load();
  }
}
