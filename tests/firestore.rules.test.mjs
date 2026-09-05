import { readFile } from 'node:fs/promises';
import { after, before, test } from 'node:test';
import assert from 'node:assert/strict';

import {
  assertFails,
  initializeTestEnvironment,
} from '@firebase/rules-unit-testing';
import { doc, getDoc, setDoc, updateDoc } from 'firebase/firestore';

const projectId = 'issuenumber-rules-test';
const root = 'spaces/issuenumber-public/ext/issuenumber';
let environment;

before(async () => {
  const rules = await readFile(new URL('../firestore.rules', import.meta.url), 'utf8');
  environment = await initializeTestEnvironment({
    projectId,
    firestore: { rules },
  });
  await environment.withSecurityRulesDisabled(async (context) => {
    const db = context.firestore();
    await setDoc(doc(db, `${root}/categories/country`), {
      id: 'country', status: 'published', title: 'Country',
    });
    await setDoc(doc(db, `${root}/questions/ireland`), {
      id: 'ireland', status: 'published', indexable: true,
      creatorUID: 'editor-private', totalRespondents: 1,
    });
    await setDoc(doc(db, `${root}/questions/ireland/issues/housing`), {
      id: 'housing', status: 'published', supporters: 1,
      creatorUID: 'issue-author-private',
    });
    await setDoc(doc(db, `${root}/questions/ireland/answers/alice`), {
      questionId: 'ireland', issueId: 'housing', languageCode: 'en',
    });
    await setDoc(doc(db, `${root}/verification/alice`), {
      eligible: true, chargeId: 'private-charge-reference',
    });
    await setDoc(doc(db, `${root}/questions/pending-question`), {
      id: 'pending-question', status: 'pending', creatorUID: 'alice',
    });
  });
});

after(async () => {
  await environment?.cleanup();
});

async function denied(operation, label) {
  await assert.doesNotReject(assertFails(operation), label);
}

test('anonymous clients cannot read raw public or private projections', async () => {
  const db = environment.unauthenticatedContext().firestore();
  await denied(getDoc(doc(db, `${root}/categories/country`)), 'raw catalog read');
  await denied(getDoc(doc(db, `${root}/questions/ireland`)), 'raw published question read');
  await denied(getDoc(doc(db, `${root}/questions/ireland/issues/housing`)), 'raw issue read');
  await denied(getDoc(doc(db, `${root}/questions/ireland/answers/alice`)), 'answer read');
  await denied(getDoc(doc(db, `${root}/verification/alice`)), 'verification read');
});

test('authentication does not expose raw creator, answer, or verification data', async () => {
  const own = environment.authenticatedContext('alice').firestore();
  const other = environment.authenticatedContext('bob').firestore();
  await denied(getDoc(doc(own, `${root}/questions/pending-question`)), 'creator pending question read');
  await denied(getDoc(doc(own, `${root}/questions/ireland/answers/alice`)), 'own answer read');
  await denied(getDoc(doc(other, `${root}/questions/ireland/answers/alice`)), 'other answer read');
  await denied(getDoc(doc(own, `${root}/verification/alice`)), 'own verification read');
  await denied(getDoc(doc(other, `${root}/verification/alice`)), 'other verification read');
});

test('anonymous and authenticated clients cannot mutate backend-owned state', async () => {
  const anonymous = environment.unauthenticatedContext().firestore();
  const authenticated = environment.authenticatedContext('alice').firestore();
  await denied(setDoc(doc(anonymous, `${root}/questions/ireland/answers/anonymous`), {
    questionId: 'ireland', issueId: 'housing',
  }), 'anonymous answer write');
  await denied(setDoc(doc(authenticated, `${root}/questions/ireland/answers/alice`), {
    questionId: 'ireland', issueId: 'housing',
  }), 'authenticated own answer write');
  await denied(updateDoc(doc(authenticated, `${root}/questions/ireland/issues/housing`), {
    supporters: 999,
  }), 'counter update');
  await denied(updateDoc(doc(authenticated, `${root}/questions/ireland/issues/housing`), {
    status: 'hidden',
  }), 'issue moderation update');
  await denied(updateDoc(doc(authenticated, `${root}/questions/ireland`), {
    status: 'hidden', indexable: false,
  }), 'question moderation update');
  await denied(setDoc(doc(authenticated, `${root}/verification/alice`), {
    eligible: true,
  }), 'self-granted payment eligibility');
});
