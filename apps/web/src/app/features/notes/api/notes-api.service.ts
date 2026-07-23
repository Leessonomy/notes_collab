import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Note } from '../domain/note.model';

@Injectable({ providedIn: 'root' })
export class NotesApiService {
  private readonly http = inject(HttpClient);

  get() {
    return this.http.get<Note[]>('/api/notes');
  }

  create(payload: { workspaceId: string; title: string; content: string }) {
    return this.http.post<Note>('/api/notes', payload);
  }

  update(id: string, changes: Partial<Pick<Note, 'title' | 'content'>>) {
    return this.http.patch<Note>(`/api/notes/${id}`, changes);
  }
}
