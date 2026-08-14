package toolcowork

import (
	"context"
	"encoding/json"
	"fmt"

	"reasonix/internal/temper/domain"
)

// temperRegisterArtifact 登记一个交付物。
type temperRegisterArtifact struct{ toolBase }

func init() {
	Register(temperRegisterArtifact{toolBase{
		name:        "temper_register_artifact",
		description: "Register a deliverable artifact (workspace file) in the current Temper work.",
		schema: mustSchema(map[string]any{
			"type": "object",
			"properties": map[string]any{
				"work_id":       map[string]any{"type": "string"},
				"project_id":    map[string]any{"type": "string"},
				"relative_path": map[string]any{"type": "string", "description": "Path relative to the workspace root"},
				"kind":          map[string]any{"type": "string", "enum": []string{"md", "txt", "json", "csv", "html", "svg", "png", "jpeg", "source"}},
				"title":         map[string]any{"type": "string"},
				"description":   map[string]any{"type": "string"},
				"sha256":        map[string]any{"type": "string"},
				"size":          map[string]any{"type": "integer"},
			},
			"required": []string{"work_id", "project_id", "relative_path", "kind", "title", "sha256", "size"},
		}),
	}})
}

type artifactArgs struct {
	WorkID       string `json:"work_id"`
	ProjectID    string `json:"project_id"`
	RelativePath string `json:"relative_path"`
	Kind         string `json:"kind"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	SHA256       string `json:"sha256"`
	Size         int64  `json:"size"`
}

func (t temperRegisterArtifact) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	var a artifactArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return "", fmt.Errorf("temper_register_artifact: %w", err)
	}
	s, err := requireStore(ctx)
	if err != nil {
		return "", err
	}
	art := domain.Artifact{
		ProjectID:    a.ProjectID,
		WorkID:       a.WorkID,
		RelativePath: a.RelativePath,
		Kind:         domain.ArtifactKind(a.Kind),
		Title:        a.Title,
		Description:  a.Description,
		SHA256:       a.SHA256,
		Size:         a.Size,
		CreatedAt:    nowUTC(),
		UpdatedAt:    nowUTC(),
	}
	if err := s.RegisterArtifact(a.WorkID, art); err != nil {
		return "", fmt.Errorf("temper_register_artifact: %w", err)
	}
	return fmt.Sprintf("artifact registered for work %s at %s", a.WorkID, a.RelativePath), nil
}

// temperSetFinalArtifact 标记最终交付物。
type temperSetFinalArtifact struct{ toolBase }

func init() {
	Register(temperSetFinalArtifact{toolBase{
		name:        "temper_set_final_artifact",
		description: "Mark an artifact as the final deliverable of the work.",
		schema: mustSchema(map[string]any{
			"type": "object",
			"properties": map[string]any{
				"work_id":     map[string]any{"type": "string"},
				"artifact_id": map[string]any{"type": "string", "description": "ID returned by temper_register_artifact"},
			},
			"required": []string{"work_id", "artifact_id"},
		}),
	}})
}

type finalArtifactArgs struct {
	WorkID     string `json:"work_id"`
	ArtifactID string `json:"artifact_id"`
}

func (t temperSetFinalArtifact) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	var a finalArtifactArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return "", fmt.Errorf("temper_set_final_artifact: %w", err)
	}
	s, err := requireStore(ctx)
	if err != nil {
		return "", err
	}
	if err := s.SetWorkFinalArtifact(a.WorkID, a.ArtifactID); err != nil {
		return "", fmt.Errorf("temper_set_final_artifact: %w", err)
	}
	return fmt.Sprintf("final artifact %s set for work %s", a.ArtifactID, a.WorkID), nil
}

// temperReportValidation 报告一次验证结果。
type temperReportValidation struct{ toolBase }

func init() {
	Register(temperReportValidation{toolBase{
		name:        "temper_report_validation",
		description: "Report the result of validating the work deliverable against acceptance criteria.",
		schema: mustSchema(map[string]any{
			"type": "object",
			"properties": map[string]any{
				"work_id":      map[string]any{"type": "string"},
				"criterion":    map[string]any{"type": "string"},
				"status":       map[string]any{"type": "string", "enum": []string{"pass", "fail", "uncertain"}},
				"evidence_ref": map[string]any{"type": "string"},
			},
			"required": []string{"work_id", "criterion", "status"},
		}),
	}})
}

type validationArgs struct {
	WorkID      string `json:"work_id"`
	Criterion   string `json:"criterion"`
	Status      string `json:"status"`
	EvidenceRef string `json:"evidence_ref"`
}

func (t temperReportValidation) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	var a validationArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return "", fmt.Errorf("temper_report_validation: %w", err)
	}
	s, err := requireStore(ctx)
	if err != nil {
		return "", err
	}
	ar := domain.AcceptanceResult{
		Criterion:   a.Criterion,
		Status:      domain.AcceptanceStatus(a.Status),
		EvidenceRef: a.EvidenceRef,
		EvaluatedAt: nowUTC(),
	}
	if err := s.RecordAcceptance(a.WorkID, ar); err != nil {
		return "", fmt.Errorf("temper_report_validation: %w", err)
	}
	return fmt.Sprintf("validation %s recorded for work %s criterion %q", a.Status, a.WorkID, a.Criterion), nil
}

// temperCompleteWork 完成工作(触发 Completion Gate)。
type temperCompleteWork struct{ toolBase }

func init() {
	Register(temperCompleteWork{toolBase{
		name:        "temper_complete_work",
		description: "Mark the work as completed. The host Completion Gate validates the task contract, final artifact, and acceptance results before allowing completion.",
		schema: mustSchema(map[string]any{
			"type": "object",
			"properties": map[string]any{
				"work_id": map[string]any{"type": "string"},
			},
			"required": []string{"work_id"},
		}),
	}})
}

type completeWorkArgs struct {
	WorkID string `json:"work_id"`
}

func (t temperCompleteWork) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	var a completeWorkArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return "", fmt.Errorf("temper_complete_work: %w", err)
	}
	s, err := requireStore(ctx)
	if err != nil {
		return "", err
	}
	w, err := s.GetWork(a.WorkID)
	if err != nil {
		return "", fmt.Errorf("temper_complete_work: %w", err)
	}
	// Host Completion Gate:必须满足才允许完成。
	if w.TaskContract == "" {
		return "", fmt.Errorf("temper_complete_work: task contract is missing; cannot complete")
	}
	if w.FinalArtifactID == "" {
		return "", fmt.Errorf("temper_complete_work: final artifact is not set; cannot complete")
	}
	if err := s.UpdateWorkStatus(a.WorkID, domain.WorkCompleted); err != nil {
		return "", fmt.Errorf("temper_complete_work: %w", err)
	}
	return fmt.Sprintf("work %s completed", a.WorkID), nil
}
