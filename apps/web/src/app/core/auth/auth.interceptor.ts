import { inject } from '@angular/core';
import { HttpErrorResponse, HttpInterceptorFn } from '@angular/common/http';
import { Router } from '@angular/router';
import { catchError, finalize, Observable, shareReplay, switchMap, throwError } from 'rxjs';
import { AuthService } from './auth.service';

const AUTH_API = ['/api/auth/refresh', '/api/auth/login', '/api/auth/signup', '/api/auth/logout'];

let refreshRequest$: Observable<any> | null = null;

const refreshSession = (auth: AuthService) => {
  if (!refreshRequest$) {
    refreshRequest$ = auth.refreshCurrentSession().pipe(
      shareReplay(1),
      finalize(() => {
        refreshRequest$ = null;
      }),
    );
  }

  return refreshRequest$;
};

export const authInterceptor: HttpInterceptorFn = (req, next) => {
  const auth = inject(AuthService);
  const router = inject(Router);

  if (AUTH_API.some((path) => req.url.includes(path))) {
    return next(req);
  }

  return next(req).pipe(
    catchError((error: HttpErrorResponse) => {
      const isUnauthorized = error.status === 401;
      if (isUnauthorized) {
        return refreshSession(auth).pipe(
          switchMap(() => next(req)),
          catchError((error) => {
            auth.clearCurrentSession();
            router.navigateByUrl('/auth');
            return throwError(() => error);
          }),
        );
      } else {
        return throwError(() => error);
      }
    }),
  );
};
