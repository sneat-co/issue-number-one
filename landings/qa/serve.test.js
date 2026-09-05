import test from 'node:test';
import assert from 'node:assert/strict';
import { serveQA } from './serve.js';

const catalog = {
  categories: [],
  concepts: [],
  questions: [
    {
      id: 'community-1',
      slug: 'which-country',
      title: 'Which country concerns you most?',
      description: 'A public country-choice priority question.',
      publication: 'published',
      indexable: true,
      choiceSource: { kind: 'predefined', entityType: 'country' },
      availableLanguages: ['en', 'ru'],
    },
  ],
};

test('serves translated community question from the public backend', async (t) => {
  const previousFetch = globalThis.fetch;
  t.after(() => {
    globalThis.fetch = previousFetch;
  });
  globalThis.fetch = async (request) => {
    const url = new URL(request instanceof Request ? request.url : request);
    if (url.pathname.endsWith('/catalog')) return Response.json(catalog);
    assert.equal(url.searchParams.get('slug'), 'which-country');
    assert.equal(url.searchParams.get('lang'), 'ru');
    return Response.json({
      question: catalog.questions[0],
      translation: {
        title: 'Какая страна беспокоит вас больше всего?',
        description: 'Публичный вопрос о приоритете среди стран.',
        sourceLanguage: 'en',
        machineTranslated: true,
      },
      issues: [],
      totalRespondents: 0,
      languageCode: 'ru',
      languageRespondents: 0,
    });
  };
  const response = await serveQA(
    new Request('https://issuenumber.one/questions/which-country?lang=ru'),
    { ISSUENUMBER_API_ORIGIN: 'https://api.example' },
  );
  assert.equal(response.status, 200);
  const html = await response.text();
  assert.ok(html.includes('Какая страна'));
  assert.ok(html.includes('/questions/which-country?lang=ru'));
  assert.ok(html.includes('<html lang="ru">'));
});

test('leaves question creation and owner preview to the authenticated app', async () => {
  assert.equal(
    await serveQA(new Request('https://issuenumber.one/questions/new'), {}),
    null,
  );
  assert.equal(
    await serveQA(
      new Request('https://issuenumber.one/questions/pending?preview=1'),
      {},
    ),
    null,
  );
});
