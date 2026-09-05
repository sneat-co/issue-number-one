import { Component, inject, signal } from '@angular/core';
import { Router } from '@angular/router';
import { SneatApiService } from '@sneat/api';
import { AnalyticsService } from '@sneat/core';
import { firstValueFrom } from 'rxjs';
import { loadPending, savePending } from './pending-answer';
import { loadPendingQuestion, savePendingQuestion } from './pending-question';
// Shared account-linked phone verification; implementation belongs to Sneat Libs.
import { PhoneVerificationComponent } from '@sneat/auth-ui';

@Component({
  selector: 'app-verification-page',
  imports: [PhoneVerificationComponent],
  template: `<main>
    <a [href]="backPath">← Back to your question</a>
    <h1>A little friction. A fairer picture.</h1>
    <p>
      Verify by phone, or make a one-time €1 card payment to support
      IssueNumber.one. Either option lets you participate. Payment gives you no
      extra voting weight.
    </p>
    <p>
      These checks make mass manipulation harder; they cannot guarantee that
      every account belongs to a different person.
    </p>
    @if (!hasPending) {
      <p>
        Choose an issue or write a question first. We’ll keep that draft while
        you verify.
      </p>
    } @else {
      <div class="options">
        @for (method of order; track method) {
          <section>
            <h2>
              {{ method === 'phone' ? 'Verify by phone' : 'Support with €1' }}
            </h2>
            <p>
              {{
                method === 'phone'
                  ? 'Receive a verification code by SMS. Your phone number is not shown on public pages.'
                  : 'Make a one-time €1 card payment. It supports this project and makes repeated participation more costly.'
              }}
            </p>
            <button [disabled]="busy()" (click)="choose(method)">
              {{
                method === 'phone'
                  ? 'Choose phone verification'
                  : 'Choose €1 payment'
              }}
            </button>
          </section>
        }
      </div>
      @if (selected() === 'phone') {
        <sneat-phone-verification (verified)="phoneCompleted()" />
      }
    }
    <p role="status" aria-live="polite">{{ status() }}</p>
    <p><a href="/privacy">How we handle your information</a></p>
  </main>`,
  styles: [
    `
      main {
        max-width: 820px;
        margin: 3rem auto;
        padding: 1.5rem;
        font: 1rem/1.6 system-ui;
        color: #1c1517;
      }
      h1 {
        font-size: 2.3rem;
        line-height: 1.15;
      }
      a {
        color: #9f1239;
      }
      .options {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 1rem;
      }
      section {
        border: 1px solid #e4ced4;
        border-radius: 16px;
        padding: 1.4rem;
      }
      button {
        font: inherit;
        background: #be123c;
        color: white;
        border: 0;
        border-radius: 24px;
        padding: 0.8rem 1rem;
        cursor: pointer;
      }
      button:disabled {
        opacity: 0.5;
      }
      :focus-visible {
        outline: 3px solid #865a00;
        outline-offset: 4px;
      }
      @media (max-width: 600px) {
        .options {
          grid-template-columns: 1fr;
        }
      }
    `,
  ],
})
export class VerificationPage {
  private readonly api = inject(SneatApiService);
  private readonly analytics = inject(AnalyticsService);
  private readonly router = inject(Router);
  readonly pending = loadPending();
  readonly pendingQuestion = loadPendingQuestion();
  readonly hasPending = !!this.pending || !!this.pendingQuestion;
  readonly selected = signal<'phone' | 'payment' | null>(null);
  readonly busy = signal(false);
  readonly status = signal(
    'Neither option is preselected. Choose whichever works for you.',
  );
  readonly order: readonly ('phone' | 'payment')[];
  get backPath(): string {
    return this.pending?.returnPath || '/questions/new';
  }
  private get returnRoute(): string {
    return this.pending ? '/answer' : '/questions/submit';
  }
  constructor() {
    const p = this.pending || this.pendingQuestion;
    if (p && !p.verificationOrder) {
      const random = new Uint8Array(1);
      crypto.getRandomValues(random);
      const verificationOrder: readonly ('phone' | 'payment')[] =
        random[0] % 2 ? ['phone', 'payment'] : ['payment', 'phone'];
      const verificationAttemptId = crypto.randomUUID();
      if (this.pending)
        savePending({
          ...this.pending,
          verificationOrder,
          verificationAttemptId,
        });
      else if (this.pendingQuestion)
        savePendingQuestion({
          ...this.pendingQuestion,
          verificationOrder,
          verificationAttemptId,
        });
    }
    const saved = loadPending() || loadPendingQuestion();
    this.order = saved?.verificationOrder || ['phone', 'payment'];
    if (p) this.track('verification_options_viewed');
  }
  private track(event: string, method?: string): void {
    this.analytics.logEvent(event, {
      category_id: this.pending?.categoryId,
      question_id: this.pending?.questionId,
      intent: this.pending?.intent,
      question_choice_source: this.pendingQuestion?.choiceSource.kind,
      method,
      first_option: this.order[0],
      attempt_id:
        this.pending?.verificationAttemptId ||
        this.pendingQuestion?.verificationAttemptId,
    });
  }
  async choose(method: 'phone' | 'payment'): Promise<void> {
    if (this.busy() || !this.hasPending) return;
    this.selected.set(method);
    this.track('verification_method_selected', method);
    if (method === 'phone') return;
    this.busy.set(true);
    this.status.set('Opening secure checkout…');
    try {
      const response = await firstValueFrom(
        this.api.post<{ checkoutUrl?: string; settled: boolean }>(
          'issuenumber/verification/checkout',
          {
            categoryId: this.pending?.categoryId,
            questionId: this.pending?.questionId,
            actionId:
              this.pending?.verificationAttemptId ||
              this.pendingQuestion?.verificationAttemptId,
          },
          { retryUnauthorizedOnce: true },
        ),
      );
      if (response.settled) {
        this.track('verification_completed', 'payment');
        await this.router.navigateByUrl(this.returnRoute, { replaceUrl: true });
        return;
      }
      if (!response.checkoutUrl) throw new Error('No checkout URL');
      const url = new URL(response.checkoutUrl);
      if (url.protocol !== 'https:' || url.hostname !== 'checkout.stripe.com')
        throw new Error('Unexpected checkout destination');
      location.assign(url.href);
    } catch {
      this.status.set(
        'Checkout could not be opened. Your draft is preserved. You can retry or choose phone verification.',
      );
      this.track('verification_failed', 'payment');
    } finally {
      this.busy.set(false);
    }
  }
  async phoneCompleted(): Promise<void> {
    this.track('verification_completed', 'phone');
    this.status.set('Phone verified. Returning to your saved draft…');
    await this.router.navigateByUrl(this.returnRoute, { replaceUrl: true });
  }
}
