# Stored question translations

This package translates question title and description while preserving the question ID, slug, answer slot, choices, and totals. The default fixed allowlist is English (`en`) and Russian (`ru`). Unsupported `?lang=` values return an error instead of silently serving English as indexed localized content.

Translations live at:

`/spaces/{spaceID}/ext/issuenumber/questions/{questionID}/translations/{language}`

Each document records the source language, source content revision, deterministic content hash, machine-translation status, and timestamps. A Firestore transaction refuses a write if the source changed while the provider request was in flight. A human-corrected current translation is preserved.

`Service.TranslateQuestion` is suitable for a bounded on-demand public read: published questions may populate a missing allowlisted language, concurrent requests are deduplicated in-process, and fresh stored translations avoid provider calls. Draft or pending questions require their creator or a trusted publishing worker. `Service.TranslateEnabledQuestion` is the publish hook for eagerly storing every enabled language.

The Google adapter uses Application Default Credentials and the host Google Cloud project. Question text is passed only as plain translation input; it is never interpolated into prompts or instructions.
