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
    headers: { "content-type": "application/json; charset=utf-8" },
  });
}

async function handleWaitlist(request, env) {
  let email = "";
  let source = "landing";
  let plan = "";
  try {
    const ct = request.headers.get("content-type") || "";
    if (ct.includes("application/json")) {
      const body = await request.json();
      email = String(body.email || "");
      if (body.source) source = String(body.source);
      if (body.plan) plan = String(body.plan);
    } else {
      const form = await request.formData();
      email = String(form.get("email") || "");
      if (form.get("source")) source = String(form.get("source"));
      if (form.get("plan")) plan = String(form.get("plan"));
    }
  } catch {
    return json({ ok: false, error: "bad_request" }, 400);
  }

  email = email.trim().toLowerCase();
  if (!EMAIL_RE.test(email) || email.length > 254) {
    return json({ ok: false, error: "invalid_email" }, 422);
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
    console.log("waitlist signup (no KV binding):", JSON.stringify(entry));
  }
  return json({ ok: true });
}

export default {
  async fetch(request, env) {
    const url = new URL(request.url);

    if (url.pathname === "/waitlist") {
      if (request.method === "POST") return handleWaitlist(request, env);
      return json({ ok: false, error: "method_not_allowed" }, 405);
    }

    return env.ASSETS.fetch(request);
  },
};
