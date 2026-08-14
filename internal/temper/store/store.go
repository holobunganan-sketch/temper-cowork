// Package store 是 Temper CoWork 的 SQLite 持久层。
//
// 业务 DB:%APPDATA%\Temper\cowork\temper.db(由 ApplyTemperIdentity 注入的
// REASONIX_STATE_HOME 决定,通常为 %APPDATA%\Temper\cowork)。
//
// 要求:WAL、foreign_keys ON、busy timeout、事务、UTC、canonical paths、
// migration、crash safe。
package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

// schemaVersion 是当前 schema 迁移版本。每次 schema 变更递增,并在
// migrate 中追加迁移步骤。
const schemaVersion = 1

// Store 是 CoWork 数据库的封装。
type Store struct {
	db *sql.DB
}

// Open 打开(或创建)位于 dir 下的 temper.db,执行迁移并返回 Store。
// dir 不存在时会创建。
func Open(dir string) (*Store, error) {
	if dir == "" {
		return nil, fmt.Errorf("temper store: empty directory")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("temper store: create dir: %w", err)
	}
	path := filepath.Join(dir, "temper.db")
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("temper store: open: %w", err)
	}
	// WAL + busy timeout + foreign keys(每个连接都要,用 DSN 与 pragma)。
	db.SetMaxOpenConns(1) // SQLite 写串行化,单连接避免锁竞争
	pragmas := []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA foreign_keys=ON",
		"PRAGMA busy_timeout=5000",
		"PRAGMA synchronous=NORMAL",
	}
	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			db.Close()
			return nil, fmt.Errorf("temper store: %s: %w", p, err)
		}
	}
	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("temper store: migrate: %w", err)
	}
	return s, nil
}

// Close 关闭数据库。
func (s *Store) Close() error {
	return s.db.Close()
}

// migrate 执行 schema 迁移,记录到 schema_migrations。
func (s *Store) migrate() error {
	return s.migrateFrom(0)
}

// migrateFrom 从给定基线版本(含)开始迁移到最新。测试用它模拟从未来
// 版本回退/迁移失败场景。
func (s *Store) migrateFrom(baseline int) error {
	if _, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version    INTEGER PRIMARY KEY,
		applied_at TEXT NOT NULL
	)`); err != nil {
		return err
	}
	var current int
	if err := s.db.QueryRow(`SELECT COALESCE(MAX(version),0) FROM schema_migrations`).Scan(&current); err != nil {
		return err
	}
	for v := max(current, baseline) + 1; v <= schemaVersion; v++ {
		tx, err := s.db.Begin()
		if err != nil {
			return err
		}
		if err := applyMigration(tx, v); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("migration %d: %w", v, err)
		}
		now := time.Now().UTC().Format(time.RFC3339)
		if _, err := tx.Exec(`INSERT INTO schema_migrations(version, applied_at) VALUES(?,?)`, v, now); err != nil {
			_ = tx.Rollback()
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

// applyMigration 按版本应用 schema 变更。迁移必须是幂等的可追加步骤。
func applyMigration(tx *sql.Tx, version int) error {
	switch version {
	case 1:
		return migrationV1(tx)
	default:
		return fmt.Errorf("unknown migration version %d", version)
	}
}

func migrationV1(tx *sql.Tx) error {
	stmts := []string{
		`CREATE TABLE projects (
			id             TEXT PRIMARY KEY,
			name           TEXT NOT NULL,
			workspace_root TEXT NOT NULL UNIQUE,
			created_at     TEXT NOT NULL,
			updated_at     TEXT NOT NULL,
			last_opened_at TEXT NOT NULL DEFAULT '',
			archived       INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE works (
			id                  TEXT PRIMARY KEY,
			project_id          TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
			title               TEXT NOT NULL,
			goal                TEXT NOT NULL,
			status              TEXT NOT NULL,
			reasonix_session_ref TEXT NOT NULL DEFAULT '',
			model_ref           TEXT NOT NULL DEFAULT '',
			quality_profile     TEXT NOT NULL DEFAULT '',
			task_contract       TEXT NOT NULL DEFAULT '',
			created_at          TEXT NOT NULL,
			updated_at          TEXT NOT NULL,
			started_at          TEXT NOT NULL DEFAULT '',
			completed_at        TEXT NOT NULL DEFAULT '',
			final_artifact_id   TEXT NOT NULL DEFAULT ''
		)`,
		`CREATE INDEX idx_works_project ON works(project_id)`,
		`CREATE TABLE work_events (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			work_id    TEXT NOT NULL REFERENCES works(id) ON DELETE CASCADE,
			event_type TEXT NOT NULL,
			detail     TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL
		)`,
		`CREATE INDEX idx_work_events_work ON work_events(work_id)`,
		`CREATE TABLE evidence (
			id          TEXT PRIMARY KEY,
			work_id     TEXT NOT NULL REFERENCES works(id) ON DELETE CASCADE,
			summary     TEXT NOT NULL,
			source_type TEXT NOT NULL,
			source_ref  TEXT NOT NULL,
			supports    TEXT NOT NULL DEFAULT '',
			timestamp   TEXT NOT NULL
		)`,
		`CREATE INDEX idx_evidence_work ON evidence(work_id)`,
		`CREATE TABLE decisions (
			id           TEXT PRIMARY KEY,
			work_id      TEXT NOT NULL REFERENCES works(id) ON DELETE CASCADE,
			decision     TEXT NOT NULL,
			rationale    TEXT NOT NULL,
			alternatives TEXT NOT NULL DEFAULT '',
			evidence_ids TEXT NOT NULL DEFAULT '',
			timestamp    TEXT NOT NULL
		)`,
		`CREATE INDEX idx_decisions_work ON decisions(work_id)`,
		`CREATE TABLE artifacts (
			id            TEXT PRIMARY KEY,
			project_id    TEXT NOT NULL,
			work_id       TEXT NOT NULL REFERENCES works(id) ON DELETE CASCADE,
			relative_path TEXT NOT NULL,
			kind          TEXT NOT NULL,
			title         TEXT NOT NULL,
			description   TEXT NOT NULL DEFAULT '',
			sha256        TEXT NOT NULL,
			size          INTEGER NOT NULL,
			validation    TEXT NOT NULL DEFAULT '',
			is_final      INTEGER NOT NULL DEFAULT 0,
			created_at    TEXT NOT NULL,
			updated_at    TEXT NOT NULL
		)`,
		`CREATE INDEX idx_artifacts_work ON artifacts(work_id)`,
		`CREATE TABLE acceptance_results (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			work_id     TEXT NOT NULL REFERENCES works(id) ON DELETE CASCADE,
			criterion   TEXT NOT NULL,
			status      TEXT NOT NULL,
			evidence_ref TEXT NOT NULL DEFAULT '',
			evaluated_at TEXT NOT NULL
		)`,
		`CREATE INDEX idx_acceptance_work ON acceptance_results(work_id)`,
		`CREATE TABLE quality_runs (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			work_id    TEXT NOT NULL REFERENCES works(id) ON DELETE CASCADE,
			gate_type  TEXT NOT NULL,
			passed     INTEGER NOT NULL,
			summary    TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL
		)`,
		`CREATE INDEX idx_quality_work ON quality_runs(work_id)`,
	}
	for _, stmt := range stmts {
		if _, err := tx.Exec(stmt); err != nil {
			return err
		}
	}
	return nil
}
