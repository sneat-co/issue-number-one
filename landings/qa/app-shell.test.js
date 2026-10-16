import assert from 'node:assert/strict';
import {
  copyFile,
  mkdtemp,
  mkdir,
  readFile,
  writeFile,
  rm,
} from 'node:fs/promises';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { fileURLToPath } from 'node:url';
import test from 'node:test';
import { unstable_dev } from 'wrangler';

const landings = fileURLToPath(new URL('..', import.meta.url));

// The real assets config (wrangler.jsonc has only full-line comments).
async function realAssetsConfig() {
  const text = await readFile(join(landings, 'wrangler.jsonc'), 'utf8');
  const config = JSON.parse(
    text
      .split('\n')
      .filter((line) => !/^\s*\/\//.test(line))
      .join('\n'),
  );
  return config.assets;
}

test('routes app-shell paths statically and only dynamic paths to the Worker', async () => {
  const dir = await mkdtemp(join(tmpdir(), 'issuenumber-assets-'));
  let worker;
  try {
    await mkdir(join(dir, 'assets'));
    const shell =
      '<!doctype html><html><body><app-root></app-root></body></html>';
    const landing = '<!doctype html><html><body>landing</body></html>';
    await writeFile(join(dir, 'assets/index.app.html'), shell);
    await writeFile(join(dir, 'assets/index.html'), landing);
    for (const file of ['_redirects', '_headers', 'robots.txt'])
      await copyFile(join(landings, 'public', file), join(dir, 'assets', file));
    const assets = await realAssetsConfig();
    assert.ok(Array.isArray(assets.run_worker_first));
    assert.notEqual(assets.not_found_handling, 'none');
    await writeFile(
      join(dir, 'wrangler.json'),
      JSON.stringify({
        name: 'issuenumber-shell-test',
        compatibility_date: '2026-06-23',
        assets: { ...assets, directory: './assets' },
        // Unreachable origin: the SSR paths must answer from the Worker (503
        // fallback page) rather than from the static landing.
        vars: { ISSUENUMBER_API_ORIGIN: 'http://127.0.0.1:9' },
      }),
    );
    worker = await unstable_dev(join(landings, 'worker.js'), {
      config: join(dir, 'wrangler.json'),
      local: true,
      logLevel: 'error',
      experimental: {
        disableExperimentalWarning: true,
        disableDevRegistry: true,
        watch: false,
      },
    });

    // App shell: /answer etc. via _redirects + _headers, /questions/new and
    // ?preview=1 via the Worker.
    for (const path of [
      '/answer',
      '/answer/x',
      '/verify',
      '/login',
      '/signup',
      '/register/x',
      '/questions/new',
      '/questions/example?preview=1',
    ]) {
      const response = await worker.fetch(path, { redirect: 'manual' });
      assert.equal(response.status, 200, path);
      assert.equal(response.headers.get('location'), null, path);
      assert.equal(await response.text(), shell, path);
      assert.equal(response.headers.get('cache-control'), 'no-store', path);
      assert.equal(
        response.headers.get('x-robots-tag'),
        'noindex, nofollow',
        path,
      );
    }

    const index = await worker.fetch('/index.html', { redirect: 'manual' });
    assert.equal(index.status, 308);
    assert.equal(
      new URL(index.headers.get('location'), 'http://x').pathname,
      '/',
    );

    const robots = await worker.fetch('/robots.txt');
    assert.equal(robots.status, 200);
    assert.match(
      await robots.text(),
      /Sitemap: https:\/\/issuenumber\.one\/sitemap\.xml/,
    );

    // Server-rendered paths reach the Worker (503 because the API is down).
    for (const path of [
      '/',
      '/issues',
      '/issues/x',
      '/questions/x',
      '/sitemap.xml',
    ]) {
      const response = await worker.fetch(path, { redirect: 'manual' });
      assert.equal(response.status, 503, path);
      await response.arrayBuffer();
    }

    const waitlist = await worker.fetch('/waitlist');
    assert.equal(waitlist.status, 405);
    assert.equal((await waitlist.json()).error, 'method_not_allowed');
  } finally {
    await worker?.stop();
    await rm(dir, { recursive: true, force: true });
  }
});
