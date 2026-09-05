import { Component, DestroyRef, inject, signal } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { Router } from '@angular/router';
import { SneatApiService } from '@sneat/api';
import { SneatAuthStateService } from '@sneat/auth-core';
import { AnalyticsService } from '@sneat/core';
import { firstValueFrom } from 'rxjs';
import {
  clearPendingQuestion,
  loadPendingQuestion,
  PendingQuestion,
  savePendingQuestion,
} from './pending-question';

interface CreateQuestionResult {
  readonly question: {
    readonly id: string;
    readonly slug: string;
    readonly title: string;
    readonly translationStatus?: string;
  };
  readonly replayed?: boolean;
}

@Component({
  selector: 'app-question-submit-page',
  template: `<main>
    <a href="/questions/new">← Edit question</a>
    <h1>
      {{ created() ? 'Your question is saved' : 'Publishing your draft' }}
    </h1>
    @if (pending) {
      <blockquote>{{ pending.title }}</blockquote>
    }
    <p role="status" aria-live="polite">{{ status() }}</p>
    @if (!busy() && !created() && pending) {
      <button (click)="submit()">Try again</button>
    }
  </main>`,
  styles: [
    `
      main {
        max-width: 680px;
        margin: 3rem auto;
        padding: 1.5rem;
        font: 1rem/1.6 system-ui;
        color: #1c1517;
      }
      h1 {
        font-size: 2.4rem;
        line-height: 1.15;
      }
      a {
        color: #9f1239;
      }
      button {
        font: inherit;
        border: 0;
        border-radius: 24px;
        padding: 0.8rem 1.2rem;
        color: white;
        background: #be123c;
        cursor: pointer;
      }
      blockquote {
        padding: 1rem;
        border-left: 3px solid #be123c;
        background: #fff5f7;
      }
      :focus-visible {
        outline: 3px solid #865a00;
        outline-offset: 4px;
      }
    `,
  ],
})
export class QuestionSubmitPage {
  private readonly router = inject(Router);
  private readonly auth = inject(SneatAuthStateService);
  private readonly api = inject(SneatApiService);
  private readonly analytics = inject(AnalyticsService);
  private readonly destroyRef = inject(DestroyRef);
  readonly pending: PendingQuestion | undefined = loadPendingQuestion();
  readonly status = signal('Checking your sign-in…');
  readonly busy = signal(true);
  readonly created = signal(false);
  private started = false;

  constructor() {
    const pending = this.pending;
    if (!pending) {
      this.status.set(
        'Your saved draft is missing or expired. Return to the question form to start again.',
      );
      this.busy.set(false);
      return;
    }
    this.auth.authUser
      .pipe(takeUntilDestroyed(this.destroyRef))
      .subscribe((user) => {
        if (user === undefined || this.started) return;
        if (!user || user.isAnonymous) {
          this.started = true;
          savePendingQuestion({ ...pending, requiredAuth: true });
          this.track('auth_required_for_question');
          void this.router.navigateByUrl(
            '/login?reason=Sign%20in%20to%20submit%20your%20question#/questions/submit',
          );
          return;
        }
        this.started = true;
        if (pending.requiredAuth) this.track('auth_completed_for_question');
        void this.submit();
      });
  }

  async submit(): Promise<void> {
    if (!this.pending) return;
    this.busy.set(true);
    this.status.set('Saving your question and preparing translations…');
    try {
      const result = await firstValueFrom(
        this.api.post<CreateQuestionResult>(
          'issuenumber/question',
          {
            title: this.pending.title,
            description: this.pending.description,
            choiceSource: this.pending.choiceSource,
            allowSuggestions: this.pending.allowSuggestions,
            sourceLanguage: this.pending.sourceLanguage || 'en',
            operationId: this.pending.operationId,
          },
          { retryUnauthorizedOnce: true },
        ),
      );
      this.created.set(true);
      this.status.set(
        result.question.translationStatus === 'failed'
          ? 'Your question is saved for review. Some translations will be retried.'
          : 'Your question is saved for review. Translations are stored with it.',
      );
      this.track('question_created');
      clearPendingQuestion();
      await this.router.navigateByUrl(
        `/questions/${encodeURIComponent(result.question.slug)}?preview=1`,
        { replaceUrl: true },
      );
    } catch (error: unknown) {
      const response = error as {
        status?: number;
        error?: { code?: string; error?: string };
      };
      if (
        response.status === 403 &&
        (response.error?.code === 'verification_required' ||
          response.error?.error === 'verification_required')
      ) {
        savePendingQuestion(this.pending);
        await this.router.navigateByUrl('/verify');
      } else if (response.status === 409) {
        this.status.set(
          'A question with this address already exists. Edit the wording and try again.',
        );
      } else if (response.status === 429) {
        this.status.set(
          'You have reached today’s question creation limit. Your draft is preserved.',
        );
      } else {
        this.status.set(
          'We could not save your question. Your draft is preserved. Please try again.',
        );
      }
    } finally {
      this.busy.set(false);
    }
  }

  private track(name: string): void {
    this.analytics.logEvent(name, {
      choice_source: this.pending?.choiceSource.kind,
      source_language: this.pending?.sourceLanguage || 'en',
    });
  }
}
