import { Location } from '@angular/common';
import { TestBed } from '@angular/core/testing';
import { Router } from '@angular/router';
import { QuestionComposerPage } from './question-composer-page';
import { clearPendingQuestion, loadPendingQuestion } from './pending-question';

describe('QuestionComposerPage', () => {
  const navigateByUrl = vi.fn().mockResolvedValue(true);
  const back = vi.fn();

  beforeEach(async () => {
    sessionStorage.clear();
    navigateByUrl.mockClear();
    back.mockClear();
    await TestBed.configureTestingModule({
      imports: [QuestionComposerPage],
      providers: [
        { provide: Router, useValue: { navigateByUrl } },
        { provide: Location, useValue: { back } },
      ],
    }).compileComponents();
  });

  it('blocks a vague question before creating a pending operation', () => {
    const component =
      TestBed.createComponent(QuestionComposerPage).componentInstance;
    component.form.patchValue({
      title: 'Top issue',
      description: 'Too short',
      choiceKind: 'free',
    });
    component.continueToAnswer();
    expect(loadPendingQuestion()).toBeUndefined();
    expect(navigateByUrl).not.toHaveBeenCalled();
  });

  it('normalizes the question and carries a constrained custom source to /answer', () => {
    const component =
      TestBed.createComponent(QuestionComposerPage).componentInstance;
    component.form.patchValue({
      title: '  What should our neighbourhood fix first  ',
      description:
        '  Choose the local issue with the greatest practical effect.  ',
      choiceKind: 'custom',
      customOptions: 'Safer crossings\nMore trees\nsafer crossings',
      allowSuggestions: true,
    });
    component.continueToAnswer();
    expect(loadPendingQuestion()).toMatchObject({
      title: 'What should our neighbourhood fix first?',
      description: 'Choose the local issue with the greatest practical effect.',
      choiceSource: {
        kind: 'custom',
        options: [{ title: 'Safer crossings' }, { title: 'More trees' }],
      },
      allowSuggestions: true,
    });
    expect(navigateByUrl).toHaveBeenCalledWith('/answer');
  });

  it('keeps an unfinished form when Cancel returns to the previous screen', () => {
    const first =
      TestBed.createComponent(QuestionComposerPage).componentInstance;
    first.form.patchValue({
      title: 'What matters most in Cork',
      description: 'Still writing this explanation',
    });
    first.cancel();
    expect(back).toHaveBeenCalled();
    clearPendingQuestion();
    const restored =
      TestBed.createComponent(QuestionComposerPage).componentInstance;
    expect(restored.form.controls.title.value).toBe(
      'What matters most in Cork',
    );
    expect(restored.form.controls.description.value).toBe(
      'Still writing this explanation',
    );
  });
});
