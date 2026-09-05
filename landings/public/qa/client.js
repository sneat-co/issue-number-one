// No auth SDK is loaded until the visitor asks to participate.
const pendingKey = 'issuenumber.pending';
const status = document.querySelector('#action-status');
const view = document.querySelector('[data-view]');
const languageCode =
  new URL(location.href).searchParams.get('lang') ||
  document.documentElement.lang ||
  'en';
function track(name, extra = {}) {
  const dimensions = {
    category_id: view?.dataset.category,
    question_id: view?.dataset.question,
    issue_id: view?.dataset.issue,
    intent: view?.dataset.intent,
    language_code: languageCode,
    ...extra,
  };
  // Reuse the configured site transports when present. Never include issue text.
  window.gtag?.('event', name, dimensions);
  window.posthog?.capture(name, dimensions);
}
function begin(data) {
  try {
    const returnPath =
      languageCode !== 'en' && !data.returnPath.includes('?')
        ? `${data.returnPath}?lang=${encodeURIComponent(languageCode)}`
        : data.returnPath;
    sessionStorage.setItem(
      pendingKey,
      JSON.stringify({
        version: 1,
        operationId: crypto.randomUUID(),
        createdAt: Date.now(),
        answerKind: 'category',
        attribution: 'anonymous',
        ...data,
        languageCode,
        returnPath,
      }),
    );
    location.assign('/answer');
  } catch {
    if (status)
      status.textContent =
        'Your browser could not preserve this answer. Enable session storage and try again; your text remains here.';
  }
}
if (view) track(view.dataset.view);
for (const button of document.querySelectorAll('[data-answer]'))
  button.addEventListener('click', () => {
    track('for_me_too_clicked', { issue_id: button.dataset.answer });
    begin({
      questionId: button.dataset.question,
      categoryId: button.dataset.category,
      issueId: button.dataset.answer,
      returnPath: button.dataset.return,
      intent: button.dataset.intent,
    });
  });
for (const form of document.querySelectorAll('[data-compose]')) {
  form
    .querySelector('input')
    .addEventListener('input', () => track('freeform_started'), { once: true });
  form.addEventListener('submit', (event) => {
    event.preventDefault();
    if (!form.reportValidity()) return;
    const title = new FormData(form).get('title').trim().replace(/\s+/g, ' ');
    if (title.length < 3) return;
    track('freeform_submitted');
    begin({
      questionId: form.dataset.question,
      categoryId: form.dataset.category,
      title,
      returnPath: form.dataset.return,
      intent: form.dataset.intent,
      attribution:
        new FormData(form).get('attribution') === 'authored'
          ? 'authored'
          : 'anonymous',
    });
  });
}
for (const button of document.querySelectorAll('[data-share]'))
  button.addEventListener('click', async () => {
    const url = new URL(button.dataset.share, 'https://issuenumber.one').href;
    const shareUrl = new URL(url);
    if (languageCode !== 'en') shareUrl.searchParams.set('lang', languageCode);
    try {
      if (navigator.share)
        await navigator.share({
          title: button.dataset.shareTitle,
          text: button.dataset.shareTitle,
          url: shareUrl.href,
        });
      else {
        await navigator.clipboard.writeText(shareUrl.href);
        if (status)
          status.textContent =
            'Link copied. Share it with someone whose opinion matters to you.';
      }
      track(button.dataset.event);
    } catch (error) {
      if (error.name !== 'AbortError' && status)
        status.textContent =
          'Could not share automatically. Copy this page’s address to share it.';
    }
  });
for (const a of document.querySelectorAll('a[data-event]'))
  a.addEventListener('click', () => track(a.dataset.event));
try {
  const recent = JSON.parse(
    sessionStorage.getItem('issuenumber.saved') || 'null',
  );
  if (
    recent?.questionId === view?.dataset.question &&
    Date.now() - recent.createdAt < 86400000
  ) {
    document
      .querySelectorAll('[data-conversion]')
      .forEach((el) => (el.hidden = false));
    if (status)
      status.textContent = recent.personalTop
        ? 'Your personal #1 is saved.'
        : 'Your category choice is saved.';
  }
} catch {
  /* Public content works without session storage. */
}
