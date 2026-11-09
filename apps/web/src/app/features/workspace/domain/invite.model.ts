import { Role } from './role.model';

export interface Invite {
  id: string;
  role: Role;
  expiresAt: string;
  createdAt: string;
}

export interface CreatedInvite extends Invite {
  token: string;
}
