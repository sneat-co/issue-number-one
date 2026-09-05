import { convertToParamMap } from '@angular/router';
import { TestBed } from '@angular/core/testing';
import { ActivatedRoute, Router } from '@angular/router';
import { SneatApiService } from '@sneat/api-public';
import { SNEAT_FIREBASE_AUTH } from '@sneat/core';
import { of } from 'rxjs';
import {
  QUESTION_PREVIEW_AUTH_OBSERVER,
  QuestionPreviewPage,
} from './question-preview-page';

vi.mock('@sneat/api-public', () => ({
  SneatApiService: class SneatApiService {},
}));
vi.mock('@sneat/core', async () => {
  const { InjectionToken } = await import('@angular/core');
  return { SNEAT_FIREBASE_AUTH: new InjectionToken('SNEAT_FIREBASE_AUTH') };
});

describe('QuestionPreviewPage', () => {
  const navigate = vi.fn().mockResolvedValue(true);
  const get = vi.fn();

  async function configure(authenticated: boolean): Promise<void> {
    navigate.mockClear();
    get.mockReset();
    await TestBed.configureTestingModule({
      imports: [QuestionPreviewPage],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: {
            snapshot: { paramMap: convertToParamMap({ slug: 'safer-dublin' }) },
          },
        },
        { provide: Router, useValue: { navigate } },
        { provide: SneatApiService, useValue: { get } },
        { provide: SNEAT_FIREBASE_AUTH, useValue: {} },
        {
          provide: QUESTION_PREVIEW_AUTH_OBSERVER,
          useValue: (_auth: unknown, observer: (user: unknown) => void) => {
            observer(authenticated ? { uid: 'owner-1' } : null);
            return () => undefined;
          },
        },
      ],
    }).compileComponents();
  }

  it('sends an anonymous owner to sign in with the exact preview return route', async () => {
    await configure(false);
    const fixture = TestBed.createComponent(QuestionPreviewPage);
    fixture.detectChanges();
    expect(navigate).toHaveBeenCalledWith(['/login'], {
      fragment: '/questions/safer-dublin?preview=1',
      queryParams: { reason: 'Sign in to preview your question' },
    });
    expect(get).not.toHaveBeenCalled();
  });

  it('renders escaped owner-only pending content and approved choices', async () => {
    await configure(true);
    get.mockReturnValue(
      of({
        question: {
          id: 'q1',
          slug: 'safer-dublin',
          title: 'What should <Dublin> fix first?',
          description: 'An owner preview before publication.',
          status: 'pending',
          choiceSource: { kind: 'custom' },
          allowSuggestions: true,
        },
        issues: [
          {
            id: 'i1',
            title: '<script>Safer crossings</script>',
            description: 'Near schools.',
          },
        ],
        totalRespondents: 0,
      }),
    );
    const fixture = TestBed.createComponent(QuestionPreviewPage);
    fixture.detectChanges();
    expect(get).toHaveBeenCalledOnce();
    const params = get.mock.calls[0][1];
    expect(params.get('slug')).toBe('safer-dublin');
    const html = fixture.nativeElement as HTMLElement;
    expect(html.textContent).toContain('Pending review · visible only to you');
    expect(html.textContent).toContain('<script>Safer crossings</script>');
    expect(html.querySelector('script')).toBeNull();
    expect(html.querySelector('meta[name="robots"]')).toBeNull();
    expect(
      document.head
        .querySelector('meta[name="robots"]')
        ?.getAttribute('content'),
    ).toBe('noindex, nofollow');
    fixture.destroy();
  });
});
