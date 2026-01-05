import { Injectable, signal } from '@angular/core';
import { User } from './user.model';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private readonly currentUser = signal<User>({
    id: 'user-1',
    name: 'Current User',
    email: 'user@example.com',
    color: '#3b82f6',
  });

  readonly user = this.currentUser.asReadonly();

  getCurrentUserId(): string {
    return this.currentUser().id;
  }
}
