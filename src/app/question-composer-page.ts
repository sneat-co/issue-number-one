import { Location } from '@angular/common';
import {
  ChangeDetectionStrategy,
  Component,
  inject,
  signal,
} from '@angular/core';
import {
  FormControl,
  FormGroup,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { Router } from '@angular/router';
import { IonButton } from '@ionic/angular/ion-button';
import { IonCheckbox } from '@ionic/angular/ion-checkbox';
import { IonContent } from '@ionic/angular/ion-content';
import { IonHeader } from '@ionic/angular/ion-header';
import { IonInput } from '@ionic/angular/ion-input';
import { IonItem } from '@ionic/angular/ion-item';
import { IonList } from '@ionic/angular/ion-list';
import { IonNote } from '@ionic/angular/ion-note';
import { IonRadio } from '@ionic/angular/ion-radio';
import { IonRadioGroup } from '@ionic/angular/ion-radio-group';
import { IonSelect } from '@ionic/angular/ion-select';
import { IonSelectOption } from '@ionic/angular/ion-select-option';
import { IonTextarea } from '@ionic/angular/ion-textarea';
import { IonTitle } from '@ionic/angular/ion-title';
import { IonToolbar } from '@ionic/angular/ion-toolbar';
import {
  loadPendingQuestion,
  PendingQuestion,
  PendingQuestionChoiceSource,
  savePendingQuestion,
} from './pending-question';

type ChoiceKind = '' | PendingQuestionChoiceSource['kind'];

@Component({
  selector: 'app-question-composer-page',
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [
    ReactiveFormsModule,
    IonHeader,
    IonToolbar,
    IonTitle,
    IonContent,
    IonList,
    IonItem,
    IonInput,
    IonTextarea,
    IonRadioGroup,
    IonRadio,
    IonSelect,
    IonSelectOption,
    IonCheckbox,
    IonNote,
    IonButton,
  ],
  template: `
    <ion-header
      ><ion-toolbar
        ><ion-title>Ask what matters most</ion-title></ion-toolbar
      ></ion-header
    >
    <ion-content class="ion-padding">
      <main class="composer-shell">
        <header>
          <p class="eyebrow">IssueNumber.one</p>
          <h1>Create a #1 issue question</h1>
          <p>
            Ask people to name the single issue that matters most in a clear
            scope.
          </p>
        </header>

        <form [formGroup]="form" (ngSubmit)="continueToAnswer()" novalidate>
          <ion-list inset>
            <ion-item>
              <ion-input
                label="Question"
                labelPlacement="stacked"
                formControlName="title"
                maxlength="180"
                placeholder="What is the top issue for remote workers?"
              ></ion-input>
            </ion-item>
            @if (submitted() && titleError()) {
              <ion-note color="danger">{{ titleError() }}</ion-note>
            }
            <ion-item>
              <ion-textarea
                label="Why this question matters"
                labelPlacement="stacked"
                formControlName="description"
                maxlength="1000"
                autoGrow="true"
                placeholder="Explain the scope and what respondents should consider."
              ></ion-textarea>
            </ion-item>
            @if (submitted() && descriptionError()) {
              <ion-note color="danger">{{ descriptionError() }}</ion-note>
            }
          </ion-list>

          <section aria-labelledby="answer-source-heading">
            <h2 id="answer-source-heading">Where should answers come from?</h2>
            <p>Choose the structure that fits this question.</p>
            <ion-radio-group formControlName="choiceKind">
              <ion-item
                ><ion-radio value="predefined"
                  >A known set of places or currencies</ion-radio
                ></ion-item
              >
              <ion-item
                ><ion-radio value="custom"
                  >A focused list for this question</ion-radio
                ></ion-item
              >
              <ion-item
                ><ion-radio value="free"
                  >People write their own issue</ion-radio
                ></ion-item
              >
            </ion-radio-group>
            @if (submitted() && !form.controls.choiceKind.value) {
              <ion-note color="danger">Choose how people will answer.</ion-note>
            }

            @if (form.controls.choiceKind.value === 'predefined') {
              <ion-item>
                <ion-select
                  label="Use"
                  labelPlacement="stacked"
                  formControlName="entityType"
                  placeholder="Choose a source"
                >
                  <ion-select-option value="country"
                    >Countries</ion-select-option
                  >
                  <ion-select-option value="city">Cities</ion-select-option>
                  <ion-select-option value="currency"
                    >Currencies</ion-select-option
                  >
                </ion-select>
              </ion-item>
              @if (submitted() && !form.controls.entityType.value) {
                <ion-note color="danger"
                  >Choose countries, cities, or currencies.</ion-note
                >
              }
            }

            @if (form.controls.choiceKind.value === 'custom') {
              <ion-item>
                <ion-textarea
                  label="Candidate issues"
                  labelPlacement="stacked"
                  formControlName="customOptions"
                  autoGrow="true"
                  placeholder="One issue per line"
                ></ion-textarea>
              </ion-item>
              <ion-note>Enter between 2 and 30 distinct issues.</ion-note>
              @if (submitted() && customOptionsError()) {
                <ion-note color="danger">{{ customOptionsError() }}</ion-note>
              }
            }

            @if (form.controls.choiceKind.value !== 'free') {
              <ion-item>
                <ion-checkbox formControlName="allowSuggestions"
                  >Let people suggest an issue that is missing</ion-checkbox
                >
              </ion-item>
            }
          </section>

          <div class="actions">
            <ion-button
              type="button"
              fill="clear"
              color="medium"
              (click)="cancel()"
              >Cancel</ion-button
            >
            <ion-button type="submit">Continue to answer</ion-button>
          </div>
        </form>
      </main>
    </ion-content>
  `,
  styles: [
    `
      .composer-shell {
        max-width: 720px;
        margin: 0 auto;
      }
      .eyebrow {
        color: var(--ion-color-primary);
        font-weight: 700;
        letter-spacing: 0.04em;
      }
      h1 {
        margin-bottom: 0.5rem;
      }
      h2 {
        font-size: 1.15rem;
        margin: 2rem 1rem 0.25rem;
      }
      section > p {
        margin: 0 1rem 1rem;
        color: var(--ion-color-medium);
      }
      ion-note {
        display: block;
        margin: 0.45rem 1rem;
      }
      .actions {
        display: flex;
        justify-content: flex-end;
        gap: 0.5rem;
        margin-top: 1.5rem;
      }
      @media (max-width: 520px) {
        .actions {
          flex-direction: column-reverse;
        }
        .actions ion-button {
          width: 100%;
        }
      }
    `,
  ],
})
export class QuestionComposerPage {
  private readonly router = inject(Router);
  private readonly location = inject(Location);
  protected readonly submitted = signal(false);

  readonly form = new FormGroup({
    title: new FormControl('', {
      nonNullable: true,
      validators: [Validators.required],
    }),
    description: new FormControl('', {
      nonNullable: true,
      validators: [Validators.required],
    }),
    choiceKind: new FormControl<ChoiceKind>('', {
      nonNullable: true,
      validators: [Validators.required],
    }),
    entityType: new FormControl<'country' | 'city' | 'currency' | ''>('', {
      nonNullable: true,
    }),
    customOptions: new FormControl('', { nonNullable: true }),
    allowSuggestions: new FormControl(false, { nonNullable: true }),
  });

  constructor() {
    const pending = loadPendingQuestion();
    if (pending) this.restore(pending);
    else this.restoreComposerDraft();
  }

  protected titleError(): string {
    const title = normalizeQuestion(this.form.controls.title.value);
    if (title.length < 10)
      return 'Write a specific question using at least 10 characters.';
    if (title.length > 180) return 'Keep the question to 180 characters.';
    return '';
  }

  protected descriptionError(): string {
    const length = this.form.controls.description.value.trim().length;
    if (length < 20) return 'Explain the question in at least 20 characters.';
    if (length > 1000) return 'Keep the explanation to 1,000 characters.';
    return '';
  }

  protected customOptionsError(): string {
    if (this.form.controls.choiceKind.value !== 'custom') return '';
    const count = parseCustomOptions(
      this.form.controls.customOptions.value,
    ).length;
    return count < 2 || count > 30
      ? 'Enter between 2 and 30 distinct issues.'
      : '';
  }

  continueToAnswer(): void {
    this.submitted.set(true);
    const source = this.choiceSource();
    if (
      !source ||
      this.titleError() ||
      this.descriptionError() ||
      this.customOptionsError()
    )
      return;
    const existing = loadPendingQuestion();
    savePendingQuestion({
      operationId: existing?.operationId ?? crypto.randomUUID(),
      createdAt: existing?.createdAt ?? new Date().toISOString(),
      title: normalizeQuestion(this.form.controls.title.value),
      description: this.form.controls.description.value.trim(),
      choiceSource: source,
      allowSuggestions:
        source.kind === 'free'
          ? true
          : this.form.controls.allowSuggestions.value,
    });
    sessionStorage.removeItem(composerDraftStorageKey);
    void this.router.navigateByUrl('/answer');
  }

  cancel(): void {
    sessionStorage.setItem(
      composerDraftStorageKey,
      JSON.stringify(this.form.getRawValue()),
    );
    this.location.back();
  }

  private choiceSource(): PendingQuestionChoiceSource | undefined {
    switch (this.form.controls.choiceKind.value) {
      case 'free':
        return { kind: 'free' };
      case 'predefined': {
        const entityType = this.form.controls.entityType.value;
        return entityType ? { kind: 'predefined', entityType } : undefined;
      }
      case 'custom': {
        const options = parseCustomOptions(
          this.form.controls.customOptions.value,
        );
        return options.length >= 2 && options.length <= 30
          ? { kind: 'custom', options }
          : undefined;
      }
      default:
        return undefined;
    }
  }

  private restore(question: PendingQuestion): void {
    const source = question.choiceSource;
    this.form.patchValue({
      title: question.title,
      description: question.description,
      choiceKind: source.kind,
      entityType: source.kind === 'predefined' ? source.entityType : '',
      customOptions:
        source.kind === 'custom'
          ? source.options.map((option) => option.title).join('\n')
          : '',
      allowSuggestions:
        source.kind === 'free' ? true : question.allowSuggestions,
    });
  }

  private restoreComposerDraft(): void {
    try {
      const raw = sessionStorage.getItem(composerDraftStorageKey);
      if (!raw) return;
      const draft = JSON.parse(raw) as Partial<typeof this.form.value>;
      this.form.patchValue(draft);
    } catch {
      sessionStorage.removeItem(composerDraftStorageKey);
    }
  }
}

const composerDraftStorageKey = 'issuenumber.one.question-composer-draft.v1';

export function normalizeQuestion(value: string): string {
  const title = value.trim().replace(/\s+/g, ' ');
  return title.endsWith('?') ? title : `${title}?`;
}

export function parseCustomOptions(value: string): { title: string }[] {
  const seen = new Set<string>();
  return value
    .split(/\r?\n/)
    .map((title) => title.trim().replace(/\s+/g, ' '))
    .filter((title) => {
      const key = title.toLocaleLowerCase();
      if (!title || seen.has(key)) return false;
      seen.add(key);
      return true;
    })
    .map((title) => ({ title }));
}
