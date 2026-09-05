import assert from 'node:assert/strict';
import { readFile } from 'node:fs/promises';
import test from 'node:test';

const seed = JSON.parse(
  await readFile(new URL('./seed.json', import.meta.url), 'utf8'),
);

const stableSlug = /^[a-z0-9]+(?:-[a-z0-9]+)*$/;

test('catalog has the expected version, taxonomy breadth, and candidate depth', () => {
  assert.equal(seed.schemaVersion, 1);
  assert.equal(seed.categories.length, 23);
  assert.ok(seed.concepts.length > 100);
  assert.ok(seed.questions.length >= 34);

  const conceptIds = new Set(seed.concepts.map(({ id }) => id));
  const categoryIds = new Set(seed.categories.map(({ id }) => id));
  for (const category of seed.categories) {
    assert.ok(category.defaultConceptIds.length >= 5, category.id);
    assert.equal(category.publication, 'published');
    assert.equal(category.indexable, true);
    assert.ok(['consumer', 'business'].includes(category.intent), category.id);
    if (category.parentCategoryId) {
      assert.ok(categoryIds.has(category.parentCategoryId), category.id);
      assert.notEqual(category.parentCategoryId, category.id);
    }
    for (const conceptId of category.defaultConceptIds) {
      assert.ok(conceptIds.has(conceptId), `${category.id}: ${conceptId}`);
    }
  }
  assert.equal(
    seed.categories.find(({ id }) => id === 'county').parentCategoryId,
    'country',
  );
  assert.equal(
    seed.categories.find(({ id }) => id === 'city').parentCategoryId,
    'county',
  );
});

test('all identities and slugs are stable and unique', () => {
  for (const [name, records] of Object.entries({
    categories: seed.categories,
    concepts: seed.concepts,
    questions: seed.questions,
  })) {
    const ids = records.map(({ id }) => id);
    assert.equal(new Set(ids).size, ids.length, `${name} IDs`);
    for (const record of records) {
      assert.match(record.id, /^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$/, record.id);
      assert.match(record.slug, stableSlug, record.slug);
    }
    const slugKeys = records.map((record) =>
      name === 'questions'
        ? `${record.categoryId}/${record.slug}`
        : record.slug,
    );
    assert.equal(new Set(slugKeys).size, records.length, `${name} slugs`);
  }
});

test('question references resolve and every category has a concrete question', () => {
  const categoryIds = new Set(seed.categories.map(({ id }) => id));
  const conceptIds = new Set(seed.concepts.map(({ id }) => id));
  const questionIds = new Set(seed.questions.map(({ id }) => id));
  const representedCategories = new Set();

  for (const question of seed.questions) {
    if (question.categoryId) {
      assert.ok(categoryIds.has(question.categoryId), question.id);
      representedCategories.add(question.categoryId);
    } else {
      assert.ok(question.slug, question.id);
    }
    if (question.choiceSource?.kind === 'predefined') {
      assert.ok(
        ['country', 'city', 'currency'].includes(
          question.choiceSource.entityType,
        ),
        question.id,
      );
      assert.equal(question.answerTargetType, question.choiceSource.entityType);
      assert.equal(question.allowSuggestions, false);
      assert.equal(question.conceptIds.length, 0);
    } else {
      assert.ok(question.conceptIds.length >= 5, question.id);
    }
    for (const conceptId of question.conceptIds) {
      assert.ok(conceptIds.has(conceptId), `${question.id}: ${conceptId}`);
    }
    for (const relatedId of question.relatedQuestionIds) {
      assert.ok(questionIds.has(relatedId), `${question.id}: ${relatedId}`);
      assert.notEqual(relatedId, question.id);
    }
    if (question.parentQuestionId) {
      assert.ok(questionIds.has(question.parentQuestionId), question.id);
      assert.notEqual(question.parentQuestionId, question.id);
    }
  }

  assert.deepEqual(representedCategories, categoryIds);
});

test('country-axis example uses canonical entities rather than fake issue concepts', () => {
  const question = seed.questions.find(
    ({ id }) => id === 'world-country-actions',
  );
  assert.equal(question.choiceSource.kind, 'predefined');
  assert.equal(question.choiceSource.entityType, 'country');
  assert.equal(question.answerTargetType, 'country');
  assert.equal(question.categoryId, undefined);
  assert.equal(question.slug, 'country-government-actions');
  assert.deepEqual(question.conceptIds, []);
  assert.match(question.description, /government/);
  assert.match(question.description, /not on the identity of its people/);
});

test('required country, county, and city hierarchy keeps separate answer slots', () => {
  const questions = new Map(
    seed.questions.map((question) => [question.id, question]),
  );
  for (const countryId of ['IE', 'GB', 'US', 'DE', 'FR', 'PL']) {
    assert.ok(
      seed.questions.some(
        ({ categoryId, scope }) =>
          categoryId === 'country' &&
          scope.type === 'country' &&
          scope.id === countryId,
      ),
      countryId,
    );
  }

  const places = [
    ['dublin', 'IE-D', 'IE-D-DUBLIN'],
    ['cork', 'IE-C', 'IE-C-CORK'],
    ['limerick', 'IE-LK', 'IE-LK-LIMERICK'],
    ['galway', 'IE-G', 'IE-G-GALWAY'],
  ];
  for (const [place, countyScopeId, cityScopeId] of places) {
    const county = questions.get(`county-${place}`);
    assert.equal(county.categoryId, 'county');
    assert.equal(county.scope.type, 'county');
    assert.equal(county.scope.id, countyScopeId);
    assert.equal(county.scope.parentId, 'IE');
    assert.equal(county.parentQuestionId, undefined);

    const cityId = `city-${place}`;
    const city = questions.get(cityId);
    assert.equal(city.categoryId, 'city');
    assert.equal(city.scope.type, 'city');
    assert.equal(city.scope.id, cityScopeId);
    assert.equal(city.scope.parentId, countyScopeId);
    assert.equal(city.parentQuestionId, undefined);
  }
});

test('seed contains candidates only and never fabricates response data', () => {
  const forbidden = new Set([
    'answerCount',
    'answers',
    'count',
    'percentage',
    'rank',
    'respondentCount',
    'supporterCount',
    'supporters',
    'votes',
  ]);

  const visit = (value, path = '$') => {
    if (Array.isArray(value)) {
      value.forEach((item, index) => visit(item, `${path}[${index}]`));
      return;
    }
    if (!value || typeof value !== 'object') return;
    for (const [key, child] of Object.entries(value)) {
      assert.ok(!forbidden.has(key), `${path}.${key}`);
      visit(child, `${path}.${key}`);
    }
  };

  visit(seed);
});
