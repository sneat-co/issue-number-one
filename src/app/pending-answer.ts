export interface PendingAnswer {
  version: 1;
  operationId: string;
  createdAt: number;
  questionId: string;
  categoryId?: string;
  languageCode?: string;
  issueId?: string;
  title?: string;
  answerKind: 'category' | 'personal';
  attribution: 'anonymous' | 'authored';
  returnPath: string;
  intent: 'consumer' | 'business';
  requiredAuth?: boolean;
  verificationOrder?: readonly ('phone' | 'payment')[];
  verificationAttemptId?: string;
}
export const pendingKey = 'issuenumber.pending';
export function loadPending(
  storage: Pick<Storage, 'getItem'> = sessionStorage,
): PendingAnswer | null {
  try {
    const p = JSON.parse(storage.getItem(pendingKey) || 'null');
    if (
      p?.version !== 1 ||
      typeof p.operationId !== 'string' ||
      !/^[a-zA-Z0-9_-]{1,100}$/.test(p.questionId) ||
      (!p.returnPath?.startsWith('/issues/') &&
        !p.returnPath?.startsWith('/questions/')) ||
      p.returnPath.includes('//') ||
      p.returnPath.includes('\\') ||
      Date.now() - p.createdAt > 24 * 60 * 60 * 1000
    )
      return null;
    if (typeof p.issueId !== 'string' && typeof p.title !== 'string')
      return null;
    return p;
  } catch {
    return null;
  }
}
export function savePending(pending: PendingAnswer): void {
  sessionStorage.setItem(pendingKey, JSON.stringify(pending));
}
