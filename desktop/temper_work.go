package main

import (
	"reasonix/internal/temper/domain"
)

// TemperWorkView 是前端可用的 Work 视图。
type TemperWorkView struct {
	ID                 string `json:"id"`
	ProjectID          string `json:"projectId"`
	Title              string `json:"title"`
	Goal               string `json:"goal"`
	Status             string `json:"status"`
	ReasonixSessionRef string `json:"reasonixSessionRef,omitempty"`
	ModelRef           string `json:"modelRef,omitempty"`
	TaskContract       string `json:"taskContract,omitempty"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
	FinalArtifactID    string `json:"finalArtifactId,omitempty"`
}

// ListTemperWorks 列出项目下的全部工作。
func (a *App) ListTemperWorks(projectID string) ([]TemperWorkView, error) {
	s, err := a.temperCoWork.get()
	if err != nil {
		return nil, err
	}
	works, err := s.ListWorksByProject(projectID)
	if err != nil {
		return nil, err
	}
	out := make([]TemperWorkView, 0, len(works))
	for _, w := range works {
		out = append(out, temperWorkView(w))
	}
	return out, nil
}

// CreateTemperWork 创建一份正式工作(draft)。
func (a *App) CreateTemperWork(projectID, title, goal, modelRef, qualityProfile string) (*TemperWorkView, error) {
	s, err := a.temperCoWork.get()
	if err != nil {
		return nil, err
	}
	w, err := s.CreateWork(projectID, title, goal, modelRef, qualityProfile)
	if err != nil {
		return nil, err
	}
	v := temperWorkView(w)
	return &v, nil
}

// UpdateTemperWorkStatus 更新工作状态。
func (a *App) UpdateTemperWorkStatus(workID, status string) error {
	s, err := a.temperCoWork.get()
	if err != nil {
		return err
	}
	return s.UpdateWorkStatus(workID, domain.WorkStatus(status))
}

func temperWorkView(w *domain.Work) TemperWorkView {
	v := TemperWorkView{
		ID:                 w.ID,
		ProjectID:          w.ProjectID,
		Title:              w.Title,
		Goal:               w.Goal,
		Status:             string(w.Status),
		ReasonixSessionRef: w.ReasonixSessionRef,
		ModelRef:           w.ModelRef,
		TaskContract:       w.TaskContract,
		CreatedAt:          w.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
		UpdatedAt:          w.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z"),
		FinalArtifactID:    w.FinalArtifactID,
	}
	return v
}
