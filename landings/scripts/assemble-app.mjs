// Merge the Angular app build into the public Astro/Worker distribution.
// The public homepage keeps dist/index.html; authenticated routes receive
// dist/index.app.html from the Worker.
import { access, cp } from 'node:fs/promises';
import { sep } from 'node:path';

const browser = '../dist/in1app/browser';
const dist = './dist';

try {
  await access(`${browser}/index.html`);
} catch {
  throw new Error(
    `App build not found at ${browser}/index.html; run the app build first.`,
  );
}

await cp(browser, dist, {
  recursive: true,
  force: true,
  filter: (source) => !source.endsWith(`${sep}index.html`),
});
await cp(`${browser}/index.html`, `${dist}/index.app.html`, { force: true });
