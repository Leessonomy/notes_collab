import { inject } from '@angular/core';
import { HttpErrorResponse, HttpInterceptorFn } from '@angular/common/http';
import { catchError, switchMap, throwError } from 'rxjs';
import { AuthService } from './auth.service';

const AUTH_API = ['/api/auth/refresh', '/api/auth/login', '/api/auth/signup', '/api/auth/logout'];

export const authInterceptor: HttpInterceptorFn = (req, next) => {
  const auth = inject(AuthService);

  if (AUTH_API.some((path) => req.url.includes(path))) {
    return next(req);
  }

  return next(req).pipe(
    catchError((error: HttpErrorResponse) => {
      if (error.status !== 401) {
        return throwError(() => error);
      }

      return auth.refreshCurrentSession().pipe(
        switchMap(() => next(req)),
        catchError((sessionError) => {
          auth.setEndSession();
          return throwError(() => sessionError);
        }),
      );
    }),
  );
};
