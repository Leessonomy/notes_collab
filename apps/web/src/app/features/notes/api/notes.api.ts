import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Note } from '../domain/note.model';

@Injectable({ providedIn: 'root' })
export class NotesApi {
  private readonly http = inject(HttpClient);

  create(payload: { workspaceId: string; title: string; content: string }) {
    return this.http.post<Note>('/api/notes', payload);
  }

  update(id: string, changes: Partial<Pick<Note, 'title' | 'content'>>) {
    return this.http.patch<Note>(`/api/notes/${id}`, changes);
  }

  delete(id: string) {
    return this.http.delete(`/api/notes/${id}`);
  }
}
