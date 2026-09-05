import {
  clearPendingQuestion,
  loadPendingQuestion,
  PendingQuestion,
  savePendingQuestion,
} from './pending-question';

class MemoryStorage {
  private readonly values = new Map<string, string>();
  getItem(key: string): string | null {
    return this.values.get(key) ?? null;
  }
  setItem(key: string, value: string): void {
    this.values.set(key, value);
  }
  removeItem(key: string): void {
    this.values.delete(key);
  }
}

describe('pending question continuation', () => {
  const question: PendingQuestion = {
    operationId: 'operation-1',
    createdAt: '2026-09-05T12:00:00.000Z',
    title: 'What should our neighbourhood fix first?',
    description:
      'Choose the issue that would make the largest difference locally.',
    choiceSource: {
      kind: 'custom',
      options: [{ title: 'Safer crossings' }, { title: 'More trees' }],
    },
    allowSuggestions: true,
  };

  it('survives an auth-style storage round trip and can be cleared after posting', () => {
    const storage = new MemoryStorage();
    savePendingQuestion(question, storage);
    expect(loadPendingQuestion(storage)).toEqual(question);
    clearPendingQuestion(storage);
    expect(loadPendingQuestion(storage)).toBeUndefined();
  });

  it('does not resume malformed or unrelated session data', () => {
    const storage = new MemoryStorage();
    storage.setItem(
      'issuenumber.one.pending-question.v1',
      JSON.stringify({ title: 'missing operation' }),
    );
    expect(loadPendingQuestion(storage)).toBeUndefined();
  });
});
