import assert from 'node:assert/strict';
import { mkdtemp, mkdir, writeFile, rm } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { fileURLToPath } from 'node:url';
import test from 'node:test';
import { unstable_dev } from 'wrangler';

test('Cloudflare serves continuation HTML without redirecting to the shell filename', async () => {
  const dir = await mkdtemp(join(tmpdir(), 'issuenumber-assets-'));
  let worker;
  try {
    await mkdir(join(dir, 'assets'));
    const shell =
      '<!doctype html><html><body><app-root></app-root></body></html>';
    await writeFile(join(dir, 'assets/index.app.html'), shell);
    await writeFile(
      join(dir, 'wrangler.json'),
      JSON.stringify({
        name: 'issuenumber-shell-test',
        compatibility_date: '2026-06-23',
        assets: {
          directory: './assets',
          binding: 'ASSETS',
          run_worker_first: true,
          not_found_handling: 'none',
        },
      }),
    );
    worker = await unstable_dev(
      fileURLToPath(new URL('../worker.js', import.meta.url)),
      {
        config: join(dir, 'wrangler.json'),
        local: true,
        logLevel: 'error',
        experimental: {
          disableExperimentalWarning: true,
          disableDevRegistry: true,
          watch: false,
        },
      },
    );
    for (const path of [
      '/answer',
      '/verify',
      '/login',
      '/questions/new',
      '/questions/example?preview=1',
    ]) {
      const response = await worker.fetch(path, { redirect: 'manual' });
      assert.equal(response.status, 200, path);
      assert.equal(response.headers.get('location'), null, path);
      assert.equal(await response.text(), shell, path);
      assert.equal(response.headers.get('cache-control'), 'no-store');
      assert.equal(response.headers.get('x-robots-tag'), 'noindex, nofollow');
    }
  } finally {
    await worker?.stop();
    await rm(dir, { recursive: true, force: true });
  }
});
