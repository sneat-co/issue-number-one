import { Route } from '@angular/router';

export const appRoutes: Route[] = [
  {
    path: 'questions/new',
    loadComponent: () =>
      import('./question-composer-page').then((m) => m.QuestionComposerPage),
  },
  {
    path: 'questions/submit',
    loadComponent: () =>
      import('./question-submit-page').then((m) => m.QuestionSubmitPage),
  },
  {
    path: 'questions/:slug',
    loadComponent: () =>
      import('./question-preview-page').then((m) => m.QuestionPreviewPage),
  },
  {
    path: 'verify',
    loadComponent: () =>
      import('./verification-page').then((m) => m.VerificationPage),
  },
  {
    path: 'answer',
    loadComponent: () => import('./answer-page').then((m) => m.AnswerPage),
  },
];
