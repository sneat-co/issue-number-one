import assert from 'node:assert/strict';
import { spawnSync } from 'node:child_process';
import { readFile } from 'node:fs/promises';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import test from 'node:test';
import { initializeTestEnvironment } from '@firebase/rules-unit-testing';

const projectId = 'issuenumber-seed-test';
const repositoryRoot = path.resolve(
  path.dirname(fileURLToPath(import.meta.url)),
  '..',
);
const catalog = JSON.parse(
  await readFile(path.join(repositoryRoot, 'catalog/seed.json'), 'utf8'),
);
const [host, portText] = (process.env.FIRESTORE_EMULATOR_HOST || '').split(':');

test(
  'seed reruns preserve operational counts and do not duplicate records',
  { skip: !host || !portText },
  async () => {
    const environment = await initializeTestEnvironment({
      projectId,
      firestore: { host, port: Number(portText) },
    });
    const runSeed = () => {
      const result = spawnSync(
        'go',
        [
          'run',
          './cmd/seed',
          '--project',
          projectId,
          '--file',
          '../catalog/seed.json',
        ],
        {
          cwd: path.join(repositoryRoot, 'backend'),
          encoding: 'utf8',
          env: { ...process.env, GCLOUD_PROJECT: projectId },
        },
      );
      assert.equal(result.status, 0, result.stderr || result.stdout);
    };

    try {
      runSeed();
      const rootOf = (database) =>
        database
          .collection('spaces')
          .doc('issuenumber-public')
          .collection('ext')
          .doc('issuenumber');
      await environment.withSecurityRulesDisabled(async (context) => {
        const issue = rootOf(context.firestore())
          .collection('questions')
          .doc('country-ie')
          .collection('issues')
          .doc('housing-affordability');
        await issue.set({ supporters: 7, weightedScore: 16 }, { merge: true });
      });

      runSeed();

      await environment.withSecurityRulesDisabled(async (context) => {
        const root = rootOf(context.firestore());
        const issue = root
          .collection('questions')
          .doc('country-ie')
          .collection('issues')
          .doc('housing-affordability');
        const [categories, concepts, questions, preserved, axis] =
          await Promise.all([
            root.collection('categories').get(),
            root.collection('concepts').get(),
            root.collection('questions').get(),
            issue.get(),
            root.collection('questions').doc('world-country-actions').get(),
          ]);
        assert.equal(categories.size, catalog.categories.length);
        assert.equal(concepts.size, catalog.concepts.length);
        assert.equal(questions.size, catalog.questions.length);
        assert.equal(preserved.get('supporters'), 7);
        assert.equal(preserved.get('weightedScore'), 16);
        assert.deepEqual(axis.get('choiceSource'), {
          kind: 'predefined',
          entityType: 'country',
        });
        assert.equal(axis.get('answerTargetType'), 'country');
      });
    } finally {
      await environment.cleanup();
    }
  },
);
