# CLAUDE.md - systemforge

SystemForge provides reusable identity, session, authorization, OAuth, BFF, and
multi-app infrastructure for Go SaaS applications.

## Authorization Notes

- `authz.Authorizer` is the stable bridge interface. Keep providers compatible
  with `Can`, `CanAll`, `CanAny`, and `Filter`.
- The SpiceDB provider checks concrete resource type/id/permission tuples.
  Applications that need field or column policy should model those as concrete
  resources rather than relying on attributes.
- Query languages such as GuardSQL (formerly GrokifyQL) should own parsing and
  AST safety. SystemForge should receive typed principal/action/resource checks.
- Saved queries, dashboard widgets, alerts, and scheduled reports should be
  re-checked against current authorization before execution.

## GuardSQL Adapter (`authzguardsql`)

- `authzguardsql` is SystemForge's adapter from GuardSQL schemas/policies to
  SystemForge authorization (`PolicyBuilder`, `ResourceBuilder`,
  `DefaultResourceBuilder`). It imports `github.com/grokify/guardsql`; the
  dependency direction is **systemforge → guardsql** (guardsql never imports
  systemforge, so the query language stays dependency-light).
- It previously lived as the nested `github.com/grokify/guardsql/authzsystemforge`
  module. It was consolidated into systemforge to remove the nested-module
  maintenance burden (fragile cross-module tagging). Consumers import
  `github.com/grokify/systemforge/authzguardsql`.

## Cross-Repo Consumers

- UIForge/DashForge uses SystemForge as the bridge from GuardSQL policy checks
  to SpiceDB, via the `authzguardsql` adapter.
