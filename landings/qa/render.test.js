import { test } from 'node:test';
import assert from 'node:assert/strict';
import {
  page,
  aggregate,
  questionPage,
  questionsPage,
  categoryPath,
  questionPath,
  sitemap,
  escapeHTML,
} from './render.js';
const cat = {
  id: 'country',
  slug: 'country',
  name: 'Country',
  question: 'What is the top issue in your country?',
  seoDescription: 'Explore what matters to people across countries.',
  description: 'Choose a country and see its priorities.',
  publication: 'published',
  indexable: true,
  defaultConceptIds: ['housing'],
  intent: 'consumer',
};
const q = {
  id: 'ie',
  slug: 'ireland',
  categoryId: 'country',
  title: 'What is the top issue in Ireland?',
  description: 'Housing, healthcare and other concerns in Ireland.',
  scope: { id: 'IE', name: 'Ireland' },
  publication: 'published',
  indexable: true,
  conceptIds: ['a', 'b', 'c', 'd', 'e'],
};
const issues = Array.from({ length: 5 }, (_, i) => ({
  id: 'i' + i,
  slug: 'issue-' + i,
  title: 'Concern ' + i,
  description: 'A meaningful issue description for the community.',
  supporters: i === 0 ? 2 : 1,
  weightedScore: i === 0 ? 2 : i === 1 ? 10 : 1,
  personalTopSupporters: i === 1 ? 1 : 0,
  status: 'published',
  conceptId: 'c' + i,
}));
const catalog = { categories: [cat], questions: [q], concepts: [] };
test('question HTML has crawlable weighted ranking with truthful real counts', () => {
  const html = questionPage(catalog, cat, {
    question: q,
    issues,
    totalRespondents: 6,
  });
  assert.equal((html.match(/<h1>/g) || []).length, 1);
  assert.ok(html.includes('<title>What is the top issue in Ireland?'));
  assert.ok(
    html.includes(
      'rel="canonical" href="https://issuenumber.one/issues/country/ireland"',
    ),
  );
  assert.ok(html.indexOf('Concern 1') < html.indexOf('Concern 0'));
  assert.ok(html.includes('10 priority points · 1 personal #1 choices'));
  assert.ok(html.includes('2 supporters'));
  assert.ok(html.includes('For me too'));
  assert.ok(!html.includes('firebase'));
});
test('empty candidates do not claim measured rank or winner', () => {
  const html = questionPage(catalog, cat, {
    question: q,
    issues: issues.map((i) => ({
      ...i,
      supporters: 0,
      weightedScore: 0,
      personalTopSupporters: 0,
    })),
    totalRespondents: 0,
  });
  assert.ok(html.includes('starting points, not measured rankings'));
  assert.ok(!html.includes('class="rank">#'));
});
test('question pages hide unset update timestamps and show real dates', () => {
  const emptyDateHtml = questionPage(catalog, cat, {
    question: q,
    issues,
    totalRespondents: 0,
    updatedAt: '0001-01-01T00:00:00Z',
  });
  assert.ok(!emptyDateHtml.includes('Updated 1/1/1'));
  assert.ok(!emptyDateHtml.includes(' · Updated '));

  const realDateHtml = questionPage(catalog, cat, {
    question: q,
    issues,
    totalRespondents: 6,
    updatedAt: '2026-09-05T16:30:00Z',
  });
  assert.ok(realDateHtml.includes(' · Updated 5/9/2026'));
});
test('candidate issue pages are noindex and hidden text is not rendered', () => {
  const snapshot = {
    question: q,
    issues: [
      ...issues,
      {
        id: 'hidden',
        slug: 'hidden',
        title: 'PRIVATE TEXT',
        status: 'pending',
      },
    ],
    totalRespondents: 6,
  };
  assert.equal(questionPage(catalog, cat, snapshot, 'hidden'), null);
  assert.ok(!questionPage(catalog, cat, snapshot).includes('PRIVATE TEXT'));
  assert.ok(
    questionPage(catalog, cat, snapshot, 'issue-0').includes(
      'content="noindex,follow"',
    ),
  );
});
test('aggregation keeps people counts separate from weighting and excludes zero winners', () => {
  const a = aggregate([
    { issues, totalRespondents: 6 },
    {
      issues: issues.map((i) => ({ ...i, supporters: 0, weightedScore: 0 })),
      totalRespondents: 0,
    },
  ]);
  assert.equal(a.total, 6);
  assert.equal(a.participating, 1);
  assert.equal(a.issues[0].id, 'c1');
  assert.equal(a.issues[0].supporters, 1);
  assert.equal(a.issues[0].weightedScore, 10);
  assert.equal(a.issues[0].leadingScopes, 1);
  assert.equal(a.issues.find((i) => i.id === 'c0').leadingScopes, 0);
});
test('country county and city paths nest independently of category membership', () => {
  const county = {
    ...cat,
    id: 'county',
    slug: 'county',
    parentCategoryId: 'country',
  };
  const city = { ...cat, id: 'city', slug: 'city', parentCategoryId: 'county' };
  const co = {
    ...q,
    id: 'lk',
    slug: 'limerick',
    categoryId: 'county',
    scope: { id: 'IE-LK', parentId: 'IE' },
  };
  const ci = {
    ...q,
    id: 'limerick',
    slug: 'limerick',
    categoryId: 'city',
    scope: { id: 'IE-LK-LIMERICK', parentId: 'IE-LK' },
  };
  const c = {
    ...catalog,
    categories: [cat, county, city],
    questions: [q, co, ci],
  };
  assert.equal(categoryPath(city, c), '/issues/country/county/city');
  assert.equal(
    questionPath(city, ci, c),
    '/issues/country/ireland/county/limerick/city/limerick',
  );
  assert.ok(
    sitemap(c).includes(
      '/issues/country/ireland/county/limerick/city/limerick',
    ),
  );
});
test('generated user question canonicals preserve founder-required trailing question mark', () => {
  const html = page({
    title: 'A question?',
    description: 'Question description',
    path: '/questions/a-question',
    body: '<h1>A question?</h1>',
  });
  assert.ok(
    html.includes(
      'rel="canonical" href="https://issuenumber.one/questions/a-question?"',
    ),
  );
  assert.ok(
    html.includes(
      'property="og:url" content="https://issuenumber.one/questions/a-question?"',
    ),
  );
});

test('translated questions use lang canonical and expose language alternatives', () => {
  const translated = questionPage(catalog, cat, {
    question: { ...q, availableLanguages: ['en', 'ru'] },
    issues,
    totalRespondents: 6,
    languageCode: 'ru',
    languageRespondents: 2,
  });
  assert.ok(
    translated.includes(
      'rel="canonical" href="https://issuenumber.one/issues/country/ireland?lang=ru"',
    ),
  );
  assert.ok(translated.includes('<html lang="ru">'));
  assert.ok(translated.includes('hreflang="en"'));
  assert.ok(translated.includes('hreflang="ru"'));
  assert.ok(
    translated.includes('6</strong> people have answered across all languages'),
  );
  assert.ok(translated.includes('2</strong> answered in Russian'));
  assert.ok(translated.includes('answers in Russian'));
  assert.ok(translated.includes('Russian priority points'));
});

test('community discovery only lists published user-created questions', () => {
  const html = questionsPage({
    ...catalog,
    questions: [
      ...catalog.questions,
      {
        id: 'public-question',
        slug: 'what-matters',
        title: 'What matters?',
        publication: 'published',
        indexable: true,
      },
      {
        id: 'private-question',
        slug: 'private-draft',
        title: 'PRIVATE DRAFT',
        publication: 'pending',
      },
    ],
  });
  assert.ok(html.includes('/questions/what-matters?'));
  assert.ok(!html.includes('PRIVATE DRAFT'));
  assert.ok(html.includes('/questions/new'));
});
test('sitemap trusts reviewed indexability without requiring seed-only concept ids', () => {
  const reviewed = {
    ...q,
    id: 'reviewed-custom',
    slug: 'reviewed-custom',
    conceptIds: undefined,
  };
  const xml = sitemap({ ...catalog, questions: [reviewed] });
  assert.ok(xml.includes('/issues/country/reviewed-custom'));
});
test('escaping prevents user text breaking attributes or structured script', () => {
  assert.equal(escapeHTML('<"&'), '&lt;&quot;&amp;');
  const html = page({
    title: '</script><script>alert(1)</script>',
    description: '" unsafe',
    path: '/test',
    body: '<h1>Test</h1>',
  });
  assert.ok(!html.includes('<script>alert(1)</script>'));
});
