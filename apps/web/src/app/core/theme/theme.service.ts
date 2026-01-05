import { Injectable, signal, effect } from '@angular/core';

@Injectable({ providedIn: 'root' })
export class ThemeService {
  private readonly storageKey = 'notes-collab-app-theme';

  private readonly _isDark = signal(
    (() => {
      const theme = localStorage.getItem(this.storageKey);
      switch (theme) {
        case 'dark':
          return true;
        case 'light':
          return false;
        default: {
          return window.matchMedia('(prefers-color-scheme: dark)').matches;
        }
      }
    })()
  );

  constructor() {
    effect(() => {
      document.documentElement.classList.toggle('dark', this._isDark());
    });

    effect(() => {
      localStorage.setItem(this.storageKey, this._isDark() ? 'dark' : 'light');
    });
  }

  get isDarkModeOn() {
    return this._isDark();
  }

  toggle() {
    this._isDark.update((d) => !d);
  }
}
