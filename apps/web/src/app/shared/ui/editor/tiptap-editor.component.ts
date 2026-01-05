import {
  Component,
  ElementRef,
  OnDestroy,
  OnInit,
  effect,
  input,
  output,
  viewChild,
} from '@angular/core';
import { NgIcon, provideIcons } from '@ng-icons/core';
import {
  lucideBold,
  lucideItalic,
  lucideUnderline,
  lucideStrikethrough,
  lucideCode,
  lucideHeading1,
  lucideHeading2,
  lucideHeading3,
  lucideList,
  lucideListOrdered,
  lucideListTodo,
  lucideQuote,
  lucideCode2,
  lucideMinus,
  lucideUndo,
  lucideRedo,
  lucideLink,
  lucideImage,
  lucideTable,
  lucideHighlighter,
} from '@ng-icons/lucide';
import { HlmButton } from '@spartan-ng/helm/button';
import { HlmIcon } from '@spartan-ng/helm/icon';
import { HlmSeparator } from '@spartan-ng/helm/separator';
import { Editor } from '@tiptap/core';
import StarterKit from '@tiptap/starter-kit';
import Placeholder from '@tiptap/extension-placeholder';
import Highlight from '@tiptap/extension-highlight';
import TaskList from '@tiptap/extension-task-list';
import TaskItem from '@tiptap/extension-task-item';
import Typography from '@tiptap/extension-typography';
import Link from '@tiptap/extension-link';
import { Underline } from '@tiptap/extension-underline';
import { TextAlign } from '@tiptap/extension-text-align';
import { Table } from '@tiptap/extension-table';
import { TableRow } from '@tiptap/extension-table-row';
import { TableCell } from '@tiptap/extension-table-cell';
import { TableHeader } from '@tiptap/extension-table-header';
import { Image } from '@tiptap/extension-image';
import { CodeBlockLowlight } from '@tiptap/extension-code-block-lowlight';
import { common, createLowlight } from 'lowlight';

@Component({
  selector: 'app-tiptap-editor',
  imports: [NgIcon, HlmButton, HlmIcon, HlmSeparator],
  providers: [
    provideIcons({
      lucideBold,
      lucideItalic,
      lucideUnderline,
      lucideStrikethrough,
      lucideCode,
      lucideHeading1,
      lucideHeading2,
      lucideHeading3,
      lucideList,
      lucideListOrdered,
      lucideListTodo,
      lucideQuote,
      lucideCode2,
      lucideMinus,
      lucideUndo,
      lucideRedo,
      lucideLink,
      lucideImage,
      lucideTable,
      lucideHighlighter,
    }),
  ],
  template: `
    <div class="tiptap-wrapper flex flex-col h-full">
      @if (editor) {
      <div
        class="toolbar border-b border-border bg-background sticky top-0 z-10 p-2 flex flex-wrap gap-1 items-center"
      >
        <button
          hlmBtn
          variant="ghost"
          size="sm"
          class="h-8 w-8 p-0"
          (click)="editor.chain().focus().toggleBold().run()"
          [class.bg-accent]="editor.isActive('bold')"
        >
          <ng-icon hlm name="lucideBold" size="sm" />
        </button>
        <button
          hlmBtn
          variant="ghost"
          size="sm"
          class="h-8 w-8 p-0"
          (click)="editor.chain().focus().toggleItalic().run()"
          [class.bg-accent]="editor.isActive('italic')"
        >
          <ng-icon hlm name="lucideItalic" size="sm" />
        </button>
        <button
          hlmBtn
          variant="ghost"
          size="sm"
          class="h-8 w-8 p-0"
          (click)="editor.chain().focus().toggleUnderline().run()"
          [class.bg-accent]="editor.isActive('underline')"
        >
          <ng-icon hlm name="lucideUnderline" size="sm" />
        </button>
        <button
          hlmBtn
          variant="ghost"
          size="sm"
          class="h-8 w-8 p-0"
          (click)="editor.chain().focus().toggleStrike().run()"
          [class.bg-accent]="editor.isActive('strike')"
        >
          <ng-icon hlm name="lucideStrikethrough" size="sm" />
        </button>
        <button
          hlmBtn
          variant="ghost"
          size="sm"
          class="h-8 w-8 p-0"
          (click)="editor.chain().focus().toggleCode().run()"
          [class.bg-accent]="editor.isActive('code')"
        >
          <ng-icon hlm name="lucideCode" size="sm" />
        </button>
        <button
          hlmBtn
          variant="ghost"
          size="sm"
          class="h-8 w-8 p-0"
          (click)="editor.chain().focus().toggleHighlight().run()"
          [class.bg-accent]="editor.isActive('highlight')"
        >
          <ng-icon hlm name="lucideHighlighter" size="sm" />
        </button>

        <hr hlmSeparator orientation="vertical" class="h-6" />

        <button
          hlmBtn
          variant="ghost"
          size="sm"
          class="h-8 w-8 p-0"
          (click)="editor.chain().focus().toggleHeading({ level: 1 }).run()"
          [class.bg-accent]="editor.isActive('heading', { level: 1 })"
        >
          <ng-icon hlm name="lucideHeading1" size="sm" />
        </button>
        <button
          hlmBtn
          variant="ghost"
          size="sm"
          class="h-8 w-8 p-0"
          (click)="editor.chain().focus().toggleHeading({ level: 2 }).run()"
          [class.bg-accent]="editor.isActive('heading', { level: 2 })"
        >
          <ng-icon hlm name="lucideHeading2" size="sm" />
        </button>
        <button
          hlmBtn
          variant="ghost"
          size="sm"
          class="h-8 w-8 p-0"
          (click)="editor.chain().focus().toggleHeading({ level: 3 }).run()"
          [class.bg-accent]="editor.isActive('heading', { level: 3 })"
        >
          <ng-icon hlm name="lucideHeading3" size="sm" />
        </button>

        <hr hlmSeparator orientation="vertical" class="h-6" />

        <button
          hlmBtn
          variant="ghost"
          size="sm"
          class="h-8 w-8 p-0"
          (click)="editor.chain().focus().toggleBulletList().run()"
          [class.bg-accent]="editor.isActive('bulletList')"
        >
          <ng-icon hlm name="lucideList" size="sm" />
        </button>
        <button
          hlmBtn
          variant="ghost"
          size="sm"
          class="h-8 w-8 p-0"
          (click)="editor.chain().focus().toggleOrderedList().run()"
          [class.bg-accent]="editor.isActive('orderedList')"
        >
          <ng-icon hlm name="lucideListOrdered" size="sm" />
        </button>
        <button
          hlmBtn
          variant="ghost"
          size="sm"
          class="h-8 w-8 p-0"
          (click)="editor.chain().focus().toggleTaskList().run()"
          [class.bg-accent]="editor.isActive('taskList')"
        >
          <ng-icon hlm name="lucideListTodo" size="sm" />
        </button>

        <hr hlmSeparator orientation="vertical" class="h-6" />

        <button
          hlmBtn
          variant="ghost"
          size="sm"
          class="h-8 w-8 p-0"
          (click)="editor.chain().focus().toggleBlockquote().run()"
          [class.bg-accent]="editor.isActive('blockquote')"
        >
          <ng-icon hlm name="lucideQuote" size="sm" />
        </button>
        <button
          hlmBtn
          variant="ghost"
          size="sm"
          class="h-8 w-8 p-0"
          (click)="editor.chain().focus().toggleCodeBlock().run()"
          [class.bg-accent]="editor.isActive('codeBlock')"
        >
          <ng-icon hlm name="lucideCode2" size="sm" />
        </button>
        <button
          hlmBtn
          variant="ghost"
          size="sm"
          class="h-8 w-8 p-0"
          (click)="editor.chain().focus().setHorizontalRule().run()"
        >
          <ng-icon hlm name="lucideMinus" size="sm" />
        </button>

        <hr hlmSeparator orientation="vertical" class="h-6" />

        <button
          hlmBtn
          variant="ghost"
          size="sm"
          class="h-8 w-8 p-0"
          (click)="toggleLink()"
          [class.bg-accent]="editor.isActive('link')"
        >
          <ng-icon hlm name="lucideLink" size="sm" />
        </button>
        <button hlmBtn variant="ghost" size="sm" class="h-8 w-8 p-0" (click)="addImage()">
          <ng-icon hlm name="lucideImage" size="sm" />
        </button>
        <button hlmBtn variant="ghost" size="sm" class="h-8 w-8 p-0" (click)="insertTable()">
          <ng-icon hlm name="lucideTable" size="sm" />
        </button>

        <hr hlmSeparator orientation="vertical" class="h-6" />

        <button
          hlmBtn
          variant="ghost"
          size="sm"
          class="h-8 w-8 p-0"
          (click)="editor.chain().focus().undo().run()"
          [disabled]="!editor.can().undo()"
        >
          <ng-icon hlm name="lucideUndo" size="sm" />
        </button>
        <button
          hlmBtn
          variant="ghost"
          size="sm"
          class="h-8 w-8 p-0"
          (click)="editor.chain().focus().redo().run()"
          [disabled]="!editor.can().redo()"
        >
          <ng-icon hlm name="lucideRedo" size="sm" />
        </button>
      </div>
      }
      <div class="flex-1 overflow-auto">
        <div #editorContainer class="tiptap-editor"></div>
      </div>
    </div>
  `,
})
export class TiptapEditorComponent implements OnInit, OnDestroy {
  readonly content = input<string>('');
  readonly placeholder = input<string>('Start typing or press "/" for commands...');
  readonly contentChange = output<string>();

  private readonly editorContainer =
    viewChild.required<ElementRef<HTMLDivElement>>('editorContainer');
  editor: Editor | null = null;

  constructor() {
    effect(() => {
      const newContent = this.content();
      if (this.editor && this.editor.getHTML() !== newContent) {
        this.editor.commands.setContent(newContent, { emitUpdate: false });
      }
    });
  }

  ngOnInit(): void {
    this.initEditor();
  }

  ngOnDestroy(): void {
    this.editor?.destroy();
  }

  private initEditor(): void {
    const container = this.editorContainer();
    if (!container) return;

    const lowlight = createLowlight(common);

    this.editor = new Editor({
      element: container.nativeElement,
      extensions: [
        StarterKit.configure({
          heading: { levels: [1, 2, 3, 4, 5, 6] },
          codeBlock: false,
        }),
        CodeBlockLowlight.configure({
          lowlight,
        }),
        Placeholder.configure({
          placeholder: this.placeholder(),
        }),
        Highlight.configure({
          multicolor: true,
        }),
        TaskList,
        TaskItem.configure({
          nested: true,
        }),
        Typography,
        Link.configure({
          openOnClick: false,
          HTMLAttributes: {
            rel: 'noopener noreferrer',
            target: '_blank',
          },
        }),
        Underline,
        TextAlign.configure({
          types: ['heading', 'paragraph'],
        }),
        Table.configure({
          resizable: true,
        }),
        TableRow,
        TableCell,
        TableHeader,
        Image.configure({
          inline: true,
          allowBase64: true,
        }),
      ],
      content: this.content(),
      editorProps: {
        attributes: {
          class:
            'tiptap prose prose-sm sm:prose lg:prose-lg xl:prose-xl max-w-none focus:outline-none',
        },
      },
      onUpdate: ({ editor }) => {
        this.contentChange.emit(editor.getHTML());
      },
    });
  }

  toggleLink(): void {
    if (!this.editor) return;

    if (this.editor.isActive('link')) {
      this.editor.chain().focus().unsetLink().run();
    } else {
      const url = prompt('Enter URL:');
      if (url) {
        this.editor.chain().focus().setLink({ href: url }).run();
      }
    }
  }

  addImage(): void {
    if (!this.editor) return;

    const url = prompt('Enter image URL:');
    if (url) {
      this.editor.chain().focus().setImage({ src: url }).run();
    }
  }

  insertTable(): void {
    if (!this.editor) return;

    this.editor.chain().focus().insertTable({ rows: 3, cols: 3, withHeaderRow: true }).run();
  }
}
