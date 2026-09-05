import { loadPending, PendingAnswer, pendingKey } from './pending-answer';

function storageFor(value: unknown): Pick<Storage, 'getItem'> {
  return {
    getItem: (key: string) =>
      key === pendingKey ? JSON.stringify(value) : null,
  };
}

function pending(overrides: Partial<PendingAnswer> = {}): PendingAnswer {
  return {
    version: 1,
    operationId: 'operation-1',
    createdAt: Date.now(),
    questionId: 'question-1',
    issueId: 'issue-1',
    answerKind: 'category',
    attribution: 'anonymous',
    returnPath: '/issues/country/ireland',
    intent: 'consumer',
    ...overrides,
  };
}

describe('pending answer continuation', () => {
  it('preserves a free-form answer for a community question through auth', () => {
    const value = pending({
      issueId: undefined,
      title: 'Safer cycling routes',
      returnPath: '/questions/what-should-dublin-fix?lang=ga',
      languageCode: 'ga',
    });

    expect(loadPending(storageFor(value))).toEqual(value);
  });

  it('rejects expired or external return paths', () => {
    expect(
      loadPending(
        storageFor(
          pending({ createdAt: Date.now() - 25 * 60 * 60 * 1000 }),
        ),
      ),
    ).toBeNull();
    expect(
      loadPending(storageFor(pending({ returnPath: '//example.com' }))),
    ).toBeNull();
  });
});
