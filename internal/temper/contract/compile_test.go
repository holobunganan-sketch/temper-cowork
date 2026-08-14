package contract

import (
	"strings"
	"testing"
)

func TestCompileFullForm(t *testing.T) {
	out := Compile(WorkForm{
		Project:     "docs",
		Goal:        "Write a release summary for v0.3.0",
		Materials:   []string{"docs/RELEASE.md", "CHANGELOG.md"},
		Deliverable: "docs/release-summary.md",
		Audience:    "maintainers",
		Constraints: []string{"keep under 500 words", "no marketing fluff"},
		AcceptanceCriteria: []string{
			"summary covers all v0.3.0 features",
			"links to RELEASE.md",
		},
		SourcePolicy: "cite sources inline",
		PausePolicy:  "pause before deleting anything",
		Quality:      "strict",
		Model:        "deepseek-r1",
	})
	for _, want := range []string{
		"# Task Contract",
		"## Context",
		"Project: docs",
		"- docs/RELEASE.md",
		"## Request",
		"Write a release summary for v0.3.0",
		"## Output format",
		"Deliverable: docs/release-summary.md",
		"## Constraints",
		"- keep under 500 words",
		"Source policy: cite sources inline",
		"Quality profile: strict",
		"Model: deepseek-r1",
		"## Acceptance criteria",
		"- [ ] AC1: summary covers all v0.3.0 features",
		"- [ ] AC2: links to RELEASE.md",
		"## Pause policy",
		"pause before deleting anything",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("compiled contract missing %q\n---\n%s", want, out)
		}
	}
}

func TestCompileEmptyFormHasDefaults(t *testing.T) {
	out := Compile(WorkForm{})
	for _, want := range []string{
		"(no goal provided)",
		"(unspecified)",
		"(no explicit criteria",
		"pause on any request for user decision",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("default contract missing %q\n---\n%s", want, out)
		}
	}
}

func TestCompileSkipsBlankMaterials(t *testing.T) {
	out := Compile(WorkForm{
		Goal:      "g",
		Materials: []string{"", "a.md", "  "},
	})
	if strings.Count(out, "- a.md") != 1 {
		t.Fatalf("expected exactly one material entry:\n%s", out)
	}
	if strings.Contains(out, "- \n") {
		t.Fatalf("blank material leaked:\n%s", out)
	}
}
