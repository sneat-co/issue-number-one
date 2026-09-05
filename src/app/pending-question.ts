export type PendingQuestionChoiceSource =
  | {
      readonly kind: 'predefined';
      readonly entityType: 'country' | 'city' | 'currency';
    }
  | {
      readonly kind: 'custom';
      readonly options: readonly {
        readonly title: string;
        readonly description?: string;
      }[];
    }
  | { readonly kind: 'free' };

export interface PendingQuestion {
  readonly operationId: string;
  readonly createdAt: string;
  readonly title: string;
  readonly description: string;
  readonly choiceSource: PendingQuestionChoiceSource;
  readonly allowSuggestions: boolean;
  readonly sourceLanguage?: string;
  readonly requiredAuth?: boolean;
  readonly verificationOrder?: readonly ('phone' | 'payment')[];
  readonly verificationAttemptId?: string;
}

const pendingQuestionStorageKey = 'issuenumber.one.pending-question.v1';

export function loadPendingQuestion(
  storage: Pick<Storage, 'getItem'> = sessionStorage,
): PendingQuestion | undefined {
  try {
    const raw = storage.getItem(pendingQuestionStorageKey);
    if (!raw) return undefined;
    const value: unknown = JSON.parse(raw);
    return isPendingQuestion(value) ? value : undefined;
  } catch {
    return undefined;
  }
}

export function savePendingQuestion(
  question: PendingQuestion,
  storage: Pick<Storage, 'setItem'> = sessionStorage,
): void {
  storage.setItem(pendingQuestionStorageKey, JSON.stringify(question));
}

export function clearPendingQuestion(
  storage: Pick<Storage, 'removeItem'> = sessionStorage,
): void {
  storage.removeItem(pendingQuestionStorageKey);
}

function isPendingQuestion(value: unknown): value is PendingQuestion {
  if (!value || typeof value !== 'object') return false;
  const question = value as Partial<PendingQuestion>;
  const createdAt = Date.parse(question.createdAt || '');
  return (
    typeof question.operationId === 'string' &&
    typeof question.createdAt === 'string' &&
    typeof question.title === 'string' &&
    typeof question.description === 'string' &&
    typeof question.allowSuggestions === 'boolean' &&
    Number.isFinite(createdAt) &&
    Date.now() - createdAt <= 24 * 60 * 60 * 1000 &&
    isChoiceSource(question.choiceSource)
  );
}

function isChoiceSource(value: unknown): value is PendingQuestionChoiceSource {
  if (!value || typeof value !== 'object') return false;
  const source = value as Partial<PendingQuestionChoiceSource> & {
    entityType?: unknown;
    options?: unknown;
  };
  if (source.kind === 'free') return true;
  if (source.kind === 'predefined') {
    return ['country', 'city', 'currency'].includes(String(source.entityType));
  }
  if (source.kind !== 'custom' || !Array.isArray(source.options)) return false;
  return source.options.every(
    (option: unknown) =>
      !!option &&
      typeof option === 'object' &&
      typeof (option as { title?: unknown }).title === 'string' &&
      ((option as { description?: unknown }).description === undefined ||
        typeof (option as { description?: unknown }).description === 'string'),
  );
}
