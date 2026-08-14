package store

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"reasonix/internal/temper/domain"
)

func TestWorkLifecycle(t *testing.T) {
	s := openTestStore(t)

	ws := filepath.Join(t.TempDir(), "workspace")
	if err := os.MkdirAll(ws, 0o755); err != nil {
		t.Fatal(err)
	}
	proj, err := s.CreateProject("P", ws)
	if err != nil {
		t.Fatal(err)
	}

	w, err := s.CreateWork(proj.ID, "  Release summary  ", "Write it", "deepseek-r1", "strict")
	if err != nil {
		t.Fatalf("CreateWork: %v", err)
	}
	if w.Title != "Release summary" {
		t.Fatalf("title = %q, want trimmed", w.Title)
	}
	if w.Status != domain.WorkDraft {
		t.Fatalf("status = %q, want draft", w.Status)
	}

	// 状态流转:ready → running → completed
	if err := s.UpdateWorkStatus(w.ID, domain.WorkReady); err != nil {
		t.Fatal(err)
	}
	if err := s.UpdateWorkStatus(w.ID, domain.WorkRunning); err != nil {
		t.Fatal(err)
	}
	got, err := s.GetWork(w.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != domain.WorkRunning || got.StartedAt == nil {
		t.Fatalf("running work = %+v, want started", got)
	}

	if err := s.SetWorkSessionRef(w.ID, "sess-1"); err != nil {
		t.Fatal(err)
	}
	if err := s.SetWorkTaskContract(w.ID, "# Task Contract\n..."); err != nil {
		t.Fatal(err)
	}
	if err := s.SetWorkFinalArtifact(w.ID, "art-1"); err != nil {
		t.Fatal(err)
	}
	if err := s.UpdateWorkStatus(w.ID, domain.WorkCompleted); err != nil {
		t.Fatal(err)
	}

	got, err = s.GetWork(w.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != domain.WorkCompleted || got.CompletedAt == nil {
		t.Fatalf("completed work = %+v", got)
	}
	if got.ReasonixSessionRef != "sess-1" || got.TaskContract == "" || got.FinalArtifactID != "art-1" {
		t.Fatalf("work fields not persisted: %+v", got)
	}

	// 事件已记录(ready/running/completed 3 个状态变更)
	events, err := s.ListWorkEvents(w.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 3 {
		t.Fatalf("events = %d, want 3 (ready/running/completed)", len(events))
	}

	// 项目内列表
	works, err := s.ListWorksByProject(proj.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(works) != 1 || works[0].ID != w.ID {
		t.Fatalf("works = %+v", works)
	}

	// 未找到
	if _, err := s.GetWork("wk-missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing work: %v, want ErrNotFound", err)
	}
}

func TestUpdateWorkStatusUnknownReturnsNotFound(t *testing.T) {
	s := openTestStore(t)
	if err := s.UpdateWorkStatus("wk-nope", domain.WorkRunning); !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}
