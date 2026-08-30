# PRD — Unified Forge Platform — Authorization Contract, IAM Console, and Shell Composition

**Initiative:** `INIT-SYSTEMFORGE-001`
**Status:** Draft
**Home repo:** `github.com/plexusone/systemforge`

## Problem

Three Forge products (ActionForge, AgentForge, DashForge) are converging on
SystemForge for identity/tenancy/authorization, but each is inventing its own
authorization surface:

- **DashForge** defines a SpiceDB schema with *unprefixed* definitions
  (`organization`, `platform`, `dashboard`) that collide with any other Forge
  app's schema in a composed deployment.
- **ActionForge** just defined app-namespaced types (`actionforge.workflow`),
  `{app}:{resource}:{verb}` OAuth scopes, and domain roles (`approver`,
  `developer`) — conventions that exist only in ActionForge's repo.
- **AgentForge** has a runtime tool-policy layer but no org-role vocabulary
  at all.

There is also no administration UI: role assignments, permission
introspection, and token/scope management have no console (the AWS IAM
equivalent), and no formal contract exists for mounting Forge UIs into the
`@plexusone/shell` (which already provides AppShell/nav/auth/session and is
consumed by agentforge-web).

The intended product motion — *start a SaaS app with SystemForge, then
optionally add ActionForge/AgentForge/DashForge as embedded capabilities,
composed under one menu, one login, one session* — requires all three gaps
closed in SystemForge itself, not per-app.

## Goals

1. **One authorization contract, owned by SystemForge.** A documented
   convention plus helpers every Forge app adopts: app-namespaced resource
   types (`{app}.{resource}`), OAuth scopes as `{app}:{resource}:{verb}` with
   `{app}:admin` implying all, a canonical role set (owner/admin/editor/
   viewer) with rules for app domain roles (e.g. ActionForge's `approver`),
   the two-gate model (role permission ∩ token scope), and a shared SpiceDB
   base schema where SystemForge owns `principal`/`organization`/`platform`
   and apps contribute `{app}_`-prefixed resource definitions referencing
   them.
2. **An IAM console in systemforge-web.** AWS-IAM-style administration:
   org members with role assignment, a role→permission matrix driven by each
   app's registered vocabulary, and API key/token management with a scope
   picker. ReBAC relationship browsing is a follow-on.
3. **A shell app-registration contract.** Each Forge UI package exports a
   typed registration (id, nav sections gated by role/scope, routes, required
   scopes); `@plexusone/shell` mounts registered apps under one chrome.
   Composition is build-time by default; the contract is transport-agnostic
   so runtime micro-frontends can be adopted later without changing app code.

## Non-Goals (this initiative)

- Runtime micro-frontends (Module Federation/import maps) — deliberately
  deferred until independent deploy cadence demands it; the registration
  contract is designed not to preclude it.
- The customer-embedded mode (GoodData-style iframe/web-component embedding
  of a single Forge capability into a third-party SaaS with scoped tokens) —
  specced as follow-on; it is a per-app embed SDK, not shell composition.
- Migrating every existing Forge screen into the shell — the reference
  composition proves the contract with two apps; full migrations are each
  app's own roadmap.

## Users and Experience

- **Product Builder starting a new SaaS**: scaffold on SystemForge +
  systemforge-web, get login/org-switcher/settings out of the box; add a
  Forge capability by installing its UI package and registering it — it
  appears in the shared nav, gated by the user's role and the app's scopes.
- **Org admin**: opens Administration → IAM: sees members and their roles per
  app, changes a member from `editor` to `approver`, issues a CI token scoped
  to `actionforge:runs:write` (and sees it cannot approve gates).
- **Forge app developer**: implements the vocabulary interface once; their
  roles/permissions/scopes appear in the IAM console automatically.

## Success Criteria

- The convention doc exists in SystemForge and all three Forge apps' authz
  vocabularies conform (DashForge schema re-namespaced; AgentForge vocabulary
  created; ActionForge aligned to the shared org definitions).
- A composed SpiceDB deployment loads SystemForge base + ≥2 app schemas with
  no definition collisions.
- The IAM console lists members, edits role assignments, and issues a token
  with a subset of an app's scopes — against real SystemForge APIs.
- A reference host app composes `@plexusone/shell` + two Forge UI packages
  with one login/session and role/scope-gated nav.
