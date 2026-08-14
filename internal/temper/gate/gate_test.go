package gate

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"reasonix/internal/temper/domain"
	"reasonix/internal/temper/store"
)

func openStore(t *testing.T) *store.Store {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "cowork")
	s, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func setupWork(t *testing.T, s *store.Store) (string, string) {
	t.Helper()
	ws := filepath.Join(t.TempDir(), "workspace")
	if err := os.MkdirAll(ws, 0o755); err != nil {
		t.Fatal(err)
	}
	proj, err := s.CreateProject("P", ws)
	if err != nil {
		t.Fatal(err)
	}
	w, err := s.CreateWork(proj.ID, "W", "g", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if err := s.SetWorkTaskContract(w.ID, "# Task Contract"); err != nil {
		t.Fatal(err)
	}
	return w.ID, ws
}

func TestCompletionGatePasses(t *testing.T) {
	s := openStore(t)
	workID, ws := setupWork(t, s)

	// 真实文件 + artifact
	content := []byte("hello temper")
	full := filepath.Join(ws, "out.md")
	if err := os.WriteFile(full, content, 0o644); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(content)
	art := domain.Artifact{
		ProjectID:    "prj",
		WorkID:       workID,
		RelativePath: "out.md",
		Kind:         domain.ArtifactMarkdown,
		Title:        "out",
		SHA256:       hex.EncodeToString(sum[:]),
		Size:         int64(len(content)),
	}
	if err := s.RegisterArtifact(workID, art); err != nil {
		t.Fatal(err)
	}
	// 从 store 拿回带 ID 的记录
	arts, err := s.ListArtifacts(workID)
	if err != nil || len(arts) != 1 {
		t.Fatalf("artifacts = %d, err=%v", len(arts), err)
	}
	if err := s.SetWorkFinalArtifact(workID, arts[0].ID); err != nil {
		t.Fatal(err)
	}
	if err := s.RecordAcceptance(workID, domain.AcceptanceResult{Criterion: "c1", Status: domain.AcceptancePass}); err != nil {
		t.Fatal(err)
	}
	if err := s.UpdateWorkStatus(workID, domain.WorkRunning); err != nil {
		t.Fatal(err)
	}

	report, err := CompletionGate(s, workID, ws)
	if err != nil {
		t.Fatal(err)
	}
	if !report.Passed {
		t.Fatalf("gate should pass: %+v", report.Checks)
	}
}

func TestCompletionGateFailsWithoutFinalArtifact(t *testing.T) {
	s := openStore(t)
	workID, ws := setupWork(t, s)
	report, err := CompletionGate(s, workID, ws)
	if err != nil {
		t.Fatal(err)
	}
	if report.Passed {
		t.Fatal("gate should fail without final artifact")
	}
	for _, c := range report.Checks {
		if c.Item == "final_artifact" && c.Passed {
			t.Fatalf("final_artifact should not pass: %+v", c)
		}
	}
}

func TestCompletionGateDetectsStaleHash(t *testing.T) {
	s := openStore(t)
	workID, ws := setupWork(t, s)

	full := filepath.Join(ws, "out.md")
	if err := os.WriteFile(full, []byte("version-1"), 0o644); err != nil {
		t.Fatal(err)
	}
	art := domain.Artifact{
		ProjectID:    "prj",
		WorkID:       workID,
		RelativePath: "out.md",
		Kind:         domain.ArtifactMarkdown,
		Title:        "out",
		SHA256:       "deadbeef", // 错误 hash
		Size:         9,
	}
	if err := s.RegisterArtifact(workID, art); err != nil {
		t.Fatal(err)
	}
	arts, _ := s.ListArtifacts(workID)
	if err := s.SetWorkFinalArtifact(workID, arts[0].ID); err != nil {
		t.Fatal(err)
	}
	if err := s.RecordAcceptance(workID, domain.AcceptanceResult{Criterion: "c1", Status: domain.AcceptancePass}); err != nil {
		t.Fatal(err)
	}
	if err := s.UpdateWorkStatus(workID, domain.WorkRunning); err != nil {
		t.Fatal(err)
	}

	report, err := CompletionGate(s, workID, ws)
	if err != nil {
		t.Fatal(err)
	}
	if report.Passed {
		t.Fatal("gate should fail on stale hash")
	}
	found := false
	for _, c := range report.Checks {
		if c.Item == "hash_current" {
			found = true
			if c.Passed {
				t.Fatal("hash_current should not pass")
			}
		}
	}
	if !found {
		t.Fatal("hash_current check missing")
	}
}

func TestCompletionGateFailsOnValidationFail(t *testing.T) {
	s := openStore(t)
	workID, ws := setupWork(t, s)

	full := filepath.Join(ws, "out.md")
	if err := os.WriteFile(full, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	art := domain.Artifact{ProjectID: "prj", WorkID: workID, RelativePath: "out.md",
		Kind: domain.ArtifactMarkdown, Title: "out", SHA256: "abc", Size: 1}
	if err := s.RegisterArtifact(workID, art); err != nil {
		t.Fatal(err)
	}
	arts, _ := s.ListArtifacts(workID)
	if err := s.SetWorkFinalArtifact(workID, arts[0].ID); err != nil {
		t.Fatal(err)
	}
	// 一条 fail 验收
	if err := s.RecordAcceptance(workID, domain.AcceptanceResult{Criterion: "c1", Status: domain.AcceptanceFail}); err != nil {
		t.Fatal(err)
	}

	report, err := CompletionGate(s, workID, ws)
	if err != nil {
		t.Fatal(err)
	}
	if report.Passed {
		t.Fatal("gate should fail on acceptance fail")
	}
}

func TestQualityGateSummarizes(t *testing.T) {
	s := openStore(t)
	workID, _ := setupWork(t, s)
	if err := s.RecordQualityRun(workID, domain.QualityRun{GateType: "validation", Passed: true, Summary: "ok"}); err != nil {
		t.Fatal(err)
	}
	if err := s.RecordQualityRun(workID, domain.QualityRun{GateType: "review", Passed: false, Summary: "needs work"}); err != nil {
		t.Fatal(err)
	}
	sum, err := QualityGate(s, workID)
	if err != nil {
		t.Fatal(err)
	}
	if len(sum.Validation) != 1 || len(sum.Reviews) != 1 {
		t.Fatalf("validation=%d reviews=%d", len(sum.Validation), len(sum.Reviews))
	}
	if sum.Passed {
		t.Fatal("quality gate should fail when a run did not pass")
	}
}
