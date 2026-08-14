package store

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"reasonix/internal/temper/domain"
)

// ErrNotFound 表示记录不存在。
var ErrNotFound = errors.New("temper store: not found")

// ErrDuplicateWorkspace 表示同一 workspace 已注册。
var ErrDuplicateWorkspace = errors.New("temper store: workspace already registered")

// CreateProject 注册一个新项目。workspaceRoot 必须是存在的目录;
// 存储 canonical 绝对路径(Windows 折叠大小写)。
func (s *Store) CreateProject(name, workspaceRoot string) (*domain.Project, error) {
	canon, err := canonicalPath(workspaceRoot)
	if err != nil {
		return nil, err
	}
	if st, err := os.Stat(canon); err != nil || !st.IsDir() {
		return nil, fmt.Errorf("temper store: workspace root %q is not an accessible directory", workspaceRoot)
	}
	now := time.Now().UTC()
	p := &domain.Project{
		ID:            newID("prj"),
		Name:          strings.TrimSpace(name),
		WorkspaceRoot: canon,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	_, err = s.db.Exec(`INSERT INTO projects(id,name,workspace_root,created_at,updated_at,last_opened_at,archived)
		VALUES(?,?,?,?,?,?,0)`,
		p.ID, p.Name, p.WorkspaceRoot, fmtTime(p.CreatedAt), fmtTime(p.UpdatedAt), "")
	if isUniqueViolation(err) {
		return nil, ErrDuplicateWorkspace
	}
	if err != nil {
		return nil, fmt.Errorf("temper store: create project: %w", err)
	}
	return p, nil
}

// GetProject 按 ID 读取项目。
func (s *Store) GetProject(id string) (*domain.Project, error) {
	row := s.db.QueryRow(`SELECT id,name,workspace_root,created_at,updated_at,last_opened_at,archived
		FROM projects WHERE id=?`, id)
	return scanProject(row)
}

// ListProjects 列出全部项目(按最近打开倒序)。
func (s *Store) ListProjects() ([]*domain.Project, error) {
	rows, err := s.db.Query(`SELECT id,name,workspace_root,created_at,updated_at,last_opened_at,archived
		FROM projects ORDER BY last_opened_at DESC, updated_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.Project
	for rows.Next() {
		p, err := scanProject(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// TouchProject 更新项目最近打开时间。
func (s *Store) TouchProject(id string) error {
	_, err := s.db.Exec(`UPDATE projects SET last_opened_at=?, updated_at=? WHERE id=?`,
		fmtTime(time.Now().UTC()), fmtTime(time.Now().UTC()), id)
	return err
}

// RenameProject 更新项目显示名。
func (s *Store) RenameProject(id, name string) error {
	res, err := s.db.Exec(`UPDATE projects SET name=?, updated_at=? WHERE id=?`,
		strings.TrimSpace(name), fmtTime(time.Now().UTC()), id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// RemoveProject 从 Temper 移除项目注册(绝不删除 workspace 目录)。
func (s *Store) RemoveProject(id string) error {
	res, err := s.db.Exec(`DELETE FROM projects WHERE id=?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// ProjectByWorkspace 按 canonical workspace 路径查找项目(供去重)。
func (s *Store) ProjectByWorkspace(workspaceRoot string) (*domain.Project, error) {
	canon, err := canonicalPath(workspaceRoot)
	if err != nil {
		return nil, err
	}
	row := s.db.QueryRow(`SELECT id,name,workspace_root,created_at,updated_at,last_opened_at,archived
		FROM projects WHERE workspace_root=?`, canon)
	return scanProject(row)
}

type scanner interface{ Scan(dest ...any) error }

func scanProject(row scanner) (*domain.Project, error) {
	var p domain.Project
	var lastOpened, createdAt, updatedAt string
	var archived int
	err := row.Scan(&p.ID, &p.Name, &p.WorkspaceRoot, &createdAt, &updatedAt, &lastOpened, &archived)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	p.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	p.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	if lastOpened != "" {
		p.LastOpenedAt, _ = time.Parse(time.RFC3339, lastOpened)
	}
	p.Archived = archived != 0
	return &p, nil
}

// canonicalPath 返回规范化绝对路径(Windows 折叠大小写,便于去重)。
func canonicalPath(p string) (string, error) {
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", fmt.Errorf("temper store: resolve path: %w", err)
	}
	abs = filepath.Clean(abs)
	if resolved, rerr := filepath.EvalSymlinks(abs); rerr == nil {
		abs = resolved
	}
	if filepath.Separator == '\\' {
		abs = strings.ToLower(abs)
	}
	return abs, nil
}

// fmtTime 序列化 UTC 时间为 RFC3339(所有时间统一 UTC)。
func fmtTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}
