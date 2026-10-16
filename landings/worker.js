import { serveQA } from './qa/serve.js';
// IssueNumber.one landing worker.
//
// wrangler.jsonc `assets.run_worker_first` sends ONLY these paths here; every
// other path is a free static asset (see landings/public/_redirects and
// _headers for the app-shell rewrites, their noindex/no-store headers and the
// /index.html redirect):
//
//   POST /waitlist  -> validate an email, store {email, ts, source, plan} in KV,
//                      return a small JSON body for the inline thank-you state.
//                      `plan` is optional (set by the pricing-card CTAs).
//   /, /issues, /issues/*, /questions, /questions/*, /sitemap.xml
//                   -> server-rendered public Q&A from the backend (qa/serve.js),
//                      except /questions/new and ?preview=1, which get the app
//                      shell (a query condition _redirects cannot express).
//   anything else   -> the static assets (env.ASSETS); reached only by
//                      non-navigation asset misses.

const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[^\s@]{2,}$/;

function json(body, status = 200) {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'content-type': 'application/json; charset=utf-8' },
  });
}

async function handleWaitlist(request, env) {
  let email = '';
  let source = 'landing';
  let plan = '';
  try {
    const ct = request.headers.get('content-type') || '';
    if (ct.includes('application/json')) {
      const body = await request.json();
      email = String(body.email || '');
      if (body.source) source = String(body.source);
      if (body.plan) plan = String(body.plan);
    } else {
      const form = await request.formData();
      email = String(form.get('email') || '');
      if (form.get('source')) source = String(form.get('source'));
      if (form.get('plan')) plan = String(form.get('plan'));
    }
  } catch {
    return json({ ok: false, error: 'bad_request' }, 400);
  }

  email = email.trim().toLowerCase();
  if (!EMAIL_RE.test(email) || email.length > 254) {
    return json({ ok: false, error: 'invalid_email' }, 422);
  }
  source = source.slice(0, 64);
  plan = plan.trim().toLowerCase().slice(0, 32);

  const entry = { email, ts: new Date().toISOString(), source };
  if (plan) entry.plan = plan;
  if (env.WAITLIST) {
    // Keyed by email: re-submitting the same address just refreshes the entry.
    await env.WAITLIST.put(`waitlist:${email}`, JSON.stringify(entry));
  } else {
    // Fallback if the KV binding is ever missing (e.g. local preview without KV).
    console.log('waitlist signup (no KV binding):', JSON.stringify(entry));
  }
  return json({ ok: true });
}

// The authenticated Angular app shell, with noindex/no-store. ASSETS
// canonicalizes .html URLs with a 307 redirect, so fetch the clean asset URL
// internally and the browser stays on the requested path.
async function serveAppShell(request, env, url) {
  if (request.method !== 'GET' && request.method !== 'HEAD')
    return new Response('Method not allowed', { status: 405 });
  const shell = await env.ASSETS.fetch(
    new Request(new URL('/index.app', url), request),
  );
  const headers = new Headers(shell.headers);
  headers.set('X-Robots-Tag', 'noindex, nofollow');
  headers.set('Cache-Control', 'no-store');
  return new Response(shell.body, { status: shell.status, headers });
}

export default {
  async fetch(request, env) {
    const url = new URL(request.url);

    if (url.pathname === '/waitlist') {
      if (request.method === 'POST') return handleWaitlist(request, env);
      return json({ ok: false, error: 'method_not_allowed' }, 405);
    }

    const publicResponse = await serveQA(request, env);
    if (publicResponse) return publicResponse;
    if (
      url.pathname === '/questions/new' ||
      (url.pathname.startsWith('/questions/') &&
        url.searchParams.get('preview') === '1')
    )
      return serveAppShell(request, env, url);
    return env.ASSETS.fetch(request);
  },
};
