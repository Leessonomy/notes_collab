import { Injectable, signal } from '@angular/core';
import { User } from '../../../core/auth/user.model';
import { OnlineUser } from '../domain/presence.model';
import { MOCK_USERS, MOCK_ONLINE_USERS } from '../data/mock-users';

@Injectable()
export class PresenceService {
  private readonly onlineUsers = signal<OnlineUser[]>([...MOCK_ONLINE_USERS]);
  private readonly currentUser = signal<User>({ ...MOCK_USERS[0] });

  readonly current = this.currentUser.asReadonly();

  getOnlineUsersForNote() {
    return this.onlineUsers;
  }

  joinNote(noteId: string) {
    const user = this.currentUser();
    const existing = this.onlineUsers().find((u) => u.user.id === user.id && u.noteId === noteId);
    if (!existing) {
      this.onlineUsers.update((users) => [...users, { user, noteId }]);
    }
  }

  leaveNote(noteId: string) {
    const userId = this.currentUser().id;
    this.onlineUsers.update((users) =>
      users.filter((u) => !(u.user.id === userId && u.noteId === noteId)),
    );
  }
}
