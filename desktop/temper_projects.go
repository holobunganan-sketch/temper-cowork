package main

import (
	"errors"

	"reasonix/internal/temper/store"
)

// TemperProjectView 是前端可用的项目视图(CoWork store 元数据 + workspace)。
type TemperProjectView struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	WorkspaceRoot string `json:"workspaceRoot"`
	CreatedAt     string `json:"createdAt"`
	LastOpenedAt  string `json:"lastOpenedAt"`
}

// ListTemperProjects 列出 CoWork store 中的项目(前端 Projects 页)。
func (a *App) ListTemperProjects() ([]TemperProjectView, error) {
	s, err := a.temperCoWork.get()
	if err != nil {
		return nil, err
	}
	projects, err := s.ListProjects()
	if err != nil {
		return nil, err
	}
	out := make([]TemperProjectView, 0, len(projects))
	for _, p := range projects {
		out = append(out, TemperProjectView{
			ID:            p.ID,
			Name:          p.Name,
			WorkspaceRoot: p.WorkspaceRoot,
			CreatedAt:     p.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
			LastOpenedAt:  p.LastOpenedAt.UTC().Format("2006-01-02T15:04:05Z"),
		})
	}
	return out, nil
}

// AddTemperProject 注册一个 workspace 为 Temper 项目:先复用 Reasonix 的
// addProject(建立 workspace 会话注册),再镜像到 CoWork store。
func (a *App) AddTemperProject(workspaceRoot, name string) error {
	root := normalizeProjectRoot(workspaceRoot)
	if root == "" {
		return errors.New("project root is required")
	}
	if err := addProject(root, name); err != nil {
		return err
	}
	return a.mirrorProjectToStore(root, name)
}

// RemoveTemperProject 从 Temper 移除项目注册(不删除 workspace 目录):
// 先移除 Reasonix 注册,再从 CoWork store 移除。
func (a *App) RemoveTemperProject(workspaceRoot string) error {
	root := normalizeProjectRoot(workspaceRoot)
	if err := removeProject(root); err != nil {
		// 未注册的 workspace 不算错误——store 可能已有该记录。
		if !errors.Is(err, store.ErrNotFound) {
			return err
		}
	}
	s, err := a.temperCoWork.get()
	if err != nil {
		return err
	}
	proj, err := s.ProjectByWorkspace(root)
	if errors.Is(err, store.ErrNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	return s.RemoveProject(proj.ID)
}

// mirrorProjectToStore 把 workspace 镜像为 CoWork store 的 Project 记录
// (幂等:已存在则更新名称)。
func (a *App) mirrorProjectToStore(root, name string) error {
	s, err := a.temperCoWork.get()
	if err != nil {
		return err
	}
	existing, err := s.ProjectByWorkspace(root)
	if err == nil {
		if name != "" {
			_ = s.RenameProject(existing.ID, name)
		}
		return nil
	}
	if !errors.Is(err, store.ErrNotFound) {
		return err
	}
	_, err = s.CreateProject(name, root)
	return err
}
