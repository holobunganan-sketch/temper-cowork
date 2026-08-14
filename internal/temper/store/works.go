package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"reasonix/internal/temper/domain"
)

// CreateWork 在项目下创建一份正式工作(draft 状态)。
func (s *Store) CreateWork(projectID, title, goal, modelRef, qualityProfile string) (*domain.Work, error) {
	now := time.Now().UTC()
	w := &domain.Work{
		ID:             newID("wk"),
		ProjectID:      projectID,
		Title:          strings.TrimSpace(title),
		Goal:           strings.TrimSpace(goal),
		Status:         domain.WorkDraft,
		ModelRef:       modelRef,
		QualityProfile: qualityProfile,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	_, err := s.db.Exec(`INSERT INTO works(id,project_id,title,goal,status,reasonix_session_ref,model_ref,quality_profile,task_contract,created_at,updated_at,started_at,completed_at,final_artifact_id)
		VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		w.ID, w.ProjectID, w.Title, w.Goal, string(w.Status), "", w.ModelRef, w.QualityProfile, "",
		fmtTime(w.CreatedAt), fmtTime(w.UpdatedAt), "", "", "")
	if err != nil {
		return nil, fmt.Errorf("temper store: create work: %w", err)
	}
	return w, nil
}

// GetWork 按 ID 读取工作。
func (s *Store) GetWork(id string) (*domain.Work, error) {
	row := s.db.QueryRow(`SELECT id,project_id,title,goal,status,reasonix_session_ref,model_ref,quality_profile,task_contract,created_at,updated_at,started_at,completed_at,final_artifact_id
		FROM works WHERE id=?`, id)
	return scanWork(row)
}

// ListWorksByProject 列出项目的全部工作(按更新时间倒序)。
func (s *Store) ListWorksByProject(projectID string) ([]*domain.Work, error) {
	rows, err := s.db.Query(`SELECT id,project_id,title,goal,status,reasonix_session_ref,model_ref,quality_profile,task_contract,created_at,updated_at,started_at,completed_at,final_artifact_id
		FROM works WHERE project_id=? ORDER BY updated_at DESC`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.Work
	for rows.Next() {
		w, err := scanWork(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, w)
	}
	return out, rows.Err()
}

// UpdateWorkStatus 更新工作状态并记录事件。
func (s *Store) UpdateWorkStatus(id string, status domain.WorkStatus) error {
	now := time.Now().UTC()
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	var started, completed string
	switch status {
	case domain.WorkRunning:
		started = fmtTime(now)
	case domain.WorkCompleted, domain.WorkFailed, domain.WorkCancelled:
		completed = fmtTime(now)
	}
	res, err := tx.Exec(`UPDATE works SET status=?, updated_at=?, started_at=CASE WHEN started_at='' THEN ? ELSE started_at END, completed_at=CASE WHEN ?<>'' THEN ? ELSE completed_at END WHERE id=?`,
		string(status), fmtTime(now), started, completed, completed, id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		_ = tx.Rollback()
		return ErrNotFound
	}
	if _, err := tx.Exec(`INSERT INTO work_events(work_id,event_type,detail,created_at) VALUES(?,?,?,?)`,
		id, "status_change", string(status), fmtTime(now)); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

// SetWorkSessionRef 关联 Reasonix Session。
func (s *Store) SetWorkSessionRef(id, sessionRef string) error {
	return s.updateWorkField(id, "reasonix_session_ref", sessionRef)
}

// SetWorkTaskContract 写入编译后的 Task Contract(供 Completion Gate 检查)。
func (s *Store) SetWorkTaskContract(id, contract string) error {
	return s.updateWorkField(id, "task_contract", contract)
}

// SetWorkFinalArtifact 标记最终交付物。
func (s *Store) SetWorkFinalArtifact(id, artifactID string) error {
	return s.updateWorkField(id, "final_artifact_id", artifactID)
}

// ListWorkEvents 列出工作事件。
func (s *Store) ListWorkEvents(workID string) ([]domain.WorkEvent, error) {
	rows, err := s.db.Query(`SELECT id,work_id,event_type,detail,created_at FROM work_events WHERE work_id=? ORDER BY id`, workID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.WorkEvent
	for rows.Next() {
		var e domain.WorkEvent
		var createdAt string
		if err := rows.Scan(&e.ID, &e.WorkID, &e.EventType, &e.Detail, &createdAt); err != nil {
			return nil, err
		}
		e.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		out = append(out, e)
	}
	return out, rows.Err()
}

func (s *Store) updateWorkField(id, field, value string) error {
	q := fmt.Sprintf(`UPDATE works SET %s=?, updated_at=? WHERE id=?`, field)
	res, err := s.db.Exec(q, value, fmtTime(time.Now().UTC()), id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func scanWork(row scanner) (*domain.Work, error) {
	var w domain.Work
	var createdAt, updatedAt, startedAt, completedAt string
	err := row.Scan(&w.ID, &w.ProjectID, &w.Title, &w.Goal, &w.Status,
		&w.ReasonixSessionRef, &w.ModelRef, &w.QualityProfile, &w.TaskContract,
		&createdAt, &updatedAt, &startedAt, &completedAt, &w.FinalArtifactID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	w.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	w.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	if startedAt != "" {
		if t, err := time.Parse(time.RFC3339, startedAt); err == nil {
			w.StartedAt = &t
		}
	}
	if completedAt != "" {
		if t, err := time.Parse(time.RFC3339, completedAt); err == nil {
			w.CompletedAt = &t
		}
	}
	return &w, nil
}
