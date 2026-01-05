import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Note } from '../domain/note.model';

@Injectable({ providedIn: 'root' })
export class NotesApiService {
  private readonly http = inject(HttpClient);

  getByWorkspace(workspaceId: string) {
    return this.http.get<Note[]>(`/api/workspaces/${workspaceId}/notes`);
  }

  create(payload: { workspaceId: string; title: string; content: string; authorId: string }) {
    return this.http.post<Note>('/api/notes', payload);
  }

  update(id: string, changes: Partial<Pick<Note, 'title' | 'content'>>) {
    return this.http.patch<Note>(`/api/notes/${id}`, changes);
  }
}
