# Query Policy Engines

SystemForge authorization can be used by application query languages to enforce
dataset, table, column, field, and semantic-model permissions. The query
language remains responsible for parsing and structural safety; SystemForge
provides the authorization decision bridge.

UIForge and GrokifyQL are the current reference pattern:

```text
GrokifyQL text
  -> GrokifyQL AST
  -> application schema/catalog
  -> SystemForge authz.Authorizer
  -> GrokifyQL policy
  -> backend query provider
```

## Responsibilities

| Layer | Responsibility |
|-------|----------------|
| Query language | Parse text, build AST, reject unsupported syntax |
| Application catalog | Define queryable sources, datasets, fields, and capabilities |
| SystemForge | Check principal/action/resource authorization |
| SpiceDB | Store relationships and evaluate permissions |
| Backend provider | Compile and execute approved ASTs with service-owned limits |

SystemForge should not parse SQL-like text. It should receive typed resources
and actions from the application layer.

## Resource Modeling

For query engines, prefer concrete resources over attribute-only checks:

```text
analytics_source:<source_id>
analytics_dataset:<dataset_id>
analytics_field:<field_id>
```

Typical relationships:

```text
analytics_source:<source_id>#org@organization:<org_id>
analytics_dataset:<dataset_id>#source@analytics_source:<source_id>
analytics_field:<field_id>#dataset@analytics_dataset:<dataset_id>
```

Typical permissions:

```zed
definition analytics_dataset {
    relation source: analytics_source
    relation owner: principal

    permission manage = owner + source->manage
    permission read = manage + source->query
    permission list = read
    permission sort = read
}

definition analytics_field {
    relation dataset: analytics_dataset
    relation owner: principal

    permission manage = owner + dataset->manage
    permission read = manage + dataset->read
    permission list = read
    permission sort = read
}
```

## Action Mapping

A query language can map AST usage to SystemForge actions:

| Query usage | SystemForge action |
|-------------|--------------------|
| Read dataset/table | `read` |
| Select/project field | `read` |
| Filter field | `list` |
| Sort or group field | `sort` |
| Manage source/model | `manage` |

Applications can define stricter action names if needed, but they should stay
stable because saved queries and dashboard widgets may depend on them.

## GrokifyQL Adapter

GrokifyQL provides an optional adapter module:

```bash
go get github.com/grokify/grokifyql/authzsystemforge
```

The adapter accepts any SystemForge `authz.Authorizer`, including the SpiceDB
provider, and returns a `grokifyql.Policy`. Applications supply the resource
mapper that converts query entities and fields into their product-specific
SystemForge resources.

## Enforcement Points

Check policy at every API boundary that accepts user-authored query text:

- when saving a query or Question
- when executing an ad-hoc query
- when rendering a dashboard widget backed by a saved query
- when refreshing alerts or scheduled reports

Persisted ASTs and compiled metadata are useful for audit and performance, but
they do not replace current authorization checks. Relationships and roles can
change after a query is saved.
