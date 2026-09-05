// Public HTML is independent of the authentication bundle.
export const ORIGIN = 'https://issuenumber.one';
export const escapeHTML = (value = '') =>
  String(value).replace(
    /[&<>"']/g,
    (c) =>
      ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' })[
        c
      ],
  );
const e = escapeHTML;
const languageNames = {
  de: 'German',
  en: 'English',
  es: 'Spanish',
  fr: 'French',
  ga: 'Irish',
  it: 'Italian',
  nl: 'Dutch',
  pl: 'Polish',
  pt: 'Portuguese',
  ru: 'Russian',
  uk: 'Ukrainian',
};
const languageName = (code) => languageNames[code] || code.toUpperCase();
export function categoryPath(category, catalog) {
  const parent = catalog?.categories.find(
    (c) => c.id === category.parentCategoryId,
  );
  return parent
    ? categoryPath(parent, catalog) + '/' + category.slug
    : '/issues/' + category.slug;
}
export function questionPath(category, question, catalog, visited = new Set()) {
  if (!question.categoryId) return '/questions/' + question.slug;
  if (visited.has(question.id)) throw new Error('Cyclic scope hierarchy');
  visited.add(question.id);
  const parent = catalog?.questions.find(
    (q) => q.scope.id === question.scope.parentId,
  );
  const parentCategory =
    parent && catalog.categories.find((c) => c.id === parent.categoryId);
  return parentCategory
    ? questionPath(parentCategory, parent, catalog, visited) +
        '/' +
        category.slug +
        '/' +
        question.slug
    : `/issues/${category.slug}/${question.slug}`;
}
export const isPublished = (item) =>
  item.publication === 'published' &&
  (!item.visibility || item.visibility === 'public');
export const isIndexable = (item) =>
  isPublished(item) && item.indexable === true;
const link = (href, text) => `<a href="${e(href)}">${e(text)}</a>`;
function basePage({
  title,
  description,
  path,
  body,
  crumbs = [],
  indexable = true,
  items = [],
  language = 'en',
  availableLanguages = [],
}) {
  const languageQuery =
    language === 'en' ? '' : `lang=${encodeURIComponent(language)}`;
  const canonical =
    ORIGIN +
    path +
    (languageQuery
      ? `?${languageQuery}`
      : path.startsWith('/questions/') && !path.endsWith('?')
        ? '?'
        : '');
  const structured = [
    {
      '@context': 'https://schema.org',
      '@type': items.length ? 'CollectionPage' : 'WebPage',
      name: title,
      url: canonical,
    },
    {
      '@context': 'https://schema.org',
      '@type': 'BreadcrumbList',
      itemListElement: [{ name: 'Home', path: '/' }, ...crumbs].map((c, i) => ({
        '@type': 'ListItem',
        position: i + 1,
        name: c.name,
        item: ORIGIN + c.path,
      })),
    },
  ];
  if (items.length)
    structured.push({
      '@context': 'https://schema.org',
      '@type': 'ItemList',
      itemListElement: items.map((x, i) => ({
        '@type': 'ListItem',
        position: i + 1,
        name: x.title,
        url: ORIGIN + x.path,
      })),
    });
  return `<!doctype html><html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>${e(title)} | IssueNumber.one</title><meta name="description" content="${e(description)}"><link rel="canonical" href="${e(canonical)}"><meta name="robots" content="${indexable ? 'index,follow' : 'noindex,follow'}"><meta property="og:type" content="website"><meta property="og:site_name" content="IssueNumber.one"><meta property="og:title" content="${e(title)}"><meta property="og:description" content="${e(description)}"><meta property="og:url" content="${e(canonical)}"><meta property="og:image" content="${ORIGIN}/og.png"><meta name="twitter:card" content="summary_large_image"><link rel="stylesheet" href="/qa/style.css"><script type="application/ld+json">${JSON.stringify(structured).replace(/</g, '\\u003c')}</script><script async src="https://www.googletagmanager.com/gtag/js?id=G-TYBDTV738R"></script><script>window.dataLayer=window.dataLayer||[];window.gtag=function(){dataLayer.push(arguments)};gtag("js",new Date());gtag("config","G-TYBDTV738R",{send_page_view:false});</script><script type="module" src="/qa/client.js"></script></head><body><a class="skip" href="#main">Skip to content</a><header><a class="brand" href="/">IssueNumber<span>.one</span></a><nav aria-label="Main">${link('/issues', 'Explore issues')} ${link('/for-work/', 'For work')} ${link('/login', 'Sign in')}</nav></header><main id="main"><nav class="breadcrumbs" aria-label="Breadcrumb">${link('/', 'Home')}${crumbs.map((c) => ' / ' + link(c.path, c.name)).join('')}</nav>${body}</main><footer><p>One choice per question. One personal #1 across them all.</p><p>These are the views of people who choose to participate, not a representative survey.</p>${link('/issues', 'Find a question')} · ${link('/privacy', 'Privacy')}</footer></body></html>`;
}

export function page(props) {
  const language = props.language || 'en';
  const alternates = (props.availableLanguages || [])
    .map(
      (code) =>
        `<link rel="alternate" hreflang="${e(code)}" href="${e(
          `${ORIGIN}${props.path}${code === 'en' ? '?' : `?lang=${encodeURIComponent(code)}`}`,
        )}">`,
    )
    .join('');
  return basePage(props)
    .replace('<html lang="en">', `<html lang="${e(language)}">`)
    .replace('<meta name="robots"', `${alternates}<meta name="robots"`);
}

export function discovery(catalog, path = '/') {
  const categories = catalog.categories.filter(isPublished);
  return page({
    title: 'What is the #1 issue?',
    description:
      'Explore what matters most in your country, city, work and everyday life. Choose your #1 issue and invite people whose views matter to you.',
    path,
    body: `<section class="hero"><p class="eyebrow">What matters most?</p><h1>So many problems.<br>What is your <em>#1?</em></h1><p>Start with a place, a community or a part of life. See what others prioritise, then add your voice.</p><a class="button" href="#categories">Find your question ↓</a></section><section id="categories"><h2>A different question for every part of life</h2><div class="grid">${categories.map((c) => `<a class="category" href="${e(categoryPath(c, catalog))}"><h3>${e(c.name)}</h3><p>${e(c.question)}</p><span>Explore →</span></a>`).join('')}</div></section>`,
  });
}

export function questionsPage(catalog) {
  const questions = catalog.questions.filter(
    (question) => !question.categoryId && isPublished(question),
  );
  return page({
    title: 'Questions from the community',
    description:
      'Explore public priority questions, or ask a focused question with predefined, custom or free answers.',
    path: '/questions',
    items: questions.map((question) => ({
      title: question.title,
      path: `/questions/${question.slug}?`,
    })),
    body: `<section class="hero"><p class="eyebrow">Community questions</p><h1>What should be #1?</h1><p>Every question asks people to choose one current priority. Answers can come from a trusted entity list, custom choices, or free suggestions.</p><a class="button" href="/questions/new">Ask a question</a></section><section><h2>Published questions</h2>${questions.length ? `<ul class="ranking">${questions.map((question) => `<li>${link(`/questions/${question.slug}?`, question.title)}</li>`).join('')}</ul>` : '<p>No community questions are published yet.</p>'}</section>`,
  });
}
export function aggregate(snapshots) {
  const concepts = new Map();
  let total = 0,
    participating = 0;
  for (const snapshot of snapshots) {
    total += snapshot.totalRespondents;
    if (snapshot.totalRespondents > 0) participating++;
    const eligible = snapshot.issues.filter((i) => i.status === 'published');
    const max = Math.max(
      0,
      ...eligible.map((i) => i.weightedScore ?? i.supporters),
    );
    const won = new Set();
    for (const issue of eligible) {
      if (!issue.conceptId) continue;
      const item = concepts.get(issue.conceptId) || {
        id: issue.conceptId,
        title: issue.title,
        supporters: 0,
        weightedScore: 0,
        leadingScopes: 0,
      };
      item.supporters += issue.supporters;
      item.weightedScore += issue.weightedScore ?? issue.supporters;
      if (
        max > 0 &&
        (issue.weightedScore ?? issue.supporters) === max &&
        !won.has(issue.conceptId)
      ) {
        item.leadingScopes++;
        won.add(issue.conceptId);
      }
      concepts.set(item.id, item);
    }
  }
  return {
    total,
    participating,
    issues: [...concepts.values()].sort(
      (a, b) =>
        (b.weightedScore ?? b.supporters) - (a.weightedScore ?? a.supporters) ||
        a.id.localeCompare(b.id),
    ),
  };
}
export function categoryPage(
  catalog,
  category,
  snapshots,
  unavailable = false,
) {
  const children = catalog.questions.filter(
    (q) => q.categoryId === category.id && isPublished(q),
  );
  const sums = aggregate(snapshots);
  const candidateIds = new Set(
    category.defaultConceptIds || category.conceptIds || [],
  );
  const candidates = catalog.concepts.filter((c) => candidateIds.has(c.id));
  const ranking = sums.total
    ? `<h2>Responses across participating scopes</h2><p>${sums.total} current answers across ${sums.participating} participating scopes. Each person has one choice per concrete question and can answer multiple scopes. Priority points rank the issues; supporter percentages count people equally. Larger scopes have more influence here.</p><ol class="ranking">${sums.issues.map((i) => `<li><strong>${e(i.title)}</strong><span>${i.weightedScore} priority points · ${i.supporters} answers · ${Math.round((i.supporters / sums.total) * 100)}%</span><small>Joint or sole #1 in ${i.leadingScopes} participating scopes</small></li>`).join('')}</ol>`
    : `<h2>Issues to consider</h2><p>These are curated candidate options, not a survey ranking. Choose a scope to give your answer.</p><ul class="ranking">${candidates.map((i) => `<li><strong>${e(i.title)}</strong><p>${e(i.description)}</p></li>`).join('')}</ul>`;
  return page({
    title: category.question,
    description: category.seoDescription,
    path: categoryPath(category, catalog),
    indexable: isIndexable(category) && children.length > 0,
    crumbs: [
      { name: 'Issues', path: '/issues' },
      { name: category.name, path: categoryPath(category, catalog) },
    ],
    body: `<section class="hero"><p class="eyebrow">${e(category.name)}</p><h1>${e(category.question)}</h1><p>${e(category.description)}</p><a class="button" href="#scopes">Choose your ${e(category.expectedChildScopeType || 'scope')}</a></section>${unavailable ? '<p role="status">Live results are temporarily unavailable. Candidate options and questions remain available below.</p>' : ''}<div class="columns"><section>${ranking}<p class="method">Only published issue counts are displayed. Pending suggestions are excluded from public rankings.</p></section><aside id="scopes"><h2>Choose a scope</h2><p>Your answer belongs to the specific question you choose.</p><ul class="links">${children.map((q) => `<li>${link(questionPath(category, q, catalog), q.title)}</li>`).join('')}</ul></aside></div>`,
  });
}
export function issueCard(
  issue,
  question,
  category,
  total,
  rank,
  catalog,
  language = 'en',
  languageTotal = total,
) {
  const path = questionPath(category, question, catalog) + '/' + issue.slug;
  const languageSupporters = issue.languageSupporters || 0;
  return `<li class="issue"><div class="issue-heading"><span class="rank">${issue.supporters > 0 ? '#' + rank : '—'}</span><div><h3>${link(path, issue.title)}</h3><p>${e(issue.description)}</p><p class="count">${issue.supporters} ${issue.supporters === 1 ? 'supporter' : 'supporters'}${total ? ' · ' + Math.round((issue.supporters / total) * 100) + '% overall' : ' · Candidate option'}</p>${languageTotal ? `<p class="count">${languageSupporters} ${languageSupporters === 1 ? 'answer' : 'answers'} in ${e(languageName(language))} · ${Math.round((languageSupporters / languageTotal) * 100)}% of ${e(languageName(language))} answers</p>` : ''}${issue.weightedScore !== undefined ? `<p class="count">${issue.weightedScore} priority points · ${issue.personalTopSupporters || 0} personal #1 choices${languageTotal ? ` · ${issue.languageWeightedScore || 0} ${e(languageName(language))} priority points` : ''}</p>` : ''}</div></div><div class="actions"><button data-answer="${e(issue.id)}" data-question="${e(question.id)}" data-category="${e(category.id)}" data-return="${e(path)}" data-intent="${e(category.intent)}">For me too</button><button class="quiet" data-share="${e(path)}" data-share-title="${e(issue.title + ' — ' + question.title)}" data-event="issue_shared">Share</button></div></li>`;
}
export function questionPage(catalog, category, snapshot, selectedSlug) {
  const { question, issues, totalRespondents } = snapshot;
  const language = snapshot.languageCode || 'en';
  const languageRespondents = snapshot.languageRespondents || 0;
  const userCreated = !question.categoryId;
  category = category || {
    id: 'user-questions',
    slug: 'questions',
    name: 'Community questions',
    intent: 'consumer',
  };
  question.scope = question.scope || {
    id: question.scopeId || question.id,
    name: question.scopeName || 'Community question',
  };
  const path = questionPath(category, question, catalog);
  const published = issues
    .filter((i) => i.status === 'published')
    .sort(
      (a, b) =>
        (b.weightedScore ?? b.supporters) - (a.weightedScore ?? a.supporters) ||
        a.id.localeCompare(b.id),
    );
  const selected = selectedSlug
    ? published.find((i) => i.slug === selectedSlug)
    : null;
  if (selectedSlug && !selected) return null;
  const title = selected
    ? `${selected.title} — ${question.scope.name}`
    : question.title;
  const description = selected
    ? `${selected.description} See its current support in “${question.title}” and choose your own #1 issue.`
    : question.description;
  const visible = selected ? [selected] : published;
  const related = catalog.questions
    .filter(
      (q) =>
        isPublished(q) &&
        q.id !== question.id &&
        (question.relatedQuestionIds?.includes(q.id) ||
          q.categoryId === question.categoryId ||
          q.scope?.parentId === question.scope.id ||
          q.scope?.id === question.scope.parentId),
    )
    .slice(0, 8);
  const issuePath = selected ? path + '/' + selected.slug : path;
  return page({
    title,
    description,
    path: issuePath,
    indexable: selected
      ? selected.indexable === true && selected.description?.length >= 100
      : isIndexable(question) && published.length >= 5,
    crumbs: [
      { name: 'Issues', path: '/issues' },
      {
        name: category.name,
        path: userCreated ? '/questions' : categoryPath(category, catalog),
      },
      { name: question.scope.name, path },
      ...(selected ? [{ name: selected.title, path: issuePath }] : []),
    ],
    items: visible.map((i) => ({ title: i.title, path: path + '/' + i.slug })),
    language,
    availableLanguages: question.availableLanguages || ['en'],
    body: `<section class="hero" data-view="${selected ? 'issue_view' : 'question_view'}" data-question="${e(question.id)}" data-category="${e(category.id)}" data-issue="${e(selected?.id || '')}" data-intent="${e(category.intent)}"><p class="eyebrow">${e(question.scope.name)}</p><h1>${e(title)}</h1><p>${e(description)}</p><p><strong>${totalRespondents}</strong> ${totalRespondents === 1 ? 'person has' : 'people have'} answered across all languages · <strong>${languageRespondents}</strong> answered in ${e(languageName(language))}${snapshot.updatedAt ? ' · Updated ' + e(new Date(snapshot.updatedAt).toLocaleDateString('en-IE', { timeZone: 'UTC' })) : ''}</p><button class="quiet" data-share="${e(path)}" data-share-title="${e(question.title)}" data-event="question_shared">Ask someone whose view matters to you ↗</button></section><div class="columns"><section aria-label="Issues"><h2>${selected ? 'This issue' : totalRespondents ? 'Current priorities across all languages' : 'What would you put first?'}</h2>${!totalRespondents ? '<p>Be among the first to answer. These options are starting points, not measured rankings.</p>' : ''}<ol class="ranking">${visible.map((i) => issueCard(i, question, category, totalRespondents, 1 + published.filter((other) => (other.weightedScore ?? other.supporters) > (i.weightedScore ?? i.supporters)).length, catalog, language, languageRespondents)).join('')}</ol>${selected ? `<p>${link(path, 'See all issues in this question →')}</p>` : ''}${question.allowSuggestions !== false || question.choiceSource?.kind === 'free' ? `<section class="composer"><h2>Don’t see your issue?</h2><form data-compose data-question="${e(question.id)}" data-category="${e(category.id)}" data-return="${e(path)}" data-intent="${e(category.intent)}"><label for="issue-title">What is your #1 issue?</label><input id="issue-title" name="title" required minlength="3" maxlength="140" autocomplete="off" list="candidate-titles" placeholder="Describe the issue in a few words"><datalist id="candidate-titles">${published.map((i) => `<option value="${e(i.title)}"></option>`).join('')}</datalist><p class="hint">Suggestions help avoid duplicates. New issues await review before joining public results.</p><label class="hint"><input type="checkbox" name="attribution" value="authored"> Show me as the author (otherwise anonymous)</label><button type="submit">Make this my choice</button></form></section>` : ''}<p class="method">The ranking above uses all languages; language figures show the selected language without changing its meaning. One current choice per person in each concrete question, plus one personal #1 across all your choices. Ireland, France, Dublin and Cork can each have their own answer. Priority points and real supporter counts are shown separately. Percentages use all current answers, including suggestions awaiting review; public shares may therefore total less than 100%. Ties share rank. This is voluntary participation, not representative polling.</p></section><aside><section class="circle"><p class="eyebrow">Better with another perspective</p><h2>${category.intent === 'business' ? 'Whose work does this affect?' : 'Who would see it differently?'}</h2><p>${category.intent === 'business' ? 'Ask a colleague what they would fix first.' : 'Invite someone whose priorities you want to understand.'}</p><button data-share="${e(path)}" data-share-title="${e(question.title)}" data-event="question_shared">${category.intent === 'business' ? 'Ask your team' : 'Ask people you know'}</button><div data-conversion hidden><hr><p>${category.intent === 'business' ? 'Explore a shared workspace for your team.' : 'Bring your people together in Sneat.app.'}</p><a data-event="${category.intent === 'business' ? 'sneat_work_cta_clicked' : 'sneat_app_cta_clicked'}" href="${category.intent === 'business' ? 'https://sneat.work/team/' : 'https://sneat.app/main'}?utm_source=issuenumber.one&amp;utm_medium=answer&amp;utm_campaign=public-issues">${category.intent === 'business' ? 'Explore Sneat.work' : 'Continue in Sneat.app'} →</a></div><p class="hint">Your contacts and individual answers are never shown on this public page.</p></section><h2>Related questions</h2><ul class="links">${related
      .map((q) => {
        const c = catalog.categories.find((c) => c.id === q.categoryId);
        return `<li>${link(
          questionPath(c || category, q, catalog) + (!q.categoryId ? '?' : ''),
          q.title,
        )}</li>`;
      })
      .join(
        '',
      )}</ul></aside></div><p id="action-status" role="status" aria-live="polite"></p>`,
  });
}
export function sitemap(catalog) {
  const paths = ['/', '/issues', '/for-work/'];
  paths.push(
    ...catalog.questions
      .filter((q) => !q.categoryId && isIndexable(q))
      .map((q) => '/questions/' + q.slug + '?'),
  );
  for (const category of catalog.categories.filter(isIndexable)) {
    const children = catalog.questions.filter(
      (q) => q.categoryId === category.id && isIndexable(q),
    );
    if (children.length)
      paths.push(
        categoryPath(category, catalog),
        ...children.map((q) => questionPath(category, q, catalog)),
      );
  }
  return `<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">${paths.map((p) => `<url><loc>${ORIGIN}${e(p)}</loc></url>`).join('')}</urlset>`;
}
