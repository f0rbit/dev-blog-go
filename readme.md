# dev-blog (archived)

> **This project has been migrated and is no longer maintained.**
>
> The blog now lives as a TypeScript-first app inside the [`devpad`](https://github.com/f0rbit/devpad) monorepo at [`apps/blog`](https://github.com/f0rbit/devpad/tree/main/apps/blog).

## Migration summary

The original Go server + React SPA was rewritten as a Cloudflare-first stack and folded into the `devpad` monorepo. The new architecture:

| Layer | Old (this repo) | New (devpad) |
|---|---|---|
| API | Go + custom HTTP handlers | Hono on Cloudflare Workers |
| Database | local SQLite | Cloudflare D1 (Drizzle ORM) |
| Post storage | DB rows | [`@f0rbit/corpus`](https://github.com/f0rbit/corpus) version-controlled stores |
| Frontend | React SPA | Astro + SolidJS |
| Auth | GitHub OAuth (custom) | shared `@devpad/api` session model |
| Schema | hand-written types | shared `@devpad/schema` (Zod + Drizzle) |
| Deployment | self-hosted VPS | Cloudflare Workers + Pages |

A VPS → Cloudflare migration script (in the `devpad` repo) was used to move existing data from the Go server's SQLite database into D1 + Corpus.

The full migration plan that drove the rewrite is preserved at [`documents/rewrite.md`](./documents/rewrite.md) for reference.

## Where to find things now

- **Blog app**: [`devpad/apps/blog`](https://github.com/f0rbit/devpad/tree/main/apps/blog)
- **API client**: [`@devpad/api`](https://github.com/f0rbit/devpad/tree/main/packages/api)
- **Schema**: [`@devpad/schema`](https://github.com/f0rbit/devpad/tree/main/packages/schema)
- **Live site**: blog.devpad.tools

## Repository status

This repo is kept for git history and the migration plan. No further development happens here.
