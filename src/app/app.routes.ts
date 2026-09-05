import { Route } from '@angular/router';

export const appRoutes: Route[] = [
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
