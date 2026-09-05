import { TestBed } from '@angular/core/testing';
import { Router } from '@angular/router';
import { SneatApiService } from '@sneat/api';
import { SneatAuthStateService } from '@sneat/auth-core';
import { AnalyticsService } from '@sneat/core';
import { of, throwError } from 'rxjs';
import { savePendingQuestion } from './pending-question';
import { QuestionSubmitPage } from './question-submit-page';

vi.mock('@sneat/api', () => ({
  SneatApiService: class SneatApiService {},
}));
vi.mock('@sneat/auth-core', () => ({
  SneatAuthStateService: class SneatAuthStateService {},
}));
vi.mock('@sneat/core', () => ({
  AnalyticsService: class AnalyticsService {},
}));

describe('QuestionSubmitPage', () => {
  const navigateByUrl = vi.fn().mockResolvedValue(true);
  const post = vi.fn();
  const logEvent = vi.fn();

  function saveDraft(): void {
    savePendingQuestion({
      operationId: 'question-operation-1',
      createdAt: new Date().toISOString(),
      title: 'What should Dublin improve first?',
      description: 'Choose the change with the greatest practical effect.',
      choiceSource: { kind: 'predefined', entityType: 'city' },
      allowSuggestions: true,
      sourceLanguage: 'en',
    });
  }

  async function configure(user: unknown): Promise<void> {
    navigateByUrl.mockClear();
    post.mockReset();
    logEvent.mockClear();
    await TestBed.configureTestingModule({
      imports: [QuestionSubmitPage],
      providers: [
        { provide: Router, useValue: { navigateByUrl } },
        { provide: SneatApiService, useValue: { post } },
        { provide: SneatAuthStateService, useValue: { authUser: of(user) } },
        { provide: AnalyticsService, useValue: { logEvent } },
      ],
    }).compileComponents();
  }

  beforeEach(() => sessionStorage.clear());

  it('preserves the draft and returns through authentication', async () => {
    saveDraft();
    await configure(null);
    TestBed.createComponent(QuestionSubmitPage).detectChanges();
    expect(post).not.toHaveBeenCalled();
    expect(navigateByUrl).toHaveBeenCalledWith(
      '/login?reason=Sign%20in%20to%20submit%20your%20question#/questions/submit',
    );
  });

  it('submits the stored choice source and opens the private preview', async () => {
    saveDraft();
    await configure({ uid: 'user-1', isAnonymous: false });
    post.mockReturnValue(
      of({
        question: {
          id: 'q1',
          slug: 'what-should-dublin-improve-first',
          title: 'What should Dublin improve first?',
          translationStatus: 'ready',
        },
      }),
    );
    TestBed.createComponent(QuestionSubmitPage).detectChanges();
    await vi.waitFor(() => expect(post).toHaveBeenCalledOnce());
    expect(post.mock.calls[0][1]).toMatchObject({
      choiceSource: { kind: 'predefined', entityType: 'city' },
      sourceLanguage: 'en',
      operationId: 'question-operation-1',
    });
    expect(navigateByUrl).toHaveBeenCalledWith(
      '/questions/what-should-dublin-improve-first?preview=1',
      { replaceUrl: true },
    );
  });

  it('continues the same draft through participation verification', async () => {
    saveDraft();
    await configure({ uid: 'user-1', isAnonymous: false });
    post.mockReturnValue(
      throwError(() => ({
        status: 403,
        error: { code: 'verification_required' },
      })),
    );
    TestBed.createComponent(QuestionSubmitPage).detectChanges();
    await vi.waitFor(() =>
      expect(navigateByUrl).toHaveBeenCalledWith('/verify'),
    );
  });
});
