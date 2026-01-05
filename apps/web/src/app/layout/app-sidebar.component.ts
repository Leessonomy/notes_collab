import { Component, inject, signal } from '@angular/core';
import { HlmSidebarImports } from '@spartan-ng/helm/sidebar';
import { HlmIcon } from '@spartan-ng/helm/icon';
import { HlmSeparator } from '@spartan-ng/helm/separator';
import { HlmTooltipImports } from '@spartan-ng/helm/tooltip';
import { BrnCollapsibleImports } from '@spartan-ng/brain/collapsible';
import { HlmCollapsibleImports } from '@spartan-ng/helm/collapsible';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { NotesFacade, NoteTabsService } from '../features/notes';
import { WorkspaceFacade, WorkspaceInviteService } from '../features/workspace';
import {
  lucidePlus,
  lucideFileText,
  lucideChevronDown,
  lucideChevronRight,
  lucideSettings,
  lucideLink,
  lucideCheck,
  lucideUsers,
  lucideFolderOpen,
  lucideUserPlus,
  lucidePlusCircle,
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
    HlmSeparator,
  ],
  providers: [
    provideIcons({
      lucidePlus,
      lucideFileText,
      lucideChevronDown,
      lucideChevronRight,
      lucideSettings,
      lucideLink,
      lucideCheck,
      lucideUsers,
      lucideFolderOpen,
      lucideUserPlus,
      lucidePlusCircle,
    }),
  ],
  template: `
    <div hlmSidebar side="left" collapsible="icon" class="w-2xs">
      <div hlmSidebarHeader>
        <div hlmSidebarGroup>
          <div hlmSidebarGroupContent>
            <ul hlmSidebarMenu>
              <hlm-collapsible [expanded]="true" class="group/collapsible">
                @let currentWorkspace = this.currentWorkspace();
                <li hlmSidebarMenuItem>
                  <div class="flex items-center gap-1">
                    <button
                      hlmCollapsibleTrigger
                      hlmSidebarMenuButton
                      class="flex w-full items-center justify-between"
                    >
                      <span class="flex items-center gap-2 truncate">
                        <ng-icon hlm name="lucideFolderOpen" size="sm" />
                        <span class="font-medium truncate group-data-[collapsible=icon]:hidden">
                          {{ currentWorkspace?.name ?? 'Select Workspace' }}
                        </span>
                      </span>
                      <ng-icon
                        hlm
                        name="lucideChevronRight"
                        size="sm"
                        class="ml-auto transition-transform group-data-[state=open]/collapsible:rotate-90"
                      />
                    </button>
                    @if (isAdmin) {
                      <button
                        hlmSidebarMenuButton
                        size="sm"
                        class="h-8 w-8 p-0 flex justify-center"
                        (click)="createWorkspace()"
                      >
                        <ng-icon hlm name="lucidePlusCircle" size="sm" />
                      </button>
                    }
                  </div>
                  <hlm-collapsible-content>
                    <ul hlmSidebarMenuSub>
                      @for (workspace of workspaces(); track workspace.id) {
                        <li hlmSidebarMenuSubItem>
                          <button
                            hlmSidebarMenuSubButton
                            class="w-full"
                            [class.bg-accent]="workspace.id === currentWorkspace?.id"
                            (click)="selectWorkspace(workspace.id)"
                          >
                            <span class="truncate">{{ workspace.name }}</span>
                            @if (workspace.id === currentWorkspace?.id) {
                              <ng-icon
                                hlm
                                name="lucideCheck"
                                size="xs"
                                class="ml-auto text-primary"
                              />
                            }
                          </button>
                        </li>
                      }
                      @if (isAdmin) {
                        <li hlmSidebarMenuSubItem>
                          <hr hlmSeparator class="my-1" />
                        </li>
                        <li hlmSidebarMenuSubItem>
                          <button
                            hlmSidebarMenuSubButton
                            class="w-full text-muted-foreground"
                            (click)="copyWorkspaceInviteLink()"
                          >
                            <ng-icon hlm name="lucideUserPlus" size="xs" />
                            <span>Invite to Workspace</span>
                          </button>
                        </li>
                      }
                    </ul>
                  </hlm-collapsible-content>
                </li>
              </hlm-collapsible>
              <hlm-collapsible [expanded]="notesOpen()" class="group/collapsible">
                <li hlmSidebarMenuItem>
                  <div class="flex items-center gap-1">
                    <button hlmCollapsibleTrigger hlmSidebarMenuButton class="flex-1">
                      <ng-icon hlm name="lucideFileText" size="sm" />
                      <span>Notes</span>
                      <ng-icon
                        hlm
                        name="lucideChevronRight"
                        size="sm"
                        class="ml-auto transition-transform group-data-[state=open]/collapsible:rotate-90"
                      />
                    </button>
                    @if (isAdmin) {
                      <button
                        hlmSidebarMenuButton
                        size="sm"
                        class="h-8 w-8 p-0 flex justify-center"
                        (click)="createNoteFromSidebar()"
                      >
                        <ng-icon hlm name="lucidePlusCircle" size="sm" />
                      </button>
                    }
                  </div>
                  <hlm-collapsible-content>
                    <ul hlmSidebarMenuSub>
                      @for (note of workspaceNotes(); track note.id) {
                        <li hlmSidebarMenuSubItem>
                          <button
                            hlmSidebarMenuSubButton
                            class="w-full h-auto min-h-[2rem] py-1.5"
                            [class.bg-accent]="note.id === currentNote()?.id"
                            (click)="selectNote(note.id)"
                          >
                            <span class="whitespace-normal break-words text-left leading-tight">{{
                              note.title || 'Untitled'
                            }}</span>
                          </button>
                        </li>
                      }
                      @if (!workspaceNotes().length) {
                        <li hlmSidebarMenuSubItem>
                          <div class="px-2 py-4 text-center text-sm text-muted-foreground">
                            No notes yet
                          </div>
                        </li>
                      }
                    </ul>
                  </hlm-collapsible-content>
                </li>
              </hlm-collapsible>
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
  private readonly workspaceFacade = inject(WorkspaceFacade);
  private readonly notesFacade = inject(NotesFacade);
  private readonly noteTabs = inject(NoteTabsService);
  private readonly inviteService = inject(WorkspaceInviteService);

  readonly workspaces = this.workspaceFacade.workspaces;
  readonly currentWorkspace = this.workspaceFacade.currentWorkspace;
  readonly currentNote = this.noteTabs.currentNote;
  readonly workspaceNotes = this.notesFacade.notes;

  readonly isAdmin = true;
  readonly notesOpen = signal(true);

  selectWorkspace(id: string) {
    this.workspaceFacade.switch(id);
    this.noteTabs.closeAll();
  }

  createWorkspace() {
    const name = prompt('Workspace name:');
    if (name) {
      this.workspaceFacade.create(name).subscribe(() => this.noteTabs.closeAll());
    }
  }

  selectNote(id: string) {
    const note = this.workspaceNotes().find((n) => n.id === id);
    if (note) this.notesFacade.openNote(note);
  }

  createNoteFromSidebar() {
    const ws = this.currentWorkspace();
    if (ws) this.notesFacade.createAndOpen(ws.id);
  }

  copyWorkspaceInviteLink() {
    const ws = this.currentWorkspace();
    if (!ws) return;
    const link = this.inviteService.generateInviteLink(ws.id);
    navigator.clipboard.writeText(link);
  }
}
