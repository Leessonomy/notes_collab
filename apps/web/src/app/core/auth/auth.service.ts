import { computed, inject, Injectable, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { catchError, of, tap } from 'rxjs';
import { User } from './user.model';

interface LoginCredentials {
  email: string;
  password: string;
}

interface SignUpCredentials {
  name: string;
  email: string;
  password: string;
}

interface AuthState {
  me: User | null;
}

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private readonly http = inject(HttpClient);

  private readonly state = signal<AuthState>({
    me: null,
  });

  readonly me = computed(() => this.state().me);
  readonly isAuthenticated = computed(() => this.state().me !== null);

  getMe() {
    return this.http.get<User>('/api/auth/me').pipe(tap((user) => this.state.set({ me: user })));
  }

  login(credentials: LoginCredentials) {
    return this.http
      .post<User>('/api/auth/login', credentials)
      .pipe(tap((user) => this.state.set({ me: user })));
  }

  signup(credentials: SignUpCredentials) {
    return this.http
      .post<User>('/api/auth/signup', credentials)
      .pipe(tap((user) => this.state.set({ me: user })));
  }

  logout() {
    return this.http.post('/api/auth/logout', {}).pipe(tap(() => this.state.set({ me: null })));
  }
}
