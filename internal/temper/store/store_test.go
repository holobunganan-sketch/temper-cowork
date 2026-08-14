package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"reasonix/internal/temper/domain"
)

func openTestStore(t *testing.T) *Store {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "cowork")
	s, err := Open(dir)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

// fresh + migration:新库自动建 schema。
func TestOpenFreshCreatesSchema(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "cowork")
	s, err := Open(dir)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer s.Close()

	var version int
	if err := s.db.QueryRow(`SELECT COALESCE(MAX(version),0) FROM schema_migrations`).Scan(&version); err != nil {
		t.Fatalf("schema_migrations: %v", err)
	}
	if version != schemaVersion {
		t.Fatalf("schema version = %d, want %d", version, schemaVersion)
	}
	// 校验 8 张业务表存在。
	for _, table := range []string{"projects", "works", "work_events", "artifacts", "evidence", "decisions", "acceptance_results", "quality_runs"} {
		var name string
		err := s.db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name=?`, table).Scan(&name)
		if err != nil {
			t.Fatalf("table %s missing: %v", table, err)
		}
	}
}

// migration:已是最新版本时再 Open 不重复迁移。
func TestOpenIdempotentMigration(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "cowork")
	s1, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	s1.Close()
	s2, err := Open(dir)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer s2.Close()
	var version int
	if err := s2.db.QueryRow(`SELECT COALESCE(MAX(version),0) FROM schema_migrations`).Scan(&version); err != nil {
		t.Fatal(err)
	}
	if version != schemaVersion {
		t.Fatalf("version = %d, want %d", version, schemaVersion)
	}
}

// future schema:高于当前版本的表不应破坏打开。
func TestOpenWithFutureSchemaTable(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "cowork")
	s1, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s1.db.Exec(`CREATE TABLE future_table (id INTEGER PRIMARY KEY)`); err != nil {
		t.Fatal(err)
	}
	s1.Close()

	s2, err := Open(dir)
	if err != nil {
		t.Fatalf("reopen with future table: %v", err)
	}
	defer s2.Close()
}

// rollback:迁移失败时不应留下半成品版本记录。
func TestMigrationRollbackOnFailure(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "cowork")
	s, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	// 模拟未来版本记录存在(如降级后的旧库),再触发一个未知迁移版本。
	// applyMigration 对未知版本返回错误,事务应回滚,不写入版本记录。
	future := schemaVersion + 1
	tx, err := s.db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	if err := applyMigration(tx, future); err == nil {
		tx.Rollback()
		t.Fatal("expected applyMigration error for unknown version")
	}
	_ = tx.Rollback()

	// 再次迁移:版本记录不应包含 future。
	if err := s.migrate(); err != nil {
		t.Fatal(err)
	}
	var maxV int
	if err := s.db.QueryRow(`SELECT COALESCE(MAX(version),0) FROM schema_migrations`).Scan(&maxV); err != nil {
		t.Fatal(err)
	}
	if maxV != schemaVersion {
		t.Fatalf("schema_migrations max = %d, want %d", maxV, schemaVersion)
	}
}

// restart:关闭后重开,数据仍在。
func TestRestartPreservesData(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "cowork")
	s1, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	ws := filepath.Join(t.TempDir(), "workspace")
	if err := os.MkdirAll(ws, 0o755); err != nil {
		t.Fatal(err)
	}
	p, err := s1.CreateProject("Demo", ws)
	if err != nil {
		t.Fatal(err)
	}
	s1.Close()

	s2, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer s2.Close()
	got, err := s2.GetProject(p.ID)
	if err != nil {
		t.Fatalf("GetProject after restart: %v", err)
	}
	if got.Name != "Demo" {
		t.Fatalf("name = %q, want Demo", got.Name)
	}
}

// concurrency:并发创建互不干扰。
func TestConcurrentCreateProjects(t *testing.T) {
	s := openTestStore(t)
	done := make(chan error, 8)
	for i := 0; i < 8; i++ {
		go func(n int) {
			ws := filepath.Join(t.TempDir(), "ws", string(rune('a'+n)))
			if err := os.MkdirAll(ws, 0o755); err != nil {
				done <- err
				return
			}
			_, err := s.CreateProject(string(rune('A'+n)), ws)
			done <- err
		}(i)
	}
	for i := 0; i < 8; i++ {
		if err := <-done; err != nil {
			t.Fatalf("concurrent create: %v", err)
		}
	}
	projects, err := s.ListProjects()
	if err != nil {
		t.Fatal(err)
	}
	if len(projects) != 8 {
		t.Fatalf("projects = %d, want 8", len(projects))
	}
}

// Project CRUD 全流程。
func TestProjectCRUD(t *testing.T) {
	s := openTestStore(t)
	ws := filepath.Join(t.TempDir(), "workspace")
	if err := os.MkdirAll(ws, 0o755); err != nil {
		t.Fatal(err)
	}

	p, err := s.CreateProject("  My Project  ", ws)
	if err != nil {
		t.Fatalf("CreateProject: %v", err)
	}
	if p.Name != "My Project" {
		t.Fatalf("name = %q, want trimmed", p.Name)
	}
	if p.WorkspaceRoot == "" {
		t.Fatal("workspace_root must be canonical")
	}

	// 重复注册同一 workspace → ErrDuplicateWorkspace
	if _, err := s.CreateProject("Dup", ws); !errors.Is(err, ErrDuplicateWorkspace) {
		t.Fatalf("dup = %v, want ErrDuplicateWorkspace", err)
	}

	// Touch 更新 last_opened
	if err := s.TouchProject(p.ID); err != nil {
		t.Fatal(err)
	}
	got, err := s.GetProject(p.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.LastOpenedAt.IsZero() {
		t.Fatal("last_opened_at not updated")
	}

	// List
	projects, err := s.ListProjects()
	if err != nil {
		t.Fatal(err)
	}
	if len(projects) != 1 || projects[0].ID != p.ID {
		t.Fatalf("list = %+v", projects)
	}

	// Remove
	if err := s.RemoveProject(p.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := s.GetProject(p.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("after remove: %v, want ErrNotFound", err)
	}
	// workspace 目录仍在(Remove 不删目录)
	if _, err := os.Stat(ws); err != nil {
		t.Fatalf("workspace dir removed: %v", err)
	}
}

// domain 模型可序列化(JSON 契约稳定性)。
func TestDomainModelsSerialize(t *testing.T) {
	w := domain.Work{
		ID: "wk-1", ProjectID: "prj-1", Title: "T", Goal: "G", Status: domain.WorkDraft,
	}
	js, err := json.Marshal(w)
	if err != nil {
		t.Fatal(err)
	}
	if len(js) == 0 {
		t.Fatal("empty json")
	}
	if string(js)[0] != '{' {
		t.Fatalf("expected object json, got %s", js)
	}
}
