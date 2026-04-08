// This file is replaced during prod builds — see project.json fileReplacements.
import { appSpecificConfig, emulatorEnvironmentConfig } from '@sneat/app';
import { IEnvironmentConfig } from '@sneat/core';

export const in1AppEnvironmentConfig: IEnvironmentConfig = appSpecificConfig({
  ...emulatorEnvironmentConfig,
  firebaseConfig: {
    ...emulatorEnvironmentConfig.firebaseConfig,
    projectId: 'sneat-work',
  },
});
