# Nest.js rules

The Nest.js microservices under `apps/microservices/`. Loaded in every session. The general, TypeScript, localization and testing rules apply first; this file adds what is specific to Nest.

1. **The tree.** `apps/microservices/<service>/src/` holds `main.ts` (one `bootstrap` function), `app.module.ts` (the root module) and one folder per concern (`health/`, `event-log/`, `demo-requests/`, `archive/`, …). Framework files keep their Nest names and meaning: `main.ts`, `app.module.ts`, `nest-cli.json`, `tsconfig.build.json`. An infrastructure component (`redis`, `minio`) is a folder under `apps/microservices/` with a Dockerfile and a README and no package.
2. **Suffixes.** In a Nest app the suffix table of the TypeScript rules reads: `.module` one `@Module()` class; `.controller` one `@Controller()` class, of HTTP routes or of `@EventPattern` handlers (the client-side widget meaning of `.controller` in the TypeScript rules does not apply here); `.strategy` one custom transport strategy; `.pipe` one `PipeTransform`; `.filter` one `ExceptionFilter`; `.guard` one `CanActivate`; `.token` one injection token; `.context` one `BaseRpcContext`. `.service` stays a class; `.factory`, `.adapter`, `.compute`, `.map`, `.fmt`, `.data`, `.types` and `.validation` keep their meaning.
3. **Injection by token.** `experimentalDecorators` is on and `emitDecoratorMetadata` is off. Every constructor parameter carries `@Inject(token)`: the class itself for a class, a `Symbol` from a `.token.ts` file for a port. `reflect-metadata` is imported once, in `main.ts`. A missing `@Inject` fails at startup, never silently, and the unit tests need no decorator plugin.
4. **Validation is Zod.** A request body passes through one validation pipe built on the Zod schema of a `.validation.ts` file; no `class-validator`, no DTO classes. A refused body answers 400 with field errors.
5. **No configuration module.** No `@nestjs/config`. A provider that needs the environment is a factory (`.factory.ts`) that reads `process.env.NAME` at creation and throws when the value is missing, as the TypeScript rules require. The app's README documents every key.
6. **Events over a durable log.** Services exchange events through the Redis Streams transport of `@langdrift/events` (a custom strategy and client), never through Nest's pub/sub transporters. Handlers are `@EventPattern` methods on a controller, acknowledged only after they complete; an envelope is validated with Zod on both sides.
7. **Copy is JSON.** Catalogs live in `src/i18n/messages/*.json` and loaders in `src/i18n/loaders/`, with ICU formatted through `intl-messageformat`. The localization rules apply with that location in place of the Next tree.
8. **Scripts call Nest directly.** `dev: nest start --watch …`, `build: nest build`, `start: node dist/main.js`, `typecheck: tsc --noEmit`, `test: vitest run`. `nest build` bundles with Nest's webpack configuration so the app's alias resolves at runtime; one alias per app, `@app/*`, in `tsconfig.paths.json`.
9. **Handlers in the controller, helpers in their own files.** A controller method dispatches; a computation, a mapping, a schema or an adapter lives in its own single-responsibility file, the Nest analog of the route-handler rule.
10. **Tests never boot Nest.** A unit test constructs the class with fakes: no `@nestjs/testing`, no application bootstrap, nothing reaching Redis, MinIO or a provider. Descriptions name the path (`describe('demo-requests/demo-requests/demo-requests.controller')`).
11. **Dockerfiles.** Every component has a Dockerfile with a `HEALTHCHECK`; `ex` derives readiness from it. A service builds from the repository root with pnpm as a checksum-pinned static binary and Turborepo's `prune --docker`, runs as the `node` user and exposes its port. CI builds no image.

```ts
// Do:
@Controller('v1/demo-requests')
export class DemoRequestsController {
	constructor(
		@Inject(EVENT_LOG_PRODUCER) private readonly producer: ClientProxy,
		@Inject(EventLogReaderService) private readonly reader: EventLogReaderService,
	) {}
}

// Don't:
@Controller('v1/demo-requests')
export class DemoRequestsController {
	constructor(private readonly reader: EventLogReaderService) {}
}
```
