// 生产路径 E2E(Go 侧):验证 Temper CoWork 完整工作流通过真实 store +
// Temper 身份 env,覆盖 Master PHASE M 的生产路径(Project → Work →
// Evidence → Artifact → Validation → Completion → 重启恢复)。
//
// 运行:go test -run TestTemperProductionE2E ./internal/temper/e2e/
package e2e

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"reasonix/internal/temper/domain"
	"reasonix/internal/temper/store"
)

func sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// TestTemperProductionE2E 模拟真实用户路径:
// 1. 设置 Temper 身份 env(ApplyTemperIdentity 的等价物)
// 2. 创建 Project(真实 workspace)
// 3. 创建 Work + Task Contract
// 4. 登记 Evidence / Decision / Artifact
// 5. Validation + Quality
// 6. Completion Gate 通过 → completed
// 7. 关闭 store(模拟退出)
// 8. 重启 store → 全部数据恢复
func TestTemperProductionE2E(t *testing.T) {
	// 1. Temper 数据目录隔离
	home := t.TempDir()
	coWorkDir := filepath.Join(home, "cowork")
	ws := filepath.Join(t.TempDir(), "temp-workspace")
	if err := os.MkdirAll(ws, 0o755); err != nil {
		t.Fatal(err)
	}

	// 2-6. 第一次运行:完整工作流
	s1, err := store.Open(coWorkDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	proj, err := s1.CreateProject("E2E Project", ws)
	if err != nil {
		t.Fatalf("create project: %v", err)
	}
	work, err := s1.CreateWork(proj.ID, "E2E Work", "Produce a summary", "deepseek-r1", "strict")
	if err != nil {
		t.Fatalf("create work: %v", err)
	}
	if err := s1.SetWorkTaskContract(work.ID, "# Task Contract\n## Request\nProduce a summary"); err != nil {
		t.Fatal(err)
	}
	if err := s1.SetWorkSessionRef(work.ID, "sess-e2e"); err != nil {
		t.Fatal(err)
	}
	if err := s1.UpdateWorkStatus(work.ID, domain.WorkRunning); err != nil {
		t.Fatal(err)
	}

	// Evidence + Decision
	if err := s1.CreateEvidence(work.ID, domain.Evidence{Summary: "verified output", SourceType: "test", SourceRef: "e2e_test.go"}); err != nil {
		t.Fatal(err)
	}
	if err := s1.CreateDecision(work.ID, domain.Decision{Decision: "use sqlite", Rationale: "no cgo"}); err != nil {
		t.Fatal(err)
	}

	// Artifact(真实文件 + hash)
	content := []byte("# E2E Summary\n")
	relPath := "summary.md"
	if err := os.WriteFile(filepath.Join(ws, relPath), content, 0o644); err != nil {
		t.Fatal(err)
	}
	art := domain.Artifact{
		ProjectID:    proj.ID,
		WorkID:       work.ID,
		RelativePath: relPath,
		Kind:         domain.ArtifactMarkdown,
		Title:        "Summary",
		SHA256:       sha256Hex(content),
		Size:         int64(len(content)),
	}
	if err := s1.RegisterArtifact(work.ID, art); err != nil {
		t.Fatal(err)
	}
	arts, err := s1.ListArtifacts(work.ID)
	if err != nil || len(arts) != 1 {
		t.Fatalf("artifacts: %d err=%v", len(arts), err)
	}
	if err := s1.SetWorkFinalArtifact(work.ID, arts[0].ID); err != nil {
		t.Fatal(err)
	}

	// Validation + Quality
	if err := s1.RecordAcceptance(work.ID, domain.AcceptanceResult{Criterion: "c1", Status: domain.AcceptancePass}); err != nil {
		t.Fatal(err)
	}
	if err := s1.RecordQualityRun(work.ID, domain.QualityRun{GateType: "validation", Passed: true, Summary: "ok"}); err != nil {
		t.Fatal(err)
	}

	// Completion Gate 检查(生产路径)
	if err := s1.UpdateWorkStatus(work.ID, domain.WorkCompleted); err != nil {
		t.Fatal(err)
	}
	got, err := s1.GetWork(work.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != domain.WorkCompleted || got.CompletedAt == nil || got.FinalArtifactID == "" {
		t.Fatalf("work not properly completed: %+v", got)
	}

	// 7. 关闭(模拟退出)
	if err := s1.Close(); err != nil {
		t.Fatal(err)
	}

	// 8. 重启 → 恢复
	s2, err := store.Open(coWorkDir)
	if err != nil {
		t.Fatalf("reopen store: %v", err)
	}
	defer s2.Close()

	proj2, err := s2.GetProject(proj.ID)
	if err != nil || proj2.Name != "E2E Project" {
		t.Fatalf("project after restart: %+v err=%v", proj2, err)
	}
	work2, err := s2.GetWork(work.ID)
	if err != nil {
		t.Fatalf("work after restart: %v", err)
	}
	if work2.Status != domain.WorkCompleted || work2.ReasonixSessionRef != "sess-e2e" {
		t.Fatalf("work state lost after restart: %+v", work2)
	}
	arts2, err := s2.ListArtifacts(work.ID)
	if err != nil || len(arts2) != 1 || arts2[0].SHA256 == "" {
		t.Fatalf("artifacts after restart: %d err=%v", len(arts2), err)
	}
	ev2, err := s2.ListEvidence(work.ID)
	if err != nil || len(ev2) != 1 {
		t.Fatalf("evidence after restart: %d err=%v", len(ev2), err)
	}
	dec2, err := s2.ListDecisions(work.ID)
	if err != nil || len(dec2) != 1 {
		t.Fatalf("decisions after restart: %d err=%v", len(dec2), err)
	}
	acc2, err := s2.ListAcceptance(work.ID)
	if err != nil || len(acc2) != 1 {
		t.Fatalf("acceptance after restart: %d err=%v", len(acc2), err)
	}
	// workspace 文件仍在
	if _, err := os.Stat(filepath.Join(ws, relPath)); err != nil {
		t.Fatalf("workspace file lost: %v", err)
	}
}
