package authz

import (
	"strings"
	"testing"
)

// testVocab is a configurable AppVocabulary for tests, modeled on the
// ActionForge facet pattern.
type testVocab struct {
	app    string
	roles  RolePermissions
	hier   RoleHierarchy
	scopes []string
	schema string
}

func (v testVocab) App() string              { return v.app }
func (v testVocab) Roles() RolePermissions   { return v.roles }
func (v testVocab) Hierarchy() RoleHierarchy { return v.hier }
func (v testVocab) Scopes() []string         { return v.scopes }
func (v testVocab) SpiceDBSchema() string    { return v.schema }

func actionforgeVocab() testVocab {
	return testVocab{
		app: "actionforge",
		roles: RolePermissions{
			"viewer":   {"workflow.read", "run.read"},
			"approver": {"workflow.read", "run.read", "run.approve"},
			"admin":    {"workflow.read", "workflow.write", "run.read", "run.approve"},
		},
		hier:   RoleHierarchy{"viewer": 20, "approver": 50, "admin": 80},
		scopes: []string{"actionforge:runs:read", "actionforge:runs:approve", "actionforge:admin"},
		schema: `
definition actionforge_org {
    relation org: organization
    relation approver: principal
    relation developer: principal

    permission workflow_read = org->view
    permission workflow_write = org->edit
    permission workflow_trigger = org->edit + developer
    permission run_approve = org->manage + approver
}

definition actionforge_workflow {
    relation app_org: actionforge_org
    permission read = app_org->workflow_read
    permission write = app_org->workflow_write
    permission trigger = app_org->workflow_trigger
}
`,
	}
}

func dashforgeVocab() testVocab {
	return testVocab{
		app:    "dashforge",
		roles:  RolePermissions{"viewer": {"dashboard.read"}},
		hier:   RoleHierarchy{"viewer": 20},
		scopes: []string{"dashforge:dashboards:read"},
		schema: `
definition dashforge_org {
    relation org: organization
    permission dashboard_read = org->view
}
`,
	}
}

func TestRegisterValidVocabularies(t *testing.T) {
	r := NewVocabularyRegistry()
	if err := r.Register(actionforgeVocab()); err != nil {
		t.Fatalf("register actionforge: %v", err)
	}
	if err := r.Register(dashforgeVocab()); err != nil {
		t.Fatalf("register dashforge: %v", err)
	}
	if apps := r.Apps(); len(apps) != 2 || apps[0] != "actionforge" || apps[1] != "dashforge" {
		t.Fatalf("Apps() = %v", apps)
	}
	if _, ok := r.Get("actionforge"); !ok {
		t.Fatal("Get(actionforge) missing")
	}
	scopes := r.AllScopes()
	if len(scopes) != 4 {
		t.Fatalf("AllScopes = %v", scopes)
	}
}

func TestRegisterRejectsDuplicateApp(t *testing.T) {
	r := NewVocabularyRegistry()
	if err := r.Register(actionforgeVocab()); err != nil {
		t.Fatal(err)
	}
	if err := r.Register(actionforgeVocab()); err == nil {
		t.Fatal("duplicate app registration should fail")
	}
}

func TestValidationFailures(t *testing.T) {
	base := actionforgeVocab()

	badName := base
	badName.app = "Action-Forge"
	if err := ValidateVocabulary(badName); err == nil {
		t.Error("bad app name should fail")
	}

	wrongPrefix := base
	wrongPrefix.scopes = []string{"dashforge:runs:read"}
	if err := ValidateVocabulary(wrongPrefix); err == nil {
		t.Error("scope with foreign app prefix should fail")
	}

	badFormat := base
	badFormat.scopes = []string{"actionforge:RunsRead"}
	if err := ValidateVocabulary(badFormat); err == nil {
		t.Error("malformed scope should fail")
	}

	dupScope := base
	dupScope.scopes = []string{"actionforge:runs:read", "actionforge:runs:read"}
	if err := ValidateVocabulary(dupScope); err == nil {
		t.Error("duplicate scope should fail")
	}

	missingHier := base
	missingHier.hier = RoleHierarchy{"viewer": 20} // approver/admin missing
	if err := ValidateVocabulary(missingHier); err == nil {
		t.Error("role missing from hierarchy should fail")
	}

	unprefixed := base
	unprefixed.schema = "definition organization { relation owner: principal }"
	if err := ValidateVocabulary(unprefixed); err == nil {
		t.Error("unprefixed definition (base collision) should fail")
	}
}

func TestComposeSpiceDBSchema(t *testing.T) {
	r := NewVocabularyRegistry()
	if err := r.Register(actionforgeVocab()); err != nil {
		t.Fatal(err)
	}
	if err := r.Register(dashforgeVocab()); err != nil {
		t.Fatal(err)
	}
	composed, err := r.ComposeSpiceDBSchema()
	if err != nil {
		t.Fatalf("compose: %v", err)
	}
	// Base first, then apps in registration order.
	for _, want := range []string{
		"definition principal", "definition organization", "definition platform",
		"definition actionforge_org", "definition actionforge_workflow",
		"definition dashforge_org",
	} {
		if !strings.Contains(composed, want) {
			t.Errorf("composed schema missing %q", want)
		}
	}
	if strings.Index(composed, "definition organization") > strings.Index(composed, "definition actionforge_org") {
		t.Error("base must precede app fragments")
	}
	// Exactly one organization definition (apps attach, never redefine).
	if n := strings.Count(composed, "definition organization "); n != 1 {
		t.Errorf("organization defined %d times, want 1", n)
	}
}

func TestComposeSkipsSchemalessApps(t *testing.T) {
	r := NewVocabularyRegistry()
	simpleOnly := dashforgeVocab()
	simpleOnly.app = "agentforge"
	simpleOnly.scopes = []string{"agentforge:agents:read"}
	simpleOnly.schema = ""
	if err := r.Register(simpleOnly); err != nil {
		t.Fatal(err)
	}
	composed, err := r.ComposeSpiceDBSchema()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(composed, "agentforge") {
		t.Error("schemaless app should contribute nothing to the composed schema")
	}
	if !strings.Contains(composed, "definition organization") {
		t.Error("base must always be present")
	}
}
