import { User } from '../../../core/auth/user.model';

export interface OnlineUser {
  user: User;
  noteId: string;
  cursorPosition?: { from: number; to: number };
}
