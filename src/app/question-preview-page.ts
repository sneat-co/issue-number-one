import { HttpParams } from '@angular/common/http';
import {
  ChangeDetectionStrategy,
  Component,
  inject,
  InjectionToken,
  OnDestroy,
  OnInit,
  signal,
} from '@angular/core';
import { Meta } from '@angular/platform-browser';
import { ActivatedRoute, Router } from '@angular/router';
import { SneatApiService } from '@sneat/api-public';
import { SNEAT_FIREBASE_AUTH } from '@sneat/core';
import { IonButton } from '@ionic/angular/ion-button';
import { IonContent } from '@ionic/angular/ion-content';
import { IonHeader } from '@ionic/angular/ion-header';
import { IonItem } from '@ionic/angular/ion-item';
import { IonLabel } from '@ionic/angular/ion-label';
import { IonList } from '@ionic/angular/ion-list';
import { IonNote } from '@ionic/angular/ion-note';
import { IonSpinner } from '@ionic/angular/ion-spinner';
import { IonTitle } from '@ionic/angular/ion-title';
import { IonToolbar } from '@ionic/angular/ion-toolbar';
import { onAuthStateChanged, Unsubscribe } from 'firebase/auth';

type PreviewAuthObserver = typeof onAuthStateChanged;
export const QUESTION_PREVIEW_AUTH_OBSERVER =
  new InjectionToken<PreviewAuthObserver>('QUESTION_PREVIEW_AUTH_OBSERVER', {
    providedIn: 'root',
    factory: () => onAuthStateChanged,
  });

interface QuestionPreviewIssue {
  readonly id: string;
  readonly slug?: string;
  readonly title: string;
  readonly description?: string;
}

interface QuestionPreviewResponse {
  readonly question: {
    readonly id: string;
    readonly slug: string;
    readonly title: string;
    readonly description: string;
    readonly status: 'pending' | 'published' | 'hidden';
    readonly choiceSource: {
      readonly kind: 'predefined' | 'custom' | 'free';
      readonly entityType?: 'country' | 'city' | 'currency';
    };
    readonly allowSuggestions: boolean;
  };
  readonly issues: readonly QuestionPreviewIssue[];
  readonly totalRespondents: number;
}

type PreviewState = 'checking-auth' | 'loading' | 'ready' | 'error';

@Component({
  selector: 'app-question-preview-page',
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [
    IonHeader,
    IonToolbar,
    IonTitle,
    IonContent,
    IonSpinner,
    IonNote,
    IonList,
    IonItem,
    IonLabel,
    IonButton,
  ],
  template: `
    <ion-header
      ><ion-toolbar
        ><ion-title>Question preview</ion-title></ion-toolbar
      ></ion-header
    >
    <ion-content class="ion-padding">
      <main class="preview-shell">
        @if (state() === 'checking-auth' || state() === 'loading') {
          <div class="center" role="status">
            <ion-spinner></ion-spinner>
            <p>Loading your question preview…</p>
          </div>
        } @else if (state() === 'error') {
          <section class="notice" role="alert">
            <h1>Preview unavailable</h1>
            <p>{{ errorMessage() }}</p>
            <ion-button fill="outline" (click)="retry()">Try again</ion-button>
          </section>
        } @else if (preview(); as result) {
          @if (result.question.status === 'pending') {
            <ion-note color="warning"
              >Pending review · visible only to you</ion-note
            >
          } @else if (result.question.status === 'hidden') {
            <ion-note color="danger"
              >This question is not publicly listed.</ion-note
            >
          } @else {
            <ion-note color="success">Published</ion-note>
          }
          <h1>{{ result.question.title }}</h1>
          <p class="description">{{ result.question.description }}</p>

          <section aria-labelledby="choices-heading">
            <h2 id="choices-heading">Answer choices</h2>
            @if (result.issues.length) {
              <ion-list>
                @for (issue of result.issues; track issue.id) {
                  <ion-item>
                    <ion-label
                      ><h3>{{ issue.title }}</h3>
                      @if (issue.description) {
                        <p>{{ issue.description }}</p>
                      }
                    </ion-label>
                  </ion-item>
                }
              </ion-list>
            } @else if (result.question.choiceSource.kind === 'free') {
              <p>People can describe the issue that matters most to them.</p>
            } @else {
              <p>No approved choices are available yet.</p>
            }
          </section>

          @if (result.question.status === 'published') {
            <a
              class="public-link"
              [href]="publicQuestionURL(result.question.slug)"
              >Open the public question</a
            >
          } @else {
            <p class="privacy-copy">
              This preview is not a public link. Publication follows review.
            </p>
          }
        }
      </main>
    </ion-content>
  `,
  styles: [
    `
      .preview-shell {
        max-width: 720px;
        margin: 0 auto;
      }
      .center,
      .notice {
        min-height: 45vh;
        display: grid;
        place-content: center;
        text-align: center;
      }
      h1 {
        margin-top: 1rem;
      }
      h2 {
        font-size: 1.15rem;
        margin-top: 2rem;
      }
      .description,
      .privacy-copy {
        color: var(--ion-color-medium);
        line-height: 1.55;
      }
      .public-link {
        display: inline-block;
        margin-top: 1.5rem;
        font-weight: 700;
      }
    `,
  ],
})
export class QuestionPreviewPage implements OnInit, OnDestroy {
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly api = inject(SneatApiService);
  private readonly auth = inject(SNEAT_FIREBASE_AUTH);
  private readonly observeAuth = inject(QUESTION_PREVIEW_AUTH_OBSERVER);
  private readonly meta = inject(Meta);
  private unsubscribeAuth?: Unsubscribe;
  private previousRobots?: string;

  protected readonly state = signal<PreviewState>('checking-auth');
  protected readonly preview = signal<QuestionPreviewResponse | undefined>(
    undefined,
  );
  protected readonly errorMessage = signal(
    'We could not load this owner preview.',
  );

  ngOnInit(): void {
    this.previousRobots = this.meta.getTag('name="robots"')?.content;
    this.meta.updateTag({ name: 'robots', content: 'noindex, nofollow' });
    this.unsubscribeAuth = this.observeAuth(this.auth, (user) => {
      this.unsubscribeAuth?.();
      this.unsubscribeAuth = undefined;
      if (!user) {
        const slug = this.route.snapshot.paramMap.get('slug') ?? '';
        void this.router.navigate(['/login'], {
          fragment: `/questions/${encodeURIComponent(slug)}?preview=1`,
          queryParams: { reason: 'Sign in to preview your question' },
        });
        return;
      }
      this.load();
    });
  }

  ngOnDestroy(): void {
    this.unsubscribeAuth?.();
    if (this.previousRobots)
      this.meta.updateTag({ name: 'robots', content: this.previousRobots });
    else this.meta.removeTag('name="robots"');
  }

  retry(): void {
    this.load();
  }

  protected publicQuestionURL(slug: string): string {
    return `/questions/${encodeURIComponent(slug)}?`;
  }

  private load(): void {
    const slug = this.route.snapshot.paramMap.get('slug');
    if (!slug) {
      this.state.set('error');
      this.errorMessage.set('The question address is missing its slug.');
      return;
    }
    this.state.set('loading');
    this.api
      .get<QuestionPreviewResponse>(
        'issuenumber/question',
        new HttpParams().set('slug', slug),
      )
      .subscribe({
        next: (response) => {
          this.preview.set(response);
          this.state.set('ready');
        },
        error: () => {
          this.state.set('error');
          this.errorMessage.set(
            'This question is unavailable or does not belong to the signed-in account.',
          );
        },
      });
  }
}
