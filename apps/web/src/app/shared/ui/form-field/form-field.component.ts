import { Component, input } from '@angular/core';
import { Field, type FieldTree } from '@angular/forms/signals';
import { HlmInputImports } from '@spartan-ng/helm/input';
import { HlmLabelImports } from '@spartan-ng/helm/label';

@Component({
  selector: 'form-field',
  imports: [Field, HlmInputImports, HlmLabelImports],
  template: `
    @let fieldState = field()();
    <div class="grid gap-2">
      <div class="flex items-center">
        <label hlmLabel [for]="id()">{{ label() }}</label>
        <ng-content select="[labelAction]" />
      </div>
      <input [id]="id()" [type]="type()" [placeholder]="placeholder()" hlmInput [field]="field()" />
      @if (fieldState.touched()) {
        @for (error of fieldState.errors(); track error.kind) {
          <p class="text-destructive text-sm">{{ error.message }}</p>
        }
      }
    </div>
  `,
})
export class FormFieldComponent<T extends string | number | boolean = string> {
  readonly field = input.required<FieldTree<T>>();
  readonly label = input.required<string>();
  readonly id = input.required<string>();
  readonly type = input<'text' | 'email' | 'password'>('text');
  readonly placeholder = input('');
}
