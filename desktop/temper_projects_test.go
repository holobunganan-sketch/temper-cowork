package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestMirrorProjectToStore 验证 workspace 注册会镜像到 CoWork store。
func TestMirrorProjectToStore(t *testing.T) {
	home := t.TempDir()
	t.Setenv("TEMPER_HOME", home)
	unsetAll(t, "REASONIX_HOME", "REASONIX_STATE_HOME", "REASONIX_CACHE_HOME")
	ApplyTemperIdentity()

	a := &App{}
	t.Cleanup(a.temperCoWork.close)
	ws := filepath.Join(t.TempDir(), "workspace")
	if err := os.MkdirAll(ws, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := a.mirrorProjectToStore(ws, "My Project"); err != nil {
		t.Fatalf("mirror: %v", err)
	}

	projects, err := a.ListTemperProjects()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(projects) != 1 {
		t.Fatalf("projects = %d, want 1", len(projects))
	}
	if projects[0].Name != "My Project" {
		t.Fatalf("name = %q, want 'My Project'", projects[0].Name)
	}
	if projects[0].WorkspaceRoot == "" {
		t.Fatal("workspace_root must be set")
	}
}

// TestMirrorProjectIdempotent 验证重复注册同一 workspace 幂等(不产生重复)。
func TestMirrorProjectIdempotent(t *testing.T) {
	home := t.TempDir()
	t.Setenv("TEMPER_HOME", home)
	unsetAll(t, "REASONIX_HOME", "REASONIX_STATE_HOME", "REASONIX_CACHE_HOME")
	ApplyTemperIdentity()

	a := &App{}
	t.Cleanup(a.temperCoWork.close)
	ws := filepath.Join(t.TempDir(), "workspace")
	if err := os.MkdirAll(ws, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := a.mirrorProjectToStore(ws, "A"); err != nil {
		t.Fatal(err)
	}
	if err := a.mirrorProjectToStore(ws, "B"); err != nil {
		t.Fatal(err)
	}
	projects, err := a.ListTemperProjects()
	if err != nil {
		t.Fatal(err)
	}
	if len(projects) != 1 {
		t.Fatalf("projects = %d, want 1 (idempotent)", len(projects))
	}
	if projects[0].Name != "B" {
		t.Fatalf("name = %q, want 'B' (rename applied)", projects[0].Name)
	}
}

// TestRemoveTemperProject 验证移除不删除 workspace 目录。
func TestRemoveTemperProject(t *testing.T) {
	home := t.TempDir()
	t.Setenv("TEMPER_HOME", home)
	unsetAll(t, "REASONIX_HOME", "REASONIX_STATE_HOME", "REASONIX_CACHE_HOME")
	ApplyTemperIdentity()

	a := &App{}
	t.Cleanup(a.temperCoWork.close)
	ws := filepath.Join(t.TempDir(), "workspace")
	if err := os.MkdirAll(ws, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := a.mirrorProjectToStore(ws, "P"); err != nil {
		t.Fatal(err)
	}

	// 未在 Reasonix 注册时,移除只清 store 记录。
	if err := a.RemoveTemperProject(ws); err != nil {
		t.Fatalf("remove: %v", err)
	}
	projects, err := a.ListTemperProjects()
	if err != nil {
		t.Fatal(err)
	}
	if len(projects) != 0 {
		t.Fatalf("projects after remove = %d, want 0", len(projects))
	}
	if _, err := os.Stat(ws); err != nil {
		t.Fatalf("workspace dir was deleted: %v", err)
	}
}
