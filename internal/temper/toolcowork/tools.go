package toolcowork

import (
	"context"
	"encoding/json"
	"fmt"

	"reasonix/internal/temper/domain"
)

// temperRecordEvidence 登记一条证据。
type temperRecordEvidence struct{ toolBase }

func init() {
	Register(temperRecordEvidence{toolBase{
		name:        "temper_record_evidence",
		description: "Record an evidence item supporting a conclusion in the current Temper work.",
		schema: mustSchema(map[string]any{
			"type": "object",
			"properties": map[string]any{
				"work_id":     map[string]any{"type": "string", "description": "Target work ID"},
				"summary":     map[string]any{"type": "string"},
				"source_type": map[string]any{"type": "string", "enum": []string{"file", "command", "tool_result", "observation", "test"}},
				"source_ref":  map[string]any{"type": "string", "description": "Path, command, or reference"},
				"supports":    map[string]any{"type": "string"},
			},
			"required": []string{"work_id", "summary", "source_type", "source_ref"},
		}),
	}})
}

type evidenceArgs struct {
	WorkID     string `json:"work_id"`
	Summary    string `json:"summary"`
	SourceType string `json:"source_type"`
	SourceRef  string `json:"source_ref"`
	Supports   string `json:"supports"`
}

func (t temperRecordEvidence) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	var a evidenceArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return "", fmt.Errorf("temper_record_evidence: %w", err)
	}
	s, err := requireStore(ctx)
	if err != nil {
		return "", err
	}
	e := domain.Evidence{
		Summary:    a.Summary,
		SourceType: a.SourceType,
		SourceRef:  a.SourceRef,
		Supports:   a.Supports,
		Timestamp:  nowUTC(),
	}
	if err := s.CreateEvidence(a.WorkID, e); err != nil {
		return "", fmt.Errorf("temper_record_evidence: %w", err)
	}
	return fmt.Sprintf("evidence recorded for work %s", a.WorkID), nil
}

// temperRecordDecision 记录一条决策。
type temperRecordDecision struct{ toolBase }

func init() {
	Register(temperRecordDecision{toolBase{
		name:        "temper_record_decision",
		description: "Record a decision made during the current Temper work.",
		schema: mustSchema(map[string]any{
			"type": "object",
			"properties": map[string]any{
				"work_id":      map[string]any{"type": "string"},
				"decision":     map[string]any{"type": "string"},
				"rationale":    map[string]any{"type": "string"},
				"alternatives": map[string]any{"type": "string"},
				"evidence_ids": map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
			},
			"required": []string{"work_id", "decision", "rationale"},
		}),
	}})
}

type decisionArgs struct {
	WorkID       string   `json:"work_id"`
	Decision     string   `json:"decision"`
	Rationale    string   `json:"rationale"`
	Alternatives string   `json:"alternatives"`
	EvidenceIDs  []string `json:"evidence_ids"`
}

func (t temperRecordDecision) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	var a decisionArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return "", fmt.Errorf("temper_record_decision: %w", err)
	}
	s, err := requireStore(ctx)
	if err != nil {
		return "", err
	}
	d := domain.Decision{
		Decision:     a.Decision,
		Rationale:    a.Rationale,
		Alternatives: a.Alternatives,
		EvidenceIDs:  a.EvidenceIDs,
		Timestamp:    nowUTC(),
	}
	if err := s.CreateDecision(a.WorkID, d); err != nil {
		return "", fmt.Errorf("temper_record_decision: %w", err)
	}
	return fmt.Sprintf("decision recorded for work %s", a.WorkID), nil
}
