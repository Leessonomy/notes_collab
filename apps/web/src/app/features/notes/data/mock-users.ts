import { User } from '../../../core/auth/user.model';
import { OnlineUser } from '../domain/presence.model';

export const MOCK_USERS: User[] = [
  {
    id: 'user-1',
    name: 'Alex Johnson',
    email: 'alex@example.com',
  },
  {
    id: 'user-2',
    name: 'Sarah Chen',
    email: 'sarah@example.com',
  },
  {
    id: 'user-3',
    name: 'Mike Wilson',
    email: 'mike@example.com',
  },
];

export const MOCK_ONLINE_USERS: OnlineUser[] = [
  { user: MOCK_USERS[0], noteId: 'note-1' },
  { user: MOCK_USERS[1], noteId: 'note-1' },
];
