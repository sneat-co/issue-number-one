import { serveQA } from './qa/serve.js';
// IssueNumber.one landing worker.
//
//   POST /waitlist  -> validate an email, store {email, ts, source, plan} in KV,
//                      return a small JSON body for the inline thank-you state.
//                      `plan` is optional (set by the pricing-card CTAs).
//   everything else -> the static Astro landing (env.ASSETS).
//
// Landing-only worker: static site + a tiny waitlist probe. When the
// IssueNumber.one app ships, fold it in following the root-mount pattern used
// by the other Sneat landings (index.app.html + reserved-path router).

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

export default {
  async fetch(request, env) {
    const url = new URL(request.url);

    if (url.pathname === '/waitlist') {
      if (request.method === 'POST') return handleWaitlist(request, env);
      return json({ ok: false, error: 'method_not_allowed' }, 405);
    }

    if (url.pathname === '/index.html')
      return Response.redirect(new URL('/', url), 308);
    const publicResponse = await serveQA(request, env);
    if (publicResponse) return publicResponse;
    if (
      /^\/(answer|verify|login|signup|register)(\/|$)/.test(url.pathname) ||
      url.pathname === '/questions/new' ||
      (url.pathname.startsWith('/questions/') &&
        url.searchParams.get('preview') === '1')
    ) {
      if (request.method !== 'GET' && request.method !== 'HEAD')
        return new Response('Method not allowed', { status: 405 });
      const shell = await env.ASSETS.fetch(
        new Request(new URL('/index.app.html', url), request),
      );
      const headers = new Headers(shell.headers);
      headers.set('X-Robots-Tag', 'noindex, nofollow');
      headers.set('Cache-Control', 'no-store');
      return new Response(shell.body, { status: shell.status, headers });
    }
    return env.ASSETS.fetch(request);
  },
};
