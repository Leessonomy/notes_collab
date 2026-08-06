import { computed, inject, Injectable, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { catchError, finalize, Observable, shareReplay, of, tap } from 'rxjs';
import { User } from './user.model';
import { Router } from '@angular/router';

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
  private readonly router = inject(Router);

  private readonly state = signal<AuthState>({ me: null });
  readonly me = computed(() => this.state().me);
  readonly isAuthenticated = computed(() => this.state().me !== null);

  private refreshInFlight: Observable<unknown> | null = null;

  getMe() {
    return this.http.get<User>('/api/auth/me').pipe(
      tap((user) => this.state.set({ me: user })),
      catchError(() => {
        this.setEndSession();
        return of(null);
      }),
    );
  }

  login(credentials: LoginCredentials) {
    return this.http
      .post<User>('/api/auth/login', credentials)
      .pipe(tap((user) => this.state.set({ me: user })));
  }

  logout() {
    return this.http.post('/api/auth/logout', {}).pipe(tap(() => this.setEndSession()));
  }

  signup(credentials: SignUpCredentials) {
    return this.http
      .post<User>('/api/auth/signup', credentials)
      .pipe(tap((user) => this.state.set({ me: user })));
  }

  refreshCurrentSession() {
    this.refreshInFlight ??= this.http.post('/api/auth/refresh', {}).pipe(
      finalize(() => (this.refreshInFlight = null)),
      shareReplay({ bufferSize: 1, refCount: false }),
    );

    return this.refreshInFlight;
  }

  setEndSession() {
    this.state.set({ me: null });
    this.router.navigateByUrl('/auth');
  }
}
