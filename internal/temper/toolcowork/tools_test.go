package toolcowork

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"reasonix/internal/temper/domain"
	"reasonix/internal/temper/store"
	"reasonix/internal/tool"
)

// testStore 打开临时 CoWork store 并注入 context。
func testCtx(t *testing.T) (context.Context, *store.Store, *domain.Work) {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "cowork")
	s, err := store.Open(dir)
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	t.Cleanup(func() { s.Close() })

	ws := filepath.Join(t.TempDir(), "workspace")
	if err := os.MkdirAll(ws, 0o755); err != nil {
		t.Fatal(err)
	}
	proj, err := s.CreateProject("P", ws)
	if err != nil {
		t.Fatal(err)
	}
	w, err := s.CreateWork(proj.ID, "W", "goal", "model", "profile")
	if err != nil {
		t.Fatal(err)
	}
	// 工具执行所需的 task contract
	if err := s.SetWorkTaskContract(w.ID, "# Task Contract\n"); err != nil {
		t.Fatal(err)
	}
	return WithStore(context.Background(), s), s, w
}

func runTool(t *testing.T, name string, args map[string]any, ctx context.Context) (string, error) {
	t.Helper()
	tool, ok := lookupTool(name)
	if !ok {
		t.Fatalf("tool %s not registered", name)
	}
	raw, err := json.Marshal(args)
	if err != nil {
		t.Fatal(err)
	}
	return tool.Execute(ctx, raw)
}

// lookupTool 从注册表取工具(测试辅助)。
func lookupTool(name string) (tool.Tool, bool) {
	for _, t := range All() {
		if t.Name() == name {
			return t, true
		}
	}
	return nil, false
}

func TestAllTemperToolsRegistered(t *testing.T) {
	names := map[string]bool{}
	for _, t := range All() {
		names[t.Name()] = true
	}
	for _, want := range []string{
		"temper_record_evidence",
		"temper_record_decision",
		"temper_register_artifact",
		"temper_set_final_artifact",
		"temper_report_validation",
		"temper_complete_work",
	} {
		if !names[want] {
			t.Fatalf("tool %s not registered", want)
		}
	}
}

func TestRecordEvidenceAndDecision(t *testing.T) {
	ctx, s, w := testCtx(t)
	if _, err := runTool(t, "temper_record_evidence", map[string]any{
		"work_id": w.ID, "summary": "verified", "source_type": "test", "source_ref": "a_test.go",
	}, ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := runTool(t, "temper_record_decision", map[string]any{
		"work_id": w.ID, "decision": "use sqlite", "rationale": "no cgo",
	}, ctx); err != nil {
		t.Fatal(err)
	}
	ev, err := s.ListEvidence(w.ID)
	if err != nil || len(ev) != 1 {
		t.Fatalf("evidence = %d, err=%v", len(ev), err)
	}
	dec, err := s.ListDecisions(w.ID)
	if err != nil || len(dec) != 1 {
		t.Fatalf("decisions = %d, err=%v", len(dec), err)
	}
}

func TestCompleteWorkRequiresFinalArtifact(t *testing.T) {
	ctx, _, w := testCtx(t)
	// 无 final artifact:必须失败。
	if _, err := runTool(t, "temper_complete_work", map[string]any{"work_id": w.ID}, ctx); err == nil {
		t.Fatal("expected completion to fail without final artifact")
	}
}

func TestFullWorkCompletionFlow(t *testing.T) {
	ctx, s, w := testCtx(t)
	if _, err := runTool(t, "temper_register_artifact", map[string]any{
		"work_id": w.ID, "project_id": "prj-x", "relative_path": "docs/out.md",
		"kind": "md", "title": "out", "sha256": "abc", "size": 12,
	}, ctx); err != nil {
		t.Fatal(err)
	}
	arts, err := s.ListArtifacts(w.ID)
	if err != nil || len(arts) != 1 {
		t.Fatalf("artifacts = %d, err=%v", len(arts), err)
	}
	if _, err := runTool(t, "temper_set_final_artifact", map[string]any{
		"work_id": w.ID, "artifact_id": arts[0].ID,
	}, ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := runTool(t, "temper_report_validation", map[string]any{
		"work_id": w.ID, "criterion": "c1", "status": "pass",
	}, ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := runTool(t, "temper_complete_work", map[string]any{"work_id": w.ID}, ctx); err != nil {
		t.Fatalf("completion failed: %v", err)
	}
	got, err := s.GetWork(w.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != domain.WorkCompleted {
		t.Fatalf("status = %q, want completed", got.Status)
	}
	if got.FinalArtifactID != arts[0].ID {
		t.Fatalf("final artifact = %q, want %q", got.FinalArtifactID, arts[0].ID)
	}
}

func TestToolsFailWithoutStore(t *testing.T) {
	// 无 store 注入:工具必须报错。
	if _, err := runTool(t, "temper_record_evidence", map[string]any{
		"work_id": "x", "summary": "s", "source_type": "file", "source_ref": "f",
	}, context.Background()); err == nil {
		t.Fatal("expected error without store in context")
	}
}
