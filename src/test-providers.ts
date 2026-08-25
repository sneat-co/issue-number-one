import { provideZonelessChangeDetection } from '@angular/core';

/**
 * Default providers for every unit test's TestBed, mirroring the
 * zoneless setup registered in `app.config.ts`. Wired via project.json's
 * `test` target `providersFile` option so individual specs don't need to
 * repeat it.
 *
 * Note on the `zone.js` devDependency: the app itself is zoneless (empty
 * `polyfills` in project.json's `build` target, `provideZonelessChangeDetection()`
 * here and in app.config.ts) and never imports zone.js. It stays a
 * devDependency only so it is hoisted into the root `node_modules` — without
 * it, `@angular/build:unit-test`'s generated test harness cannot statically
 * resolve the `zone.js/testing` subpath export under pnpm's strict linking,
 * even though the harness only imports it behind an `if (typeof Zone !==
 * 'undefined')` guard that is always false here. Do not add zone.js to
 * `polyfills` or import it anywhere in application code.
 */
export default [provideZonelessChangeDetection()];
