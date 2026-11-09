import { Component, inject, signal } from '@angular/core';
import { Router } from '@angular/router';
import { HlmSidebarImports } from '@spartan-ng/helm/sidebar';
import { HlmIcon } from '@spartan-ng/helm/icon';
import { HlmTooltipImports } from '@spartan-ng/helm/tooltip';
import { BrnCollapsibleImports } from '@spartan-ng/brain/collapsible';
import { HlmCollapsibleImports } from '@spartan-ng/helm/collapsible';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { NotesFacade, NoteTabsService } from '../features/notes';
import { Workspace, WorkspaceFacade } from '../features/workspace';
import { AppFacade } from '../app.facade';
import { appUrls } from '../core/app-urls';
import {
  lucidePlus,
  lucideFileText,
  lucideChevronRight,
  lucideSettings,
  lucideFolderOpen,
  lucideDoorOpen,
  lucideTrash2,
  lucideX,
} from '@ng-icons/lucide';

@Component({
  selector: 'app-sidebar',
  imports: [
    NgIcon,
    HlmIcon,
    HlmSidebarImports,
    HlmTooltipImports,
    BrnCollapsibleImports,
    HlmCollapsibleImports,
  ],
  providers: [
    provideIcons({
      lucidePlus,
      lucideFileText,
      lucideChevronRight,
      lucideSettings,
      lucideFolderOpen,
      lucideDoorOpen,
      lucideTrash2,
      lucideX,
    }),
  ],
  template: `
    <div hlmSidebar side="left" collapsible="icon" class="w-2xs">
      <div hlmSidebarHeader>
        <div hlmSidebarGroup>
          <div
            hlmSidebarGroupLabel
            class="text-[13px] font-semibold uppercase tracking-wider text-sidebar-foreground/60"
          >
            <span>Workspaces</span>
            <span
              class="ml-3 inline-flex items-center leading-none text-[12px] font-bold tabular-nums text-sidebar-foreground"
            >
              {{ workspaces().length }}
            </span>
          </div>

          <button hlmSidebarGroupAction (click)="createWorkspace()">
            <ng-icon hlm name="lucidePlus" size="sm" />
          </button>

          <div hlmSidebarGroupContent>
            <ul hlmSidebarMenu>
              @for (workspace of workspaces(); track workspace.id) {
                <hlm-collapsible
                  [expanded]="isOpen(workspace.id)"
                  (expandedChange)="setOpen(workspace.id, $event)"
                  class="group/collapsible"
                >
                  <li hlmSidebarMenuItem>
                    <button
                      hlmCollapsibleTrigger
                      hlmSidebarMenuButton
                      class="group-has-[[data-sidebar=menu-action]]/menu-item:pr-16"
                      [isActive]="workspace.id === currentWorkspace()?.id"
                      (click)="selectWorkspace(workspace.id)"
                    >
                      <ng-icon hlm name="lucideFolderOpen" size="sm" />
                      <span>{{ workspace.name }}</span>
                      <ng-icon
                        hlm
                        name="lucideChevronRight"
                        size="sm"
                        class="ml-auto transition-transform group-data-[state=open]/collapsible:rotate-90"
                      />
                    </button>

                    @if (workspace.rights.write) {
                      <button
                        hlmSidebarMenuAction
                        class="right-8 !top-1 !w-6"
                        [attr.aria-label]="'New note in ' + workspace.name"
                        (click)="createNote(workspace.id)"
                      >
                        <ng-icon hlm name="lucidePlus" size="sm" />
                      </button>
                    }

                    @if (workspace.rights.manage) {
                      <button
                        hlmSidebarMenuAction
                        class="hover:text-destructive !top-1 !w-6"
                        [attr.aria-label]="'Delete ' + workspace.name"
                        (click)="deleteWorkspace(workspace.id)"
                      >
                        <ng-icon hlm name="lucideTrash2" size="sm" />
                      </button>
                    } @else {
                      <button
                        hlmSidebarMenuAction
                        class="hover:text-destructive !top-1 !w-6"
                        [attr.aria-label]="'Leave ' + workspace.name"
                        (click)="leaveWorkspace(workspace)"
                      >
                        <ng-icon hlm name="lucideDoorOpen" size="sm" />
                      </button>
                    }

                    <hlm-collapsible-content>
                      <ul hlmSidebarMenuSub>
                        @for (note of notesFacade.byWorkspace(workspace.id); track note.id) {
                          <li hlmSidebarMenuSubItem class="flex items-center gap-0.5">
                            <button
                              hlmSidebarMenuSubButton
                              [isActive]="note.id === currentNote()?.id"
                              (click)="notesFacade.openNote(note.id)"
                            >
                              <ng-icon hlm name="lucideFileText" size="xs" />
                              <span>{{ note.title || 'Untitled' }}</span>
                            </button>

                            @if (workspace.rights.write) {
                              <button
                                type="button"
                                [attr.aria-label]="'Delete ' + (note.title || 'Untitled')"
                                (click)="deleteNote(note.id)"
                                class="text-sidebar-foreground/60 ring-sidebar-ring hover:bg-sidebar-accent hover:text-destructive flex aspect-square w-6 shrink-0 items-center justify-center rounded-md outline-none hover:cursor-pointer focus-visible:ring-2"
                              >
                                <ng-icon hlm name="lucideX" size="sm" />
                              </button>
                            }
                          </li>
                        } @empty {
                          <li hlmSidebarMenuSubItem>
                            <div class="px-2 py-3 text-center text-sm text-muted-foreground">
                              No notes yet
                            </div>
                          </li>
                        }
                      </ul>
                    </hlm-collapsible-content>
                  </li>
                </hlm-collapsible>
              } @empty {
                <li hlmSidebarMenuItem>
                  <div class="px-2 py-4 text-center text-sm text-muted-foreground">
                    @if (isLoading()) {
                      Loading workspaces…
                    } @else {
                      No workspaces yet
                    }
                  </div>
                </li>
              }
            </ul>
          </div>
        </div>
      </div>

      <div hlmSidebarFooter>
        <div hlmSidebarGroup>
          <ul hlmSidebarMenu>
            <li hlmSidebarMenuItem>
              <button hlmSidebarMenuButton>
                <ng-icon hlm name="lucideSettings" size="sm" />
                <span>Settings</span>
              </button>
            </li>
          </ul>
        </div>
      </div>
    </div>
  `,
})
export class AppSidebarComponent {
  private readonly router = inject(Router);
  private readonly workspaceFacade = inject(WorkspaceFacade);
  private readonly noteTabs = inject(NoteTabsService);
  private readonly appFacade = inject(AppFacade);

  protected readonly notesFacade = inject(NotesFacade);

  private readonly openWorkspaces = signal<Record<string, boolean>>({});

  readonly workspaces = this.workspaceFacade.workspaces;
  readonly currentWorkspace = this.workspaceFacade.currentWorkspace;
  readonly currentNote = this.noteTabs.currentNote;
  readonly isLoading = this.appFacade.isLoading;

  isOpen(workspaceId: string) {
    return this.openWorkspaces()[workspaceId] ?? workspaceId === this.currentWorkspace()?.id;
  }

  setOpen(workspaceId: string, open: boolean) {
    this.openWorkspaces.update((state) => ({ ...state, [workspaceId]: open }));
  }

  selectWorkspace(id: string) {
    this.workspaceFacade.switch(id);
    this.router.navigateByUrl(appUrls.workspace);
  }

  createWorkspace() {
    const name = prompt('Workspace name:');
    if (name) {
      this.workspaceFacade.create(name).subscribe((workspace) => {
        this.setOpen(workspace.id, true);
        this.router.navigateByUrl(appUrls.workspace);
      });
    }
  }

  deleteNote(noteId: string) {
    this.notesFacade.deleteNote(noteId);
  }

  deleteWorkspace(workspaceId: string) {
    const losesOpenNote = this.currentNote()?.workspaceId === workspaceId;

    this.workspaceFacade.delete(workspaceId).subscribe(() => {
      this.appFacade.load();
      if (losesOpenNote) this.router.navigateByUrl(appUrls.workspace);
    });
  }

  leaveWorkspace(workspace: Workspace) {
    if (!confirm(`Leave ${workspace.name}? You will need a new invite to come back.`)) return;

    const losesOpenNote = this.currentNote()?.workspaceId === workspace.id;

    this.workspaceFacade.leave(workspace.id).subscribe(() => {
      this.appFacade.load();
      if (losesOpenNote) this.router.navigateByUrl(appUrls.workspace);
    });
  }

  createNote(workspaceId: string) {
    this.setOpen(workspaceId, true);
    this.notesFacade.createAndOpen(workspaceId);
  }
}
