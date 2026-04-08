import {
  ApplicationConfig,
  provideBrowserGlobalErrorListeners,
} from '@angular/core';
import { provideRouter } from '@angular/router';
import {
  getStandardSneatProviders,
  provideAppInfo,
  provideRolesByType,
} from '@sneat/app';
import { authRoutes } from '@sneat/auth-ui';
import { SneatApp } from '@sneat/core';
import { in1AppEnvironmentConfig } from '../environments/environment';
import { appRoutes } from './app.routes';

export const appConfig: ApplicationConfig = {
  providers: [
    provideBrowserGlobalErrorListeners(),
    ...getStandardSneatProviders(in1AppEnvironmentConfig),
    // NOTE: 'in1app' is not yet in the SneatApp union type in @sneat/core@0.4.0;
    // follow-up: extend SneatApp in sneat-libs/libs/core/src/lib/app.service.ts and republish.
    provideAppInfo({ appId: 'in1app' as SneatApp, appTitle: 'issue-number-one' }),
    provideRolesByType(undefined),
    provideRouter([...appRoutes, ...authRoutes]),
  ],
};
