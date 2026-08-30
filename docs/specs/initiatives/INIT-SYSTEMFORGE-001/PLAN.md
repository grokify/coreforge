# PLAN — Unified Forge Platform — Authorization Contract, IAM Console, and Shell Composition

**Initiative:** `INIT-SYSTEMFORGE-001`
**Status:** Draft

## Entry Gate

1. Repos registered: systemforge, systemforge-web, actionforge, agentforge,
   dashforge (done).
2. ActionForge's authz package accepted as the convention seed (it already
   implements namespaced types, `{app}:{resource}:{verb}` scopes, the
   two-gate Guard, and a prefixed SpiceDB schema).

## Build Order

### Phase 1 — Authorization Contract

The contract must land before the console (which introspects vocabularies)
and before per-app migrations.

1. `RMI-SYSTEMFORGE-003` Convention doc + `AppVocabulary` interface and
   vocabulary registry (scope-format and duplicate-definition validation) in
   systemforge/authz.
2. `RMI-SYSTEMFORGE-004` Shared SpiceDB base schema (`principal`,
   `organization`, `platform`) + schema-composition helper.
3. `RMI-DASHFORGE-016` DashForge adoption: re-namespace schema to
   `dashforge_*` referencing the shared base; conform scopes/roles;
   conformance test.
4. `RMI-AGENTFORGE-001` AgentForge adoption: create org-role vocabulary
   (roles/permissions/scopes/schema) on the convention; conformance test.
5. `RMI-ACTIONFORGE-037` ActionForge alignment: switch to the shared base
   `organization` definition and register its AppVocabulary.

**Exit criteria:** composed schema of base + 3 apps validates with no
collisions; all three apps pass conformance tests.

### Phase 2 — IAM Console

6. `RMI-SYSTEMFORGE-005` Admin APIs: role-assignment CRUD per
   (org, principal, app); vocabulary introspection; token issuance bounded by
   AllowedScopes.
7. `RMI-SYSTEMFORGEWEB-001` IAM pages: Members + role assignment; Roles &
   Permissions matrix (vocabulary-driven); API Keys & Tokens with scope
   picker.

**Exit criteria:** PRD success criterion — assign a role and issue a
scope-limited token entirely through the console.

### Phase 3 — Shell Composition

8. `RMI-SYSTEMFORGEWEB-002` `ForgeApp` registration contract in
   `@plexusone/shell` (nav/route/scope gating), versioned, documented.
9. `RMI-SYSTEMFORGEWEB-003` Reference composed host: shell + auth + two
   Forge UI surfaces under one login/menu/session; role/scope-gated nav
   proven.

**Exit criteria:** PRD success criterion — the reference host runs with one
session across both apps.

## Risks and Mitigations

See TRD: DashForge schema migration timing (pre-production), vocabulary
drift (conformance tests + startup validation), contract over-fit (minimal
versioned `ForgeApp`).

## Follow-On

- ReBAC relationship browser in the IAM console.
- Customer-embedded mode per app (iframe/web-component + scoped-token embed
  SDK) — the GoodData-style story.
- Runtime micro-frontends (Module Federation) if independent deploy cadence
  emerges; `ForgeApp` is the remote-module contract when it does.
- ActionForge UI as a React package for shell composition (its standalone
  single-file UI remains for `actionforge serve`).
