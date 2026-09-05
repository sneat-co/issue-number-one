import { Component, DestroyRef, inject, signal } from '@angular/core';
import { Router } from '@angular/router';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { SneatAuthStateService } from '@sneat/auth-core';
import { SneatApiService } from '@sneat/api';
import { AnalyticsService } from '@sneat/core';
import { firstValueFrom } from 'rxjs';
import {
  loadPending,
  pendingKey,
  savePending,
  PendingAnswer,
} from './pending-answer';
import { loadPendingQuestion } from './pending-question';

interface AnswerResult {
  changed?: boolean;
  issueId?: string;
  answer?: { issueId: string };
  personalTop?: boolean;
}
@Component({
  selector: 'app-answer-page',
  template: `<main class="answer-shell">
    <a href="/issues">← Explore issues</a>
    <h1>{{ saved() ? 'Your choice is saved' : 'Your voice matters' }}</h1>
    <p role="status" aria-live="polite">{{ status() }}</p>
    @if (pending?.title) {
      <blockquote>{{ pending.title }}</blockquote>
    }
    @if (saved()) {
      @if (pending?.answerKind !== 'personal') {
        <p>Is this your most important issue across every part of life?</p>
        <button [disabled]="busy()" (click)="makePersonal()">
          Make this my personal #1
        </button>
      }
      <p>
        <a [href]="pending?.returnPath || '/issues'"
          >See the public results →</a
        >
      </p>
      <h2>
        {{
          pending?.intent === 'business'
            ? 'Whose work does this affect?'
            : 'Whose perspective matters to you?'
        }}
      </h2>
      <p>Share the question and give them space to answer for themselves.</p>
      <button (click)="share()">Share this question</button>
      <p>
        <a [href]="destination" (click)="conversion()"
          >{{
            pending?.intent === 'business'
              ? 'Explore IssueNumber.one for work'
              : 'Continue in Sneat.app'
          }}
          →</a
        >
      </p>
    } @else if (!busy() && pending) {
      <button (click)="submit()">Try again</button>
    }
    <p>
      <a [href]="pending?.returnPath || '/issues'">{{
        saved() ? 'Back to question' : 'Cancel and return to the question'
      }}</a>
    </p>
  </main>`,
  styles: [
    `
      .answer-shell {
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
      button:disabled {
        opacity: 0.5;
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
export class AnswerPage {
  private readonly router = inject(Router);
  private readonly auth = inject(SneatAuthStateService);
  private readonly api = inject(SneatApiService);
  private readonly analytics = inject(AnalyticsService);
  private readonly destroyRef = inject(DestroyRef);
  pending: PendingAnswer | null = loadPending();
  readonly status = signal('Checking your sign-in…');
  readonly busy = signal(true);
  readonly saved = signal(false);
  private started = false;
  get destination(): string {
    return this.pending?.intent === 'business'
      ? '/for-work/'
      : 'https://sneat.app/main?utm_source=issuenumber.one&utm_medium=answer&utm_campaign=public-issues';
  }
  constructor() {
    if (!this.pending) {
      if (loadPendingQuestion()) {
        void this.router.navigateByUrl('/questions/submit', {
          replaceUrl: true,
        });
        return;
      }
      this.status.set(
        'Choose an issue first. If your saved draft expired, return to the question to start again.',
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
          if (this.pending) {
            this.pending.requiredAuth = true;
            savePending(this.pending);
          }
          this.track('auth_required_for_answer');
          void this.router.navigateByUrl(
            '/login?reason=Sign%20in%20to%20save%20your%20choice#/answer',
          );
        } else {
          this.started = true;
          if (this.pending?.requiredAuth)
            this.track('auth_completed_for_answer');
          void this.submit();
        }
      });
  }
  private track(name: string): void {
    this.analytics.logEvent(name, {
      category_id: this.pending?.categoryId,
      question_id: this.pending?.questionId,
      intent: this.pending?.intent,
    });
  }
  async submit(): Promise<void> {
    const p = this.pending;
    if (!p) return;
    this.busy.set(true);
    this.status.set('Saving your choice…');
    try {
      const result = await firstValueFrom(
        this.api.post<AnswerResult>(
          'issuenumber/answer',
          {
            questionId: p.questionId,
            issueId: p.issueId,
            title: p.title,
            operationId: p.operationId,
            answerKind: p.answerKind,
            attribution: p.attribution,
            languageCode: p.languageCode || 'en',
          },
          { retryUnauthorizedOnce: true },
        ),
      );
      p.issueId = result.issueId || result.answer?.issueId || p.issueId;
      if (p.issueId) delete p.title;
      this.saved.set(true);
      this.status.set(
        p.answerKind === 'personal'
          ? 'This is now your one personal #1 across all your questions.'
          : 'This is now your choice for this question. Your choices in other questions remain separate.',
      );
      sessionStorage.setItem(
        'issuenumber.saved',
        JSON.stringify({
          questionId: p.questionId,
          createdAt: Date.now(),
          personalTop: p.answerKind === 'personal',
        }),
      );
      sessionStorage.removeItem(pendingKey);
      this.track(result.changed ? 'answer_changed' : 'answer_submitted');
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
        savePending(p);
        await this.router.navigateByUrl('/verify');
      } else
        this.status.set(
          'We couldn’t save your choice. Your draft is preserved. Please try again.',
        );
    } finally {
      this.busy.set(false);
    }
  }
  async makePersonal(): Promise<void> {
    if (!this.pending?.issueId || this.busy()) return;
    this.pending = {
      ...this.pending,
      answerKind: 'personal',
      operationId: crypto.randomUUID(),
      createdAt: Date.now(),
    };
    savePending(this.pending);
    await this.submit();
  }
  async share(): Promise<void> {
    if (!this.pending) return;
    const url = 'https://issuenumber.one' + this.pending.returnPath;
    try {
      if (navigator.share)
        await navigator.share({ title: 'What is your #1 issue?', url });
      else {
        await navigator.clipboard.writeText(url);
        this.status.set('Link copied.');
      }
      this.track('question_shared');
    } catch (error) {
      if ((error as Error).name !== 'AbortError')
        this.status.set('Copy the question’s address to share it.');
    }
  }
  conversion(): void {
    this.track(
      this.pending?.intent === 'business'
        ? 'sneat_work_cta_clicked'
        : 'sneat_app_cta_clicked',
    );
  }
}
