import { expect, test } from '@playwright/test';

test('question author can fill before auth and keeps the draft at sign-in', async ({
  page,
}) => {
  await page.goto('/questions/new');

  await expect(
    page.getByRole('heading', { name: 'Create a #1 issue question' }),
  ).toBeVisible();
  await expect(page.locator('ion-radio[aria-checked="true"]')).toHaveCount(0);

  await page
    .locator('ion-input[formcontrolname="title"] input')
    .fill('What should Dublin improve first');
  await page
    .locator('ion-textarea[formcontrolname="description"] textarea')
    .fill(
      'Choose the change with the greatest practical effect for residents.',
    );
  await page
    .getByRole('radio', { name: 'A focused list for this question' })
    .click();
  await page
    .locator('ion-textarea[formcontrolname="customOptions"] textarea')
    .fill('Safer crossings\nMore frequent buses\nSafer crossings');
  await page.getByRole('button', { name: 'Submit question' }).click();

  await expect(page).toHaveURL(/\/login/);
  const stored = await page.evaluate(() =>
    JSON.parse(
      sessionStorage.getItem('issuenumber.one.pending-question.v1') || 'null',
    ),
  );
  expect(stored).toMatchObject({
    title: 'What should Dublin improve first?',
    requiredAuth: true,
    choiceSource: {
      kind: 'custom',
      options: [{ title: 'Safer crossings' }, { title: 'More frequent buses' }],
    },
  });
});

test('verification has no default and keeps one random order for the attempt', async ({
  page,
}) => {
  await page.goto('/');
  await page.evaluate(() => {
    sessionStorage.setItem(
      'issuenumber.pending',
      JSON.stringify({
        version: 1,
        operationId: 'answer-operation-1',
        createdAt: Date.now(),
        questionId: 'ireland',
        categoryId: 'country',
        title: 'Housing affordability',
        answerKind: 'category',
        attribution: 'anonymous',
        returnPath: '/issues/country/ireland',
        intent: 'consumer',
      }),
    );
  });

  await page.goto('/verify');
  await expect(
    page.getByRole('heading', {
      name: 'A little friction. A fairer picture.',
    }),
  ).toBeVisible();
  await expect(page.getByText('Neither option is preselected.')).toBeVisible();
  const firstOrder = await page.locator('.options h2').allTextContents();
  const normalizedOrder = firstOrder.map((value) => value.trim());
  expect([...normalizedOrder].sort()).toEqual(
    ['Support with €1', 'Verify by phone'].sort(),
  );
  await expect(page.locator('sneat-phone-verification')).toHaveCount(0);

  await page.reload();
  await expect(page.locator('.options h2')).toHaveText(normalizedOrder);
  const stored = await page.evaluate(() =>
    JSON.parse(sessionStorage.getItem('issuenumber.pending') || 'null'),
  );
  expect(stored.title).toBe('Housing affordability');
  expect(stored.verificationOrder).toHaveLength(2);
  expect(stored.verificationAttemptId).toBeTruthy();
});
