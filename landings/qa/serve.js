import {
  discovery,
  categoryPage,
  questionPage,
  sitemap,
  page,
  isPublished,
  categoryPath,
  questionPath,
  questionsPage,
} from './render.js';
async function api(env, path) {
  const base = env.ISSUENUMBER_API_ORIGIN || 'https://api.sneat.cloud';
  const result = await fetch(new URL('/v0/issuenumber/' + path, base), {
    headers: { accept: 'application/json' },
    signal: AbortSignal.timeout(5000),
  });
  if (!result.ok) throw new Error('Public data unavailable');
  const length = Number(result.headers.get('content-length') || 0);
  if (length > 2_000_000) throw new Error('Public data too large');
  // The owning backend bounds these DTOs; never fetch private answers here.
  return result.json();
}
export async function serveQA(request, env) {
  const url = new URL(request.url);
  const path = url.pathname.replace(/\/$/, '') || '/';
  if (path === '/questions/new' || url.searchParams.get('preview') === '1')
    return null;
  if (
    !['/', '/issues', '/sitemap.xml', '/robots.txt'].includes(path) &&
    !path.startsWith('/issues/') &&
    !path.startsWith('/questions/') &&
    path !== '/questions'
  )
    return null;
  if (request.method !== 'GET' && request.method !== 'HEAD')
    return new Response('Method not allowed', {
      status: 405,
      headers: { allow: 'GET, HEAD' },
    });
  if (path === '/robots.txt')
    return new Response(
      'User-agent: *\nAllow: /\nDisallow: /answer\nDisallow: /verify\nDisallow: /login\nSitemap: https://issuenumber.one/sitemap.xml\n',
      { headers: { 'content-type': 'text/plain' } },
    );
  const headers = {
    'content-type': 'text/html; charset=utf-8',
    'cache-control': 'public, max-age=30',
    'x-content-type-options': 'nosniff',
    'referrer-policy': 'strict-origin-when-cross-origin',
  };
  try {
    const catalog = await api(env, 'catalog');
    if (path === '/sitemap.xml')
      return new Response(request.method === 'HEAD' ? null : sitemap(catalog), {
        headers: { ...headers, 'content-type': 'application/xml' },
      });
    let html;
    if (path === '/questions') {
      html = questionsPage(catalog);
    } else if (path === '/' || path === '/issues') {
      html = discovery(catalog, path);
    } else if (path.startsWith('/questions/')) {
      const slug = path.slice('/questions/'.length);
      const language = url.searchParams.get('lang') || 'en';
      const snapshot = await api(
        env,
        `question?slug=${encodeURIComponent(slug)}&lang=${encodeURIComponent(language)}`,
      );
      html = questionPage(catalog, undefined, snapshot);
    } else {
      const category = catalog.categories.find(
        (c) => categoryPath(c, catalog) === path && isPublished(c),
      );
      if (category) {
        const children = catalog.questions.filter(
          (q) => q.categoryId === category.id && isPublished(q),
        );
        const snapshots = [];
        for (let i = 0; i < children.length; i += 6)
          snapshots.push(
            ...(await Promise.all(
              children
                .slice(i, i + 6)
                .map((q) =>
                  api(env, 'question?questionId=' + encodeURIComponent(q.id)),
                ),
            )),
          );
        html = categoryPage(catalog, category, snapshots);
      } else {
        for (const question of catalog.questions.filter(isPublished)) {
          const c = catalog.categories.find(
            (c) => c.id === question.categoryId && isPublished(c),
          );
          if (!c) continue;
          const canonical = questionPath(c, question, catalog);
          const tail = path.startsWith(canonical + '/')
            ? path.slice(canonical.length + 1)
            : '';
          if (path === canonical || (tail && !tail.includes('/'))) {
            const language = url.searchParams.get('lang') || 'en';
            html = questionPage(
              catalog,
              c,
              await api(
                env,
                `question?questionId=${encodeURIComponent(question.id)}&lang=${encodeURIComponent(language)}`,
              ),
              tail || undefined,
            );
            if (html) break;
          }
          const legacy = '/issues/' + c.slug + '/' + question.slug;
          if (path === legacy && legacy !== canonical)
            return Response.redirect(
              new URL(canonical, 'https://issuenumber.one'),
              308,
            );
        }
      }
    }
    if (!html)
      return new Response(
        page({
          title: 'Question not found',
          description: 'Explore the published questions on IssueNumber.one.',
          path,
          body: '<h1>Question not found</h1><p><a href="/issues">Explore published questions</a></p>',
          indexable: false,
        }),
        { status: 404, headers },
      );
    return new Response(request.method === 'HEAD' ? null : html, { headers });
  } catch {
    return new Response(
      page({
        title: 'Results temporarily unavailable',
        description: 'Please try again shortly.',
        path,
        body: '<h1>We couldn’t load the latest issues</h1><p>Your answers are safe. Please try again shortly.</p><p><a href="/for-work/">Explore IssueNumber.one for work</a></p>',
        indexable: false,
      }),
      {
        status: 503,
        headers: {
          ...headers,
          'cache-control': 'no-store',
          'retry-after': '30',
        },
      },
    );
  }
}
