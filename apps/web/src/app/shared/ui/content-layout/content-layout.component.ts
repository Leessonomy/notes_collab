import { Component } from '@angular/core';

@Component({
  selector: 'content-layout',
  template: `
    <div class="w-full max-w-7xl mx-auto p-8">
      <ng-content />
    </div>
  `,
})
export class ContentLayoutComponent {}
