import { Role } from './role.model';

export interface Membership {
  role: Role;
  joinedAt: string;
}

export interface Rights {
  read: boolean;
  write: boolean;
  manage: boolean;
}

export interface Workspace {
  id: string;
  name: string;
  membership: Membership;
  rights: Rights;
  createdAt: Date;
  updatedAt: Date;
}
