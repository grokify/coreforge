package authzguardsql

import (
	"context"
	"testing"

	"github.com/grokify/guardsql"
	"github.com/plexusone/systemforge/authz"
)

type mockAuthorizer struct {
	allowed map[string]bool
}

func (m mockAuthorizer) Can(_ context.Context, _ authz.Principal, action authz.Action, resource authz.Resource) (bool, error) {
	return m.allowed[string(resource.Type)+":"+string(action)+":"+attr(resource, "entity")+":"+attr(resource, "field")], nil
}

func (m mockAuthorizer) CanAll(ctx context.Context, principal authz.Principal, actions []authz.Action, resource authz.Resource) (bool, error) {
	for _, action := range actions {
		ok, err := m.Can(ctx, principal, action, resource)
		if err != nil || !ok {
			return ok, err
		}
	}
	return true, nil
}

func (m mockAuthorizer) CanAny(ctx context.Context, principal authz.Principal, actions []authz.Action, resource authz.Resource) (bool, error) {
	for _, action := range actions {
		ok, err := m.Can(ctx, principal, action, resource)
		if err != nil || ok {
			return ok, err
		}
	}
	return false, nil
}

func (m mockAuthorizer) Filter(ctx context.Context, principal authz.Principal, action authz.Action, resources []authz.Resource) ([]authz.Resource, error) {
	var out []authz.Resource
	for _, resource := range resources {
		ok, err := m.Can(ctx, principal, action, resource)
		if err != nil {
			return nil, err
		}
		if ok {
			out = append(out, resource)
		}
	}
	return out, nil
}

func TestPolicyBuilderBuildsGuardSQLPolicy(t *testing.T) {
	builder := PolicyBuilder{
		Authorizer: mockAuthorizer{allowed: map[string]bool{
			"guardsql_entity:read:items:":      true,
			"guardsql_field:read:items:name":   true,
			"guardsql_field:list:items:status": true,
			"guardsql_field:sort:items:score":  true,
		}},
		Schema: guardsql.Schema{Entities: map[string]guardsql.Entity{
			"items": {
				Fields: map[string]guardsql.Field{
					"name":   {Type: guardsql.FieldString},
					"status": {Type: guardsql.FieldString},
					"score":  {Type: guardsql.FieldNumber},
				},
			},
			"hidden": {
				Fields: map[string]guardsql.Field{"secret": {Type: guardsql.FieldString}},
			},
		}},
		RequireLimit: true,
		MaxLimit:     100,
	}
	policy, err := builder.Build(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(policy.AllowedEntities) != 1 || policy.AllowedEntities[0] != "items" {
		t.Fatalf("allowed entities = %#v, want items", policy.AllowedEntities)
	}
	if !policy.Fields["items"]["name"].Selectable {
		t.Fatalf("name should be selectable: %#v", policy.Fields)
	}
	if !policy.Fields["items"]["status"].Filterable {
		t.Fatalf("status should be filterable: %#v", policy.Fields)
	}
	if !policy.Fields["items"]["score"].Sortable {
		t.Fatalf("score should be sortable: %#v", policy.Fields)
	}
}

func attr(resource authz.Resource, key string) string {
	value, _ := resource.Attributes[key].(string)
	return value
}
